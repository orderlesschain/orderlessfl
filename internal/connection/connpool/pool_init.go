package connpool

import (
	"context"
	"crypto/tls"
	"gitlab.lrz.de/orderless/orderlessfl/internal/config"
	"gitlab.lrz.de/orderless/orderlessfl/internal/customcrypto/certificates"
	"gitlab.lrz.de/orderless/orderlessfl/internal/logger"
	protos "gitlab.lrz.de/orderless/orderlessfl/protos/goprotos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"net"
	"os"
	"strings"
	"time"
)

var KACP = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             10 * time.Second, // wait 30 second for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

func factoryNode(ip string) (*grpc.ClientConn, error) {
	configTLS := &tls.Config{
		InsecureSkipVerify: false,
		RootCAs:            certificates.CAs,
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(credentials.NewTLS(configTLS)), grpc.WithKeepaliveParams(KACP),
		grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip")),
		// https://chromium.googlesource.com/external/github.com/grpc/grpc-go/+/HEAD/Documentation/encoding.md
	}
	conn, err := grpc.Dial(ip+config.Config.TransactionServerPort, opts...)
	if err != nil {
		logger.FatalLogger.Fatalln("Failed to start gRPC connection:", err)
	}
	return conn, err
}

func GetNodeConnectionsForFederated(participatingNodes []string) map[string]*Pool {
	serverPool := make(map[string]*Pool)
	for _, name := range participatingNodes {
		factory := factoryNode
		pool, err := NewPoolWithIP(factory, config.Config.Nodes[name], config.Config.InsideContractConnectionPoolCount, config.Config.InsideContractConnectionPoolCount, 10*time.Second)
		if err != nil {
			logger.FatalLogger.Fatalln("Failed to create gRPC pool:", err)
		}
		serverPool[name] = pool
	}
	return serverPool
}

func GetNodeConnectionsStreamProposal(participatingNodes []string) map[string]*Pool {
	serverPool := make(map[string]*Pool)
	for _, name := range participatingNodes {
		factory := factoryNode
		pool, err := NewPoolWithIP(factory, config.Config.Nodes[name], config.Config.ProposalConnectionStreamPoolCount, config.Config.ProposalConnectionStreamPoolCount, 10*time.Second)
		if err != nil {
			logger.FatalLogger.Fatalln("Failed to create gRPC pool:", err)
		}
		serverPool[name] = pool
	}
	return serverPool
}

func GetSingleNodeConnectionsStreamProposal(participatingNode string) *Pool {
	factory := factoryNode
	pool, err := NewPoolWithIP(factory, config.Config.Nodes[participatingNode], config.Config.ProposalConnectionStreamPoolCount, config.Config.ProposalConnectionStreamPoolCount, 10*time.Second)
	if err != nil {
		logger.FatalLogger.Fatalln("Failed to create gRPC pool:", err)
	}
	return pool
}

func GetNodeConnectionsStreamTransactions(participatingNodes []string) map[string]*Pool {
	serverPool := make(map[string]*Pool)
	for _, name := range participatingNodes {
		factory := factoryNode
		pool, err := NewPoolWithIP(factory, config.Config.Nodes[name], config.Config.TransactionConnectionStreamPoolCount, config.Config.TransactionConnectionStreamPoolCount, 10*time.Second)
		if err != nil {
			logger.FatalLogger.Fatalln("Failed to create gRPC pool:", err)
		}
		serverPool[name] = pool
	}
	return serverPool
}

func GetSingleNodeConnectionsStreamTransactions(participatingNode string) *Pool {
	factory := factoryNode
	pool, err := NewPoolWithIP(factory, config.Config.Nodes[participatingNode], config.Config.TransactionConnectionStreamPoolCount, config.Config.TransactionConnectionStreamPoolCount, 10*time.Second)
	if err != nil {
		logger.FatalLogger.Fatalln("Failed to create gRPC pool:", err)
	}
	return pool
}

func GetNodeConnectionsWatchProposalEvent(participatingNodes []string) map[string]*Pool {
	serverPool := make(map[string]*Pool)
	for _, name := range participatingNodes {
		factory := factoryNode
		pool, err := NewPoolWithIP(factory, config.Config.Nodes[name], config.Config.ProposalEventConnectionPoolCount, config.Config.ProposalEventConnectionPoolCount, 10*time.Second)
		if err != nil {
			logger.FatalLogger.Fatalln("Failed to create gRPC pool:", err)
		}
		serverPool[name] = pool
	}
	return serverPool
}

func GetSingleNodeConnectionsWatchProposalEvent(participatingNode string) *Pool {
	factory := factoryNode
	pool, err := NewPoolWithIP(factory, config.Config.Nodes[participatingNode], config.Config.ProposalEventConnectionPoolCount, config.Config.ProposalEventConnectionPoolCount, 10*time.Second)
	if err != nil {
		logger.FatalLogger.Fatalln("Failed to create gRPC pool:", err)
	}
	return pool
}

func GetNodeConnectionsWatchTransactionEvent(participatingNodes []string) map[string]*Pool {
	serverPool := make(map[string]*Pool)
	for _, name := range participatingNodes {
		factory := factoryNode
		pool, err := NewPoolWithIP(factory, config.Config.Nodes[name], config.Config.TransactionEventConnectionPoolCount, config.Config.TransactionEventConnectionPoolCount, 10*time.Second)
		if err != nil {
			logger.FatalLogger.Fatalln("Failed to create gRPC pool:", err)
		}
		serverPool[name] = pool
	}
	return serverPool
}

func GetSingleNodeConnectionsWatchTransactionEvent(participatingNode string) *Pool {
	factory := factoryNode
	pool, err := NewPoolWithIP(factory, config.Config.Nodes[participatingNode], config.Config.TransactionEventConnectionPoolCount, config.Config.TransactionEventConnectionPoolCount, 10*time.Second)
	if err != nil {
		logger.FatalLogger.Fatalln("Failed to create gRPC pool:", err)
	}
	return pool
}

func GetNodeConnectionsWatchEventWithoutSelf(participatingNodes []string) map[string]*Pool {
	currentIP := GetOutboundIP()
	serverPool := make(map[string]*Pool)
	for _, name := range participatingNodes {
		nodeIP := config.Config.Nodes[name]
		if nodeIP != currentIP {
			factory := factoryNode
			pool, err := NewPoolWithIP(factory, nodeIP, config.Config.TransactionEventConnectionPoolCount, config.Config.TransactionEventConnectionPoolCount, 10*time.Second)
			if err != nil {
				logger.FatalLogger.Fatalln("Failed to create gRPC pool:", err)
			}
			serverPool[name] = pool
		}
	}
	return serverPool
}

func GetAllNodesConnections() map[string]*Pool {
	serverPool := make(map[string]*Pool)
	for name, nodeIP := range config.Config.Nodes {
		factory := factoryNode
		pool, err := NewPoolWithIP(factory, nodeIP, config.Config.GossipConnectionPoolCount, config.Config.GossipConnectionPoolCount, 10*time.Second)
		if err != nil {
			logger.FatalLogger.Fatalln("Failed to create gRPC pool:", err)
		}
		serverPool[name] = pool
	}
	return serverPool
}

func GetAllClientsConnections() map[string]*Pool {
	clientConnectionPool := make(map[string]*Pool)
	for name, clientIP := range config.Config.Clients {
		factory := func(ip string) (*grpc.ClientConn, error) {
			opts := []grpc.DialOption{grpc.WithKeepaliveParams(KACP), grpc.WithInsecure()}
			conn, err := grpc.Dial(ip+config.Config.TransactionServerPort, opts...)
			if err != nil {
				logger.FatalLogger.Fatalln("Failed to start gRPC connection:", err)
			}
			return conn, err
		}
		pool, err := NewPoolWithIP(factory, clientIP, config.Config.ClientConnectionPoolCount, config.Config.ClientConnectionPoolCount, 10*time.Second)
		if err != nil {
			logger.FatalLogger.Fatalln("Failed to create gRPC pool:", err)
		}
		clientConnectionPool[name] = pool
	}
	return clientConnectionPool
}

func GetFederatedLearningLocalHostConnections() *Pool {
	factory := func(ip string) (*grpc.ClientConn, error) {
		opts := []grpc.DialOption{grpc.WithKeepaliveParams(KACP), grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.UseCompressor("gzip"))}
		conn, err := grpc.Dial(ip+config.Config.FederatedServerPort, opts...)
		if err != nil {
			logger.FatalLogger.Fatalln("Failed to start gRPC connection:", err)
		}
		return conn, err
	}
	pool, err := NewPoolWithIP(factory, "localhost", config.Config.FederatedConnectionPoolCount, config.Config.FederatedConnectionPoolCount, 10*time.Second)
	if err != nil {
		logger.FatalLogger.Fatalln("Failed to create gRPC pool:", err)
	}
	return pool
}

func SleepAndReconnect() {
	time.Sleep(5 * time.Second)
}

func SleepAndContinue() {
	time.Sleep(2 * time.Second)
}

func GetOutboundIP() string {
	if len(config.Config.IPAddress) != 0 {
		return config.Config.IPAddress
	}
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
	defer func(conn net.Conn) {
		if errCon := conn.Close(); errCon != nil {
			logger.ErrorLogger.Println(errCon)
		}
	}(conn)
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func GetComponentPseudoName() string {
	currentIP := GetOutboundIP()
	for name, nodeIP := range config.Config.Nodes {
		if nodeIP == currentIP {
			return name
		}
	}
	for name, clientIP := range config.Config.Clients {
		if clientIP == currentIP {
			return name
		}
	}
	//logger.WarningLogger.Println("Current IP does not match. If local deployment, this is alright.")
	name, err := os.Hostname()
	if err != nil {
		logger.ErrorLogger.Println(err)
	}
	return strings.ReplaceAll(name, "-", "")
}

func RestartFederated() {
	coon := GetFederatedLearningLocalHostConnections()
	conn, err := coon.Get(context.Background())
	if conn == nil || err != nil {
		return
	}
	client := protos.NewFederatedServiceClient(conn.ClientConn)
	_, err = client.ChangeModeRestart(context.Background(), &protos.Empty{})
	if err != nil {
		return
	}
	if errCon := conn.Close(); errCon != nil {
		logger.ErrorLogger.Println(errCon)
	}
}
