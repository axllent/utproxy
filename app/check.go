package app

import (
	"fmt"
	"time"

	"github.com/axllent/utproxy/checks/exec"
	"github.com/axllent/utproxy/checks/http"
	"github.com/axllent/utproxy/checks/mysql"
	"github.com/axllent/utproxy/checks/ping"
	"github.com/axllent/utproxy/checks/tcp"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
)

var (
	tmpCache = cache.New(55*time.Second, time.Minute)
)

// Check global function
// Caches responses to reduce load from multiple checkers.
// Success responses are cached for 55 seconds, failures 25 seconds.
func Check(test string) error {
	cached, found := tmpCache.Get(test)
	if found {
		if cached == nil {
			return nil
		}
		return cached.(error)
	}

	err := testCheck(test)

	cacheFor := cache.DefaultExpiration
	if err != nil {
		// reduce cache for errors
		cacheFor = 25 * time.Second
	}

	tmpCache.Set(test, err, cacheFor)

	return err
}

// realCheck function
func testCheck(test string) error {
	if !viper.Sub("Services").IsSet(test) {
		return fmt.Errorf("unknown")
	}

	model := viper.Sub("Services").Sub(test)

	if !model.IsSet("Type") {
		return fmt.Errorf("unknown")
	}

	disabled := model.GetBool("Disabled")
	if disabled {
		return fmt.Errorf("disabled")
	}

	mtype := model.GetString("Type")

	switch mtype {
	case "http":
		return http.Check(model)
	case "tcp":
		return tcp.Check(model)
	case "exec":
		return exec.Check(model)
	case "mysql":
		return mysql.Check(model)
	case "ping":
		return ping.Check(model)
	default:
		return fmt.Errorf("unknown")
	}
}
