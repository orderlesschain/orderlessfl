package mnistfederatedchain

import (
	"context"
	"github.com/golang/protobuf/proto"
	"gitlab.lrz.de/orderless/orderlessfl/internal/connection/connpool"
	"gitlab.lrz.de/orderless/orderlessfl/internal/contract/contractinterface"
	"gitlab.lrz.de/orderless/orderlessfl/internal/logger"
	protos "gitlab.lrz.de/orderless/orderlessfl/protos/goprotos"
	"io/ioutil"
	"path"
	"strconv"
)

const ModelID1 = "mnistfederatedchain1"
const ModelID2 = "mnistfederatedchain2"
const ModelID3 = "mnistfederatedchain3"
const ModelID4 = "mnistfederatedchain4"
const BatchSize = 10000

func GetGlobalModelAndTrainAndSubmit(options *contractinterface.BenchmarkFunctionOptions) *contractinterface.BenchmarkFunctionOutputs {
	totalClients := options.NumberOfKeysSecond
	rangePerClient := options.NumberOfKeys / totalClients
	min := options.CurrentClientPseudoId * rangePerClient
	//max := min + rangePerClient
	//clientId := options.BenchmarkUtils.GetRandomIntWithMinAndMax(min, max)
	clientId := min + (options.Counter % rangePerClient)

	modelId := ""
	modelType := int32(0)
	if options.CrdtObjectCount == "1" {
		modelId = ModelID1
		modelType = 1
	} else if options.CrdtObjectCount == "2" {
		modelId = ModelID2
		modelType = 2
	} else if options.CrdtObjectCount == "3" {
		modelId = ModelID3
		modelType = 3
	} else if options.CrdtObjectCount == "4" {
		modelId = ModelID4
		modelType = 4
	}

	failedResponse := &contractinterface.BenchmarkFunctionOutputs{
		Outputs:            []string{modelId, "Failed"},
		ContractMethodName: "getglobalmodelandtrainandsubmit",
		WriteReadType:      protos.ProposalRequest_Write,
	}

	options.SingleFunctionCounter.Lock.Lock()
	counter := options.SingleFunctionCounter.Counter[clientId]
	options.SingleFunctionCounter.Lock.Unlock()

	modelWeightsFromNode, err := getLatestModelWeightFromOneNode(options, modelId, clientId)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return failedResponse
	}
	if modelWeightsFromNode.ClientLatestClock != int32(counter) {
		logger.InfoLogger.Println("Previous update has not been applied, clientId:", clientId, " , ", modelWeightsFromNode.ClientLatestClock, " , ", counter)
		return failedResponse
	}

	options.SingleFunctionCounter.Lock.Lock()
	options.SingleFunctionCounter.Counter[clientId]++
	counter = options.SingleFunctionCounter.Counter[clientId]
	options.SingleFunctionCounter.Lock.Unlock()

	if int32(counter)-modelWeightsFromNode.ClientLatestClock != 1 {
		logger.InfoLogger.Println("Too fast, clientId:", clientId, " , ", modelWeightsFromNode.ClientLatestClock, " , ", counter)
		options.SingleFunctionCounter.Lock.Lock()
		options.SingleFunctionCounter.Counter[clientId]--
		options.SingleFunctionCounter.Lock.Unlock()
		return failedResponse
	}

	conn, err := options.FederatedConnPool.Get(context.Background())
	if conn == nil || err != nil {
		logger.ErrorLogger.Println(err)
		return failedResponse
	}
	modelUpdateResponse, err := protos.NewFederatedServiceClient(conn.ClientConn).TrainUpdateModel(context.Background(), &protos.ModelUpdateRequest{
		ModelId:           modelId,
		ModelType:         modelType,
		RequestCounter:    int32(counter),
		ClientId:          int32(clientId),
		BatchSize:         BatchSize,
		TotalClientsCount: int32(options.NumberOfKeys),
		ClientLatestClock: int32(counter),
		LayersWeight:      modelWeightsFromNode.LayersWeight,
	})
	if err != nil {
		logger.ErrorLogger.Println(err)
		return failedResponse
	}
	if err = conn.Close(); err != nil {
		logger.ErrorLogger.Println(err)
		return failedResponse
	}
	if modelUpdateResponse.Status != protos.ModelUpdateRequest_SUCCEEDED {
		return failedResponse
	}
	return &contractinterface.BenchmarkFunctionOutputs{
		Outputs:              []string{modelId},
		ContractMethodName:   "getglobalmodelandtrainandsubmit",
		FederatedModelUpdate: modelUpdateResponse.LayersWeight,
		WriteReadType:        protos.ProposalRequest_Write,
		ClientClock:          strconv.Itoa(counter),
		ShouldBeSubmitted:    true,
	}
}

func GetGlobalModelAndSubmitPretrainedUpdate(options *contractinterface.BenchmarkFunctionOptions) *contractinterface.BenchmarkFunctionOutputs {
	totalClients := options.NumberOfKeysSecond
	rangePerClient := options.NumberOfKeys / totalClients
	min := options.CurrentClientPseudoId * rangePerClient
	//max := min + rangePerClient
	//clientId := options.BenchmarkUtils.GetRandomIntWithMinAndMax(min, max)
	clientId := min + (options.Counter % rangePerClient)

	modelId := ""

	if options.CrdtObjectCount == "1" {
		modelId = ModelID1
	} else if options.CrdtObjectCount == "2" {
		modelId = ModelID2
	} else if options.CrdtObjectCount == "3" {
		modelId = ModelID3
	} else if options.CrdtObjectCount == "4" {
		modelId = ModelID4
	}

	failedResponse := &contractinterface.BenchmarkFunctionOutputs{
		Outputs:            []string{modelId, "Failed"},
		ContractMethodName: "getglobalmodelandtrainandsubmit",
		WriteReadType:      protos.ProposalRequest_Write,
	}

	options.SingleFunctionCounter.Lock.Lock()
	options.SingleFunctionCounter.Counter[clientId]++
	counter := options.SingleFunctionCounter.Counter[clientId]
	options.SingleFunctionCounter.Lock.Unlock()

	preTrainedUpdate, err := readPreTrainedModelUpdate(&protos.ModelUpdateRequest{
		ModelId:           modelId,
		RequestCounter:    int32(counter),
		ClientId:          int32(clientId),
		BatchSize:         BatchSize,
		TotalClientsCount: int32(options.NumberOfKeys),
		ClientLatestClock: int32(counter),
	})
	if err != nil {
		logger.ErrorLogger.Println(err)
		return failedResponse
	}
	return &contractinterface.BenchmarkFunctionOutputs{
		Outputs:              []string{modelId},
		ContractMethodName:   "getglobalmodelandtrainandsubmit",
		FederatedModelUpdate: preTrainedUpdate,
		WriteReadType:        protos.ProposalRequest_Write,
		ClientClock:          strconv.Itoa(counter),
		ShouldBeSubmitted:    true,
	}
}

var preTrainedLayers [][]byte

func readPreTrainedModelUpdate(modelUpdate *protos.ModelUpdateRequest) ([]byte, error) {
	if preTrainedLayers == nil || len(preTrainedLayers) == 0 {
		logger.InfoLogger.Println("Initializing the preTrained")
		modelWeightsByte, err := ioutil.ReadFile(path.Join("../federated/federated_models/", "pretrained-"+modelUpdate.ModelId))
		if err != nil {
			modelWeightsByte, err = ioutil.ReadFile(path.Join("../../federatedpython/federated_models/", "pretrained-"+modelUpdate.ModelId))
			if err != nil {
				return nil, err
			}
		}
		modelWeights := &protos.ModelWeights{}
		if err = proto.Unmarshal(modelWeightsByte, modelWeights); err != nil {
			logger.ErrorLogger.Println(err)
		}
		preTrainedLayers = modelWeights.Layers
	}
	modelWeights := &protos.ModelWeights{
		Layers:            preTrainedLayers,
		ModelClock:        0,
		ClientId:          modelUpdate.ClientId,
		ClientLatestClock: modelUpdate.ClientLatestClock,
	}
	modelWeightsByte, err := proto.Marshal(modelWeights)
	if err != nil {
		return nil, err
	}
	return modelWeightsByte, nil
}

func getLatestModelWeightFromOneNode(options *contractinterface.BenchmarkFunctionOptions, modelId string, clientId int) (*protos.ModelUpdateRequest, error) {
	var nodeConnection *connpool.Pool

	selectedNodeIdForModel := "node" + strconv.Itoa(options.CurrentClientPseudoId%len(options.NodesConnPool))
	nodeConnection = options.NodesConnPool[selectedNodeIdForModel]

	conn, err := nodeConnection.Get(context.Background())
	if conn == nil || err != nil {
		return nil, err
	}
	resp, err := protos.NewTransactionServiceClient(conn.ClientConn).GetLatestFederatedModel(context.Background(), &protos.ModelUpdateRequest{
		ModelId:  modelId,
		ClientId: int32(clientId),
	})
	if err != nil {
		return nil, err
	}
	if err = conn.Close(); err != nil {
		return nil, err
	}
	return resp, nil
}
