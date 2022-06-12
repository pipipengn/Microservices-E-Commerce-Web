package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"user_srv/proto"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	userClient = proto.NewUserClient(conn)
}

func TestGetUserList() {
	reply, err := userClient.GetUserList(context.Background(), &proto.PageInfo{Pn: 1, PSize: 5})
	if err != nil {
		panic(err)
	}

	for _, user := range reply.Data {
		fmt.Println(user)
		checkResponse, _ := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.Password,
		})
		fmt.Println(checkResponse.Success)
	}
}

func TestCreateUser() {
	reply, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Password: "qwe123",
		Mobile:   "2136897847",
		NickName: "pipipengn",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(reply)
}

func main() {
	//testing.Init()
	//defer func(conn *grpc.ClientConn) {
	//	err := conn.Close()
	//	if err != nil {
	//
	//	}
	//}(conn)
	////TestGetUserList()
	//TestCreateUser()

	m := map[int]int{1: 1}
	if _, ok := m[1]; ok {
		println("cun zai")
	} else {
		println("bu cun zai")
	}

}
