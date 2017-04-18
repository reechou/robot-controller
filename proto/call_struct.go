package proto

const (
	RESPONSE_OK = iota
	RESPONSE_ERR
)

type Response struct {
	Code int64       `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type CreateRobotHost struct {
	Host     string `json:"host"`
	RobotNum int    `json:"robotNum"`
}

type GetRobotReq struct {
	RobotWx string `json:"robotWx"`
}

type GetRobotListFromTypeReq struct {
	RobotType int `json:"robotType"`
}

type GetLoginRobotListFromTypeReq struct {
	RobotType int `json:"robotType"`
}
