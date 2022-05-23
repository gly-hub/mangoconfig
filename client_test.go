/**
* @Author vangogh
* @Description 客户端TEST
* @File:  client_test.go
* @Datetime 2022/5/18 14:17
**/
package mangoconfig

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("192.168.199.135:9456", []string{""})
	client.Start()
	for{
		select {

		}
	}
}