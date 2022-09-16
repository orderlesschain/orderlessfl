package contracts

import (
	"errors"
	"gitlab.lrz.de/orderless/orderlessfl/contractsbenchmarks/contracts/mnistfederatedchaincontract"
	"gitlab.lrz.de/orderless/orderlessfl/internal/contract/contractinterface"
	"strings"
)

func GetContract(contractName string) (contractinterface.ContractInterface, error) {
	contractName = strings.ToLower(contractName)
	switch contractName {
	case "mnistfederatedchaincontract":
		return mnistfederatedchaincontract.NewContract(), nil
	default:
		return nil, errors.New("contract not found")
	}
}

func GetContractNames() []string {
	return []string{
		"mnistfederatedchaincontract",
	}
}
