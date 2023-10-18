package service

import (
	"HiChat/common"
	"HiChat/dao"
	"HiChat/middleware"
	"HiChat/models"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func List(ctx *gin.Context) {
	list, err := dao.GetUserList()

	if err != nil {
		zap.S().Info("获取用户列表失败", err)
		ctx.JSON(200, gin.H{
			"code": -1, //0 表示成功， -1 表示失败
			"msg":  "获取用户列表失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, list)
}

func LoginByNameAndPassWord(ctx *gin.Context) {
	name := ctx.PostForm("name")
	password := ctx.PostForm("password")
	data, err := dao.FindUserByName(name)

	if err != nil {
		zap.S().Info("登录失败:", err)
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "登录失败",
		})
		return
	}

	if data.Name == "" {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "用户名不存在",
		})
		return
	}

	//由于数据库密码保存是使用md5密文的， 所以验证密码时，是将密码再次加密，然后进行对比，后期会讲解md:common.CheckPassWord
	ok := common.CheckPassWord(password, data.Salt, data.Password)

	if !ok {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "密碼錯誤",
		})
		return
	}

	resp , err := dao.FindUserByNameAndPwd(name, data.Password)
	if err!=nil{
		zap.S().Info("登录失败", err)
	}

	token, err := middleware.GenerateToken(data.ID, data.Identify)
	if err != nil {
		zap.S().Info("生成token失敗", err)
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "生成token失敗",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":   0,
		"msg":    "登录成功",
		"token":  token,
		"userId": resp.ID,
	})
}

func NewUser(ctx *gin.Context) {
	user := models.UserBasic{}
	user.Name = ctx.Request.FormValue("name")
	password := ctx.Request.FormValue("password")
	repassword := ctx.Request.FormValue("repassword")

	_, err := dao.FindUser(user.Name)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "用戶已存在",
		})
		return
	}

	if password != repassword {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "密碼不一致",
		})
		return
	}

	//生成盐值
	salt := fmt.Sprintf("%d", rand.Int31())

	//password 加密
	user.Password = common.SaltPassword(password, salt)
	user.Salt = salt
	t := time.Now()
	user.LoginTime = &t
	user.HeartBeatTime = &t
	user.LoginOutTime = &t
	data, err := dao.CreateUser(user)
	if err != nil {
		zap.S().Info("用戶創建失敗", err)
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "用戶創建失敗",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 0,
		"msg":  "注册成功",
		"data": data,
	})
}

func UpdataUser(ctx *gin.Context) {
	user := models.UserBasic{}

	id, err := strconv.Atoi(ctx.Request.FormValue("id"))
	if err != nil {
		zap.S().Info("id 類型轉換失敗", err)
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "修改帳號失敗",
		})
		return
	}

	user.ID = uint(id)
	user.Name = ctx.Request.FormValue("name")
	password := ctx.Request.FormValue("password")
	user.Avatar = ctx.Request.FormValue("avatar")
	user.Email = ctx.Request.FormValue("email")
	user.Phone = ctx.Request.FormValue("phone")
	user.Gender = ctx.Request.FormValue("gender")

	salt := fmt.Sprintf("%d", rand.Int31())
	user.Salt = salt
	user.Password = common.SaltPassword(password, salt)

	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		zap.S().Info("參數錯誤", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  "參數錯誤",
		})
		return
	}

	data, err := dao.UpdateUser(user)
	if err != nil {
		zap.S().Info("更新用戶失敗", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "更新用戶失敗",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "更新用戶成功",
		"data": data,
	})
}

func DeleteUser(ctx *gin.Context) {
	user := models.UserBasic{}

	id, err := strconv.Atoi(ctx.Request.FormValue("id"))
	if err != nil {
		zap.S().Info("id 類型轉換失敗", err)
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "註銷帳號失敗",
		})
		return
	}

	user.ID = uint(id)
	err = dao.DeleteUser(user)
	if err != nil {
		zap.S().Info("註銷用戶失敗", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "註銷用戶失敗",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "註銷用戶成功",
	})
}

//SendUserMsg 发送消息
func SendUserMsg(ctx *gin.Context) {
    models.Chat(ctx.Writer, ctx.Request)
}
