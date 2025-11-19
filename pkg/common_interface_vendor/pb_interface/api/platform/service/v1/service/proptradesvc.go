package service

import (
	"context"

	pb "github.com/suhanyujie/go_utils/pkg/common_interface_vendor/pb_interface/api/platform/service/v1"
)

type PropTradeSvcService struct {
	pb.UnimplementedPropTradeSvcServer
}

func NewPropTradeSvcService() *PropTradeSvcService {
	return &PropTradeSvcService{}
}

func (s *PropTradeSvcService) ConsumeProp(ctx context.Context, req *pb.ConsumePropReq) (*pb.ConsumePropResp, error) {
	return &pb.ConsumePropResp{}, nil
}
func (s *PropTradeSvcService) QueryProp(ctx context.Context, req *pb.QueryPropReq) (*pb.QueryPropResp, error) {
	return &pb.QueryPropResp{}, nil
}
func (s *PropTradeSvcService) AddProp(ctx context.Context, req *pb.AddPropReq) (*pb.AddPropResp, error) {
	return &pb.AddPropResp{}, nil
}
func (s *PropTradeSvcService) QueryOrderStatus(ctx context.Context, req *pb.QueryOrderStatusReq) (*pb.QueryOrderStatusResp, error) {
	return &pb.QueryOrderStatusResp{}, nil
}
func (s *PropTradeSvcService) CreateOneGoods(ctx context.Context, req *pb.CreateOneGoodsReq) (*pb.CreateOneGoodsResp, error) {
	return &pb.CreateOneGoodsResp{}, nil
}
func (s *PropTradeSvcService) VerifyPayOrder(ctx context.Context, req *pb.VerifyPayOrderReq) (*pb.VerifyPayOrderResp, error) {
	return &pb.VerifyPayOrderResp{}, nil
}
func (s *PropTradeSvcService) OrderDeliver(ctx context.Context, req *pb.OrderDeliverReq) (*pb.OrderDeliverResp, error) {
	return &pb.OrderDeliverResp{}, nil
}
