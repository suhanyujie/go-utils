package service

import (
	"context"

	pb "github.com/suhanyujie/go_utils/pkg/common_interface_vendor/pb_interface/api/platform/service/v1"
)

type PlayerTaskSvcService struct {
	pb.UnimplementedPlayerTaskSvcServer
}

func NewPlayerTaskSvcService() *PlayerTaskSvcService {
	return &PlayerTaskSvcService{}
}

func (s *PlayerTaskSvcService) QueryTasksState(ctx context.Context, req *pb.QueryTasksStateReq) (*pb.QueryTasksStateResp, error) {
	return &pb.QueryTasksStateResp{}, nil
}
