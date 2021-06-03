package server

import (
	"fmt"
	"net"
	"time"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	mwLogging "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"eth-proxy/pkg/logger"
	"eth-proxy/proxy-service/config"
	"eth-proxy/utils"
)

// IGrpc is basic GRPC server interface
type IGrpc interface {
	Listen()
	Close()

	GetServer() *grpc.Server
}

// Server is Grpc Server instance
type Server struct {
	addr     string
	listener net.Listener
	*grpc.Server
}

// NewServer create new GRPC server instance
func NewServer(cfg *config.Grpc) (IGrpc, error) {
	server := &Server{addr: utils.CreateAddr(cfg.Host, cfg.Port)}

	if err := server.initListener(); err != nil {
		return nil, fmt.Errorf("GRPC server initializing: %w", err)
	}

	if err := server.initServer(cfg.Timeout); err != nil {
		return nil, fmt.Errorf("GRPC server iniitializing: %w", err)
	}

	return server, nil
}

// Listen open and listening incoming Tcp connections to Grpc Server port
func (s *Server) Listen() {
	logger.Log().Infof("listen and serve GRPC on %s", s.addr)
	if err := s.Serve(s.listener); err != nil {
		logger.Log().Fatalf("errored listening for grpc connections: %s", err)
	}
}

// GetServer return Grpc Server core instance
func (s *Server) GetServer() *grpc.Server {
	return s.Server
}

// Close closing Grpc Server listener and resetting connections
func (s *Server) Close() {
	if s.Server != nil {
		s.GracefulStop()
	}
	return
}

// initListener initialize GRPC server connections listener
func (s *Server) initListener() (err error) {
	s.listener, err = net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("failed to listen GRPC on addr %s: %w", s.addr, err)
	}
	return nil
}

// initServer initialize GRPC server core
func (s *Server) initServer(timeout string) error {
	duration, err := time.ParseDuration(timeout)
	if err != nil {
		return fmt.Errorf("initialize GRPC server: %w", err)
	}

	middlewares := middleware.ChainUnaryServer(
		mwLogging.UnaryServerInterceptor(logger.Log().WithField("layer", "GRPC server")),
		mwRecovery.UnaryServerInterceptor(),
	

	var kaep = keepalive.EnforcementPolicy{
		MinTime:             5 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}

	var kasp = keepalive.ServerParameters{
		MaxConnectionIdle:     15 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionAge:      30 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
		MaxConnectionAgeGrace: 5 * time.Second,  // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  5 * time.Second,  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               1 * time.Second,  // Wait 1 second for the ping ack before assuming the connection is dead
	}

	s.Server = grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(kaep),
		grpc.KeepaliveParams(kasp),
		grpc.ConnectionTimeout(duration),
		grpc.UnaryInterceptor(middlewares),
	)

	return nil
}
