//go:build !linux

package credsmgr

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vogtp/go-creds-mgr/pkg/creds"
)

var secrets_pass = "SECRETS PASSWORD"

// Get returns a credential manager instance
func Get(ctx context.Context) (creds.Manager, error) {
	return nil, fmt.Errorf("Credentials manager is only implemented for linux")
}

// Command returns a cobra.Command allowing to manage credemtials
func Command(ctx context.Context) (*cobra.Command, error) {
	return nil, fmt.Errorf("Credentials manager is only implemented for linux")
}
