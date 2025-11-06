package runner

import (
	"context"
	"fmt"
	"strings"
	"time"
	"encoding/json"
	"bytes"

	"github.com/SwampPear/argo/pkg/state"
	"github.com/SwampPear/argo/pkg/llm"
)

// Analyzer object.
type Analyzer struct {
	LLM 							 llm.LLMClient // abstract LLM interface
	MaxTokensPerBatch  int	   	 		 // batching parameter 		 (1800)
	MaxLogsPerBatch    int     	 		 // batching parameter 		 (120)
	ApproxTokensPerLog int     	 		 // heuristic for batching (16)
	WarnThreshold 		 float64 	 		 // status threshold 			 (0.60)
	ErrThreshold  		 float64 	 		 // status threshold 			 (0.85)
}

// TODO: beef up
type BugReport struct {
	Score       float64  `json:"score"`
	Explanation string   `json:"explanation"`
	Indicators  []string `json:"indicators,omitempty"`
}

// Starts the analyzer (single pass, no streaming).
func (a *Analyzer) Start(m *state.Manager) error {
	// start message
	log(m, "start", "OK", "Analyzer started.", "", "", 0)
	
	// initialize llm
	if err := a.initLLM(m); err != nil {
		log(m, "error", "ERR", err.Error(), "", "", 0)
	}

	// defaults
	if a.MaxTokensPerBatch <= 0 {
		a.MaxTokensPerBatch = 1800
	}
	if a.ApproxTokensPerLog <= 0 {
		a.ApproxTokensPerLog = 32
	}
	if a.MaxLogsPerBatch <= 0 {
		a.MaxLogsPerBatch = 8
	}
	if a.WarnThreshold <= 0 {
		a.WarnThreshold = 0.60
	}
	if a.ErrThreshold <= 0 {
		a.ErrThreshold = 0.85
	}

	// acquire logs
	logs := m.Logs()

	// no logs
	if len(logs) == 0 {
		log(m, "noLogs", "OK", "No logs to analyze.", "", "", 0)
		return nil
	}

	// make batch
	batches := a.makeBatches(m, logs)

	// process batches
	startAll := time.Now()
	ctx := context.Background()
	for bi, batch := range batches {
		bStart := time.Now()

		// llm context
		score, explanation, indicators, raw := a.callLLM(ctx, batch)

		// status
		status := "OK"

		// summary
		summary := fmt.Sprintf(
			"[Batch %d/%d] bug-likelihood=%.2f — %s | signals: %s | raw: %s",
			bi+1, len(batches), clamp(score, 0, 1),
			truncate(explanation, 120),
			truncate(strings.Join(indicators, " | "), 120),
			truncate(raw, 140),
		)

		// emit report log
		log(m, "batchComplete", status, summary, fmt.Sprintf("%d-logs", len(batch)), time.Since(bStart).String(), 
				clamp(score, 0, 1))

		fmt.Println(raw)
	}

	log(m, "complete", "OK", "Analyzation complete.", fmt.Sprintf("%d-logs", len(logs)), 
			time.Since(startAll).String(), 0)

	return nil
}

// Emits an analyzer event.
func log(m *state.Manager, action string, status string, summary string, target string, duration string, confidence float64) {
	// defaults
	if target == "" {
		target = "-"
	}
	if duration == "" {
		target = "0ms"
	}
	if confidence == 0 {
		confidence = 1.0
	}

	// append log to state
	m.AppendLog(state.LogEntry{
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
		Module:     "Analyzer",
		Action:     action,
		Target:     "-",
		Status:     status,
		Duration:   "0ms",
		Confidence: 1.0,
		Summary:    summary,
	})
}

// Computer max lines per batch.
func (a *Analyzer) maxLines() int {
	maxLines := a.MaxTokensPerBatch / a.ApproxTokensPerLog
	if maxLines <= 0 {
		maxLines = 1
	}
	if maxLines > a.MaxLogsPerBatch {
		maxLines = a.MaxLogsPerBatch
	}

	return maxLines
}

// Make batches
func (a *Analyzer) makeBatches(m *state.Manager, logs []state.LogEntry) [][]state.LogEntry {// filter
	var batches [][]state.LogEntry

	maxLines := a.maxLines()
	
	// slice batches
	var count int
	batch := make([]state.LogEntry, 0, maxLines)
	for _, e := range logs {
		batch = append(batch, e)
		count++

		// batch full
		if count >= maxLines {
			batches = append(batches, batch)
			batch = make([]state.LogEntry, 0, maxLines)
			count = 0
		}
	}

	// final batch
	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	return batches
}

// Ensures LLM is configured on analyzation start.
func (a *Analyzer) initLLM(m *state.Manager) error {
	cfg := m.GetState().Settings.LLM

	switch cfg.Provider {
	case "ollama":
		c := &llm.OllamaClient{}
		if err := c.Init(m); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Invalid config provider: %s", cfg.Provider)
	}

	return nil
}

// Removes the timestamp from a log string.
func trimTimestamp(ts string) string {
	if i := strings.IndexByte(ts, '.'); i > 0 {
		return ts[:i] + "Z"
	}

	return ts
}

// Removes newlines and trims whitspace.
func safe(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return s
}

// Truncates a string to n characters.
func truncate(s string, n int) string {
	if len(s) <= n { 
		return s 
	}
	if n <= 3 { 
		return s[:n] 
	}

	return s[:n-3] + "..."
}

// Clamps a number between a range.
func clamp(x, lo, hi float64) float64 {
	if x < lo {
		return lo
	}
	if x > hi {
		return hi
	}

	return x
}

// Cleans indicators
func cleanIndicators(in []string) []string {
	var out []string
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s != "" {
			out = append(out, s)
		}
	}

	return out
}

// Calls an LLM.
func (a *Analyzer) callLLM(ctx context.Context, logs []state.LogEntry) (score float64, explanation string, indicators []string, raw string) {
	// system prompt
	system := `
	You are SEC-REPORTER, a structured security assistant. Your job is to summarize potential vulnerabilities into a 
	single, valid JSON object. Use precise, professional language similar to HackerOne reports. Your report should contain
	a brief exploit vector. You should also give a brief score for how actionable and vulnerable the code may be. Output 
	JSON only — no explanations, no markdown. If information is missing, set the field to the default value for the data 
	type. Avoid sensitive data. You should be scrutinous in your scoring and should focus on web bugs.
	Follow this schema exactly:
	{
		"score": 0.0-1.0,
		"explanation": "<=160 chars",
		"indicators": ["short phrases"]
	}
	`

	// user prompt
	var user bytes.Buffer
	fmt.Fprintln(&user, "Summarize the following logs into one JSON object using the system schema above.\n LOGS:")
	for _, e := range logs {
		fmt.Fprintf(&user, "- ts=%s mod=%s act=%s tgt=%s st=%s dur=%s conf=%.2f :: %s\n",
			trimTimestamp(e.Timestamp), safe(e.Module), safe(e.Action), safe(e.Target),
			safe(e.Status), safe(e.Duration), e.Confidence, truncate(e.Summary, 140))
	}

	// query LLM
	resp, err := a.LLM.Complete(ctx, system, user.String())
	raw = strings.TrimSpace(resp)

	if err != nil {
		return 0.25, "LLM analysis failed", []string{"llm-call-error"}, err.Error()
	}

	// format report
	if raw != "" {
		var rep BugReport
		if json.Unmarshal([]byte(raw), &rep) == nil && rep.Explanation != "" {
			return clamp(rep.Score, 0, 1), strings.TrimSpace(rep.Explanation), cleanIndicators(rep.Indicators), raw
		}
	}

	return 0.25, "LLM analysis failed; defaulting to low risk", []string{"llm-call-error"}, raw
}