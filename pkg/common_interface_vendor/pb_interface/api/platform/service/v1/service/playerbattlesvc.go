package service

import (
	"context"

	pb "github.com/suhanyujie/go_utils/pkg/common_interface_vendor/pb_interface/api/platform/service/v1"
)

type PlayerBattleSvcService struct {
	pb.UnimplementedPlayerBattleSvcServer
}

func NewPlayerBattleSvcService() *PlayerBattleSvcService {
	return &PlayerBattleSvcService{}
}

func (s *PlayerBattleSvcService) GetUserBattleInfo(ctx context.Context, req *pb.GetUserBattleInfoReq) (*pb.GetUserBattleInfoResp, error) {
	return &pb.GetUserBattleInfoResp{}, nil
}
func (s *PlayerBattleSvcService) GetBestRecordOfFriendBattle(ctx context.Context, req *pb.GetBestRecordOfFriendBattleReq) (*pb.GetBestRecordOfFriendBattleResp, error) {
	return &pb.GetBestRecordOfFriendBattleResp{}, nil
}
func (s *PlayerBattleSvcService) GetNftTopNInfo(ctx context.Context, req *pb.GetNftTopNInfoReq) (*pb.GetNftTopNInfoResp, error) {
	return &pb.GetNftTopNInfoResp{}, nil
}
func (s *PlayerBattleSvcService) GetSurpassRateByScore(ctx context.Context, req *pb.GetSurpassRateByScoreReq) (*pb.GetSurpassRateByScoreResp, error) {
	return &pb.GetSurpassRateByScoreResp{}, nil
}
func (s *PlayerBattleSvcService) GetRuneRate(ctx context.Context, req *pb.GetRuneRateReq) (*pb.GetRuneRateResp, error) {
	return &pb.GetRuneRateResp{}, nil
}
func (s *PlayerBattleSvcService) GetDefeatNftTopNInfo(ctx context.Context, req *pb.GetDefeatNftTopNInfoReq) (*pb.GetDefeatNftTopNInfoResp, error) {
	return &pb.GetDefeatNftTopNInfoResp{}, nil
}
func (s *PlayerBattleSvcService) GetInfoOfFriendHasUsingNft(ctx context.Context, req *pb.GetInfoOfFriendHasUsingNftReq) (*pb.GetInfoOfFriendHasUsingNftResp, error) {
	return &pb.GetInfoOfFriendHasUsingNftResp{}, nil
}
func (s *PlayerBattleSvcService) GetNftUsingRateInfo(ctx context.Context, req *pb.GetNftUsingRateInfoReq) (*pb.GetNftUsingRateInfoResp, error) {
	return &pb.GetNftUsingRateInfoResp{}, nil
}
