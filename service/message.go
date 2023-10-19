package service

import (
	"HiChat/common"
	"HiChat/dao"
	"HiChat/global"
	"HiChat/models"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gopkg.in/fatih/set.v0"
)

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
		n, err := udpConn.Read(buf[:])
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
	msg := models.Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		zap.S().Info("服務端解析json消息失敗:", err)
		return
	}

	fmt.Println("解析数据:", msg, "msg.FormId", msg.FromId, "targetId:", msg.TargetId, "type:", msg.Type)

	//判断消息类型
	switch msg.Type {
	case 1: //私聊
		sendMsgAndSave(msg.TargetId, data)
	case 2: //群聊
		sendGroupMsg(uint(msg.FromId), uint(msg.TargetId), data)
	}
}

// sendMsgAndSave 向用户单聊发送消息
func sendMsgAndSave(targetId int64, msg []byte) {
	rwLocker.RLock()                //保证线程安全，上锁
	node, ok := clientMap[targetId] //对方是否在线
	rwLocker.RUnlock()              //解锁

	jsonMsg := models.Message{}
	json.Unmarshal(msg, &jsonMsg)
	ctx := context.Background()
	targetIdStr := strconv.Itoa(int(targetId))
	userIdStr := strconv.Itoa(int(jsonMsg.FromId))

	if !ok {
		zap.S().Info("userId不存在對應的node")
		return
	}

	//如果当前用户在线，将消息转发到当前用户的websocket连接中，然后进行存储
	zap.S().Info("targetId :", targetId, "node:", node)
	node.DataQueue <- msg

	//userIdStr和targetIdStr进行拼接唯一key
	var key string
	if targetId > jsonMsg.FromId {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}

	//将消息存储到redis中
	res, err := global.RedisDB.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		zap.S().Info("創建redis紀錄失敗", err)
		return
	}

	//將聊天記錄寫入redis中
	score := float64(cap(res)) + 1
	ress, err := global.RedisDB.ZAdd(ctx, key, redis.Z{score, msg}).Result()
	if err != nil {
		zap.S().Info("寫入redis紀錄失敗", err)
		return
	}
	fmt.Println("寫入redis紀錄成功", ress)
}

// sendGroupMsg 向群聊发送消息
func sendGroupMsg(fromId uint, targetId uint, data []byte) (int, error) {
	//群发的逻辑：1获取到群里所有用户，然后向除开自己的每一位用户发送消息
	userIDs, err := dao.FindUsers(fromId)
	if err != nil {
		zap.S().Info("查詢用戶失敗", err)
		return -1, err
	}

	for _, userId := range userIDs {
		if fromId != userId {
			sendMsgAndSave(int64(userId), data)
		}
	}

	return 0 ,nil
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

	}
}

func broadMsg(data []byte) {
	upSendChan <- data
}

// 從redis中獲取聊天記錄
func RedisMsgService(userId int64, targetId int64, start int64, end int64, isRev bool) []string {
	ctx := context.Background()
	userIdStr := strconv.Itoa(int(userId))
	targetIdStr := strconv.Itoa(int(targetId))

	var key string
	if userId > targetId {
		key = "msg_" + targetIdStr + "_" + userIdStr
	} else {
		key = "msg_" + userIdStr + "_" + targetIdStr
	}

	var res []string
	var err error

	if isRev {
		res, err = global.RedisDB.ZRange(ctx, key, start, end).Result()
	} else {
		res, err = global.RedisDB.ZRevRange(ctx, key, start, end).Result()
	}

	if err != nil {
		zap.S().Info("獲取redis紀錄失敗", err)
		return nil
	}

	return res
}




func RedisMsg(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
	start, _ := strconv.Atoi(c.PostForm("start"))
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev"))

	res := RedisMsgService(int64(userIdA), int64(userIdB), int64(start), int64(end), isRev)
	common.RespOKList(c.Writer, "ok", res)
}