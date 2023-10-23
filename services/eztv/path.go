package eztv

import (
	"strings"

	"github.com/abibby/eztvrss/config"
)

func eztvURL(p string) string {
	domain := config.EztvDomain
	domain, _ = strings.CutSuffix(domain, "/")
	p, _ = strings.CutPrefix(p, "/")
	return domain + "/" + p
}
