package dao

import (
	"HiChat/common"
	"HiChat/global"
	"HiChat/models"
	"errors"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func GetUserList() ([]*models.UserBasic, err) {
	var list []*models.UserBasic
	if tx := global.DB.Find(&list); tx.RowsAffected == 0 {
		return nil, errors.New("没有用户")
	}

	return list, nil
}

func FindUserByNameAndPwd(name string, password string) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if tx := global.DB.Where("name=? and password=?", name, password).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("没有用户")
	}

	//token加密
	token := strconv.Itoa(int(time.Now().Unix()))
	temp := common.Md5encoder([]byte(token))
	if tx := global.DB.Model(&user).Where("id=? ", user.ID).Update("identify", temp); tx.RowsAffected == 0 {
		return nil, errors.New("寫入identify失敗")
	}

	return &user, nil

}

func FindUserByName(name string) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if tx := global.DB.Where("name=?", name).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("没有查詢到紀錄")
	}

	return &user, nil
}

func FindUser(name string) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if tx := global.DB.Where("name=?", name).First(&user); tx.RowsAffected == 1 {
		return nil, errors.New("用戶已存在")
	}

	return &user, nil
}

func FindUserByPhone(phone string) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if tx := global.DB.Where("phone=?", phone).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("没有查詢到紀錄")
	}

	return &user, nil
}

func FindUerByEmail(email string) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if tx := global.DB.Where("email=?", email).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("没有查詢到紀錄")
	}

	return &user, nil
}

func FindUserID(ID uint) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if tx := global.DB.Where(ID).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("没有查詢到紀錄")
	}

	return &user, nil
}

func CreateUser(user models.UserBasic) (*models.UserBasic, error) {
	tx := global.DB.Create(&user)
	if tx.RowsAffected==0 {
		zap.S().Info("新建用户失敗")
        return nil, errors.New("新增用户失敗")
	}
	return &user, nil
}

func UpdateUser(user models.UserBasic) (*models.UserBasic, error) {
	tx := global.DB.Model(&user).Updates(models.UserBasic{
		Name:     user.Name,
        Password: user.Password,
        Gender:   user.Gender,
        Phone:    user.Phone,
        Email:    user.Email,
        Avatar:   user.Avatar,
        Salt:     user.Salt,
	})
	if tx.RowsAffected==0 {
		zap.S().Info("更新用户失敗")
		return nil, errors.New("更新用户失敗")
	}
	return &user, nil
}

func DeleteUser(user models.UserBasic) error {
	if tx := global.DB.Delete(&user); tx.RowsAffected == 0 {
        zap.S().Info("删除失敗")
        return errors.New("删除用户失敗")
    }
    return nil
}
