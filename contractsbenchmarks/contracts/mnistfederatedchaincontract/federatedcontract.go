package mnistfederatedchaincontract

import (
	"errors"
	"gitlab.lrz.de/orderless/orderlessfl/internal/contract/contractinterface"
	protos "gitlab.lrz.de/orderless/orderlessfl/protos/goprotos"
	"strings"
)

type MNIST struct {
}

func NewContract() *MNIST {
	return &MNIST{}
}

func (t *MNIST) Invoke(shim contractinterface.ShimInterface, proposal *protos.ProposalRequest) (*protos.ProposalResponse, error) {
	proposal.MethodName = strings.ToLower(proposal.MethodName)
	if proposal.MethodName == "getglobalmodelandtrainandsubmit" {
		return t.getGlobalModelAndTrainAndSubmit(shim, proposal.MethodParams, proposal.FederatedModelUpdateHash)
	}
	return shim.Error(), errors.New("method name not found")
}

func (t *MNIST) getGlobalModelAndTrainAndSubmit(shim contractinterface.ShimInterface, args []string, modelUpdate []byte) (*protos.ProposalResponse, error) {
	modelId := args[0]
	if modelUpdate != nil && len(modelUpdate) > 0 {
		shim.PutFederatedModelWeights(modelId, modelUpdate)
	}
	return shim.Success(), nil
}
