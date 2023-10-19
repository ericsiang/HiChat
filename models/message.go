package models

type Message struct {
	Model
	FromId      int64  `json:"userId"`   // 發送者id
	TargetId    int64  `json:"targetId"` // 接收者id
	Type        int    //聊天類型 1.群聊 2.私聊 3.廣播
	MessageType int    //信息類別 1.文字 2.圖片 3.音頻
	Content     string //消息內容
	Pic         string `json:"url"` //圖片地址
	Url         string //文件相關
	Desc        string //描述
	Amount      int    //其他數據大小
}

func (m *Message) TableName() string {
	return "message"
}