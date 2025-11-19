package service

import (
	"context"

	pb "github.com/suhanyujie/go_utils/pkg/common_interface_vendor/pb_interface/api/platform/service/v1"
)

type UserInfoSvcService struct {
	pb.UnimplementedUserInfoSvcServer
}

func NewUserInfoSvcService() *UserInfoSvcService {
	return &UserInfoSvcService{}
}

func (s *UserInfoSvcService) GetUserInfo(ctx context.Context, req *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	return &pb.GetUserInfoResp{}, nil
}
