package sign

import (
	"context"
)

// Signer 自动签到接口
type Signer interface {
	// AutoSign 自动签到
	AutoSign(ctx context.Context, cookies string)
}
