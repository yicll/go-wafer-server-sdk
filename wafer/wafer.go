package wafer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RequestInterface struct {
	Appid         string            `json:"appid"`
	InterfaceName string            `json:"interfaceName"`
	Params        map[string]string `json:"para"`
}

type RequestBody struct {
	Version       int              `json:"version"`
	ComponentName string           `json:"componentName"`
	ReqInterface  RequestInterface `json:"interface"`
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

// wafer-session-server返回数据
type ResponseBody struct {
	Code    int          `json:"returnCode"`
	Message string       `json:"returnMessage"`
	Data    ResponseData `json:"returnData"`
}

// SDK统一返回数据结构
// PrintWafer: 是否需要打印出wafer sdk需要的信息，true的话，请直接将PrintResut转成json返回前端，否则忽略，只有使用wafer sdk需要关注
// PrintResult: 需要打印输出到前端的wafer sdk可读信息
// Data: 返回数据
type ReturnResult struct {
	PrintWafer  bool
	PrintResult interface{}
	Data        ResponseData
}

type WaferSDK struct {
	Request       *http.Request
	Appid         string
	AuthServerUrl string
	Wafer         bool
	code          string
	encryptData   string
	iv            string
	id            string
	skey          string
	err           WaferError
}

// r: *http.Request，如果需要配合前端wafer-sdk使用，请传入http request对象，否则该参数为nil即认为单独使用
// appid: 小程序appid
// url: 请求wafer-session-server的地址
// wafer: 前端是否使用wafer-sdk配合请求，true：是、false：否（表示单独使用）
// rs: 当wafer为true的时候，必须传递该参数，http request对象
func NewWaferSDK(appid string, url string, wafer bool, rs ...*http.Request) *WaferSDK {
	var r *http.Request
	if len(rs) == 0 {
		r = nil
	} else {
		r = rs[0]
	}

	return &WaferSDK{
		Request:       r,
		Appid:         appid,
		AuthServerUrl: url,
		Wafer:         wafer,
	}
}

// 通过wafer-session-server进行登录
// 如果小程序前端使用了wafer sdk，则不需要传递参数
// 如果单独使用，则需要传递相应参数
//  param1: code string，微信登录临时凭证
//  param2: encrypt_data string 微信加密后的用户信息
//  param3: iv string 加密偏移向量
func (w *WaferSDK) Login(p ...string) (r ReturnResult, err error) {
	err = w.validate("login", p)

	if err != nil {
		return r, err
	}

	params := map[string]string{
		"code":         w.code,
		"encrypt_data": w.encryptData,
		"iv":           w.iv,
	}

	var result ReturnResult
	if w.Wafer {
		result.PrintWafer = true
	} else {
		result.PrintWafer = false
	}

	data, err := w.sendRequest(INTERFACE_LOGIN, params)
	if err != nil {
		if w.Wafer {
			result.PrintResult = map[string]string{
				WX_SESSION_MAGIC_ID: "1",
				"error":             ERR_LOGIN_FAILED,
				"message":           err.Error(),
			}
		}
	} else {
		result.Data = data
		if w.Wafer {
			printRet := make(map[string]interface{})
			printRet[WX_SESSION_MAGIC_ID] = "1"
			printRet["session"] = map[string]string{
				"id":   data.Id,
				"skey": data.Skey,
			}
			result.PrintResult = printRet
		}
	}

	return result, err
}

// 通过wafer-session-server进行登录状态验证
// 如果小程序前端使用了wafer sdk，则不需要传递参数
// 如果单独使用，则需要传递相应参数
//  param1: id string wafer-session-server生成的uuid
//  param2: skey string wafer-session-server生成的skey
func (w *WaferSDK) Check(p ...string) (r ReturnResult, err error) {
	err = w.validate("check", p)

	if err != nil {
		return r, err
	}

	params := map[string]string{
		"id":   w.id,
		"skey": w.skey,
	}

	var result ReturnResult
	data, err := w.sendRequest(INTERFACE_CHECK, params)
	if err != nil {
		var errCode string
		if err.(WaferError).GetCode() == RETURN_CODE_SKEY_EXPIRED || err.(WaferError).GetCode() == RETURN_CODE_WX_SESSION_FAILED {
			errCode = ERR_INVALID_SESSION
		} else {
			errCode = ERR_CHECK_LOGIN_FAILED
		}

		if w.Wafer {
			result.PrintWafer = true
			result.PrintResult = map[string]string{
				WX_SESSION_MAGIC_ID: "1",
				"error":             errCode,
				"message":           err.Error(),
			}
		}
	} else {
		result.PrintWafer = false
		result.Data = data
	}

	return result, err
}

// 发送请求到wafer-session-server
func (w *WaferSDK) sendRequest(api string, params map[string]string) (r ResponseData, err error) {
	url := w.AuthServerUrl

	reqBody := w.packReqData(api, params)

	body, err := json.Marshal(reqBody)
	if err != nil {
		err = WaferError{RETURN_CODE_SERVE_ERR, err.Error()}
		return r, err
	}

	resp, err := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(body))
	if err != nil {
		err = WaferError{RETURN_CODE_SERVE_ERR, err.Error()}
		return r, err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = WaferError{RETURN_CODE_SERVE_ERR, err.Error()}
		return r, err
	}

	var respBody ResponseBody
	err = json.Unmarshal(content, &respBody)
	if err != nil {
		err = WaferError{RETURN_CODE_SERVE_ERR, err.Error()}
		return r, err
	}

	if respBody.Code != RETURN_CODE_SUCCESS {
		err = WaferError{respBody.Code, respBody.Message}
		return r, err
	}

	r = respBody.Data

	return r, nil
}

// 组装请求数据
func (w *WaferSDK) packReqData(api string, params map[string]string) RequestBody {
	var ri RequestInterface
	ri.Appid = w.Appid
	ri.InterfaceName = api
	ri.Params = params

	var r RequestBody
	r.Version = 1
	r.ComponentName = "MA"
	r.ReqInterface = ri

	return r
}

// 校验信息
func (w *WaferSDK) validate(act string, p []string) error {

	if w.Appid == "" {
		return WaferError{RETURN_CODE_PARAM_ERR, "appid is empty"}
	}

	if w.AuthServerUrl == "" {
		return WaferError{RETURN_CODE_PARAM_ERR, "AuthServerUrl is empty"}
	}

	if w.Wafer && w.Request == nil {
		return WaferError{RETURN_CODE_PARAM_ERR, "use wafer sdk, http request required"}
	}

	if act == "login" {

		if w.Wafer {
			w.code = w.Request.Header.Get(WX_HEADER_CODE)
			w.encryptData = w.Request.Header.Get(WX_HEADER_ENCRYPTED_DATA)
			w.iv = w.Request.Header.Get(WX_HEADER_IV)
			if w.code == "" {
				return WaferError{RETURN_CODE_HEADER_ERR,
					fmt.Sprintf("use wafer sdk, request header %s is empty", WX_HEADER_CODE)}
			}
			if w.encryptData == "" {
				return WaferError{RETURN_CODE_HEADER_ERR,
					fmt.Sprintf("use wafer sdk, request header %s is empty", WX_HEADER_ENCRYPTED_DATA)}
			}
			if w.iv == "" {
				return WaferError{RETURN_CODE_HEADER_ERR,
					fmt.Sprintf("use wafer sdk, request header %s is empty", WX_HEADER_IV)}
			}
		} else {
			if len(p) < 3 {
				return WaferError{RETURN_CODE_PARAM_ERR,
					fmt.Sprintf("login func require 3 params, %d params given", len(p))}
			}
			w.code = p[0]
			w.encryptData = p[1]
			w.iv = p[2]
			if w.code == "" {
				return WaferError{RETURN_CODE_PARAM_ERR, "login func params[0] code is empty"}
			}
			if w.encryptData == "" {
				return WaferError{RETURN_CODE_PARAM_ERR, "login func params[1] encryptData is empty"}
			}
			if w.iv == "" {
				return WaferError{RETURN_CODE_PARAM_ERR, "login func params[2] iv is empty"}
			}
		}
	} else if act == "check" {

		if w.Wafer {
			w.id = w.Request.Header.Get(WX_HEADER_ID)
			w.skey = w.Request.Header.Get(WX_HEADER_SKEY)
			if w.id == "" {
				return WaferError{RETURN_CODE_HEADER_ERR,
					fmt.Sprintf("use wafer sdk, request header %s is empty", WX_HEADER_ID)}
			}
			if w.skey == "" {
				return WaferError{RETURN_CODE_HEADER_ERR,
					fmt.Sprintf("use wafer sdk, request header %s is empty", WX_HEADER_SKEY)}
			}
		} else {
			if len(p) < 2 {
				return WaferError{RETURN_CODE_PARAM_ERR,
					fmt.Sprintf("login func require 2 params, %d params given", len(p))}
			}
			w.id = p[0]
			w.skey = p[1]
			if w.id == "" {
				return WaferError{RETURN_CODE_PARAM_ERR, "check func params[0] id is empty"}
			}
			if w.skey == "" {
				return WaferError{RETURN_CODE_PARAM_ERR, "check func params[1] skey is empty"}
			}
		}
	} else {
		return WaferError{RETURN_CODE_PARAM_ERR, "invalid act, must in [login,check]"}
	}

	return nil
}
