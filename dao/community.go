package dao

import (
	"HiChat/global"
	"HiChat/models"
	"errors"
)

func CreateCommunity(community models.Community) (int, error) {
	com := models.Community{}
	if tx := global.DB.Where("name = ?", community.Name).First(&com); tx.RowsAffected == 1 {
		return -1, errors.New("community name already exists")
	}

	tx := global.DB.Begin()
	if tx := tx.Create(&community); tx.RowsAffected == 0 {
		tx.Rollback()
		return -1, errors.New("create community failed")
	}

	relation := models.Relation{}
	relation.OwnerId = community.OwnerId //群主id
	relation.TargetId = community.ID     //群id
	relation.Type = 2                    //群

	if tx := tx.Create(&relation); tx.RowsAffected == 0 {
		tx.Rollback()
		return -1, errors.New("create community relation failed")
	}

	tx.Commit()
	return 0, nil
}

func GetCommunityList(ownerid uint) (*[]models.Community, error) {
	//获取我加入的群列表
	relation := make([]models.Relation, 0)
	if tx := global.DB.Where("owner_id = ? and type = 2", ownerid).Find(&relation); tx.RowsAffected == 0 {
		return nil, errors.New("沒有加入任何群")
	}

	communityID := make([]uint, 0)
	for _, v := range relation {
		cid := v.TargetId
		communityID = append(communityID, cid)
	}

	community := make([]models.Community, 0)
	if tx := global.DB.Where("id in ?", communityID).Find(&community); tx.RowsAffected == 0 {
		return nil, errors.New("獲取群數據失敗")
	}

	return &community, nil
}

// JoinCommunity 根据群昵称搜索并加入群
func JoinCommunity(ownerId uint, cname string) (int, error) {
	community := models.Community{}
	if tx := global.DB.Where("name = ?", cname).First(&community); tx.RowsAffected == 0 {
		return -1, errors.New("群紀錄不存在")
	}

	relation := models.Relation{}
	if tx := global.DB.Where("owner_id = ? and target_id =? and type=2", ownerId, community.ID).First(&relation); tx.RowsAffected == 1 {
		return -1, errors.New("已加入群")
	}

	relation = models.Relation{}
	relation.OwnerId = ownerId
	relation.TargetId = community.ID
	relation.Type = 2

	if tx := global.DB.Create(&relation); tx.RowsAffected == 0 {
		return -1, errors.New("加入群失敗")
	}

	return 0, nil
}

// FindUsers 获取群成员id
func FindUsers(groupId uint) ([]uint, error) {
	relation := make([]models.Relation, 0)

	if tx := global.DB.Where("target_id = ? and type = 2", groupId).Find(&relation); tx.RowsAffected == 0 {
		return nil, errors.New("未在群中找到該成員")
	}

	userIDs := make([]uint, 0)
	for _, v := range relation {
		userId := v.OwnerId
		userIDs = append(userIDs, userId)
	}
	return userIDs, nil
}
