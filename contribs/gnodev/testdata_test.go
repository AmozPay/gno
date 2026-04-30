package main

import (
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gnolang/gno/gnovm/pkg/integration"
	"github.com/rogpeppe/go-internal/testscript"
	"github.com/stretchr/testify/require"
)

func TestScripts(t *testing.T) {
	p := integration.NewTestingParams(t, "testdata")

	buildDir := t.TempDir()
	gnodevBin := filepath.Join(buildDir, "gnodev")
	buildCmd := exec.Command(
		"go",
		"build",
		"-ldflags",
		"-X github.com/gnolang/gno/tm2/pkg/version.Version=testscript-version",
		"-o",
		gnodevBin,
		".",
	)
	buildCmd.Dir = "."
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, string(output))

	if p.Cmds == nil {
		p.Cmds = make(map[string]func(ts *testscript.TestScript, neg bool, args []string))
	}

	p.Cmds["gnodev"] = func(ts *testscript.TestScript, neg bool, args []string) {
		err := ts.Exec(gnodevBin, args...)
		if err != nil {
			ts.Logf("gnodev command error: %+v", err)
		}

		if (err == nil) == neg {
			ts.Fatalf("unexpected gnodev command outcome (err=%t expected=%t)", err == nil, !neg)
		}
	}

	testscript.Run(t, p)
}
