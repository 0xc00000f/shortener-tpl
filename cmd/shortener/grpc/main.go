package main

import (
	"os"
	"os/signal"

	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/cmd/shortener/grpc/grpc"
	pb "github.com/0xc00000f/shortener-tpl/cmd/shortener/grpc/grpc/proto"
)

func main() {
	s := grpc.New("status", ":3090", zap.L())
	pb.RegisterShortenerServiceServer(s.Server, grpc.ShortenerService{}) //nolint:exhaustruct

	go s.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	s.Server.GracefulStop()
}
