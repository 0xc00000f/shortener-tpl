package grpc

import (
	"context"
	"time"

	pb "github.com/0xc00000f/shortener-tpl/cmd/shortener/grpc/grpc/proto"
)

type ShortenerService struct {
	pb.UnimplementedShortenerServiceServer
}

func (ShortenerService) Status(context.Context, *pb.StatusReq) (*pb.StatusResp, error) {
	return &pb.StatusResp{Uptime: time.Now().Unix()}, nil
}
