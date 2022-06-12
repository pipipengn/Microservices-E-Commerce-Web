package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"strings"
	"time"
	"user_srv/global"
	"user_srv/model"
	"user_srv/proto"
)

type UserServer struct {
	*proto.UnimplementedUserServer
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		if page == 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func ModelToResponse(user model.User) proto.UserInfoResponse {
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		Mobile:   user.Mobile,
		Password: user.Password,
		Nickname: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
	}
	if user.Birthday != nil {
		userInfoRsp.Birthaday = uint64(user.Birthday.Unix())
	}

	return userInfoRsp
}

func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	rsp := &proto.UserListResponse{Total: uint32(result.RowsAffected)}
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users)

	var protoUsers []*proto.UserInfoResponse
	for _, user := range users {
		userInfoResponse := ModelToResponse(user)
		protoUsers = append(protoUsers, &userInfoResponse)
	}
	rsp.Data = protoUsers

	return rsp, nil
}

func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)

	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "user does not exist")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	userInfoResponse := ModelToResponse(user)
	return &userInfoResponse, nil
}

func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)

	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "user does not exist")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	userInfoResponse := ModelToResponse(user)
	return &userInfoResponse, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).Find(&user)
	if result.RowsAffected != 0 {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName

	options := &password.Options{SaltLen: 10, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt, encodedPwd := password.Encode(req.Password, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}

	userInfoResponse := ModelToResponse(user)
	return &userInfoResponse, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "user does not exist")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	birthday := time.Unix(int64(req.Birthaday), 0)
	user.NickName = req.Nickname
	user.Birthday = &birthday
	user.Gender = req.Gender

	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}
	return &empty.Empty{}, nil
}

func (s *UserServer) CheckPassword(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	options := &password.Options{SaltLen: 10, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	check := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options)
	return &proto.CheckResponse{Success: check}, nil
}
