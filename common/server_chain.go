package common

// ServerChan ServerChan推送设置
type ServerChan struct {
    Key string
}

// ServerChanRequest ServerChan调用请求体
type ServerChanRequest struct {
    Text string `json:"text"`
    Desp string `json:"desp"`
}
