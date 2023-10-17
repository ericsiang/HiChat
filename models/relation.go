package models

type Relation struct {
	Model
	OwnerId   uint //誰的關係信息
	TargetId uint //對應的誰(加入的群id)
	Type     int  //關係類型 1 好友 2 群
	Desc     string
}

func (r *Relation) TableName() string {
	return "relaction"
}