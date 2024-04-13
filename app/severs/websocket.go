package severs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"remembrance/app/common"
	"remembrance/app/controller"
	"remembrance/app/models"
	"remembrance/app/response"
	"sync"
	"time"

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
	clients            = make(map[string][]*websocket.Conn) // 存储所有客户端连接，按照 groupid 分组
	clientsLock        sync.Mutex                           // 用于 clients 的互斥锁，防止并发冲突
	messageHistoryLock sync.Mutex                           // 用于 messageHistory 的互斥锁
)

// HandleConnections 处理 WebSocket 连接的函数
func HandleConnections(c *gin.Context) {
	// 从 URL 参数中获取 groupid
	groupID := c.Query("groupid")

	// 升级 HTTP 连接到 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}
	defer conn.Close() // 确保在函数返回前关闭连接

	// 将新的连接添加到相应的组中
	clientsLock.Lock()
	clients[groupID] = append(clients[groupID], conn)
	clientsLock.Unlock()

	// 读取历史消息，发送给新连接
	sendMessageHistory(groupID, conn)

	// 循环读取消息
	for {
		var mess controller.Message
		err := conn.ReadJSON(&mess)
		if err != nil {
			// 读取出错，断开连接并从群组中移除
			clientsLock.Lock()
			removeConnection(groupID, conn)
			clientsLock.Unlock()
			break
		}

		// 处理消息
		if mess.Cloudurl != "" {
			photo := mess.GetGroupPhoto()
			messageHistoryLock.Lock()
			common.DB.Create(&photo)
			messageHistoryLock.Unlock()

			// 广播消息给同一群组内的所有连接
			broadcastToGroup(groupID, photo)
			response.OkMsg(c, "发送成功")
		}

		time.Sleep(100 * time.Microsecond)
	}
}

// sendMessageHistory 发送历史消息给连接
func sendMessageHistory(groupID string, conn *websocket.Conn) {
	messageHistoryLock.Lock()
	defer messageHistoryLock.Unlock()

	var messageHistory []models.GroupPhoto
	common.DB.Limit(20).Table("group_photos").Where("Group_id = ?", groupID).Find(&messageHistory)

	for _, msg := range messageHistory {
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("Error marshaling message:", err)
			break
		}

		if err := conn.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}
}

// broadcastToGroup 广播消息给同一群组内的所有连接
func broadcastToGroup(groupID string, photo models.GroupPhoto) {
	clientsLock.Lock()
	defer clientsLock.Unlock()

	jsonMessage, err := json.Marshal(photo)
	if err != nil {
		fmt.Println("Error marshaling GroupPhoto:", err)
		return
	}

	for _, client := range clients[groupID] {
		if err := client.WriteMessage(websocket.TextMessage, jsonMessage); err != nil {
			client.Close() // 发送消息失败时关闭连接
			removeConnection(groupID, client)
		}
	}
}

// removeConnection 从群组中移除连接
func removeConnection(groupID string, conn *websocket.Conn) {
	// 查找并移除连接
	for i, c := range clients[groupID] {
		if c == conn {
			clients[groupID] = append(clients[groupID][:i], clients[groupID][i+1:]...)
			break
		}
	}
}
