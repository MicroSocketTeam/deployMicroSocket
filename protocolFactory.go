package microSocket

import (
	"log"
	"microSocket/util"
)

type AbstractFactory interface {
	NewMsf(msfEvent MsfEventer) *Msf
}

// 协议工厂
func Factory(protocolName string) AbstractFactory {
	util.DealBaseCli()
	var abstractFactory AbstractFactory
	switch protocolName {
		case "tcp":
			abstractFactory = TcpFactory{}
		case "websocket":

		default:
			log.Println("错误的协议名称(Wrong protocol name!)")
	}
	return abstractFactory
}