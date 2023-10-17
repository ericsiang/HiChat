package service

import (
	"HiChat/common"
	"HiChat/dao"
	"HiChat/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// user 對返回數據進行遮蔽
type user struct {
	Name     string
	Avatar   string
	Gender   string
	Phone    string
	Email    string
	Identify string
}

func FriendList(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Request.FormValue("userId"))
	users, err := dao.FriendList(uint(id))
	if err != nil {
		zap.S().Info("查詢好友列表失敗：", err)
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "好友列表為空",
		})
		return
	}

	infos := make([]user, 0)
	for _, v := range *users {
		info := user{
			Name:     v.Name,
			Avatar:   v.Avatar,
			Gender:   v.Gender,
			Phone:    v.Phone,
			Email:    v.Email,
			Identify: v.Identify,
		}

		infos = append(infos, info)
	}

	common.RespOKList(ctx.Writer, infos, len(infos))
}

// AddFriendByName 通过昵称加好友
func AddFriendByName(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.PostForm("userId"))

	if err != nil {
		zap.S().Info("類型轉換失敗:", err)
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "通过昵称加好友失败",
		})
		return
	}

	tar := ctx.PostForm("targetName")
	target, err := strconv.Atoi(tar)
	if err != nil {
		code, err := dao.AddFriendByName(uint(userId), tar)
		if err != nil {
			HandleErr(code, ctx, err)
			return
		}
	} else {
		code, err := dao.AddFriend(uint(userId), uint(target))
		if err != nil {
			HandleErr(code, ctx, err)
			return
		}
	}

	ctx.JSON(200, gin.H{
		"code":    0, //  0成功   -1失败
		"message": "添加好友成功",
	})
}

func HandleErr(code int, ctx *gin.Context, err error) {
	switch code {
	case -1:
		ctx.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": err.Error(),
		})
	case 0:
		ctx.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "好友已存在",
		})
	case -2:
		ctx.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "不能添加自己",
		})
	}
}

// NewGroup 新建群聊
func NewGroup(ctx *gin.Context) {
	owner := ctx.PostForm("ownerId")
	ownerId, err := strconv.Atoi(owner)
	if err != nil {
		zap.S().Info("owner類型轉換失敗:", err)
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "新建群聊失敗",
		})
		return
	}

	ty := ctx.PostForm("cate")
	_, err = strconv.Atoi(ty)
	if err != nil {
		zap.S().Info("ty類型轉換失敗:", err)
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "新建群聊失敗",
		})
		return
	}

	img := ctx.PostForm("icon")
	name := ctx.PostForm("name")
	desc := ctx.PostForm("desc")

	community := models.Community{}
	if ownerId == 0 {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "尚未登入",
		})
		return
	}

	community.Image = img
	community.Name = name
	community.Desc = desc
	community.OwnerId = uint(ownerId)

	code, err := dao.CreateCommunity(community)
	if err != nil {
		HandleErr(code, ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"code": 0, //  0成功   -1失败
		"msg":  "新建群聊成功",
	})

}

// GroupList 获取群列表
func GroupList(ctx *gin.Context) {
	owner := ctx.PostForm("ownerId")
	ownerId, err := strconv.Atoi(owner)
	if err != nil {
		zap.S().Info("owner類型轉換失敗:", err)
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "获取群列表失敗",
		})
		return
	}

	if ownerId == 0 {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "尚未登入",
		})
		return
	}

	community, err := dao.GetCommunityList(uint(ownerId))
	if err != nil {
		zap.S().Info("获取群列表失敗:", err)
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "你还没加入任何群聊",
		})
		return
	}

	common.RespOKList(ctx.Writer, community, len(*community))
}

// JoinGroup 加入群聊
func JoinGroup(ctx *gin.Context) {
	comName := ctx.PostForm("comName")


	user := ctx.PostForm("userId")
	userId, err := strconv.Atoi(user)
	if err != nil {
		zap.S().Info("userId 類型轉換失敗:", err)
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "加入群聊失敗",
		})
		return
	}

	if userId == 0 {
		ctx.JSON(200, gin.H{
			"code": -1,
			"msg":  "尚未登入",
		})
		return
	}

	code, err := dao.JoinCommunity(uint(userId),comName)
	if err != nil {
		HandleErr(code, ctx, err)
		return
	}

	ctx.JSON(200, gin.H{
		"code": 0, //  0成功   -1失败
		"msg":  "加入群聊成功",
	})
}
