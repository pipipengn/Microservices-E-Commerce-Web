package api

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"user-web/forms"
	"user-web/global"
	"user-web/global/response"
	"user-web/middlewares"
	"user-web/proto"
	"strconv"
	"time"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err == nil {
		return
	}
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"msg": e.Message(),
			})
		case codes.Internal:
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "Internaln Error",
			})
		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "Augument Error",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": e.Code(),
				"msg":  e.Message(),
			})
		}
	}
}

// GetUserList 获取所有用户
func GetUserList(c *gin.Context) {

	// jwt auth中间件set了
	userId, _ := c.Get("userId")
	zap.S().Info("当前用户Id：", userId)

	// 获取分页参数
	pn, _ := strconv.Atoi(c.DefaultQuery("pn", "0"))
	pSize, _ := strconv.Atoi(c.DefaultQuery("psize", "100"))

	// 远程调用
	userList, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pn),
		PSize: uint32(pSize),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询用户列表失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}

	// 返回结果
	var result []interface{}
	for _, value := range userList.Data {

		data := response.UserResponse{
			Id:       value.Id,
			NickName: value.Nickname,
			Birthday: time.Unix(int64(value.Birthaday), 0).Format("2006-01-02"),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}

		result = append(result, data)
	}

	c.JSON(http.StatusOK, result)
}

// PasswordLogin 用户登录
func PasswordLogin(c *gin.Context) {
	passwordLoginForm := forms.PasswordLoginForm{}

	// 表单验证
	if err := c.ShouldBind(&passwordLoginForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证码
	//if verify := store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true); !verify {
	//	c.JSON(http.StatusBadRequest, gin.H{
	//		"captcha": "验证码错误",
	//	})
	//	return
	//}

	// 查询用户信息 远程调用
	userInfo, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	})
	if err != nil {
		zap.S().Errorw("[PasswordLogin] 查询用户手机号失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}

	// 验证密码
	checked, _ := global.UserSrvClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
		Password:          passwordLoginForm.Password,
		EncryptedPassword: userInfo.Password,
	})

	if checked.Success {
		// 生成token
		j := middlewares.NewJWT()
		claims := middlewares.CustomClaims{
			ID:          uint(userInfo.Id),
			NickName:    userInfo.Nickname,
			AuthorityId: uint(userInfo.Role),
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix(),
				ExpiresAt: time.Now().Unix() + 60*60*24*30,
				Issuer:    "pipipengn",
			},
		}
		token, err := j.CreateToken(claims)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "生成token失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"msg":        "Login Successful",
			"id":         userInfo.Id,
			"nickname":   userInfo.Nickname,
			"token":      token,
			"expired_at": (time.Now().Unix() + 60*60*24*30) * 1000,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "Login Failed"})
}

// Register 用户注册
func Register(c *gin.Context) {
	registerForm := forms.RegisterForm{}

	// 表单验证
	if err := c.ShouldBind(&registerForm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 远程调用 创建用户
	newUser, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		Password: registerForm.Password,
		Mobile:   registerForm.Mobile,
	})
	if err != nil {
		zap.S().Errorw("[Register] 创建用户失败", "msg", err.Error())
		HandleGrpcErrorToHttp(err, c)
		return
	}
	c.JSON(http.StatusOK, newUser)
}
