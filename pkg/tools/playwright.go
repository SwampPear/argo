package tools
/*
package tools

import (
	"bufio"
	"context"
	"encoding/json"
	"os/exec"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type PW struct {
	Ctx func() context.Context // inject Wails ctx getter
}

// StartInteractive launches headful Chromium and streams events until user closes it or ctx is canceled.
func (p PW) StartInteractive(ctx context.Context, args map[string]any) error {
	cmd := exec.CommandContext(ctx, "node", "./tools/playwright/pw_runner.mjs")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil { return err }
	_ = json.NewEncoder(stdin).Encode(args) // send config
	_ = stdin.Close()

	sc := bufio.NewScanner(stdout)
	for sc.Scan() {
		var m map[string]any
		if json.Unmarshal(sc.Bytes(), &m) != nil { continue }
		ch := "log:event"
		if t, ok := m["type"].(string); ok && t == "obs" { ch = "obs:event" }
		runtime.EventsEmit(p.Ctx(), ch, m)
	}
	_ = cmd.Wait()
	return ctx.Err()
}
*/