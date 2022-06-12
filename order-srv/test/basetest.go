package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"order_srv/proto"
)

var conn *grpc.ClientConn
var client proto.OrderClient

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client = proto.NewOrderClient(conn)
}

func TestCreateCartItem(userID, Nums, goodsId int32) {
	rsp, err := client.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  userID,
		GoodsId: goodsId,
		Nums:    Nums,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id, "success")
}

func TestCartItemList(userID int32) {
	rsp, err := client.CartItemList(context.Background(), &proto.UserInfo{Id: userID})
	if err != nil {
		panic(err)
	}
	for _, item := range rsp.Data {
		fmt.Println(item)
	}
}

func TestUpdateCartItem(Id int32) {
	_, err := client.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		Id:      Id,
		Checked: true,
	})
	if err != nil {
		panic(err)
	}
}

func TestCreateOrder() {
	rsp, err := client.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  1,
		Address: "Los Angeles",
		Name:    "pipipengn",
		Mobile:  "125123123",
		Post:    "qwewqewq",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id, "success")
}

func TestGetOrderDetail(orderId int32) {
	rsp, err := client.OrderDetail(context.Background(), &proto.OrderRequest{
		UserId: 1,
		Id:     orderId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp)
}

func TestOrderList() {
	rsp, err := client.OrderList(context.Background(), &proto.OrderFilterRequest{UserId: 1})
	if err != nil {
		panic(err)
	}
	for _, datum := range rsp.Data {
		fmt.Println(datum)
	}
}

func main() {
	Init()
	//TestCreateCartItem(1, 1, 422)
	//TestUpdateCartItem(1)
	//TestCartItemList(1)

	//TestCreateOrder()
	//TestGetOrderDetail(3)

	TestOrderList()
}
