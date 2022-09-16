package benchmark

import (
	"context"
	"encoding/csv"
	"gitlab.lrz.de/orderless/orderlessfl/contractsbenchmarks/benchmark/benchmarkfunctions"
	"gitlab.lrz.de/orderless/orderlessfl/internal/benchmark/benchmarkutils"
	"gitlab.lrz.de/orderless/orderlessfl/internal/benchmark/bencmarkdb"
	"gitlab.lrz.de/orderless/orderlessfl/internal/benchmark/latencymeasurment"
	"gitlab.lrz.de/orderless/orderlessfl/internal/config"
	"gitlab.lrz.de/orderless/orderlessfl/internal/connection/connpool"
	"gitlab.lrz.de/orderless/orderlessfl/internal/contract/contractinterface"
	"gitlab.lrz.de/orderless/orderlessfl/internal/customcrypto/keygenerator"
	"gitlab.lrz.de/orderless/orderlessfl/internal/customcrypto/signer"
	"gitlab.lrz.de/orderless/orderlessfl/internal/logger"
	"gitlab.lrz.de/orderless/orderlessfl/internal/transaction"
	protos "gitlab.lrz.de/orderless/orderlessfl/protos/goprotos"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TransactionResult struct {
	transaction                       *transaction.ClientTransaction
	transactionCounter                int
	receivedProposalCount             int
	receivedProposalExpected          int
	receivedTransactionCommitCount    int
	receivedTransactionCommitExpected int
	latencyMeasurement                latencymeasurment.LatencyMeasurement
	latencyMeasurementInstance        latencymeasurment.LatencyMeasurementInstance
	startTime                         time.Time
	endTime                           time.Time
	latency                           time.Duration
	readWriteType                     protos.ProposalRequest_WriteReadTransaction
	sentProposalBytes                 int
	receivedProposalBytes             int
	sentTransactionBytes              int
	receivedTransactionBytes          int
	federatedModelUpdate              []byte
	modelId                           string
}

func MakeNewTransactionResultOrderlessFL(counter, endorsementPolicyWithExtra, endorsementPolicy int, signer *signer.Signer) *TransactionResult {
	return &TransactionResult{
		transaction:                       transaction.NewClientTransaction(signer),
		transactionCounter:                counter,
		latencyMeasurement:                latencymeasurment.Latency,
		receivedProposalCount:             0,
		receivedProposalExpected:          endorsementPolicyWithExtra,
		receivedTransactionCommitCount:    0,
		receivedTransactionCommitExpected: endorsementPolicy,
	}
}

func MakeNewTransactionResultFederated(counter, endorsementPolicyWithExtra, endorsementPolicy int) *TransactionResult {
	return &TransactionResult{
		transaction:                       transaction.NewClientTransaction(nil),
		transactionCounter:                counter,
		latencyMeasurement:                latencymeasurment.Latency,
		receivedProposalCount:             0,
		receivedProposalExpected:          endorsementPolicyWithExtra,
		receivedTransactionCommitCount:    0,
		receivedTransactionCommitExpected: endorsementPolicy,
	}
}

func (t *TransactionResult) EndTransactionMeasurements() {
	if t.latencyMeasurementInstance != nil {
		t.startTime, t.latency, t.endTime = t.latencyMeasurementInstance()
	}
}

type TransactionsResult struct {
	transactions map[string]*TransactionResult
	lock         *sync.Mutex
}

type ExecutorDriver func()

type proposalStream struct {
	streamOrderlessFL protos.TransactionService_ProcessProposalOrderlessFLStreamClient
}

type transactionStream struct {
	streamOrderlessFL protos.TransactionService_CommitOrderlessFLTransactionStreamClient
}

type Executor struct {
	doneChannel                           chan bool
	dbOps                                 *bencmarkdb.ExperimentOperations
	nodesConnectionsStreamProposals       map[string]*connpool.Pool
	nodesConnectionsStreamTransactions    map[string]*connpool.Pool
	nodesConnectionsWatchProposalEvent    map[string]*connpool.Pool
	nodesConnectionsWatchTransactionEvent map[string]*connpool.Pool
	clientProposalStreamLock              *sync.RWMutex
	clientProposalStream                  map[string]*proposalStream
	clientTransactionStreamLock           *sync.RWMutex
	clientTransactionStream               map[string]*transactionStream
	transactionsResult                    *TransactionsResult
	currentClientPseudoId                 int
	inExperimentParticipatingNodes        []string
	inExperimentParticipatingClients      []string
	isThisClientParticipating             bool
	clientsCorrespondingNodes             []string
	roundExecutor                         *RoundExecutor
	roundNotDone                          bool
	PublicPrivateKey                      *keygenerator.RSAKey
	federatedConnPool                     *connpool.Pool
	nodesConnectionsForFederated          map[string]*connpool.Pool
}

type transactionInitiator struct {
	counter   int
	startTime time.Time
}

type RoundExecutor struct {
	experimentID                                string
	experimentDBID                              uint
	benchmarkConfig                             *protos.BenchmarkConfig
	transactionPerSecondFromBenchmark           int
	transactionsSendDurationSecondFromBenchmark int
	totalTransactionsFromBenchmark              int
	totalDesiredTransaction                     int
	sentTransactionsCount                       int
	doneTransactionsCount                       int
	benchmarkFunction                           contractinterface.BenchmarkFunction
	benchmarkUtils                              *benchmarkutils.BenchmarkUtils
	endorsementPolicyOrgs                       int
	endorsementPolicyOrgsWithExtraEndorsement   int
	selectedOrgsEndorsementPolicy               []string
	selectedOrgsEndorsementPolicyCount          int
	baseContractOptions                         *contractinterface.BaseContractOptions
	roundDone                                   chan bool
	executor                                    *Executor
	signer                                      *signer.Signer
	transactionInitiators                       chan *transactionInitiator
}

func NewExecutor() *Executor {
	tempEx := &Executor{
		doneChannel:      make(chan bool),
		PublicPrivateKey: keygenerator.LoadPublicPrivateKeyFromFile(),
	}
	tempEx.currentClientPseudoId = tempEx.getCurrentClientPseudoId()
	tempEx.setInExperimentParticipatingComponents()
	tempEx.setIsThisClientParticipating()
	if !tempEx.isThisClientParticipating {
		logger.InfoLogger.Println("This client in NOT participating in the experiment")
		return tempEx
	}
	tempEx.dbOps = bencmarkdb.NewExperimentOperations()
	tempEx.nodesConnectionsWatchTransactionEvent = connpool.GetNodeConnectionsWatchTransactionEvent(tempEx.clientsCorrespondingNodes)
	if config.Config.IsOrderlessFL || config.Config.IsFederatedLearning {
		go tempEx.subscriberForTransactionEventsOrderlessFL()
		tempEx.federatedConnPool = connpool.GetFederatedLearningLocalHostConnections()
		tempEx.nodesConnectionsForFederated = connpool.GetNodeConnectionsForFederated(tempEx.clientsCorrespondingNodes)
	} else {
		logger.FatalLogger.Fatalln("target system not set")
	}
	tempEx.nodesConnectionsStreamProposals = connpool.GetNodeConnectionsStreamProposal(tempEx.clientsCorrespondingNodes)
	tempEx.nodesConnectionsStreamTransactions = connpool.GetNodeConnectionsStreamTransactions(tempEx.clientsCorrespondingNodes)
	tempEx.nodesConnectionsWatchProposalEvent = connpool.GetNodeConnectionsWatchProposalEvent(tempEx.clientsCorrespondingNodes)
	tempEx.clientProposalStreamLock = &sync.RWMutex{}
	tempEx.clientTransactionStreamLock = &sync.RWMutex{}
	tempEx.clientProposalStream = map[string]*proposalStream{}
	tempEx.clientTransactionStream = map[string]*transactionStream{}
	go tempEx.makeAllStreamProposal()
	if config.Config.IsOrderlessFL || config.Config.IsFederatedLearning {
		go tempEx.subscriberForProposalEventsOrderlessFL()
		go tempEx.makeAllStreamTransactionOrderlessFL()
	}
	connpool.SleepAndContinue()
	return tempEx
}

func NewRoundExecutor(bconfig *protos.BenchmarkConfig, executor *Executor) *RoundExecutor {
	benchmarkFunction, _ := benchmarkfunctions.GetBenchmarkFunctions(bconfig.ContractName, bconfig.BenchmarkFunctionName)
	tempEx := &RoundExecutor{
		experimentID:      bconfig.Base.ExperimentId,
		benchmarkConfig:   bconfig,
		benchmarkFunction: benchmarkFunction,
		transactionsSendDurationSecondFromBenchmark: int(bconfig.TransactionSendDurationSecond),
		totalTransactionsFromBenchmark:              int(bconfig.TotalTransactions),
		transactionPerSecondFromBenchmark:           int(bconfig.TransactionPerSecond),
		sentTransactionsCount:                       0,
		doneTransactionsCount:                       0,
		executor:                                    executor,
		endorsementPolicyOrgs:                       int(bconfig.EndorsementPolicyOrgs),
		endorsementPolicyOrgsWithExtraEndorsement:   int(bconfig.EndorsementPolicyOrgs), // + config.Config.ExtraEndorsementOrgs,
		selectedOrgsEndorsementPolicy:               []string{},
		roundDone:                                   make(chan bool),
		signer:                                      signer.NewSigner(executor.PublicPrivateKey),
		transactionInitiators:                       make(chan *transactionInitiator, 100000),
	}
	tempEx.setEndorsementPolicyIds()
	if tempEx.transactionsSendDurationSecondFromBenchmark > 0 {
		tempEx.totalDesiredTransaction = tempEx.transactionsSendDurationSecondFromBenchmark * tempEx.transactionPerSecondFromBenchmark
	}
	if tempEx.totalTransactionsFromBenchmark > 0 {
		tempEx.totalDesiredTransaction = tempEx.totalTransactionsFromBenchmark
	}
	if tempEx.totalDesiredTransaction == 0 {
		logger.FatalLogger.Fatalln("no number of transactions is sent to send")
	}
	tempEx.executor.transactionsResult = &TransactionsResult{
		transactions: map[string]*TransactionResult{},
		lock:         &sync.Mutex{},
	}
	tempEx.benchmarkUtils = benchmarkutils.NewBenchmarkDist(0, tempEx.executor.currentClientPseudoId, tempEx.totalDesiredTransaction, bconfig.NumberOfKeys)
	tempEx.baseContractOptions = &contractinterface.BaseContractOptions{
		BenchmarkUtils:        tempEx.benchmarkUtils,
		Bconfig:               tempEx.benchmarkConfig,
		BenchmarkFunction:     tempEx.benchmarkFunction,
		CurrentClientPseudoId: tempEx.executor.currentClientPseudoId,
		TotalTransactions:     tempEx.totalDesiredTransaction,
		SingleFunctionCounter: contractinterface.NewSingleFunctionCounter(bconfig.NumberOfKeys),
		NodesConnectionPool:   tempEx.executor.nodesConnectionsForFederated,
	}
	tempEx.baseContractOptions.FederatedConnPool = tempEx.executor.federatedConnPool

	executor.roundExecutor = tempEx
	executor.roundNotDone = true
	return tempEx
}

func (rex *RoundExecutor) ExecuteBenchmark() {
	expId, err := rex.executor.dbOps.AddBenchmarkToDB(&bencmarkdb.ExperimentsORM{
		ExperimentID: rex.experimentID,
		Status:       bencmarkdb.Running,
		ResultPath:   "",
	})
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
	rex.experimentDBID = expId
	transactionInterval := time.NewTicker(time.Second / time.Duration(rex.transactionPerSecondFromBenchmark))
	for {
		rex.sentTransactionsCount++
		rex.transactionInitiators <- &transactionInitiator{counter: rex.sentTransactionsCount, startTime: <-transactionInterval.C}
		if rex.sentTransactionsCount == rex.totalDesiredTransaction {
			transactionInterval.Stop()
			break
		}
	}
	timer := time.NewTimer(config.Config.TransactionTimeoutSecondConverted)
	select {
	case <-rex.roundDone:
		timer.Stop()
	case <-timer.C:
		rex.executor.roundNotDone = false
	}
	rex.CleanUpTransactions()
	rex.ExportReportFile()
	err = rex.executor.dbOps.UpdateExperimentStatus(rex.experimentDBID, bencmarkdb.Done)
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
}

func (rex *RoundExecutor) runExecuteBenchmarkStarter() {
	if config.Config.IsOrderlessFL {
		for {
			starter := <-rex.transactionInitiators
			go rex.executeTransactionPart1OrderlessFL(starter.counter, starter.startTime)
		}
	} else if config.Config.IsFederatedLearning {
		for {
			starter := <-rex.transactionInitiators
			go rex.executeTransactionPart1Federated(starter.counter, starter.startTime)
		}
	}
}

func (ex *Executor) getCurrentClientPseudoId() int {
	currentClientId := connpool.GetComponentPseudoName()
	currentClientId = strings.ReplaceAll(currentClientId, "client", "")
	currentClientIdInt, err := strconv.Atoi(currentClientId)
	if err != nil {
		logger.FatalLogger.Fatalln(err)
	}
	return currentClientIdInt
}

func (ex *Executor) setInExperimentParticipatingComponents() {
	for i := 0; i < config.Config.TotalNodeCount; i++ {
		ex.inExperimentParticipatingNodes = append(ex.inExperimentParticipatingNodes, "node"+strconv.Itoa(i))
	}
	for i := 0; i < config.Config.TotalClientCount; i++ {
		ex.inExperimentParticipatingClients = append(ex.inExperimentParticipatingClients, "client"+strconv.Itoa(i))
	}
	currentClientId := connpool.GetComponentPseudoName()
	currentClientId = strings.ReplaceAll(currentClientId, "client", "")
	currentClientIdInt, err := strconv.Atoi(currentClientId)
	if err != nil {
		logger.FatalLogger.Fatalln(err)
	}
	nodeCount := len(ex.inExperimentParticipatingNodes)
	endorsingOrgsCount := config.Config.EndorsementPolicy + config.Config.ExtraEndorsementOrgs
	if endorsingOrgsCount > nodeCount {
		endorsingOrgsCount = nodeCount
	}
	inExperimentParticipatingClientsCount := len(ex.inExperimentParticipatingClients)
	if config.Config.OrgsPercentageIncreasedLoad == 0 {
		totalNodeConnectionsCounter := 0
		for clientCount := 0; clientCount < inExperimentParticipatingClientsCount; clientCount++ {
			for endorser := 0; endorser < endorsingOrgsCount; endorser++ {
				nodeId := totalNodeConnectionsCounter % nodeCount
				totalNodeConnectionsCounter++
				if clientCount == currentClientIdInt {
					ex.clientsCorrespondingNodes = append(ex.clientsCorrespondingNodes, "node"+strconv.Itoa(nodeId))
				}
			}
		}
	} else if config.Config.LoadIncreasePercentage != 0 {
		logger.InfoLogger.Println("Re-balancing the loads")
		percentageNodeCount := int(math.Ceil((float64(nodeCount * config.Config.OrgsPercentageIncreasedLoad)) / 100))
		uniformLoad := (inExperimentParticipatingClientsCount * endorsingOrgsCount) / nodeCount
		percentageIncreaseLoadCount := int(math.Ceil((float64(config.Config.LoadIncreasePercentage * uniformLoad)) / 100))

		connectionPerNodeCounter := inExperimentParticipatingClientsCount * endorsingOrgsCount
		allCorrespondingNodes := make([]string, connectionPerNodeCounter)
		positionInEndorsementSlice := 0
		endorsementRingCounter := 0
		totalLoad := percentageIncreaseLoadCount + uniformLoad
		for nodeId := 0; nodeId < percentageNodeCount; nodeId++ {
			nodeIdString := "node" + strconv.Itoa(nodeCount-nodeId-1)
			for load := 0; load < totalLoad; load++ {
				location := ((endorsementRingCounter % inExperimentParticipatingClientsCount) * endorsingOrgsCount) + positionInEndorsementSlice
				endorsementRingCounter++
				if endorsementRingCounter%inExperimentParticipatingClientsCount == 0 {
					positionInEndorsementSlice++
				}
				if location >= connectionPerNodeCounter || len(allCorrespondingNodes[location]) > 0 {
					logger.FatalLogger.Fatalln("the location for load balancing already occupied")
				}
				allCorrespondingNodes[location] = nodeIdString
			}
		}
		nodeCountOthers := nodeCount - percentageNodeCount
		var allCorrespondingNodesOthers []string
		missingLocations := connectionPerNodeCounter - endorsementRingCounter
		for i := 0; i < missingLocations; i++ {
			nodeId := i % nodeCountOthers
			allCorrespondingNodesOthers = append(allCorrespondingNodesOthers, "node"+strconv.Itoa(nodeId))
		}
		endorsementRingCounter = 0
		for i := 0; i < connectionPerNodeCounter; i++ {
			if len(allCorrespondingNodes[i]) == 0 {
				allCorrespondingNodes[i] = allCorrespondingNodesOthers[endorsementRingCounter]
				endorsementRingCounter++
			}
		}
		baseCount := currentClientIdInt * endorsingOrgsCount
		for endorser := 0; endorser < endorsingOrgsCount; endorser++ {
			if baseCount < len(allCorrespondingNodes) {
				ex.clientsCorrespondingNodes = append(ex.clientsCorrespondingNodes, allCorrespondingNodes[baseCount])
			}
			baseCount++
		}
		unique := make(map[string]bool, len(ex.clientsCorrespondingNodes))
		for _, node := range ex.clientsCorrespondingNodes {
			if _, ok := unique[node]; ok {
				logger.FatalLogger.Fatalln("non unique nodes in the list")
			}
			unique[node] = true
		}
	}
	logger.InfoLogger.Println(currentClientIdInt, ex.clientsCorrespondingNodes, len(ex.clientsCorrespondingNodes))
}

func (ex *Executor) setInExperimentParticipatingFaultyComponents(nodes *protos.FaultyNodes) {
	if len(ex.roundExecutor.selectedOrgsEndorsementPolicy) == 0 {
		return
	}
	failureRandom := rand.New(rand.NewSource(time.Now().UnixNano()))
	nodeCount := len(ex.inExperimentParticipatingNodes)
	endorsingOrgsCount := config.Config.EndorsementPolicy + config.Config.ExtraEndorsementOrgs
	if endorsingOrgsCount > nodeCount {
		endorsingOrgsCount = nodeCount
	}
	tempSelectedOrgsEndorsementPolicy := map[string]bool{}
	for i := 0; i < len(ex.roundExecutor.selectedOrgsEndorsementPolicy); i++ {
		blackListNode := false
		for j := 0; j < len(nodes.NodeId); j++ {
			if ex.roundExecutor.selectedOrgsEndorsementPolicy[i] == nodes.NodeId[j] {
				blackListNode = true
			}
		}
		if !blackListNode {
			tempSelectedOrgsEndorsementPolicy[ex.roundExecutor.selectedOrgsEndorsementPolicy[i]] = true
		}
	}
	if len(tempSelectedOrgsEndorsementPolicy) == endorsingOrgsCount {
		return
	}
	newlySelectedOrgsEndorsementPolicy := map[string]bool{}
	nodeDifference := endorsingOrgsCount - len(tempSelectedOrgsEndorsementPolicy)
	for i := 0; i < nodeDifference; i++ {
		randNodeIdString := "node" + strconv.Itoa(failureRandom.Intn(nodeCount))
		blackListed := false
		for j := 0; j < len(nodes.NodeId); j++ {
			if randNodeIdString == nodes.NodeId[j] {
				blackListed = true
				break
			}
		}
		if blackListed {
			i--
			continue
		}
		notUnique := false
		for nodeId := range tempSelectedOrgsEndorsementPolicy {
			if randNodeIdString == nodeId {
				notUnique = true
				break
			}
		}
		if notUnique {
			i--
			continue
		}
		newlySelectedOrgsEndorsementPolicy[randNodeIdString] = true
		tempSelectedOrgsEndorsementPolicy[randNodeIdString] = true
	}
	var tempSelectedOrgsEndorsementPolicyList []string
	for nodeId := range tempSelectedOrgsEndorsementPolicy {
		tempSelectedOrgsEndorsementPolicyList = append(tempSelectedOrgsEndorsementPolicyList, nodeId)
	}
	if endorsingOrgsCount != len(tempSelectedOrgsEndorsementPolicyList) {
		logger.ErrorLogger.Println("the number of nodes are not equal to endorsement policy")
	}
	unique := make(map[string]bool, len(tempSelectedOrgsEndorsementPolicyList))
	for _, node := range tempSelectedOrgsEndorsementPolicyList {
		if _, ok := unique[node]; ok {
			logger.ErrorLogger.Println("non unique nodes in the list")
		}
		unique[node] = true
	}
	logger.InfoLogger.Println("After notify", tempSelectedOrgsEndorsementPolicyList, len(tempSelectedOrgsEndorsementPolicyList))

	timer := time.NewTimer(time.Duration(nodes.StartAfterS) * time.Second)
	<-timer.C

	ex.clientProposalStreamLock.Lock()
	ex.clientTransactionStreamLock.Lock()

	for _, nodeId := range nodes.NodeId {
		delete(ex.clientProposalStream, nodeId)
		delete(ex.clientTransactionStream, nodeId)
	}
	for nodeId := range newlySelectedOrgsEndorsementPolicy {
		ex.nodesConnectionsStreamProposals[nodeId] = connpool.GetSingleNodeConnectionsStreamProposal(nodeId)
		ex.nodesConnectionsStreamTransactions[nodeId] = connpool.GetSingleNodeConnectionsStreamTransactions(nodeId)
		ex.nodesConnectionsWatchTransactionEvent[nodeId] = connpool.GetSingleNodeConnectionsWatchTransactionEvent(nodeId)
		ex.nodesConnectionsWatchProposalEvent[nodeId] = connpool.GetSingleNodeConnectionsWatchProposalEvent(nodeId)
	}
	ex.roundExecutor.selectedOrgsEndorsementPolicy = tempSelectedOrgsEndorsementPolicyList
	ex.roundExecutor.selectedOrgsEndorsementPolicyCount = len(ex.roundExecutor.selectedOrgsEndorsementPolicy)

	ex.clientTransactionStreamLock.Unlock()
	ex.clientProposalStreamLock.Unlock()

	go ex.subscriberForNewlyAddedProposalEventsOrderlessFL(newlySelectedOrgsEndorsementPolicy)
	go ex.subscriberForNewlyAddedTransactionEventsOrderlessFL(newlySelectedOrgsEndorsementPolicy)

	for nodeId := range newlySelectedOrgsEndorsementPolicy {
		if err := ex.makeSingleStreamProposal(nodeId); err != nil {
			logger.ErrorLogger.Println(err)
		}
		if err := ex.makeSingleStreamTransactionOrderlessFL(nodeId); err != nil {
			logger.ErrorLogger.Println(err)
		}
	}

}

func (ex *Executor) getNext(nodeDist []float64, desiredCount int) []int {
	var nodes []int
	for i := 0; i < len(nodeDist); i++ {
		if nodeDist[i] > 0 {
			nodes = append(nodes, i)
			nodeDist[i]--
		}
		if len(nodes) == desiredCount {
			break
		}
	}
	return nodes
}

func (rex *RoundExecutor) setEndorsementPolicyIds() {
	nodeCount := len(rex.executor.inExperimentParticipatingNodes)
	endorsingOrgsCount := config.Config.EndorsementPolicy + config.Config.ExtraEndorsementOrgs
	if endorsingOrgsCount > nodeCount {
		endorsingOrgsCount = nodeCount
	}
	for i := 0; i < endorsingOrgsCount; i++ {
		rex.selectedOrgsEndorsementPolicy = append(rex.selectedOrgsEndorsementPolicy, rex.executor.clientsCorrespondingNodes[i])
	}
	rex.selectedOrgsEndorsementPolicyCount = len(rex.selectedOrgsEndorsementPolicy)
}

func (ex *Executor) setIsThisClientParticipating() {
	currentClientId := connpool.GetComponentPseudoName()
	for _, client := range ex.inExperimentParticipatingClients {
		if currentClientId == client {
			ex.isThisClientParticipating = true
		}
	}
}

func (rex *RoundExecutor) TransactionsDoneChecker() {
	for {
		<-rex.executor.doneChannel
		rex.doneTransactionsCount++
		if rex.sentTransactionsCount == rex.totalDesiredTransaction &&
			rex.sentTransactionsCount == rex.doneTransactionsCount {
			rex.roundDone <- true
			return
		}
	}
}

func (rex *RoundExecutor) CleanUpTransactions() {
	rex.executor.transactionsResult.lock.Lock()
	for _, transactionResults := range rex.executor.transactionsResult.transactions {
		if transactionResults.transaction.Status == protos.TransactionStatus_RUNNING {
			transactionResults.EndTransactionMeasurements()
			transactionResults.transaction.Status = protos.TransactionStatus_FAILED_TIMEOUT
		}
	}
	rex.executor.transactionsResult.lock.Unlock()
}

func (rex *RoundExecutor) ExportReportFile() {
	reportPath := filepath.Join("./orderlessfl-experiments/results/",
		time.Now().Format("20060102150405"))
	err := os.MkdirAll(reportPath, os.ModePerm)
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
	reportPath = filepath.Join(reportPath, "report.csv")
	reportFile, err := os.Create(reportPath)
	if err != nil {
		logger.ErrorLogger.Println("failed creating file:", err)
	}
	csvWriter := csv.NewWriter(reportFile)
	for _, transactionResults := range rex.executor.transactionsResult.transactions {
		_ = csvWriter.Write([]string{
			transactionResults.transaction.TransactionId,
			strconv.FormatInt(transactionResults.startTime.UnixNano(), 10),
			strconv.FormatInt(transactionResults.endTime.UnixNano(), 10),
			transactionResults.latency.String(),
			transactionResults.transaction.Status.String(),
			strconv.FormatBool(rex.benchmarkConfig.ReportImportance),
			transactionResults.readWriteType.String(),
			strconv.Itoa(transactionResults.sentProposalBytes),
			strconv.Itoa(transactionResults.receivedProposalBytes),
			strconv.Itoa(transactionResults.sentTransactionBytes),
			strconv.Itoa(transactionResults.receivedTransactionBytes),
		})
	}
	csvWriter.Flush()
	err = rex.executor.dbOps.UpdateExperimentReportPath(rex.experimentDBID, reportPath)
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
	err = reportFile.Close()
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
}

func (ex *Executor) makeTransactionDone() {
	ex.doneChannel <- true
}

func (ex *Executor) makeAllStreamProposal() {
	for node := range ex.nodesConnectionsStreamProposals {
		if err := ex.makeSingleStreamProposal(node); err != nil {
			logger.ErrorLogger.Println(err)
		}
	}
}

func (ex *Executor) makeSingleStreamProposal(node string) error {
	conn, err := ex.nodesConnectionsStreamProposals[node].Get(context.Background())
	if conn == nil || err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	client := protos.NewTransactionServiceClient(conn.ClientConn)
	tempStream := &proposalStream{}
	if config.Config.IsOrderlessFL || config.Config.IsFederatedLearning {
		tempStream.streamOrderlessFL, err = client.ProcessProposalOrderlessFLStream(context.Background())
	}
	if err != nil {
		if errCon := conn.Close(); errCon != nil {
			logger.ErrorLogger.Println(errCon)
		}
		connpool.SleepAndReconnect()
		err = ex.makeSingleStreamProposal(node)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}
		return nil
	}
	ex.clientProposalStreamLock.Lock()
	ex.clientProposalStream[node] = tempStream
	ex.clientProposalStreamLock.Unlock()
	return nil
}
