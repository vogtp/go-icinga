//go:build linux

package credsmgr

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vogtp/go-creds-mgr/pkg/creds"
	"github.com/vogtp/go-creds-mgr/pkg/tpmstorage"
)

var secrets_pass = "SECRETS PASSWORD"

type commander interface {
	CobraCommand() *cobra.Command
}

var instance creds.Manager

// Get returns a credential manager instance
func Get(ctx context.Context) (creds.Manager, error) {
	if instance != nil {
		return instance, nil
	}

	tpm, err := tpmstorage.New(ctx,
		tpmstorage.SecretsPassword(secrets_pass),
		tpmstorage.TPMDevice("/dev/tpmrm0", tpmstorage.Simulator),
	)
	if err != nil {
		return nil, fmt.Errorf("could not open tpm persistent storage: %w", err)
	}
	instance, err = creds.New(secrets_pass, tpm)
	return instance, err
}

// Command returns a cobra.Command allowing to manage credemtials
func Command(ctx context.Context) (*cobra.Command, error) {
	credsManager, err := Get(ctx)
	if err != nil {
		return nil, err
	}
	cmder, ok := credsManager.(commander)
	if !ok {
		return nil, fmt.Errorf("cannot get a cobra.Command")
	}
	return cmder.CobraCommand(), nil
}
