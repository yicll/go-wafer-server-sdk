package wafer

const (
	WX_HEADER_CODE           string = "X-WX-Code"
	WX_HEADER_ENCRYPTED_DATA string = "X-WX-Encrypted-Data"
	WX_HEADER_IV             string = "X-WX-IV"
	WX_HEADER_ID             string = "X-WX-Id"
	WX_HEADER_SKEY           string = "X-WX-Skey"

	WX_SESSION_MAGIC_ID string = "F2C224D4-2BCE-4C64-AF9F-A6D872000D1A"

	ERR_LOGIN_FAILED       string = "ERR_LOGIN_FAILED"
	ERR_INVALID_SESSION    string = "ERR_INVALID_SESSION"
	ERR_CHECK_LOGIN_FAILED string = "ERR_CHECK_LOGIN_FAILED"

	INTERFACE_LOGIN string = "qcloud.cam.id_skey"
	INTERFACE_CHECK string = "qcloud.cam.auth"

	RETURN_CODE_SUCCESS           int = 0
	RETURN_CODE_PARAM_ERR         int = 1001
	RETURN_CODE_HEADER_ERR        int = 1002
	RETURN_CODE_SERVE_ERR         int = 2000
	RETURN_CODE_SKEY_EXPIRED      int = 60011
	RETURN_CODE_WX_SESSION_FAILED int = 60012
)
