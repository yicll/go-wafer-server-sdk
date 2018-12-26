Golang Wafer Server SDK
===============

本项目是 [go-wafer-session-server](https://github.com/yicll/go-wafer-session-server) 组成部分

可搭配前端的wafer-sdk使用，也可以单独引入项目中使用

### go版本wafer sdk

**使用方法**

由于公司的gitlab仓库没有支持https，所以go get的时候需要加上-insecure参数

> go get github.com/yicll/go-wafer-server-sdk

**引入方法**

```Go
import (
	"github.com/yicll/go-wafer-server-sdk/wafer"
)
```

**使用说明**

初始化sdk

```Go
sdk := wafer.NewWaferSDK(
	"wxf4daafbf1c76d304",					 //小程序的appid，不能为空
	"http://domob-206.domob-inc.cn:6789",    //wafer-session-server地址
	true,									 // 是否配合前端wafer-sdk使用，true：是、false：否
	this.Ctx.Request,						 //当使用wafer-sdk为true的时候，该参数不能为空，*http.Request
)
```

登录

```Go
//配合wafer-sdk的使用场景
data, err := sdk.Login()

//单独使用的场景
//code: 小程序调用login接口获取到的临时登录凭据
//encryptData: 加密的用户信息
//iv: 加密向量偏移量
data, err := sdk.Login(code, encryptData, iv)
```

检查登录状态

```Go
//配合wafer-sdk的使用场景
data, err := sdk.Check()

//单独使用的场景
//id: 通过调用sdk的login接口返回结果中的Id字段
//skey: 通过调用sdk的login接口返回的结果中的Skey字段
data, err := sdk.Check(id, skey)
```

Login&Check方法返回结构体`ReturnResult `详细说明

```Go
// SDK统一返回数据结构
// PrintWafer: 是否需要打印出wafer sdk需要的信息，true的话，请直接将PrintResut转成json返回前端，否则忽略，只有使用wafer sdk需要关注
// PrintResult: 需要打印输出到前端的wafer sdk可读信息
// Data: 返回数据
type ReturnResult struct {
	PrintWafer  bool
	PrintResult interface{}
	Data        ResponseData
}

// wafer-session-server接口返回数据结构体
// Id: session uuid，结合skey用于识别session
// Skey: session skey，结合id组合用于识别session
// UserInfo: 返回的用户信息结构
// Duration: 登录有效期，业务端可忽略
type ResponseData struct {
	Id       string   `json:"id"`
	Skey     string   `json:"skey"`
	UserInfo UserInfo `json:"user_info"`
	Duration int      `json:"duration"`
}

// 解析成功之后的用户信息结构
type UserInfo struct {
	OpenId    string `json:"openId"`
	UnionId   string `json:"unionId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	Language  string `json:"language"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarUrl string `json:"avatarUrl"`
}
```

**使用示例**

以beego框架为例，实现的api

```Go
package controllers

import (
	"git.domob-inc.cn/mp-lib/go-wafer-server-sdk/wafer"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

// 独立使用的示例
// 登录接口
func (this *MainController) Login() {

	sdk := wafer.NewWaferSDK(
		"wxf4daafbf1c76d304",
		"http://domob-206.domob-inc.cn:6789",
		false,
	)
	
	code := "001btntB1UQJ7d0VVasB1ZKGtB1btnth"
	encryptData := "qs8afGiRlAsjIcNuG9CqxMbMgr6tpaTqOrpa9szUSrYfObQR54ThGhmAadEhkuW/6Flyqa+r+p/4BuKnCLx81TzwqM+7gP3pdOG4rLvlvWCtDes2blsGZm2wNFOqqwj+xfVQqj25JznX75lNbObY5Ic67ZTiaszMzJym0QDy7vaBQMCwdGLfTiVPc35cpfq9ZZzGDVVewHoNGauhPrkOxdu+ec/M6/Fp39J32yEyfi/7lkUwauobdDl7ovazjoFGvfeBOjdXlmBGuF0+W5KKjdsXINLHWL1m4gZD5twLQxICC4A6W6YvXoLAHr41eslvfFvGptIJFOW4GXnEZyhzc7tgubiSvMy9cMA0NcB6o8qIh7GrZ1sp6FSdrCDaDj3zXNlHzgbXvNfX/Q7PkQ18AaofjapSnoEOUxfiHwR/yNpK05yqviCgY7UdoNUSKd3GtMXg+KJTG5yfvOfN23JiaQnJ4P30wJ15IJb07pQsEMk0C6QthDfPvRnxpU07ERgGL7FKmP1f3Z2HlrzET/Z2Jw=="
	iv := "nU6TJmoVfrz8Vt8FJbrZYA=="
	
	data, err := sdk.Login(code, encryptData, iv)
	if err != nil {
		beego.Error(err.Error())
	}

	this.Data["json"] = data.Data

	this.ServeJSON()
}

// 检查登录状态
func (this *MainController) Check() {

	sdk := wafer.NewWaferSDK(
		"wxf4daafbf1c76d304",
		"http://domob-206.domob-inc.cn:6789",
		false,
	)
	
	id := "54cddc016436921eb6da31962d88e946"
	skey := "1e9d61aa73f67407f374326cde532b13"

	data, err := sdk.Check(id, skey)
	if err != nil {
		beego.Error(err.Error())
	}
	
	this.Data["json"] = data.Data

	this.ServeJSON()
}

// 配合前端wafer sdk使用
// 登录接口
func (this *MainController) Login1() {

	sdk := wafer.NewWaferSDK(
		"wxf4daafbf1c76d304",
		"http://domob-206.domob-inc.cn:6789",
		true,
		this.Ctx.Request,
	)
	data, err := sdk.Login()
	if err != nil {
		beego.Error(err.Error())
	}

	if data.PrintWafer {
		this.Data["json"] = data.PrintResult
	} else {
		this.Data["json"] = data.Data
	}

	this.ServeJSON()
}

// 检查登录状态
func (this *MainController) Check1() {

	sdk := wafer.NewWaferSDK(
		"wxf4daafbf1c76d304",
		"http://domob-206.domob-inc.cn:6789",
		true,
		this.Ctx.Request,
	)

	data, err := sdk.Login()
	if err != nil {
		beego.Error(err.Error())
	}

	if data.PrintWafer {
		this.Data["json"] = data.PrintResult
	} else {
		this.Data["json"] = data.Data
	}

	this.ServeJSON()
}

```


