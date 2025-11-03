package runner

import (
	"os/exec"

	"github.com/SwampPear/argo/pkg/state"
	"github.com/SwampPear/argo/pkg/logs"
)

type Playwright struct{}

func (p *Playwright) Start(m *state.Manager) error {
	cmd := exec.Command("node", "./tools/playwright/pw_runner.mjs")
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	go logs.Stream(stdout, m.AppendLog)
	go logs.Stream(stderr, m.AppendLog)

	return cmd.Start()
}
