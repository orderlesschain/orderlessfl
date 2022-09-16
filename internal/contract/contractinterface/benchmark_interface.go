package contractinterface

import (
	"gitlab.lrz.de/orderless/orderlessfl/internal/benchmark/benchmarkutils"
	"gitlab.lrz.de/orderless/orderlessfl/internal/connection/connpool"
	protos "gitlab.lrz.de/orderless/orderlessfl/protos/goprotos"
	"sync"
)

type BenchmarkFunctionOptions struct {
	Counter                     int
	NodeId                      int
	BenchmarkUtils              *benchmarkutils.BenchmarkUtils
	CurrentClientPseudoId       int
	TotalTransactions           int
	SingleFunctionCounter       *SingleFunctionCounter
	CrdtObjectCount             string
	CrdtOperationPerObjectCount string
	CrdtObjectType              string
	NumberOfKeys                int
	NumberOfKeysSecond          int
	NumberOfKeysThird           int
	FederatedConnPool           *connpool.Pool
	NodesConnPool               map[string]*connpool.Pool
}

type SingleFunctionCounter struct {
	Lock    *sync.Mutex
	Counter map[int]int
}

func NewSingleFunctionCounter(numberOfKeysInt64 int64) *SingleFunctionCounter {
	numberOfKeys := int(numberOfKeysInt64)
	tempCounter := &SingleFunctionCounter{
		Lock:    &sync.Mutex{},
		Counter: map[int]int{},
	}
	for i := 0; i < numberOfKeys; i++ {
		tempCounter.Counter[i] = 0
	}
	return tempCounter
}

type BenchmarkFunctionOutputs struct {
	ContractMethodName   string
	Outputs              []string
	WriteReadType        protos.ProposalRequest_WriteReadTransaction
	FederatedModelUpdate []byte
	ClientClock          string
	ShouldBeSubmitted    bool
}

type BenchmarkFunction func(*BenchmarkFunctionOptions) *BenchmarkFunctionOutputs
