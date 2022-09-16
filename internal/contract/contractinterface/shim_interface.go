package contractinterface

import (
	"gitlab.lrz.de/orderless/orderlessfl/internal/crdtmanagerv2"
	"gitlab.lrz.de/orderless/orderlessfl/internal/transactionprocessor/transactiondb"
	protos "gitlab.lrz.de/orderless/orderlessfl/protos/goprotos"
)

type SharedShimResources struct {
	DBConnections        map[string]*transactiondb.Operations
	FederatedCRDTManager *crdtmanagerv2.FederatedCRDTManager
}

type ShimInterface interface {
	GetKeyValueWithVersion(string) ([]byte, error)
	GetKeyValueWithVersionZero(string) ([]byte, error)
	GetKeyValueNoVersion(string) ([]byte, error)
	GetKeyRangeValueWithVersion(string) (map[string][]byte, error)
	GetKeyRangeValueWithVersionZero(string) (map[string][]byte, error)
	GetKeyRangeValueNoVersion(string) (map[string][]byte, error)
	PutKeyValue(string, []byte)
	GetSharedShimResources() *SharedShimResources
	SuccessWithOutput([]byte) *protos.ProposalResponse
	Success() *protos.ProposalResponse
	Error() *protos.ProposalResponse
	PutFederatedModelWeights(string, []byte)
	GetFederatedModelWeights(string) (*protos.ModelUpdateRequest, error)
}
