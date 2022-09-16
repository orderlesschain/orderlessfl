package contract

import (
	"gitlab.lrz.de/orderless/orderlessfl/internal/config"
	"gitlab.lrz.de/orderless/orderlessfl/internal/contract/contractinterface"
	"gitlab.lrz.de/orderless/orderlessfl/internal/transactionprocessor/transactiondb"
	protos "gitlab.lrz.de/orderless/orderlessfl/protos/goprotos"
)

type Shim struct {
	proposalRequest  *protos.ProposalRequest
	proposalResponse *protos.ProposalResponse
	sharedResources  *contractinterface.SharedShimResources
	dbOp             *transactiondb.Operations
	contractName     string
}

func NewShim(proposal *protos.ProposalRequest, sharedResources *contractinterface.SharedShimResources, contractName string) *Shim {
	return &Shim{
		proposalRequest: proposal,
		proposalResponse: &protos.ProposalResponse{
			ProposalId: proposal.ProposalId,
			NodeId:     config.Config.UUID,
			ReadWriteSet: &protos.ProposalResponseReadWriteSet{
				ReadKeys:       &protos.ProposalResponseReadKeys{ReadKeys: []*protos.ReadKey{}},
				WriteKeyValues: &protos.ProposalResponseWriteKeyValues{WriteKeyValues: []*protos.WriteKeyValue{}},
			},
			NodeSignature: []byte{},
		},
		sharedResources: sharedResources,
		dbOp:            sharedResources.DBConnections[contractName],
		contractName:    contractName,
	}
}

func (s *Shim) GetSharedShimResources() *contractinterface.SharedShimResources {
	return s.sharedResources
}

func (s *Shim) GetKeyValueWithVersion(key string) ([]byte, error) {
	keyValueORM, err := s.dbOp.GetKeyValueWithVersion(key)
	if err != nil {
		return nil, err
	}
	s.proposalResponse.ReadWriteSet.ReadKeys.ReadKeys = append(s.proposalResponse.ReadWriteSet.ReadKeys.ReadKeys, &protos.ReadKey{
		Key:           key,
		VersionNumber: keyValueORM.VersionNumber,
	})
	return keyValueORM.Value, nil
}

func (s *Shim) GetKeyValueWithVersionZero(key string) ([]byte, error) {
	keyValueORM, err := s.dbOp.GetKeyValueWithVersion(key)
	if err != nil {
		return nil, err
	}
	s.proposalResponse.ReadWriteSet.ReadKeys.ReadKeys = append(s.proposalResponse.ReadWriteSet.ReadKeys.ReadKeys, &protos.ReadKey{
		Key:           key,
		VersionNumber: 0,
	})
	return keyValueORM.Value, nil
}

func (s *Shim) GetKeyValueNoVersion(key string) ([]byte, error) {
	keyValueByte, err := s.dbOp.GetKeyValueNoVersion(key)
	if err != nil {
		return nil, err
	}
	s.proposalResponse.ReadWriteSet.ReadKeys.ReadKeys = append(s.proposalResponse.ReadWriteSet.ReadKeys.ReadKeys, &protos.ReadKey{
		Key:           key,
		VersionNumber: 0,
	})
	return keyValueByte, nil
}

func (s *Shim) GetKeyRangeValueWithVersion(key string) (map[string][]byte, error) {
	keyValues := map[string][]byte{}
	keyValueORMs, err := s.dbOp.GetKeyRangeValueWithVersion(key)
	if err != nil {
		return nil, err
	}
	for keyValueORMKey, keyValueORM := range keyValueORMs {
		s.proposalResponse.ReadWriteSet.ReadKeys.ReadKeys = append(s.proposalResponse.ReadWriteSet.ReadKeys.ReadKeys, &protos.ReadKey{
			Key:           keyValueORMKey,
			VersionNumber: keyValueORM.VersionNumber,
		})
		keyValues[keyValueORMKey] = keyValueORM.Value
	}
	return keyValues, nil
}

func (s *Shim) GetKeyRangeValueWithVersionZero(key string) (map[string][]byte, error) {
	keyValues := map[string][]byte{}
	keyValueORMs, err := s.dbOp.GetKeyRangeValueWithVersion(key)
	if err != nil {
		return nil, err
	}
	for keyValueORMKey, keyValueORM := range keyValueORMs {
		s.proposalResponse.ReadWriteSet.ReadKeys.ReadKeys = append(s.proposalResponse.ReadWriteSet.ReadKeys.ReadKeys, &protos.ReadKey{
			Key:           keyValueORMKey,
			VersionNumber: 0,
		})
		keyValues[keyValueORMKey] = keyValueORM.Value
	}
	return keyValues, nil
}

func (s *Shim) GetKeyRangeValueNoVersion(key string) (map[string][]byte, error) {
	keyValueORMs, err := s.dbOp.GetKeyRangeNoVersion(key)
	if err != nil {
		return nil, err
	}
	s.proposalResponse.ReadWriteSet.ReadKeys.ReadKeys = append(s.proposalResponse.ReadWriteSet.ReadKeys.ReadKeys, &protos.ReadKey{
		Key:           key,
		VersionNumber: 0,
	})
	return keyValueORMs, nil
}

func (s *Shim) PutKeyValue(key string, value []byte) {
	s.proposalResponse.ReadWriteSet.WriteKeyValues.WriteKeyValues = append(s.proposalResponse.ReadWriteSet.WriteKeyValues.WriteKeyValues, &protos.WriteKeyValue{
		Key:       key,
		Value:     value,
		WriteType: protos.WriteKeyValue_BINARY,
	})
}

func (s *Shim) SuccessWithOutput(output []byte) *protos.ProposalResponse {
	s.proposalResponse.Status = protos.ProposalResponse_SUCCESS
	s.proposalResponse.ShimOutput = output
	return s.proposalResponse
}

func (s *Shim) Success() *protos.ProposalResponse {
	s.proposalResponse.Status = protos.ProposalResponse_SUCCESS
	return s.proposalResponse
}

func (s *Shim) Error() *protos.ProposalResponse {
	s.proposalResponse.Status = protos.ProposalResponse_FAIL
	return s.proposalResponse
}

func (s *Shim) PutFederatedModelWeights(key string, modelUpdateBinary []byte) {
	s.proposalResponse.ReadWriteSet.WriteKeyValues.WriteKeyValues = append(s.proposalResponse.ReadWriteSet.WriteKeyValues.WriteKeyValues, &protos.WriteKeyValue{
		Key:       key,
		Value:     modelUpdateBinary,
		WriteType: protos.WriteKeyValue_FEDERATED,
	})
}

func (s *Shim) GetFederatedModelWeights(key string) (*protos.ModelUpdateRequest, error) {
	modelUpdates, _, err := s.sharedResources.FederatedCRDTManager.GetModelCRDTByte(&protos.ModelUpdateRequest{ModelId: key})
	if err != nil {
		return nil, err
	}
	s.proposalResponse.ReadWriteSet.ReadKeys.ReadKeys = append(s.proposalResponse.ReadWriteSet.ReadKeys.ReadKeys, &protos.ReadKey{
		Key:           key,
		VersionNumber: 0,
	})
	return &protos.ModelUpdateRequest{
		ModelId:      key,
		LayersWeight: modelUpdates,
	}, nil
}
