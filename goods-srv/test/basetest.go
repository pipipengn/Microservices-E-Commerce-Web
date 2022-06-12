package main

import (
	"context"
	"fmt"
	"goods_srv/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

var conn *grpc.ClientConn
var client proto.GoodsClient

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client = proto.NewGoodsClient(conn)
}

func TestGetAllCategorysList() {
	reply, err := client.GetAllCategorysList(context.Background(), &emptypb.Empty{})
	if err != nil {
		panic(err)
	}

	fmt.Println(reply.JsonData)
}

func TestGetSubCategory() {
	reply, err := client.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: 135487,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(reply.SubCategorys)
}

func TestCategoryBrandList() {
	reply, err := client.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{
		Pages:       1,
		PagePerNums: 5,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(reply.Data)
}

func TestGoodsList() {
	reply, err := client.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		TopCategory: 130361,
		PriceMin:    90,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(reply.Total)

	for _, good := range reply.Data {
		fmt.Println(good.Name, good.ShopPrice)
	}
}

func TestBatchGetGoods() {
	reply, err := client.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: []int32{421, 422, 423},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(reply.Total)

	for _, good := range reply.Data {
		fmt.Println(good.Name, good.ShopPrice)
	}
}

func TestGetGoodsDetail() {
	reply, err := client.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{Id: 421})
	if err != nil {
		panic(err)
	}
	fmt.Println(reply.Name)
}

func main() {
	Init()
	//TestGetBrandList()
	TestGetAllCategorysList()
	//TestGetSubCategory()
	//TestCategoryBrandList()
	//TestGoodsList()

	//TestBatchGetGoods()
	//TestGetGoodsDetail()
}
