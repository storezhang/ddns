package dns

// Record 缓存中的解析记录
type Record struct {
    ID     string
    Domain string
    Value  string
    RR     string
}

var recordCache map[string]*Record
