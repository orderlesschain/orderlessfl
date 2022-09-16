package transaction

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"gitlab.lrz.de/orderless/orderlessfl/internal/config"
	"gitlab.lrz.de/orderless/orderlessfl/internal/contract/contractinterface"
	"gitlab.lrz.de/orderless/orderlessfl/internal/customcrypto/hasher"
	"gitlab.lrz.de/orderless/orderlessfl/internal/customcrypto/signer"
	protos "gitlab.lrz.de/orderless/orderlessfl/protos/goprotos"
	"sort"
	"strconv"
)

type ClientTransaction struct {
	TransactionId     string
	ProposalResponses map[string]*protos.ProposalResponse
	Status            protos.TransactionStatus
	signer            *signer.Signer
}

func NewClientTransaction(signer *signer.Signer) *ClientTransaction {
	return &ClientTransaction{
		TransactionId:     uuid.NewString(),
		Status:            protos.TransactionStatus_RUNNING,
		ProposalResponses: map[string]*protos.ProposalResponse{},
		signer:            signer,
	}
}

func (t *ClientTransaction) MakeProposalRequestBenchmarkExecutor(counter int, baseOptions *contractinterface.BaseContractOptions) (*protos.ProposalRequest, []byte) {
	proposalParams := baseOptions.BenchmarkFunction(&contractinterface.BenchmarkFunctionOptions{
		Counter:                     counter,
		CurrentClientPseudoId:       baseOptions.CurrentClientPseudoId,
		BenchmarkUtils:              baseOptions.BenchmarkUtils,
		TotalTransactions:           baseOptions.TotalTransactions,
		SingleFunctionCounter:       baseOptions.SingleFunctionCounter,
		CrdtOperationPerObjectCount: baseOptions.Bconfig.CrdtOperationPerObjectCount,
		CrdtObjectType:              baseOptions.Bconfig.CrdtObjectType,
		CrdtObjectCount:             baseOptions.Bconfig.CrdtObjectCount,
		NumberOfKeys:                int(baseOptions.Bconfig.NumberOfKeys),
		NumberOfKeysSecond:          int(baseOptions.Bconfig.NumberOfKeysSecond),
		FederatedConnPool:           baseOptions.FederatedConnPool,
		NodesConnPool:               baseOptions.NodesConnectionPool,
	})
	if proposalParams.ShouldBeSubmitted {
		return &protos.ProposalRequest{
			TargetSystem:         baseOptions.Bconfig.TargetSystem,
			ProposalId:           t.TransactionId,
			ClientId:             config.Config.UUID,
			ContractName:         baseOptions.Bconfig.ContractName,
			MethodName:           proposalParams.ContractMethodName,
			MethodParams:         proposalParams.Outputs,
			WriteReadTransaction: proposalParams.WriteReadType,
			FederatedModelUpdateHash: hasher.Hash(append(proposalParams.FederatedModelUpdate,
				[]byte(proposalParams.ClientClock)...)),
		}, proposalParams.FederatedModelUpdate
	}
	return &protos.ProposalRequest{
		TargetSystem:         baseOptions.Bconfig.TargetSystem,
		ProposalId:           t.TransactionId,
		ClientId:             config.Config.UUID,
		ContractName:         baseOptions.Bconfig.ContractName,
		MethodName:           proposalParams.ContractMethodName,
		MethodParams:         proposalParams.Outputs,
		WriteReadTransaction: proposalParams.WriteReadType,
	}, nil
}

func (t *ClientTransaction) MakeTransactionBenchmarkExecutorWithClientFullSign(bconfig *protos.BenchmarkConfig, endorsementPolicy int) (*protos.Transaction, error) {
	if t.Status == protos.TransactionStatus_FAILED_GENERAL {
		return nil, errors.New("transaction failed")
	}
	tempTransaction := &protos.Transaction{
		TargetSystem:      bconfig.TargetSystem,
		TransactionId:     t.TransactionId,
		ClientId:          config.Config.UUID,
		ContractName:      bconfig.ContractName,
		NodeSignatures:    map[string][]byte{},
		EndorsementPolicy: int32(endorsementPolicy),
	}
	for _, proposal := range t.ProposalResponses {
		tempTransaction.NodeSignatures[proposal.NodeId] = proposal.NodeSignature
	}
	for _, proposal := range t.ProposalResponses {
		tempTransaction.ReadWriteSet = proposal.ReadWriteSet
		break
	}
	sort.Slice(tempTransaction.ReadWriteSet.ReadKeys.ReadKeys, func(i, j int) bool {
		return tempTransaction.ReadWriteSet.ReadKeys.ReadKeys[i].Key < tempTransaction.ReadWriteSet.ReadKeys.ReadKeys[j].Key
	})
	sort.Slice(tempTransaction.ReadWriteSet.WriteKeyValues.WriteKeyValues, func(i, j int) bool {
		return tempTransaction.ReadWriteSet.WriteKeyValues.WriteKeyValues[i].Key < tempTransaction.ReadWriteSet.WriteKeyValues.WriteKeyValues[j].Key
	})
	marshalledReadWriteSet, _ := proto.Marshal(tempTransaction.ReadWriteSet)
	tempTransaction.ClientSignature = t.signer.Sign(marshalledReadWriteSet)

	return tempTransaction, nil
}

func (t *ClientTransaction) MakeTransactionBenchmarkExecutorFederated(counter int, baseOptions *contractinterface.BaseContractOptions,
	bconfig *protos.BenchmarkConfig, endorsementPolicy int, clientPseudoId int) (*protos.Transaction, protos.ProposalRequest_WriteReadTransaction, bool, error) {

	proposalParams := baseOptions.BenchmarkFunction(&contractinterface.BenchmarkFunctionOptions{
		Counter:                     counter,
		CurrentClientPseudoId:       baseOptions.CurrentClientPseudoId,
		BenchmarkUtils:              baseOptions.BenchmarkUtils,
		TotalTransactions:           baseOptions.TotalTransactions,
		SingleFunctionCounter:       baseOptions.SingleFunctionCounter,
		CrdtOperationPerObjectCount: baseOptions.Bconfig.CrdtOperationPerObjectCount,
		CrdtObjectType:              baseOptions.Bconfig.CrdtObjectType,
		CrdtObjectCount:             baseOptions.Bconfig.CrdtObjectCount,
		NumberOfKeys:                int(baseOptions.Bconfig.NumberOfKeys),
		NumberOfKeysSecond:          int(baseOptions.Bconfig.NumberOfKeysSecond),
		FederatedConnPool:           baseOptions.FederatedConnPool,
		NodesConnPool:               baseOptions.NodesConnectionPool,
	})

	tempTransaction := &protos.Transaction{
		TargetSystem:      bconfig.TargetSystem,
		TransactionId:     t.TransactionId,
		ClientId:          config.Config.UUID,
		ContractName:      bconfig.ContractName,
		NodeSignatures:    map[string][]byte{},
		EndorsementPolicy: int32(endorsementPolicy),
		IsFederated:       true,
		ClientPseudoId:    strconv.Itoa(clientPseudoId),
	}

	tempTransaction.ReadWriteSet = &protos.ProposalResponseReadWriteSet{
		ReadKeys:       &protos.ProposalResponseReadKeys{ReadKeys: []*protos.ReadKey{}},
		WriteKeyValues: &protos.ProposalResponseWriteKeyValues{WriteKeyValues: []*protos.WriteKeyValue{}},
	}
	tempTransaction.ReadWriteSet.WriteKeyValues.WriteKeyValues = append(tempTransaction.ReadWriteSet.WriteKeyValues.WriteKeyValues, &protos.WriteKeyValue{
		Key:       proposalParams.Outputs[0],
		Value:     proposalParams.FederatedModelUpdate,
		WriteType: protos.WriteKeyValue_FEDERATED,
	})

	return tempTransaction, proposalParams.WriteReadType, proposalParams.ShouldBeSubmitted, nil
}
