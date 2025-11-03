package runner

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
	"os"

	"github.com/SwampPear/argo/pkg/state"
)

// Interface TODO: back with specific LLM.
type LLMClient interface {
	Complete(ctx context.Context, system, prompt string) (string, error)
}

type Analyzer struct {
	LLM LLMClient

	// batching
	MaxTokensPerBatch int // default 1800
	MaxLogsPerBatch   int // default 120

	// heuristic
	ApproxTokensPerLog int // default 16

	// status thresholds
	WarnThreshold float64 // default 0.60
	ErrThreshold  float64 // default 0.85
}

type BugReport struct {
	Score       float64  `json:"score"`
	Explanation string   `json:"explanation"`
	Indicators  []string `json:"indicators,omitempty"`
}

// Starts the analyzer (single pass, no streaming).
func (a *Analyzer) Start(m *state.Manager) error {
	// start message
	m.AppendLog(state.LogEntry{
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
		Module:     "Analyzer",
		Action:     "start",
		Target:     "-",
		Status:     "OK",
		Duration:   "0ms",
		Confidence: 1.0,
		Summary:    "Analyzer started.",
	})

	a.ensureLLM()

	// defaults
	if a.MaxTokensPerBatch <= 0 {
		a.MaxTokensPerBatch = 1800
	}
	if a.ApproxTokensPerLog <= 0 {
		a.ApproxTokensPerLog = 16
	}
	if a.MaxLogsPerBatch <= 0 {
		a.MaxLogsPerBatch = 120
	}
	if a.WarnThreshold <= 0 {
		a.WarnThreshold = 0.60
	}
	if a.ErrThreshold <= 0 {
		a.ErrThreshold = 0.85
	}

	// acquire logs
	logs := m.Logs()
	if logs == nil {
		return fmt.Errorf("Error retrieving logs.")
	}

	// no logs
	if len(logs) == 0 {
		m.AppendLog(state.LogEntry{
			Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
			Module:     "Analyzer",
			Action:     "analyze",
			Target:     "-",
			Status:     "OK",
			Duration:   "0ms",
			Confidence: 1.0,
			Summary:    "No logs to analyze.",
		})

		return nil
	}

	// order by time
	startAll := time.Now()
	sort.SliceStable(logs, func(i, j int) bool {
		ti, tj := logs[i].Timestamp, logs[j].Timestamp
		if ti != tj {
			return ti < tj
		}

		return false
	})

	// batch
	batches := a.makeBatches(logs)

	ctx := context.Background()
	for bi, batch := range batches {
		bStart := time.Now()

		score, explanation, indicators, raw := a.callLLM(ctx, batch)

		status := "OK"
		switch {
		case score >= a.ErrThreshold:
			status = "ERROR"
		case score >= a.WarnThreshold:
			status = "WARN"
		}

		// emit per batch analysis log
		summary := fmt.Sprintf(
			"[Batch %d/%d] bug-likelihood=%.2f — %s | signals: %s | raw: %s",
			bi+1, len(batches), clamp(score, 0, 1),
			truncate(explanation, 120),
			truncate(strings.Join(indicators, " | "), 120),
			truncate(raw, 140),
		)

		m.AppendLog(state.LogEntry{
			Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
			Module:     "Analyzer",
			Action:     "analyzer-batchComplete",
			Target:     fmt.Sprintf("%d-logs", len(batch)),
			Status:     status,
			Duration:   time.Since(bStart).String(),
			Confidence: clamp(score, 0, 1),
			Summary:    summary,
		})
	}

	m.AppendLog(state.LogEntry{
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
		Module:     "Analyzer",
		Action:     "analyzer-complete",
		Target:     fmt.Sprintf("%d-logs", len(logs)),
		Status:     "OK",
		Duration:   time.Since(startAll).String(),
		Confidence: 1.0,
		Summary:    "Analyzation complete.",
	})

	return nil
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
func (a *Analyzer) makeBatches(logs []state.LogEntry) [][]state.LogEntry {
	var batches [][]state.LogEntry

	maxLines := a.maxLines()
	
	// slice by count
	cur := make([]state.LogEntry, 0, maxLines)
	var curCount int
	for _, e := range logs {
		cur = append(cur, e)
		curCount++
		if curCount >= maxLines {
			batches = append(batches, cur)
			cur = make([]state.LogEntry, 0, maxLines)
			curCount = 0
		}
	}

	if len(cur) > 0 {
		batches = append(batches, cur)
	}

	return batches
}

// Ensures LLM is configured on analyzation start.
func (a *Analyzer) ensureLLM() error {
	if a.LLM != nil {
		return nil
	}

	// Try Ollama
	ollamaURL := os.Getenv("OLLAMA_HOST")
	ollamaModel := os.Getenv("OLLAMA_MODEL")

	a.LLM = &OllamaClient{
		BaseURL:     ollamaURL,
		Model:       ollamaModel,
		Temperature: 0.0,
		Timeout:     60 * time.Second,
	}

	return nil
}

// Calls an LLM.
func (a *Analyzer) callLLM(ctx context.Context, logs []state.LogEntry) (score float64, explanation string, indicators []string, raw string) {
	// format prompt
	/*
	system := "You are a senior reliability engineer. Analyze execution logs for hidden bugs and return STRICT JSON."

	var b bytes.Buffer
	fmt.Fprintln(&b, "You will receive execution logs (compact lines).")
	fmt.Fprintln(&b, "Infer likelihood that a hidden bug exists, even if no error was thrown.")
	fmt.Fprintln(&b, "Consider: status flips, retries/timeouts, long durations, inconsistent sequences, masked failures, flaky steps.")
	fmt.Fprintln(&b, "Return ONLY JSON: {\"score\":0..1,\"explanation\":\"<=160 chars\",\"indicators\":[\"short phrases\"]}")
	fmt.Fprintln(&b, "")
	fmt.Fprintln(&b, "Logs:")
	for _, e := range logs {
		fmt.Fprintf(&b, "- ts=%s mod=%s act=%s tgt=%s st=%s dur=%s conf=%.2f :: %s\n",
			trimTS(e.Timestamp), safe(e.Module), safe(e.Action), safe(e.Target),
			safe(e.Status), safe(e.Duration), e.Confidence, truncate(e.Summary, 140))
	}
	*/

	// query LLM
	/*
	resp, err := a.LLM.Complete(ctx, system, b.String())
	raw = strings.TrimSpace(resp)
	*/

	// format report
	/*
	if err == nil && raw != "" {
		var rep BugReport
		if json.Unmarshal([]byte(raw), &rep) == nil && rep.Explanation != "" {
			return clamp(rep.Score, 0, 1), strings.TrimSpace(rep.Explanation), cleanIndicators(rep.Indicators), raw
		}
	}
	*/

	return 0.25, "LLM analysis failed; defaulting to low risk", []string{"llm-call-error"}, raw
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