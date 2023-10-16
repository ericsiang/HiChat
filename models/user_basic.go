package models

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type UserBasic struct {
	Model
	Name          string
	Password      string
	Avatar        string
	Gender        string `gorm:"gender;default:male;type:varchar(6);comment: 'male :男 , famale:女'"`
	Phone         string `valid:"match(^09\\d{8}$)"`
	Email         string `valid:"email"`
	Identify      string
	ClientIP      string `valid:"ipv4"`
	ClientPort    string
	Salt          string     //盐值
	LoginTime     *time.Time `gorm:"column:login_time"`
	HeartBeatTime *time.Time `gorm:"column:heart_beat_time"`
	LoginOutTime  *time.Time `gorm:"column:login_out_time"`
	IsLoginOut    bool
	DeviceInfo    string //登录设备
}


func (table *UserBasic) UserTableName() string {
	return "user_basic"
}