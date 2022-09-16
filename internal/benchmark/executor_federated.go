package benchmark

import (
	"github.com/golang/protobuf/proto"
	"gitlab.lrz.de/orderless/orderlessfl/internal/logger"
	"gitlab.lrz.de/orderless/orderlessfl/internal/profiling"
	protos "gitlab.lrz.de/orderless/orderlessfl/protos/goprotos"
	"time"
)

func (rex *RoundExecutor) executeTransactionPart1Federated(counter int, startTime time.Time) {
	transactionResult := MakeNewTransactionResultFederated(counter, rex.endorsementPolicyOrgsWithExtraEndorsement, rex.endorsementPolicyOrgs)
	rex.executor.transactionsResult.lock.Lock()
	rex.executor.transactionsResult.transactions[transactionResult.transaction.TransactionId] = transactionResult
	rex.executor.transactionsResult.lock.Unlock()
	transaction, readWriteType, shouldBeSubmitted, err := transactionResult.transaction.MakeTransactionBenchmarkExecutorFederated(transactionResult.transactionCounter,
		rex.baseContractOptions, rex.benchmarkConfig, rex.endorsementPolicyOrgs, rex.executor.currentClientPseudoId)
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
	transactionResult.readWriteType = readWriteType
	transactionResult.latencyMeasurementInstance = transactionResult.latencyMeasurement(startTime)
	if shouldBeSubmitted {
		rex.streamTransactionFederated(transactionResult, transaction)
	} else {
		rex.executor.roundExecutor.executeTransactionPart3OrderlessFL(transactionResult)
	}
}

func (rex *RoundExecutor) streamTransactionFederated(transactionResult *TransactionResult, transaction *protos.Transaction) {
	if profiling.IsBandwidthProfiling {
		transactionResult.sentTransactionBytes = rex.selectedOrgsEndorsementPolicyCount * proto.Size(transaction)
	}
	sentTransactions := 0
	for _, nodeId := range rex.selectedOrgsEndorsementPolicy {
		rex.executor.clientTransactionStreamLock.RLock()
		streamer, ok := rex.executor.clientTransactionStream[nodeId]
		rex.executor.clientTransactionStreamLock.RUnlock()
		if !ok {
			continue
		}
		if err := streamer.streamOrderlessFL.Send(transaction); err != nil {
			rex.executor.clientTransactionStreamLock.Lock()
			delete(rex.executor.clientTransactionStream, nodeId)
			rex.executor.clientTransactionStreamLock.Unlock()
			err = rex.executor.makeSingleStreamTransactionOrderlessFL(nodeId)
			if err != nil {
				logger.ErrorLogger.Println(err)
			}
		}
		sentTransactions++
	}
	if sentTransactions < rex.endorsementPolicyOrgs {
		transaction.Status = protos.TransactionStatus_FAILED_GENERAL
	}
}
