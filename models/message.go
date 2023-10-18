package models

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"gopkg.in/fatih/set.v0"
)

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

type Node struct {
	Conn      *websocket.Conn // socket連接
	Addr      string          //客戶端地址
	DataQueue chan []byte     //消息內容的數據管道
	GroupSets set.Interface   //群組的集合 好友,群
}

// 映射關係
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 讀寫鎖 綁定Node時需要
var rwLocker sync.RWMutex

// 全局channel
var upSendChan chan []byte = make(chan []byte, 1024)

// init 初始化
func init() {
	go UdpSendProc()
	go UdpRecieveProc()
}

// UdpSendProc 完成upd数据发送, 连接到udp服务端，将全局channel中的消息体，写入udp服务端
func UdpSendProc() {
	udpConn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 3000,
		Zone: "Asia/Taipei",
	})
	if err != nil {
		zap.S().Info("udp连接失败", err)
		return
	}

	defer udpConn.Close()

	for {
		select {
		case data := <-upSendChan:
			_, err := udpConn.Write(data)
			if err != nil {
				zap.S().Info("udp发送失败", err)
				return
			}
		}
	}
}

// UpdRecProc 完成udp数据的接收，启动udp服务，获取udp客户端的写入的消息
func UdpRecieveProc() {
	udpConn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 3000,
	})

	if err != nil {
		zap.S().Info("監聽udp端口失败", err)
		return
	}
	defer udpConn.Close()

	for {
		var buf [1024]byte
		n, addr, err := udpConn.ReadFromUDP(buf[:])
		if err != nil {
			zap.S().Info("接收udp消息失败", err)
			return
		}

		//處理發送邏輯
		dispatch(buf[0:n])
	}
}

// dispatch 解析消息，聊天类型判断
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		zap.S().Info("服務端解析json消息失敗:", err)
		return
	}

	switch msg.Type {
	case 1: //私聊
		sendMsg(msg.TargetId, data)
	case 2: //群聊
		sendGroupMsg(uint(msg.FromId), uint(msg.TargetId), data)
	}
}

// sendMs 向用户单聊发送消息
func sendMsg(targetId int64, data []byte) {
	rwLocker.RLock()
	node, ok := clientMap[targetId]
	rwLocker.RUnlock()
	if !ok {
		zap.S().Info("userId不存在對應的node")
		return
	}
	zap.S().Info("targetId :", targetId, "node:", node)
	node.DataQueue <- data
}

// sendGroupMsg 向群聊发送消息
func sendGroupMsg(fromId uint, targetId uint, data []byte) (int, error) {
	//群发的逻辑：1获取到群里所有用户，然后向除开自己的每一位用户发送消息
	userIDs, err := FindUsers(fromId)
	if err != nil {
		zap.S().Info("查詢用戶失敗", err)
		return -1, err
	}

	for _,userId := range userIDs{
		if fromId != userId{
			sendMsgAndSave(int64(userId), data)
		}
	}
}

func sendMsgAndSave(userId int64, msg []byte){

}

func Chat(w http.ResponseWriter, r *http.Request) {
	//獲取發送者id
	query := r.URL.Query()
	userId := query.Get("userId")
	sendId, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		zap.S().Info("類型轉換失敗", err)
		return
	}

	//http升级为websocket
	wsUpgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		zap.S().Info("http升级为websocket失敗", err)
		return
	}

	//获取socket连接,构造消息节点
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	//將sendId和node綁定
	rwLocker.Lock()
	clientMap[sendId] = node
	rwLocker.Unlock()

	//服務端發送消息
	go sendProc(node)

	//服務端接收消息
	go recieveProc(node)
}

// sendProc 从node中获取信息并写入websocket中
func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				zap.S().Info("服務端發送消息失敗", err)
				return
			}
			fmt.Println("服務端發送socket消息成功")
		}
	}
}

// 從websocket中将消息体拿出，然后进行解析，再进行信息类型判断， 最后将消息发送至目的用户的node中
func recieveProc(node *Node) {
	for {
		//從websocket中读取数据
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			zap.S().Info("服務端讀取消息失敗", err)
			return
		}

		//將消息放入全局channel中
		broadMsg(data)

		//移到dispatch()中
		// msg := Message{}
		// err = json.Unmarshal(data, &msg)
		// if err != nil {
		// 	zap.S().Info("服務端解析json消息失敗:", err)
		// 	return
		// }
		// //fmt.Println(msg)
		// if msg.Type == 1 {
		// 	zap.S().Info("私訊：", msg.Content)
		// 	targetNode, ok := clientMap[msg.TargetId]
		// 	if !ok {
		// 		zap.S().Info("不存在對應的node")
		// 		return
		// 	}
		// 	targetNode.DataQueue <- data
		// 	fmt.Println("發送目的用户的node成功 :", string(data))
		// }
	}
}

func broadMsg(data []byte) {
	upSendChan <- data
}
