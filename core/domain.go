package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/goexl/uda"
)

// Domain 域名
type Domain struct {
	name      string
	subdomain string
	typ       uda.Type
	value     string
	ttl       time.Duration
	prefix    string
	staff     string
}

// NewDomain 创建域名
func NewDomain(
	name string, subdomain string,
	typ uda.Type, value string, ttl time.Duration,
	prefix string, staff string,
) *Domain {
	return &Domain{
		name:      name,
		subdomain: subdomain,
		typ:       typ,
		value:     value,
		ttl:       ttl,
		prefix:    prefix,
		staff:     staff,
	}
}

func (d *Domain) Subdomain() (subdomain string) {
	subdomain = d.subdomain
	if `` != strings.TrimSpace(d.prefix) {
		subdomain = fmt.Sprintf(`%s.%s`, d.prefix, subdomain)
	}
	if `` != strings.TrimSpace(d.staff) {
		subdomain = fmt.Sprintf(`%s.%s`, subdomain, d.staff)
	}

	return
}

func (d *Domain) Final() string {
	return fmt.Sprintf(`%s.%s`, d.Subdomain(), d.name)
}

func (d *Domain) Ttl() time.Duration {
	return d.ttl
}

func (d *Domain) Value() string {
	return d.value
}

func (d *Domain) Name() string {
	return d.name
}
