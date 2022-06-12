package main

import (
	"context"
	"fmt"
	"goods_srv/proto"
)

func TestGetBrandList() {
	reply, err := client.BrandList(context.Background(), &proto.BrandFilterRequest{Pages: 1, PagePerNums: 10})
	if err != nil {
		panic(err)
	}

	fmt.Println(reply.Total)

	for _, brand := range reply.Data {
		fmt.Println(brand)
	}
}
