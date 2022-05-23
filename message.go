/**
* @Author vangogh
* @Description 消息
* @File:  message
* @Datetime 2022/5/18 10:32
**/
package mangoconfig

import (
	"encoding/json"
)

type MessageType int

const (
	HEARTBEAT MessageType = 0 // 心跳
	UPDATE_CONFIG_MOMERY MessageType = 1 // 只更新内存配置
	UPDATE_CONFIG_FILE MessageType = 2 // 更新文件配置
	RESETUP_CONFIG MessageType = 3 // 初始化普通函数
	RESTART_CONFIG MessageType = 4 // 重启服务
)


type MessageMng struct {
	Name string
	IP string
	Message Message
}

// 消息分为消息加解码
type Message struct {
	Type MessageType // 消息类型
	Data string // 数据
}

// 加码
func (m *Message) Encode() (string, error) {
	mByte, err := json.Marshal(m)
	if err != nil{
		return "", err
	}
	return string(mByte), nil
}

// 解码
func (m *Message) Decode(mStr string) error {
	return json.Unmarshal([]byte(mStr), m)
}


type ConfMsg struct {
	Name string
	Data string
}

// 加码
func (m *ConfMsg) Encode() (string, error) {
	mByte, err := json.Marshal(m)
	if err != nil{
		return "", err
	}
	return string(mByte), nil
}

// 解码
func (m *ConfMsg) Decode(mStr string) error {
	return json.Unmarshal([]byte(mStr), m)
}

// 配置注入
func (m *ConfMsg) mapConf(conf interface{}) error {
	return json.Unmarshal([]byte(m.Data), conf)
}