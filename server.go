/**
* @Author vangogh
* @Description 配置服务端
* @File:  server
* @Datetime 2022/5/17 14:48
**/
package mangoconfig

import (
	"context"
	"fmt"
	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/server"
	"net"
	"time"
)

type RPCxServer struct {
	messages chan MessageMng
	s *server.Server
	addr string
}

type RpcConfig struct {
	Addr string
	ETCDAddr []string
}

func NewRPCxServer(config RpcConfig) *RPCxServer{
	s := server.NewServer()
	rcvHandle := new(API)

	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + config.Addr,
		EtcdServers:    config.ETCDAddr,
		BasePath:       "/config",
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		panic(err)
	}
	s.Plugins.Add(r)

	// 注册服务名
	if err := s.RegisterName(CONFIG_SERVICE, rcvHandle, ""); err != nil{
		panic(err)
	}

	rpcS:=&RPCxServer{
		messages: make(chan MessageMng),
		s:    s,
		addr: config.Addr,
	}

	return rpcS
}

func (r *RPCxServer) SendMessage(m MessageMng){
	r.messages <- m
}

func (r *RPCxServer) Start(){
	// 启动服务
	r.server()
	// 启动消息处理
	r.messageProc()
	// 启动心跳检测
	r.heartbeat()
}

func (r *RPCxServer) server(){
	port := ParsePort(r.addr)
	// 启动服务
	go func() {
		err:=r.s.Serve("tcp", fmt.Sprintf("%s", port))
		if err != nil {
			panic("RPCxServer Error")
		}
	}()
}

// 消息处理
func (r *RPCxServer) messageProc(){
	go func() {
		for {
			select {
			case m := <- r.messages:
				conn := cMng.GetClient(m.Name, m.IP)
				if conn != nil && conn.Conn != nil {
					mStr, err := m.Message.Encode()
					if err != nil{
						fmt.Println(err)
						continue
					}
					err = r.s.SendMessage(conn.Conn, "", "", nil, []byte(mStr))
					if err != nil {
						cMng.RemoveClient(m.Name, m.IP)
					}
				}
			}
		}
	}()
}

// 客户端连接处理
func (r *RPCxServer) heartbeat(){
	go func() {
		for {
			for _, clients := range cMng.Clients{
				for i:=0;i<len(clients);i++{
					// 发送心跳
					m := MessageMng{
						Name:    clients[i].Name,
						IP:      clients[i].IP,
						Message: Message{
							Type: HEARTBEAT,
							Data: "",
						},
					}
					r.messages <- m

					//m :=  Message{
					//	Type: HEARTBEAT,
					//	Data: "",
					//}
					//
					//mStr, _ := m.Encode()
					//err := r.s.SendMessage(clients[i].Conn, "", "", nil, []byte(mStr))
					//if err != nil {
					//	cMng.RemoveClient(clients[i].Name, clients[i].IP)
					//}
				}
			}
			time.Sleep(time.Minute)
		}
	}()
}

func ParsePort(addr string)(port string) {
	for pos, c := range addr {
		switch c {
		case ':':
			port = addr[pos:]
		}
	}
	return
}


// 控制服务
type API struct {

}

// 注册客户端
func (a *API) RegisterClient(ctx context.Context, req RegisterReq, resp *RegisterResp) error {
	client := ctx.Value(server.RemoteConnContextKey).(net.Conn)
	cMng.AddClient(client.RemoteAddr().String(), req.Name, client)
	//cMng.AddClient("", req.Name, client)
	resp.Code = 200
	return nil
}

// 注册配置
func (a *API) RegisterConfigStruct(ctx context.Context, req RegisterConfigStruct, resp *RegisterResp) error {
	client := ctx.Value(server.RemoteConnContextKey).(net.Conn)
	cMng.AddConfigSetup(client.RemoteAddr().String(), req.ClientName, req.ConfName, req.ConfData)
	resp.Code = 200
	return nil
}

 // 注册初始化方法
func (a *API) RegisterSetupFunc(ctx context.Context, req RegisterSetupFunc, resp *RegisterResp) error {
	client := ctx.Value(server.RemoteConnContextKey).(net.Conn)
	fmt.Println(client.RemoteAddr().String())
	cMng.AddSetupFunc(client.RemoteAddr().String(), req.ClientName, req.FuncName, req.FuncLabel, req.FuncDesc)
	resp.Code = 200
	return nil
}
