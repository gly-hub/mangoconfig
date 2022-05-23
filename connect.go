/**
* @Author vangogh
* @Description 连接管理
* @File:  connect
* @Datetime 2022/5/17 15:16
**/
package mangoconfig

import (
	"net"
	"sync"
)

type SetupFunc struct {
	FuncName string
	FuncLabel string
	FuncDesc string
}

type ClientConf struct {
	ConfName string
	ConfData string
}

var cMng *ConnectionMng

func init(){
	cMng = &ConnectionMng{
		Clients: make(map[string][]*Connection),
		mu:      sync.RWMutex{},
	}
}


type Connection struct {
	Conn net.Conn // 服务器连接
	IP string // 服务器ip
	Name string // 服务器名称
	Config []ClientConf // 服务配置列表
	Setups []SetupFunc // 初始化方法列表
}

type ConnectionMng struct {
	Clients map[string][]*Connection
	mu        sync.RWMutex
}

// 添加
func (c *ConnectionMng) AddClient(ip, name string, conn net.Conn){
	client := &Connection{
		Conn: conn,
		IP:   ip,
		Name: name,
		Config: make([]ClientConf, 0),
		Setups: make([]SetupFunc, 0),
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Clients[name] = append(c.Clients[name], client)
}

// 取出一个分组客户端
func (c *ConnectionMng) GetClients(name string)[]*Connection{
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Clients[name]
}

// 取出一个客户端
func (c *ConnectionMng) GetClient(name, ip string) *Connection {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for i:=0; i<len(c.Clients[name]);i++{
		if c.Clients[name][i].IP == ip {
			return c.Clients[name][i]
		}
	}
	return nil
}

// 移除
func (c *ConnectionMng) RemoveClient(name, ip string){
	c.mu.RLock()
	defer c.mu.RUnlock()
	for i:=0; i<len(c.Clients[name]);i++{
		if c.Clients[name][i].IP == ip {
			c.Clients[name] = append(c.Clients[name][:i], c.Clients[name][i+1:]...)
		}
	}
}

// 添加初始化方法
func (c *ConnectionMng) AddSetupFunc(ip, name, funcName, funcLabel, funcDesc string){
	funcInfo := SetupFunc{
		FuncName:  funcName,
		FuncLabel: funcLabel,
		FuncDesc:  funcDesc,
	}

	conn := c.GetClient(name, ip)
	if conn != nil {
		conn.Setups = append(conn.Setups, funcInfo)
	}
}

// 添加配置方法
func (c *ConnectionMng) AddConfigSetup(ip, name, confName, confData string){
	confObj := ClientConf{
		ConfName: confName,
		ConfData: confData,
	}

	conn := c.GetClient(name, ip)
	if conn != nil {
		conn.Config = append(conn.Config, confObj)
	}
}


// 获取现有的所有连接
func GetAllConnection()map[string][]*Connection{
	return cMng.Clients
}