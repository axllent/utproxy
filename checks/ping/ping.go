package ping

import (
	"errors"
	"fmt"
	"net"
	"runtime"
	"time"

	"github.com/go-ping/ping"
	"github.com/spf13/viper"
)

type response struct {
	addr *net.IPAddr
	rtt  time.Duration
}

// Check returns a test
func Check(v *viper.Viper) error {

	ep := v.GetString("Endpoint")

	if ep == "" {
		return errors.New("No endpoint set")
	}

	pinger, err := ping.NewPinger(ep)
	if err != nil {
		return err
	}

	// required for Windows
	if runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}

	pinger.Count = 1
	pinger.Timeout = time.Second

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return err
	}

	stats := pinger.Statistics()

	if stats.PacketLoss != 0 {
		return fmt.Errorf("ping %s failed", ep)
	}

	return nil
}
