package handler

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sales-product-srv/driver"
	"sales-product-srv/model"
	"sales-product-srv/proto"
	"sales-product-srv/utils"
)

type BannerServer struct {
}

func (b *BannerServer) GetBannerList(ctx context.Context, req *proto.BannerRequest) (*proto.BannerResponseList, error) {
	var count int64
	var banners []model.Banner
	var bannerList []*proto.BannerResponse
	result := driver.DB.Find(&banners)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "查询失败")
	}

	rsp := &proto.BannerResponseList{}

	// 分页
	driver.DB.Scopes(utils.Paginate(int(req.PageIndex), int(req.PageSize))).Find(&banners)
	//查询所有总数
	driver.DB.Model(&model.Brand{}).Count(&count)
	for _, value := range banners {
		bannerList = append(bannerList, &proto.BannerResponse{
			Id:    value.Id,
			Image: value.Image,
			Index: value.Index,
			Url:   value.Url,
		})
	}
	rsp.Total = int32(count)
	rsp.Data = bannerList

	return rsp, nil
}

// 删除轮播图

func (b *BannerServer) DeleteBanner(ctx context.Context, req *proto.DeleteBannerRequest) (*empty.Empty, error) {
	var banner model.Banner
	if result := driver.DB.Where("id=?", req.Id).Find(&banner); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}
	if result := driver.DB.Where("id=?", req.Id).Delete(&banner); result.RowsAffected != 1 {
		return nil, status.Errorf(codes.Internal, "删除失败")
	}
	return &empty.Empty{}, nil
}
