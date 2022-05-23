/**
* @Author vangogh
* @Description 数据结构体
* @File:  model
* @Datetime 2022/5/18 10:18
**/
package mangoconfig

// 请求结构体
type (
	RegisterReq struct {
		Name string
	}

	RegisterConfigStruct struct {
		ClientName string
		ConfName string
		ConfData string
	}

	RegisterSetupFunc struct {
		ClientName string
		FuncName string
		FuncLabel string
		FuncDesc string
	}
)


// 响应结构体
type (
	RegisterResp struct {
		Code int
	}
)