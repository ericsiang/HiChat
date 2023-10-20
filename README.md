## IM 項目 (使用websocket)

#### 參考：[从0到1搭建一个IM项目](https://learnku.com/articles/74274)

### 主要功能
* 登入、註冊、用戶資料更新、帳號註銷
* 單聊、群聊
* 加好友、好友列表、建群組、加入群組、群組列表

### 技术 Tool
* Go、Gin、Websocket、UDP、Mysql、Redis、Viper(config setting)、Gorm(ORM)、Zap(log controller)、lumberjack(cutting log file )、Md5、Jwt

    * zap、lumberjack 參考 : [Go日志库zap使用详解 ](https://www.cnblogs.com/jiujuan/p/17304844.html) 、[高性能日志库zap配置示例](https://studygolang.com/articles/17394)、[Go 項目實現日誌](https://www.readfog.com/a/1709305422763102208)

### 系统架构
![](flowData/system.png)

### 通信流程
![](flowData/system2.png)

![](flowData/flow.png)

![](flowData/websocket_connect_flow.png)