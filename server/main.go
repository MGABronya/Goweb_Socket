// @Title  server
// @Description  服务端
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-10-08 19:26
package main

import (
	"fmt"
	"net"
	"time"
)

// 连接池
var ConnSlice map[net.Conn]*Heartbeat

// Heartbeat			定义了心跳消息的结构体
type Heartbeat struct {
	endTime int64 // 过期时间
}

// @title    main
// @description   用于启动服务
// @auth      MGAronya（张健）             2022-10-08 19:26
// @param     void
// @return    void
func main() {
	// TODO 初始化连接池
	ConnSlice = map[net.Conn]*Heartbeat{}

	// TODO 监听本地端口8806
	l, err := net.Listen("tcp", "127.0.0.1:8086")

	// TODO 启动失败
	if err != nil {
		fmt.Println("服务启动失败")
	}
	defer l.Close()

	// TODO 接收信息
	for {
		// TODO 获取客户端请求连接
		conn, err := l.Accept()

		// TODO 接收失败
		if err != nil {
			fmt.Println("Error accepting: ", err)
		}

		// TODO 打印连接的客户端地址和本地地址
		fmt.Printf("Received message %s -> %s\n", conn.RemoteAddr(), conn.LocalAddr())

		// TODO 更新心跳时间
		ConnSlice[conn] = &Heartbeat{
			endTime: time.Now().Add(time.Second * 5).Unix(),
		}

		// TODO 处理连接
		go handelConn(conn)
	}
}

// @title    handelConn
// @description   用于处理连接
// @auth      MGAronya（张健）             2022-10-08 19:26
// @param     c net.Conn			接收一个连接
// @return    void
func handelConn(c net.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := c.Read(buffer)
		// TODO 更新心跳时间
		if ConnSlice[c].endTime > time.Now().Unix() {
			ConnSlice[c].endTime = time.Now().Add(time.Second * 5).Unix()
		} else {
			fmt.Println("长时间未发送消息断开连接")
			return
		}

		// TODO 读取消息失败
		if err != nil {
			return
		}

		// TODO 如果是心跳检测，则不执行剩下的代码
		if string(buffer[0:n]) == "1" {
			// TODO 向客户端发送心跳回复，保证服务存活
			c.Write([]byte("1"))
			continue
		}

		// TODO 心跳检测，在需要发送数据时才检查规定时间内有没有数据到达
		for conn, heart := range ConnSlice {
			if conn == c {
				continue
			}
			// TODO 如果过期，从列表中删除连接，并关闭连接
			if heart.endTime < time.Now().Unix() {
				delete(ConnSlice, conn)
				conn.Close()
				fmt.Println("删除连接", conn.RemoteAddr())
				fmt.Println("现在存有连接", ConnSlice)
				continue
			}
			// TODO 写入该链接发送过来的信息
			conn.Write(buffer[0:n])
		}
	}
}
