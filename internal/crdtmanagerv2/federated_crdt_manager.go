package crdtmanagerv2

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"gitlab.lrz.de/orderless/orderlessfl/internal/logger"
	protos "gitlab.lrz.de/orderless/orderlessfl/protos/goprotos"
	"io/ioutil"
	"path"
	"sync"
	"time"
)

type ModelWeightsCRDT struct {
	Layers            []*LayerWeightsCRDT
	ModelClock        int32
	ClientId          int32
	ClientLatestClock int32
	CRDTObjectLock    *sync.RWMutex
}

type LayerWeightsCRDT struct {
	LayerRank      int32
	LayerType      protos.LayerWeights_LayerWeightsType
	LayerWeight1D  []float64
	LayerWeights2D [][]float64
}

type FederatedCRDTManager struct {
	crdtObjs               map[string]*ModelWeightsCRDT
	crdtObjsLock           *sync.RWMutex
	clientsObservation     map[int32]int32
	clientsObservationLock *sync.RWMutex
	latestModels           map[string][]byte
	latestModelsUpdated    map[string]time.Time
	latestModelLock        *sync.RWMutex
}

func NewFederatedCRDTManager() *FederatedCRDTManager {
	return &FederatedCRDTManager{
		crdtObjs:               map[string]*ModelWeightsCRDT{},
		crdtObjsLock:           &sync.RWMutex{},
		clientsObservation:     map[int32]int32{},
		clientsObservationLock: &sync.RWMutex{},
		latestModels:           map[string][]byte{},
		latestModelsUpdated:    map[string]time.Time{},
		latestModelLock:        &sync.RWMutex{},
	}
}

func (m *FederatedCRDTManager) ApplyModelUpdates(modelId string, modelUpdate []byte) error {
	modelWeightsCRDT, err := ConvertWeightProtoToCRDT(modelUpdate)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}
	clientCount := 0
	m.clientsObservationLock.RLock()
	clientObservation, ok := m.clientsObservation[modelWeightsCRDT.ClientId]
	clientCount = len(m.clientsObservation)
	m.clientsObservationLock.RUnlock()
	if !ok {
		m.clientsObservationLock.Lock()
		m.clientsObservation[modelWeightsCRDT.ClientId] = 0
		clientCount = len(m.clientsObservation)
		m.clientsObservationLock.Unlock()
	} else {
		if clientObservation == modelWeightsCRDT.ClientLatestClock {
			logger.InfoLogger.Println("repetitive model update sent by client id", modelWeightsCRDT.ClientId)
			return nil
		}
		if clientObservation > modelWeightsCRDT.ModelClock {
			logger.InfoLogger.Println("older model sent by client id", modelWeightsCRDT.ClientId)
			return nil
		}
	}
	m.crdtObjsLock.RLock()
	crdtObj, ok := m.crdtObjs[modelId]
	m.crdtObjsLock.RUnlock()
	if !ok {
		crdtObj, err = ReadInitModel(modelId)
		if err != nil {
			logger.ErrorLogger.Println(err)
			return err
		}
		m.crdtObjsLock.Lock()
		m.crdtObjs[modelId] = crdtObj
		m.crdtObjsLock.Unlock()
	}
	m.modifyFederated(crdtObj, modelWeightsCRDT, clientCount)
	m.clientsObservationLock.Lock()
	m.clientsObservation[modelWeightsCRDT.ClientId] += 1
	m.clientsObservationLock.Unlock()

	return nil
}

func (m *FederatedCRDTManager) modifyFederated(crdtObject *ModelWeightsCRDT, modelUpdate *ModelWeightsCRDT, clientCount int) {
	crdtObject.CRDTObjectLock.RLock()
	crdtObject.ModelClock++
	clockDiff := crdtObject.ModelClock - modelUpdate.ModelClock
	if clockDiff <= 0 {
		clockDiff = 1
	}
	clockPenalty := 1.0 / float64(clockDiff)
	clientAverage := 1.0 / float64(clientCount)
	totalPenalty := clockPenalty * clientAverage
	for layerIndex := range modelUpdate.Layers {
		if modelUpdate.Layers[layerIndex].LayerType == protos.LayerWeights_ONE_D {
			for i := range modelUpdate.Layers[layerIndex].LayerWeight1D {
				crdtObject.Layers[layerIndex].LayerWeight1D[i] += modelUpdate.Layers[layerIndex].LayerWeight1D[i] * totalPenalty
			}
		} else if modelUpdate.Layers[layerIndex].LayerType == protos.LayerWeights_TWO_D {
			for i := range modelUpdate.Layers[layerIndex].LayerWeights2D {
				for j := range modelUpdate.Layers[layerIndex].LayerWeights2D[i] {
					crdtObject.Layers[layerIndex].LayerWeights2D[i][j] += modelUpdate.Layers[layerIndex].LayerWeights2D[i][j] * totalPenalty
				}
			}
		}
	}
	crdtObject.CRDTObjectLock.RUnlock()
}

func (m *FederatedCRDTManager) GetModelCRDTByte(modelUpdate *protos.ModelUpdateRequest) ([]byte, int32, error) {
	crdtObjId := modelUpdate.ModelId
	modelShouldBeUpdated := false

	expirationTime := 100 * time.Millisecond

	var err error
	var crdtObjByte []byte

	if _, ok := m.latestModels[crdtObjId]; !ok {
		modelShouldBeUpdated = true
	} else {
		m.latestModelLock.Lock()
		if latestTime, ok := m.latestModelsUpdated[modelUpdate.ModelId]; !ok {
			m.latestModelsUpdated[modelUpdate.ModelId] = time.Now()
			modelShouldBeUpdated = true
		} else {
			if time.Since(latestTime) >= expirationTime {
				m.latestModelsUpdated[modelUpdate.ModelId] = time.Now()
				modelShouldBeUpdated = true
			}
		}
		m.latestModelLock.Unlock()
	}

	if modelShouldBeUpdated {
		m.crdtObjsLock.RLock()
		crdtObj, ok := m.crdtObjs[crdtObjId]
		if ok {
			crdtObjByte, err = ConvertWeightCRDTToProto(crdtObj)
		}
		m.crdtObjsLock.RUnlock()
		if err != nil {
			return nil, 0, err
		}
		if !ok {
			crdtObj, err = ReadInitModel(modelUpdate.ModelId)
			if err != nil {
				logger.ErrorLogger.Println(err)
				return nil, 0, err
			}
			m.crdtObjsLock.Lock()
			m.crdtObjs[crdtObjId] = crdtObj
			crdtObjByte, err = ConvertWeightCRDTToProto(crdtObj)
			m.crdtObjsLock.Unlock()
		}
		m.latestModels[crdtObjId] = crdtObjByte
		m.clientsObservationLock.RLock()
		_, ok = m.clientsObservation[modelUpdate.ClientId]
		m.clientsObservationLock.RUnlock()
		if !ok {
			m.clientsObservationLock.Lock()
			m.clientsObservation[modelUpdate.ClientId] = 0
			m.clientsObservationLock.Unlock()
		}
	}
	m.clientsObservationLock.RLock()
	latestClientObservation, _ := m.clientsObservation[modelUpdate.ClientId]
	m.clientsObservationLock.RUnlock()

	return m.latestModels[crdtObjId], latestClientObservation, nil
}

func ConvertWeightProtoToCRDT(weights []byte) (*ModelWeightsCRDT, error) {
	modelWeights := &protos.ModelWeights{}
	if err := proto.Unmarshal(weights, modelWeights); err != nil {
		logger.ErrorLogger.Println(err)
	}
	modelWeightsCRDT := &ModelWeightsCRDT{
		Layers:            []*LayerWeightsCRDT{},
		ModelClock:        modelWeights.ModelClock,
		ClientId:          modelWeights.ClientId,
		ClientLatestClock: modelWeights.ClientLatestClock,
		CRDTObjectLock:    &sync.RWMutex{},
	}
	for _, layer := range modelWeights.Layers {
		layerWeights := &protos.LayerWeights{}
		if err := proto.Unmarshal(layer, layerWeights); err != nil {
			logger.ErrorLogger.Println(err)
			return nil, err
		}
		layerWeightsCRDT := &LayerWeightsCRDT{
			LayerRank: layerWeights.LayerRank,
			LayerType: layerWeights.WeightType,
		}
		if layerWeights.WeightType == protos.LayerWeights_ONE_D {
			layerWeightsCRDT.LayerWeight1D = layerWeights.OneDWeights
		} else if layerWeights.WeightType == protos.LayerWeights_TWO_D {
			layerWeightsCRDT.LayerWeights2D = [][]float64{}
			for _, layer2Weight := range layerWeights.TwoDWeights {
				layer2WeightsFloat := &protos.LayerTwoDWeights{}
				if err := proto.Unmarshal(layer2Weight, layer2WeightsFloat); err != nil {
					logger.ErrorLogger.Println(err)
					return nil, err
				}
				layerWeightsCRDT.LayerWeights2D = append(layerWeightsCRDT.LayerWeights2D, layer2WeightsFloat.Weights)
			}
		} else {
			return nil, errors.New("proto to weight error, because the layer type is not found")
		}
		modelWeightsCRDT.Layers = append(modelWeightsCRDT.Layers, layerWeightsCRDT)
	}
	return modelWeightsCRDT, nil
}

func ConvertWeightCRDTToProto(modelWeightsCRDT *ModelWeightsCRDT) ([]byte, error) {
	modelWeights := &protos.ModelWeights{
		ModelClock: modelWeightsCRDT.ModelClock,
	}
	for _, layer := range modelWeightsCRDT.Layers {
		layerWeights := &protos.LayerWeights{
			WeightType: layer.LayerType,
			LayerRank:  layer.LayerRank,
		}
		if layer.LayerType == protos.LayerWeights_ONE_D {
			layerWeights.OneDWeights = layer.LayerWeight1D
		} else if layer.LayerType == protos.LayerWeights_TWO_D {
			layerWeights.TwoDWeights = [][]byte{}
			for _, layer2DWeights := range layer.LayerWeights2D {
				layer2DProtosByte, err := proto.Marshal(&protos.LayerTwoDWeights{
					Weights: layer2DWeights,
				})
				if err != nil {
					return nil, err
				}
				layerWeights.TwoDWeights = append(layerWeights.TwoDWeights, layer2DProtosByte)
			}
		} else {
			return nil, errors.New("proto to weight error, because the layer type is not found")
		}
		layerWeightsByte, err := proto.Marshal(layerWeights)
		if err != nil {
			return nil, err
		}
		modelWeights.Layers = append(modelWeights.Layers, layerWeightsByte)
	}
	modelWeightsByte, err := proto.Marshal(modelWeights)
	if err != nil {
		return nil, err
	}
	return modelWeightsByte, nil
}

func ReadInitModel(modelId string) (*ModelWeightsCRDT, error) {
	modelWeightsByte, err := ioutil.ReadFile(path.Join("../federated/federated_models/", modelId))
	if err != nil {
		modelWeightsByte, err = ioutil.ReadFile(path.Join("../../federatedpython/federated_models/", modelId))
		if err != nil {
			return nil, err
		}
	}
	return ConvertWeightProtoToCRDT(modelWeightsByte)
}
