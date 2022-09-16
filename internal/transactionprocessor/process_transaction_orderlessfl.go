package transactionprocessor

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"gitlab.lrz.de/orderless/orderlessfl/internal/blockprocessor"
	"gitlab.lrz.de/orderless/orderlessfl/internal/logger"
	"gitlab.lrz.de/orderless/orderlessfl/internal/profiling"
	"gitlab.lrz.de/orderless/orderlessfl/internal/transaction"
	protos "gitlab.lrz.de/orderless/orderlessfl/protos/goprotos"
	"reflect"
	"sort"
)

type OrderlessFLNodeTransactionResponseSubscriber struct {
	stream   protos.TransactionService_SubscribeNodeTransactionsServer
	finished chan<- bool
}

func (p *Processor) signProposalResponseOrderlessFL(proposalResponse *protos.ProposalResponse) (*protos.ProposalResponse, error) {
	sort.Slice(proposalResponse.ReadWriteSet.ReadKeys.ReadKeys, func(i, j int) bool {
		return proposalResponse.ReadWriteSet.ReadKeys.ReadKeys[i].Key < proposalResponse.ReadWriteSet.ReadKeys.ReadKeys[j].Key
	})
	sort.Slice(proposalResponse.ReadWriteSet.WriteKeyValues.WriteKeyValues, func(i, j int) bool {
		return proposalResponse.ReadWriteSet.WriteKeyValues.WriteKeyValues[i].Key < proposalResponse.ReadWriteSet.WriteKeyValues.WriteKeyValues[j].Key
	})
	marshalledReadWriteSet, err := proto.Marshal(proposalResponse.ReadWriteSet)
	if err != nil {
		return p.makeFailedProposal(proposalResponse.ProposalId), nil
	}
	proposalResponse.NodeSignature = p.signer.Sign(marshalledReadWriteSet)
	return proposalResponse, nil
}

func (p *Processor) preProcessValidateReadWriteSetOrderlessFL(tx *protos.Transaction) error {
	sort.Slice(tx.ReadWriteSet.ReadKeys.ReadKeys, func(i, j int) bool {
		return tx.ReadWriteSet.ReadKeys.ReadKeys[i].Key < tx.ReadWriteSet.ReadKeys.ReadKeys[j].Key
	})
	sort.Slice(tx.ReadWriteSet.WriteKeyValues.WriteKeyValues, func(i, j int) bool {
		return tx.ReadWriteSet.WriteKeyValues.WriteKeyValues[i].Key < tx.ReadWriteSet.WriteKeyValues.WriteKeyValues[j].Key
	})
	txDigest, err := proto.Marshal(tx.ReadWriteSet)
	if err != nil {
		return err
	}
	passedSignature := int32(0)
	for node, nodeSign := range tx.NodeSignatures {
		err = p.signer.Verify(node, txDigest, nodeSign)
		if err == nil {
			passedSignature++
		}
	}
	if passedSignature < tx.EndorsementPolicy {
		return errors.New("failed signature validation")
	}
	err = p.signer.Verify(tx.ClientId, txDigest, tx.ClientSignature)
	if err != nil {
		return err
	}
	return nil
}

func (p *Processor) ProcessProposalOrderlessFLStream(proposal *protos.ProposalRequest) {
	p.txJournal.AddProposalToQueue(proposal)
}

func (p *Processor) runProposalQueueProcessingOrderlessFL() {
	for {
		proposals := <-p.txJournal.DequeuedProposalsChan
		go p.processDequeuedProposalsOrderlessFL(proposals)
	}
}

func (p *Processor) processDequeuedProposalsOrderlessFL(proposals *transaction.DequeuedProposals) {
	for _, proposal := range proposals.DequeuedProposals {
		go p.processProposalOrderlessFL(proposal)
	}
}

func (p *Processor) processProposalOrderlessFL(proposal *protos.ProposalRequest) {
	response, err := p.executeContract(proposal)
	if err != nil {
		p.sendProposalResponseToSubscriber(proposal.ClientId, response)
		return
	}
	response, err = p.signProposalResponseOrderlessFL(response)
	if err != nil {
		p.sendProposalResponseToSubscriber(proposal.ClientId, response)
		return
	}
	p.sendProposalResponseToSubscriber(proposal.ClientId, response)
}

func (p *Processor) ProcessTransactionOrderlessFLStream(tx *protos.Transaction) {
	p.txJournal.AddTransactionToQueue(tx)
}

func (p *Processor) runTransactionQueueProcessingOrderlessFL() {
	for {
		transactions := <-p.txJournal.DequeuedTransactionsChan
		go p.processDequeuedTransactionsOrderlessFL(transactions)
	}
}

func (p *Processor) processDequeuedTransactionsOrderlessFL(transactions *transaction.DequeuedTransactions) {
	for _, tx := range transactions.DequeuedTransactions {
		go p.processTransactionOrderlessFL(tx)
	}
}

func (p *Processor) processTransactionOrderlessFL(tx *protos.Transaction) {
	if !tx.FromNode && !tx.IsFederated {
		if err := p.preProcessValidateReadWriteSetOrderlessFL(tx); err != nil {
			tx.Status = protos.TransactionStatus_FAILED_SIGNATURE_VALIDATION
		}
	}
	shouldBeAdded := p.txJournal.AddTransactionToJournalIfNotExist(tx)
	if shouldBeAdded {
		if tx.IsFederated {
			p.processTransactionInBlock(tx, nil)
		} else {
			p.blockProcessor.TransactionChan <- tx
		}
	} else {
		if !tx.FromNode {
			p.finalizeAlreadySubmittedTransactionFromClientOrderlessFL(tx)
		}
	}
}

func (p *Processor) processTransactionFromOtherNodesOrderlessFL(txs []*protos.Transaction) {
	for _, tx := range txs {
		shouldBeAdded := p.txJournal.AddTransactionToJournalNodeIfNotExist(tx)
		if shouldBeAdded {
			alreadyInJournal := p.txJournal.IsTransactionInJournal(tx)
			if !alreadyInJournal {
				tx.FromNode = true
				p.txJournal.AddTransactionToQueue(tx)
			}
		}
	}
}

func (p *Processor) runTransactionProcessorOrderlessFL() {
	for {
		block := <-p.blockProcessor.MinedBlock
		go p.processBlockOrderlessFL(block)
	}
}

func (p *Processor) processBlockOrderlessFL(block *blockprocessor.MinedBlock) {
	for _, tx := range block.Transactions {
		go p.processTransactionInBlock(tx, block)
	}
}

func (p *Processor) processTransactionInBlock(tx *protos.Transaction, block *blockprocessor.MinedBlock) {
	var blockHeader []byte
	if tx.IsFederated {
		blockHeader = []byte("federated")
	} else {
		blockHeader = block.Block.ThisBlockHash
	}
	if tx.Status == protos.TransactionStatus_FAILED_SIGNATURE_VALIDATION {
		p.finalizeTransactionResponseOrderlessFL(tx.ClientId,
			p.makeFailedTransactionResponse(tx.TransactionId, tx.Status, blockHeader), tx)
	} else {
		response, err := p.commitTransactionOrderlessFL(tx)
		if err == nil {
			response.BlockHeader = blockHeader
			p.finalizeTransactionResponseOrderlessFL(tx.ClientId, response, tx)
		} else {
			p.finalizeTransactionResponseOrderlessFL(tx.ClientId,
				p.makeFailedTransactionResponse(tx.TransactionId, protos.TransactionStatus_FAILED_DATABASE, blockHeader), tx)
		}
	}
}

func (p *Processor) commitTransactionOrderlessFL(tx *protos.Transaction) (*protos.TransactionResponse, error) {
	if len(tx.FederatedModelUpdate) > 0 {
		if err := p.federatedCRDTManager.ApplyModelUpdates(tx.ModelId, tx.FederatedModelUpdate); err != nil {
			logger.InfoLogger.Println(err)
			return nil, err
		}
		//if err := p.addTransactionToDatabaseOrderlessFL(tx); err != nil {
		//	logger.InfoLogger.Println(err)
		//	return nil, err
		//}
	}
	return p.makeSuccessTransactionResponse(tx.TransactionId, []byte{}), nil
}

func (p *Processor) addTransactionToDatabaseOrderlessFL(tx *protos.Transaction) error {
	dbOp := p.sharedShimResources.DBConnections[tx.ContractName]
	var err error
	if err = dbOp.PutKeyValueNoVersion(tx.ModelId+"-"+tx.TransactionId, tx.FederatedModelUpdate); err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	return err
}

func (p *Processor) finalizeAlreadySubmittedTransactionFromClientOrderlessFL(transaction *protos.Transaction) {
	var response *protos.TransactionResponse
	if transaction.Status == protos.TransactionStatus_FAILED_SIGNATURE_VALIDATION {
		response = p.makeFailedTransactionResponse(transaction.TransactionId, transaction.Status, []byte("BlockHeader"))
	} else {
		response = p.makeSuccessTransactionResponse(transaction.TransactionId, []byte("BlockHeader"))
	}
	response.NodeSignature = p.signer.Sign(response.BlockHeader)
	p.sendTransactionResponseToSubscriber(transaction.ClientId, response)
}

func (p *Processor) finalizeTransactionResponseOrderlessFL(clientId string, txResponse *protos.TransactionResponse, transaction *protos.Transaction) {
	if !transaction.FromNode {
		if !transaction.IsFederated {
			txResponse.NodeSignature = p.signer.Sign(txResponse.BlockHeader)
		}
		p.sendTransactionResponseToSubscriber(clientId, txResponse)
	}
}

func (p *Processor) sendTransactionBatchToSubscribedNodeOrderlessFL(transactions []*protos.Transaction) {
	gossips := &protos.NodeTransactionResponse{
		Transaction: transactions,
	}
	p.nodeSubscribersLock.RLock()
	nodes := reflect.ValueOf(p.OrderlessFLNodeTransactionResponseSubscriber).MapKeys()
	if profiling.IsBandwidthProfiling {
		_ = proto.Size(gossips)
	}
	p.nodeSubscribersLock.RUnlock()
	for _, node := range nodes {
		nodeId := node.String()
		p.nodeSubscribersLock.RLock()
		streamer, ok := p.OrderlessFLNodeTransactionResponseSubscriber[nodeId]
		p.nodeSubscribersLock.RUnlock()
		if !ok {
			logger.ErrorLogger.Println("Node was not found in the subscribers streams.", nodeId)
			return
		}
		if err := streamer.stream.Send(gossips); err != nil {
			streamer.finished <- true
			logger.ErrorLogger.Println("Could not send the response to the node " + nodeId)
			p.nodeSubscribersLock.Lock()
			delete(p.OrderlessFLNodeTransactionResponseSubscriber, nodeId)
			p.nodeSubscribersLock.Unlock()
		}
	}
}

func (p *Processor) NodeTransactionResponseSubscriptionOrderlessFL(subscription *protos.TransactionResponseEventSubscription,
	stream protos.TransactionService_SubscribeNodeTransactionsServer) error {
	finished := make(chan bool)
	p.nodeSubscribersLock.Lock()
	p.OrderlessFLNodeTransactionResponseSubscriber[subscription.ComponentId] = &OrderlessFLNodeTransactionResponseSubscriber{
		stream:   stream,
		finished: finished,
	}
	p.nodeSubscribersLock.Unlock()
	cntx := stream.Context()
	for {
		select {
		case <-finished:
			return nil
		case <-cntx.Done():
			return nil
		}
	}
}

func (p *Processor) GetLatestFederatedModel(modelUpdateRequest *protos.ModelUpdateRequest) (*protos.ModelUpdateRequest, error) {
	latestModelWeight, clientLatestClock, err := p.federatedCRDTManager.GetModelCRDTByte(modelUpdateRequest)
	if err != nil {
		return nil, err
	}
	return &protos.ModelUpdateRequest{
		ModelId:           modelUpdateRequest.ModelId,
		LayersWeight:      latestModelWeight,
		ClientId:          modelUpdateRequest.ClientId,
		ClientLatestClock: clientLatestClock,
	}, nil
}
