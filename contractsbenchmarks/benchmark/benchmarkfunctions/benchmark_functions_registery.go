package benchmarkfunctions

import (
	"errors"
	"gitlab.lrz.de/orderless/orderlessfl/contractsbenchmarks/benchmark/benchmarkfunctions/mnistfederatedchain"
	"gitlab.lrz.de/orderless/orderlessfl/internal/contract/contractinterface"
	"gitlab.lrz.de/orderless/orderlessfl/internal/logger"
	"strings"
)

func GetBenchmarkFunctions(contactName string, benchmarkFunctionName string) (contractinterface.BenchmarkFunction, error) {
	contactName = strings.ToLower(contactName)
	benchmarkFunctionName = strings.ToLower(benchmarkFunctionName)
	switch contactName {
	case "mnistfederatedchaincontract":
		switch benchmarkFunctionName {
		case "getglobalmodelandtrainandsubmit":
			return mnistfederatedchain.GetGlobalModelAndTrainAndSubmit, nil
		case "getglobalmodelandsubmitpretrainedupdate":
			return mnistfederatedchain.GetGlobalModelAndSubmitPretrainedUpdate, nil
		}
	}
	logger.ErrorLogger.Println("smart contract function was not found")
	return nil, errors.New("functions not found")
}
