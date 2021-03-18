package tcp

import (
	"errors"
	"net"
	"time"

	"github.com/spf13/viper"
)

// Check returns a test
func Check(v *viper.Viper) error {

	ep := v.GetString("Endpoint")

	if ep == "" {
		return errors.New("No endpoint set")
	}

	timeout := time.Second

	dialer := func() (net.Conn, error) {
		return net.DialTimeout("tcp", ep, timeout)
	}

	conn, err := dialer()
	if err == nil {
		conn.Close()
		return nil
	}

	return err
}
