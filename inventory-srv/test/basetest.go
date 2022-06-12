package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"inventory_srv/proto"
	"sync"
)

var conn *grpc.ClientConn
var client proto.InventoryClient

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	client = proto.NewInventoryClient(conn)
}

func TestSetInv(goodsId, Num int32) {
	if _, err := client.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
		Num:     Num,
	}); err != nil {
		panic(err)
	}
	fmt.Println("success")
}

func TestInvDetail(goodsId int32) {
	resp, err := client.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: goodsId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}

func TestSell(wg *sync.WaitGroup) {
	/*
		1.第一件扣减成功，第二件：1.没有库存信息，库存不足
		2.两件都成功
	*/
	_, err := client.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 1},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("success")
	wg.Done()
}

func TestReback() {
	_, err := client.Reback(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 10},
			{GoodsId: 422, Num: 30},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("success")
}

func main() {
	Init()
	//TestSetInv(422, 40)
	//TestInvDetail(421)
	//TestSell()
	//TestReback()

	//for i := 421; i < 841; i++ {
	//	TestSetInv(int32(i), 100)
	//}

	var wg sync.WaitGroup
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go TestSell(&wg)
	}
	wg.Wait()
}
