package powershell

import (
	"context"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	ctx := context.Background()
	s, err := New(ctx, "ITS-TEST25-05.its.unibas.ch", "d-idag-mon", os.Getenv("PW"))
	if err != nil {
		t.Fatal(err)
	}
	s.Run(ctx, "c:/programdata/syscheck/syscheck.exe net stat --remote.path c:/programdata/syscheck --remote.windows --verbose --remote.is_remote")
	t.Fail()
}
