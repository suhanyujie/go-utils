package service

import (
	"context"

	pb "github.com/suhanyujie/go_utils/pkg/common_interface_vendor/pb_interface/api/platform/service/v1"
)

type GameMixedSvcService struct {
	pb.UnimplementedGameMixedSvcServer
}

func NewGameMixedSvcService() *GameMixedSvcService {
	return &GameMixedSvcService{}
}

func (s *GameMixedSvcService) FarmEventPush(ctx context.Context, req *pb.FarmEventPushReq) (*pb.FarmEventPushResp, error) {
	return &pb.FarmEventPushResp{}, nil
}
func (s *GameMixedSvcService) GetUserNftList(ctx context.Context, req *pb.GetUserNftListReq) (*pb.GetUserNftListResp, error) {
	return &pb.GetUserNftListResp{}, nil
}
