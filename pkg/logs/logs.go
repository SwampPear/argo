package logs

import (
	"bufio"
	"encoding/json"
	"io"
	"github.com/SwampPear/argo/pkg/state"
)

// Processes stream from Playwright runner.
func Stream(r io.Reader, append func(state.LogEntry)) {
	sc := bufio.NewScanner(r)
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 1024*1024)

	for sc.Scan() {
		line := sc.Text()
		var entry state.LogEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			entry = state.LogEntry{
				Summary: line,
				Status:  "OK",
				Module:  "Playwright",
				Action:  "log",
			}
		}
		append(entry)
	}
}
