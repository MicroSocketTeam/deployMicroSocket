package main

import (
	"fmt"
	"log"
	msf "microSocket"
	"strconv"
)

var wbsSer = msf.NewWebskt(&wbsEvent{})

//框架事件
type wbsEvent struct {
}

//客户端握手成功事件
func (this wbsEvent) OnHandel(fd string) bool {
	log.Println(fd, "链接成功类")
	return true
}

//断开连接事件
func (this wbsEvent) OnClose(fd string) {
	log.Println(fd, "链接断开类")
}

//接收到消息事件
func (this wbsEvent) OnMessage(fd string, msg string) bool {
	log.Println(fd, msg)
	return true
}

//---------------------------------------------------------------------
//框架业务逻辑
type Test struct {
}

func (this Test) Default() {
	fmt.Println("is default")
}

func (this Test) BeforeRequest(data map[string]string) bool {
	log.Println("before")
	return true
}

func (this Test) AfterRequest(data map[string]string) {
	log.Println("after")
}

func (this Test) Hello(data map[string]string) {
	fd, _ := strconv.Atoi(data["fd"])
	log.Println("收到消息了")
	wbsSer.SessionMaster.WriteByid(uint32(fd), "Hello")
}

//---------------------------------------------------------------------

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.Llongfile)

	wbsSer.EventPool.Register("test", &Test{})
	wbsSer.Listening("127.0.0.1:9501")
}
