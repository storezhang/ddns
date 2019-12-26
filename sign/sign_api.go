package sign

import (
    "context"
)

// AutoSignResult 自动签到结果
type AutoSignResult struct {
    Success bool
    Before  string
    After   string
    Msg     string
}

// Signer 自动签到接口
type Signer interface {
    // AutoSign 自动签到
    AutoSign(ctx context.Context, cookies string) (result AutoSignResult, err error)
}
