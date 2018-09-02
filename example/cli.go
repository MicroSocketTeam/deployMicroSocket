package main

import (
	"fmt"
	"microSocket"
	"microSocket/util"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8565")
	if err != nil {
		return
	}
	defer conn.Close()

	data := make(map[string]string)
	data["module"] = "test"
	data["action"] = "Hello"
	data["name"] = "jd"
	//把map转化为string
	a := []byte(util.Map2String(data))
	//把string打包
	b := microSocket.Enpack(a)
	//发送数据
	conn.Write(b)

	res := make([]byte, 20)
	conn.Read(res)
	fmt.Println(res)
}
