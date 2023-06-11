package grpc

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	name   string
	addr   string
	Logger *zap.Logger

	Server   *grpc.Server
	Listener net.Listener
}

func New(
	name string,
	addr string,
	logger *zap.Logger,
) *Server {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("failed to listen: %v", zap.Error(err))
	}

	n := "grpc server"
	if name != "" {
		n = name + " " + n
	}

	return &Server{
		name:   n,
		addr:   addr,
		Logger: logger,

		Server:   grpc.NewServer(),
		Listener: lis,
	}
}

func (s *Server) Run() {
	s.Logger.Info(s.name+" is about to listen", zap.String("addr", s.addr))

	// Register reflection service on gRPC server.
	reflection.Register(s.Server)

	if err := s.Server.Serve(s.Listener); err != nil {
		s.Logger.Fatal(s.name+" listen failed", zap.Error(err))
	}
}
