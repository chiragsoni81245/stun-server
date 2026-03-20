package config

import (
	"flag"
	"strings"
	"time"

	"github.com/chiragsoni81245/stun-server/internal/stunserver"
)

func Load() stunserver.Config {
	addrs := flag.String("addrs", ":3478,:3479", "comma-separated listen addresses")
	flag.Parse()

	return stunserver.Config{
		Addrs:          strings.Split(*addrs, ","),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		RateLimitPerIP: 100,
	}
}
