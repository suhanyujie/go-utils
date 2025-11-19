package service

import (
	"context"

	pb "github.com/suhanyujie/go_utils/pkg/common_interface_vendor/pb_interface/api/game/service/v1"
)

type GameSvcService struct {
	pb.UnimplementedGameSvcServer
}

func NewGameSvcService() *GameSvcService {
	return &GameSvcService{}
}

func (s *GameSvcService) CreateMatch(ctx context.Context, req *pb.CreateMatchReq) (*pb.CreateMatchResp, error) {
	return &pb.CreateMatchResp{}, nil
}
func (s *GameSvcService) CheckCanReconnect(ctx context.Context, req *pb.CheckCanReconnectReq) (*pb.CheckCanReconnectResp, error) {
	return &pb.CheckCanReconnectResp{}, nil
}
func (s *GameSvcService) PlayerExit(ctx context.Context, req *pb.PlayerExitReq) (*pb.PlayerExitResp, error) {
	return &pb.PlayerExitResp{}, nil
}
