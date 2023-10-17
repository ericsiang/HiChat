package dao

import (
	"HiChat/global"
	"HiChat/models"
	"errors"

	"go.uber.org/zap"
)

func FriendList(userId uint) (*[]models.UserBasic, error) {
	relation := make([]models.Relation, 0)
	if tx := global.DB.Where("ower_id =? and type=1", userId).Find(&relation); tx.RowsAffected == 0 {
		return nil, errors.New("没有查詢到紀錄")
	}

	userID := make([]uint, 0)
	for _, v := range relation {
		userID = append(userID, v.TargetId)
	}

	user := make([]models.UserBasic, 0)
	if tx := global.DB.Where("id in ?", userID).Find(&user); tx.RowsAffected == 0 {
		zap.S().Info("未查詢到Relation好友關係")
		return nil, errors.New("没有查詢到好友紀錄")
	}

	return &user, nil
}

// AddFriend 加好友
func AddFriend(userID, TargetId uint) (int, error) {
	if userID == TargetId {
		return -2, errors.New("不能添加自己")
	}

	targetUser, err := FindUserByID(TargetId)
	if err != nil {
		return -1, errors.New("未查詢到用戶")
	}

	if targetUser.ID == 0 {
		zap.S().Info("未查詢到用戶")
		return -1, errors.New("未查詢到用戶")
	}

	relation := models.Relation{}

	if tx := global.DB.Where("ower_id = ? and target_id = ? and type = 1", userID, TargetId).First(&relation); tx.RowsAffected == 1 {
		zap.S().Info("好友已存在")
		return 0, errors.New("好友已存在")
	}

	if tx := global.DB.Where("ower_id = ? and target_id = ? and type = 1", TargetId, userID).First(&relation); tx.RowsAffected == 1 {
		zap.S().Info("好友已存在")
		return 0, errors.New("好友已存在")
	}

	//開啟事務
	tx := global.DB.Begin()
	relation.OwnerId = userID
	relation.TargetId = TargetId
	relation.Type = 1

	if err := tx.Create(&relation).Error; err != nil {
		zap.S().Info("添加好友失敗:", err)
		//回滾事務
		tx.Rollback()
		return -1, errors.New("添加好友失敗")
	}

	relation = models.Relation{}
	relation.OwnerId = TargetId
	relation.TargetId = userID
	relation.Type = 1

	if err := tx.Create(&relation).Error; err != nil {
		zap.S().Info("添加好友失敗:", err)
		//回滾事務
		tx.Rollback()
		return -1, errors.New("添加好友失敗")
	}

	//提交事務
	tx.Commit()
	return 1, nil
}

// AddFriendByName 昵称加好友
func AddFriendByName(userId uint, targetName string) (int, error) {
	user, err := FindUserByName(targetName)
	if err != nil {
		zap.S().Info("該用戶不存在:", err)
		return -1, errors.New("該用戶不存在")
	}

	if user.ID == 0 {
		zap.S().Info("該用戶不存在")
		return -1, errors.New("該用戶不存在")
	}

	return AddFriend(userId, user.ID)
}
