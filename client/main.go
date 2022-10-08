// @Title  client
// @Description  用户端
// @Author  MGAronya（张健）
// @Update  MGAronya（张健）  2022-10-08 19:26
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

// @title    main
// @description   用于启动客户端
// @auth      MGAronya（张健）             2022-10-08 19:26
// @param     void
// @return    void
func main() {
	// TODO 使用本地端口8086
	server := "127.0.0.1:8086"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)

	// TODO 使用失败
	if err != nil {
		fmt.Print("Fatal error:", err.Error())
		os.Exit(1)
	}

	// TODO 生成一个连接
	conn, err := net.DialTCP("tcp", nil, tcpAddr)

	// TODO 生成失败
	if err != nil {
		fmt.Print("Fatal error:", err.Error())
		return
	}

	// TODO 打印连接成功信息
	fmt.Print(conn.RemoteAddr().String(), "connect success!")

	// TODO 创建定时器，每次服务器端发送消息就刷新时间
	Sender(conn)

	// TODO 终止
	fmt.Print("end")
}

// @title    Sender
// @description   创建定时器，每次服务器端发送消息就刷新时间
// @auth      MGAronya（张健）             2022-10-08 19:26
// @param     void
// @return    void
func Sender(conn *net.TCPConn) {

	// TODO 保证连接最后关闭
	defer conn.Close()

	// TODO 读入
	sc := bufio.NewReader(os.Stdin)

	// TODO 生成定时器
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for {
			// TODO 每秒写入一个1，作为心跳包发送给服务端
			<-t.C
			_, err := conn.Write([]byte("1"))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}()

	// TODO 接收用户的聊天名称
	name := ""
	fmt.Println("请输入聊天名称")

	// TODO 写入名称
	fmt.Fscan(sc, &name)

	// TODO 消息
	msg := ""
	buffer := make([]byte, 1024)

	// TODO 建立5s的定时器
	_t := time.NewTimer(time.Second * 5)
	defer _t.Stop()

	go func() {
		// TODO 5s服务端没有反应，则报错
		<-_t.C
		fmt.Println("服务器出现故障，断开连接")
	}()

	for {
		go func() {
			for {
				// TODO 读入服务端消息
				n, err := conn.Read(buffer)
				if err != nil {
					return
				}
				// TODO 收到消息就刷新_t定时器，如果time.Second * 5时间到了，则<-_t.C就不会阻塞
				_t.Reset(time.Second * 5)

				// TODO 打印除心跳回应外的消息
				if string(buffer[0:1]) != "1" {
					fmt.Println(string(buffer[0:n]))
				}
			}
		}()

		// TODO 写入消息
		fmt.Fscan(sc, &msg)

		// TODO 打印时间
		i := time.Now().Format("2022-10-08 19:48:05")
		conn.Write([]byte(fmt.Sprintf("%s\n\t%s:%s", i, name, msg)))
	}
}
