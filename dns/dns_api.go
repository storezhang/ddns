package dns

type Resolver interface {
    Resolve(domain string, pr string, value string, dnsType string, ttl int)
}
