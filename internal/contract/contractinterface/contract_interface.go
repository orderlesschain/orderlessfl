package contractinterface

import (
	"gitlab.lrz.de/orderless/orderlessfl/internal/benchmark/benchmarkutils"
	"gitlab.lrz.de/orderless/orderlessfl/internal/connection/connpool"
	protos "gitlab.lrz.de/orderless/orderlessfl/protos/goprotos"
)

type ContractInterface interface {
	Invoke(ShimInterface, *protos.ProposalRequest) (*protos.ProposalResponse, error)
}

type BaseContractOptions struct {
	Bconfig               *protos.BenchmarkConfig
	BenchmarkFunction     BenchmarkFunction
	BenchmarkUtils        *benchmarkutils.BenchmarkUtils
	CurrentClientPseudoId int
	TotalTransactions     int
	SingleFunctionCounter *SingleFunctionCounter
	NodesConnectionPool   map[string]*connpool.Pool
	FederatedConnPool     *connpool.Pool
}
