package dns

// ResolveResult 自动签到结果
type ResolveResult struct {
    Success bool
    Before  string
    After   string
}

// Resolver 动态域名解析接口
type Resolver interface {
    // Resolve 域名解析
    Resolve(domain string, pr string, value string, dnsType string, ttl int) (result ResolveResult, err error)
}
