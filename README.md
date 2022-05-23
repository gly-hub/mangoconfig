### MangoConfig配置管理包
mangoconfig配置服务用于服务器配置修改，在服务器启动或部署中可能存在配置信息更新情况,
某些情况需要重启服务从而重新初始化某些服务，例如mysql初始化、redis初始化等。mangoconfig包
提供了修改内存配置及重启初始化方法等功能。能够满足小部分需求，相较于现有的配置中心，功能更
加间接。
### 使用技术
1. golang
2. rpc
3. etcd
### 使用方法
#### server端
```golang
config := mangoconfig.RpcConfig{
    Addr:     "192.168.199.135:9456",
    ETCDAddr: []string{"192.168.199.243:2379"},
}
global.ConfigServer = mangoconfig.NewRPCxServer(config)
global.ConfigServer.Start()
```
**服务端发送更新配置消息**
```golang
mm := mangoconfig.ConfMsg{
    Name: reqJ.ConfName,
    Data: reqJ.Config,
}
mmStr, err := mm.Encode()
if err != nil{
    return
}

m := mangoconfig.MessageMng{
    Name:    reqJ.ClientName,
    IP:      reqJ.IP,
    Message: mangoconfig.Message{
        Type: mangoconfig.UPDATE_CONFIG_MOMERY,
        Data: mmStr,
    },
}

global.ConfigServer.SendMessage(m)
```
**服务端发送重新初始化消息**
```golang
m := mangoconfig.MessageMng{
		Name:    reqJ.ClientName,
		IP:      reqJ.IP,
		Message: mangoconfig.Message{
			Type: mangoconfig.RESETUP_CONFIG,
			Data: reqJ.Label,
		},
	}

	global.ConfigServer.SendMessage(m)
```
#### client端
```golang
client := mangoconfig.NewClient("MangoApi", []string{"192.168.199.243:2379"})
client.Start()
_ = client.RegisterConfStruct("config", config.Conf)
_ = client.RegisterSetupFunc("初始化rpc", "rpcSetup", "初始化rpc方法", clientrpcx.Setup)
_ = client.RegisterSetupFunc("初始化http", "httpSetup", "初始化http方法", server.Setup)
```
### 方法
#### server提供
+ SendMessage():用于服务端发送通知消息
#### client提供
+ RegisterConfStruct():用于注册配置对象
+ RegisterSetupFunc():用于注册初始化方法
### 心跳问题
心跳由服务端主动推送，若客户端没响应则会删掉客户端。
### 关于服务
