// Package exec provides an execution check for running commands
package exec

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/spf13/viper"
)

// Check returns a test
func Check(v *viper.Viper) error {

	command := v.GetString("Command")

	if command == "" {
		return errors.New("no command set")
	}

	args := v.GetStringSlice("Args")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("error: %s\noutput: %s", err.Error(), string(output))
	}

	return nil
}
