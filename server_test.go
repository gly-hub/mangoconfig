/**
* @Author vangogh
* @Description //TODO
* @File:  server_test.go
* @Datetime 2022/5/18 11:03
**/
package mangoconfig

import (
	"fmt"
	"testing"
	"time"
)

func TestNewRPCxServer(t *testing.T) {
	config := RpcConfig{Addr:"192.168.199.135:9456"}
	rpcServer := NewRPCxServer(config)
	rpcServer.Start()

	time.Sleep(6 * time.Second)
	m := MessageMng{
		Name:    "Test",
		IP:      "",
		Message: Message{
			Type: 2,
			Data: "123",
		},
	}
	fmt.Println("发送消息")
	rpcServer.SendMessage(m)

	for {
		select {

		}
	}
}