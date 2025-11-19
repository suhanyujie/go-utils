package service

import (
	"context"

	pb "github.com/suhanyujie/go_utils/pkg/common_interface_vendor/pb_interface/api/platform/service/v1"
)

type FriendSvcService struct {
	pb.UnimplementedFriendSvcServer
}

func NewFriendSvcService() *FriendSvcService {
	return &FriendSvcService{}
}

func (s *FriendSvcService) CheckIsFriend(ctx context.Context, req *pb.CheckIsFriendReq) (*pb.CheckIsFriendResp, error) {
	return &pb.CheckIsFriendResp{}, nil
}
func (s *FriendSvcService) AddFriend(ctx context.Context, req *pb.AddFriendReq) (*pb.AddFriendResp, error) {
	return &pb.AddFriendResp{}, nil
}
func (s *FriendSvcService) GetFriendList(ctx context.Context, req *pb.GetFriendListReq) (*pb.GetFriendListResp, error) {
	return &pb.GetFriendListResp{}, nil
}
