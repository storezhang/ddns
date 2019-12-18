package common

// ServerChan ServerChan推送设置
type ServerChan struct {
    ScKey string
}

// ServerChanRequest ServerChan调用请求体
type ServerChanRequest struct {
    Text string
    Desp string
}
