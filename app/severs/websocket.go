package severs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"remembrance/app/common"
	"remembrance/app/controller"
	"remembrance/app/models"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// websocket.Upgrader 结构体，用于配置 WebSocket 连接的参数
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // 读缓冲区大小
	WriteBufferSize: 1024, // 写缓冲区大小
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源的 WebSocket 连接
	},
}

// 定义全局变量以存储当前所有已连接的客户端
var (
	clients            = make(map[*websocket.Conn]bool) // 存储所有客户端连接
	clientsLock        sync.Mutex                       // 用于 clients 的互斥锁，防止并发冲突
	messageHistoryLock sync.Mutex                       // 用于 messageHistory 的互斥锁
)

// handleConnections 处理 WebSocket 连接的函数
func HandleConnections(c *gin.Context) {
	// 升级 HTTP 连接到 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}
	defer conn.Close() // 确保在函数返回前关闭连接

	// 将新的连接添加到 clients 集合中
	clientsLock.Lock()
	clients[conn] = true
	clientsLock.Unlock()

	//读取历史消息
	var mes controller.Message
	c.BindJSON(&mes)
	var messageHistory []models.GroupPhoto
	common.DB.Limit(7).Table("groupphotos").Where("Group_id = ?", mes.GroupId).Find(&messageHistory)
	// 发送历史消息给新连接的客户端
	messageHistoryLock.Lock()
	for _, msg := range messageHistory {
		// 序列化每条消息为JSON
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("Error marshaling message:", err)
			break
		}
		// 发送序列化后的消息
		if err := conn.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}

	messageHistoryLock.Unlock()

	// 循环读取新的消息
	for {
		var mess controller.Message
		c.BindJSON(&mess)
		if err != nil {
			break // 如果读取出错，退出循环
		}
		photo := mess.GetGroupPhoto()
		// 将新消息添加到历史记录并广播给所有客户端
		messageHistoryLock.Lock()
		common.DB.Create(&photo)
		messageHistoryLock.Unlock()

		broadcast(photo) // 广播消息
	}

	// 从 clients 集合中移除断开连接的客户端
	clientsLock.Lock()
	delete(clients, conn)
	clientsLock.Unlock()
}

// broadcast 将消息广播给所有已连接的客户端
func broadcast(photo models.GroupPhoto) {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	//转化为json
	jsonMessage, err := json.Marshal(photo)
	if err != nil {
		fmt.Println("Error marshaling GroupPhoto:", err)
		return
	}

	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
			client.Close()          // 发送消息失败时关闭连接
			delete(clients, client) // 从客户端列表中移除
		}
	}
}
