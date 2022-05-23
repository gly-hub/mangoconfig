/**
* @Author vangogh
* @Description 配置客户端
* @File:  client
* @Datetime 2022/5/17 14:48
**/
package mangoconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	etcd_client "github.com/rpcxio/rpcx-etcd/client"
)

type ConfigClient struct {
	client *client.OneClient
	message chan *protocol.Message
	name string
	setups map[string]func()
	config map[string]interface{}
}

// 初始化一个客户端
func NewClient(name string, etcd []string) *ConfigClient {
	var configClient = &ConfigClient{
		message:make(chan *protocol.Message),
		name: name,
		setups: make(map[string]func()),
		config: make(map[string]interface{}),
	}


	d, err := etcd_client.NewEtcdV3Discovery("/config", "", etcd, true, nil)
	if err != nil{
		panic(err)
	}

	configClient.client = client.NewBidirectionalOneClient(client.Failtry, client.RandomSelect, d, client.DefaultOption, configClient.message)

	// 注册客户端
	req := RegisterReq{
		Name:   name,
	}
	resp := &RegisterResp{}
	err = configClient.client.Call(context.Background(), CONFIG_SERVICE, "RegisterClient", req, resp)
	if err != nil {
		panic(fmt.Sprintf("failed to call: %v", err))
	}

	return configClient
}

// 启动后开始监听
func (c *ConfigClient) Start(){
	go func() {
		for msg := range c.message{
			var message Message
			err := message.Decode(string(msg.Payload))
			if err != nil{
				continue
			}
			switch message.Type{
			case UPDATE_CONFIG_MOMERY:
				// 普通内存更新
				c.procMemoryConfigMessage(message.Data)

			case RESETUP_CONFIG:
				c.proSetup(message.Data)
			}
		}
	}()
}

/**
 * @Description: 注册配置结构
 * @receiver c
 * @param name 配置名
 * @param conf 名字对象
 */
func (c *ConfigClient) RegisterConfStruct(name string, conf interface{}) error {
	confByte, err := json.Marshal(conf)
	if err != nil{
		fmt.Println(err)
		return err
	}
	req := RegisterConfigStruct{
		ClientName: c.name,
		ConfName:   name,
		ConfData:   string(confByte),
	}

	resp := &RegisterResp{}
	err = c.client.Call(context.Background(), CONFIG_SERVICE,"RegisterConfigStruct", req, resp)
	if err != nil {
		return err
	}

	c.config[name] = conf
	return nil
}

/**
 * @Description: 注册初始化方法
 * @receiver c
 * @param name 名称
 * @param label 标签，用于标记
 * @param desc 描述
 * @param setup 无参初始化方法
 * @return error
 */
func (c *ConfigClient) RegisterSetupFunc(name, label, desc string, setup func()) error {
	req := RegisterSetupFunc{
		ClientName: c.name,
		FuncName:   name,
		FuncLabel:  label,
		FuncDesc:   desc,
	}

	resp := &RegisterResp{}
	err := c.client.Call(context.Background(), CONFIG_SERVICE,"RegisterSetupFunc", req, resp)
	if err != nil {
		return err
	}
	c.setups[label] = setup
	return nil
}

// 普通内存配置消息处理函数
func (c *ConfigClient) procMemoryConfigMessage(data string){
	var confMsg  = &ConfMsg{}
	if err := json.Unmarshal([]byte(data), confMsg); err != nil{
		fmt.Println(err)
		return
	}
	if err := confMsg.mapConf(c.config[confMsg.Name]); err != nil{
		fmt.Println(err)
		return
	}
}

// 重新初始化处理函数
func (c *ConfigClient) proSetup(data string){
	c.setups[data]()
}