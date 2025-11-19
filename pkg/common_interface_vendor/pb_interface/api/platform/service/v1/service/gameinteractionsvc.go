package service

import (
	"context"

	pb "github.com/suhanyujie/go_utils/pkg/common_interface_vendor/pb_interface/api/platform/service/v1"
)

type GameInteractionSvcService struct {
	pb.UnimplementedGameInteractionSvcServer
}

func NewGameInteractionSvcService() *GameInteractionSvcService {
	return &GameInteractionSvcService{}
}

func (s *GameInteractionSvcService) ConsumeEnergy(ctx context.Context, req *pb.ConsumeEnergyReq) (*pb.ConsumeEnergyResp, error) {
	return &pb.ConsumeEnergyResp{}, nil
}
func (s *GameInteractionSvcService) GameCallback(ctx context.Context, req *pb.GameCallbackReq) (*pb.GameCallbackResp, error) {
	return &pb.GameCallbackResp{}, nil
}
func (s *GameInteractionSvcService) GameSessCallback(ctx context.Context, req *pb.GameSessCallbackReq) (*pb.GameSessCallbackResp, error) {
	return &pb.GameSessCallbackResp{}, nil
}
func (s *GameInteractionSvcService) PvpKnockoutSettle(ctx context.Context, req *pb.PvpKnockoutSettleReq) (*pb.PvpKnockoutSettleResp, error) {
	return &pb.PvpKnockoutSettleResp{}, nil
}
func (s *GameInteractionSvcService) GamePveNext(ctx context.Context, req *pb.GamePveNextReq) (*pb.GamePveNextResp, error) {
	return &pb.GamePveNextResp{}, nil
}
