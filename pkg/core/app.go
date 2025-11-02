package core
/*
func (a *App) runPipeline(ctx context.Context, id string, cfg settings.Settings) {
	a.mu.Lock()
  start := time.Now()
  log := func(phase, module, action, target, status string, conf float64, d time.Duration, summary string) {
    a.emit(LogEvent{
      Step: a.step, ID: id, Timestamp: time.Now().Format(time.RFC3339), Module: module, Action: action, Target: target, 
			Status: status, Duration: d.String(), Confidence: conf, Summary: summary,
    })
    a.step++
  }
	a.mu.Unlock()

  // 1) Scope/policy check
  // 2) Recon (e.g., run nuclei safe set)
  // 3) Browser crawl (Playwright sidecar or HTTP runner)

  log("Triage","Run","complete","-", "OK", 1, time.Since(start), "Run finished")
	log("Triage","Run","complete","-", "OK", 1, time.Since(start), "Run finished")
}*/


// Runs the pipeline.
/*
func (a *App) Run() (string, error) {
  a.mu.Lock()
  cfg := a.settings
  a.mu.Unlock()

  id := xid.New().String()

  ctx, cancel := context.WithCancel(context.Background())

  a.mu.Lock()
  if a.active == nil {
    a.active = make(map[string]context.CancelFunc)
  }
  a.active[id] = cancel
  a.mu.Unlock()

  go a.runPipeline(ctx, id, cfg)
  return id, nil
}*/
