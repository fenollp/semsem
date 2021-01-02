package semsem_test

import (
	"testing"

	"github.com/fenollp/semsem"
	"github.com/stretchr/testify/require"
)

func TestNoBreakingAPIChanges(t *testing.T) {
	x, err := semsem.X("./...")
	require.NoError(t, err)

	t.Logf(">>> x.Pkgs: %d", len(x.Pkgs))
	for i, v := range x.Pkgs {
		t.Logf("    x.Pkgs[%d] = %+v", i, v)
	}

	t.Logf(">>> x.Self: %d", len(x.Self))
	for k, v := range x.Self {
		t.Logf("    x.Self[%+v] = %+v", k, v)
	}

	t.Logf(">>> x.Imports: %d", len(x.Imports))
	for k, v := range x.Imports {
		t.Logf("    x.Imports[%+v] = %+v", k, v)
	}

	t.Logf(">>> x.OldObjs: %d", len(x.OldObjs))
	for k, v := range x.OldObjs {
		t.Logf("    x.OldObjs[%+v] = %+v", k, v)
	}

	t.Logf(">>> x.Report: %+v", x.Report)
}
