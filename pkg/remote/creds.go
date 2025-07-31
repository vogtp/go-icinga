package remote

import (
	"context"
	"log/slog"

	"github.com/spf13/viper"

	"github.com/vogtp/go-icinga/pkg/credsmgr"
)

func GetPassword(ctx context.Context) string {
	pass := viper.GetString(passwordFlag)
	if len(pass) > 0 {
		return pass
	}
	username := viper.GetString(UserFlag)
	if len(username) < 1 {
		return ""
	}
	cm, err := credsmgr.Get(ctx)
	if err != nil {
		slog.Warn("Cannot get credential manager", "err", err)
	}
	pw, err := cm.Load(ctx, username)
	if err != nil {
		slog.Info("Cannot get credential for user", "err", err, "username", username)
	}
	return string(pw)
}
