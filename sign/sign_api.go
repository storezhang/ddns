package sign

// Signer 自动签到接口
type Signer interface {
	// AutoSign 自动签到
	AutoSign()
}
