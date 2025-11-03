import Menu from '@components/menu/Menu'
import Page from '@components/pages/Page'
import { useAppStore, type LogEntry } from '@state/state'
import { EventsOn } from '@wails/runtime'
import { useEffect } from 'react'
import './App.css'


const App = () => {
  const addLog = useAppStore(s => s.addLog)

  const handleEvent = (e: any) => {
    const entry: LogEntry = {
      id: String(e?.id ?? `${Date.now()}-${Math.random().toString(16).slice(2)}`),
      timestamp: String(e?.timestamp ?? '') || String(Date.now()),
      module: String(e?.module ?? ''),
      action: String(e?.action ?? ''),
      target: String(e?.target ?? ''),
      status: String(e?.status ?? 'OK').toUpperCase(),
      duration: String(e?.duration ?? ''),
      confidence: Number(e?.confidence ?? 0),
      summary: String(e?.summary ?? ''),
      parent_step_id: e?.parent_step_id ?? undefined
    }

    addLog(entry)
  }

  useEffect(() => {
    const off = EventsOn('log:event', handleEvent)
    return () => { off() }
  }, [])

  return (
    <div id="app">
      <Menu/>
      <div id="container">
        <Page/>
      </div>
    </div>
  )
}

export default App