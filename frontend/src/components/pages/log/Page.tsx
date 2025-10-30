import { useEffect, useRef, useState } from 'react'
import { EventsOn } from '../../../../wailsjs/runtime'
import styles from './Page.module.css'

type LogEvent = {
  timestamp: string
  run_id: string
  step_id: string
  phase?: string
  module: string
  action: string
  target: string
  status: string
  duration: string
  confidence: number
  summary: string
  parent_step_id?: string | number
}

interface ICellProps { children: React.ReactNode }

const Cell = ({ children }: ICellProps) => (
  <td>
    <div className={styles.cellScroll}>{children}</div>
  </td>
)

const MAX_ROWS = 1000

const Page = () => {
  const [rows, setRows] = useState<LogEvent[]>([])
  const tbodyRef = useRef<HTMLTableSectionElement | null>(null)

  const appendRow = (e: any) => {
    // Defensive normalize in case fields are missing/typed oddly
    const ev: LogEvent = {
      timestamp: String(e?.timestamp ?? ''),
      run_id: String(e?.run_id ?? ''),
      step_id: String(e?.step_id ?? ''),
      phase: e?.phase ? String(e.phase) : undefined,
      module: String(e?.module ?? ''),
      action: String(e?.action ?? ''),
      target: String(e?.target ?? ''),
      status: String(e?.status ?? ''),
      duration: String(e?.duration ?? ''),
      confidence: Number(e?.confidence ?? 0),
      summary: String(e?.summary ?? ''),
      parent_step_id: e?.parent_step_id ?? undefined
    }

    setRows(prev => {
      const next = [...prev, ev]
      if (next.length > MAX_ROWS) next.shift()
      return next
    })
  }

  // Subscribe to backend events
  useEffect(() => {
    const off = EventsOn('log:event', appendRow)
    return () => { off() }
  }, [])

  // Auto-scroll to bottom on new rows
  useEffect(() => {
    if (!tbodyRef.current) return
    // Find the nearest scrollable container (the table wrapper)
    const scroller = tbodyRef.current.parentElement?.parentElement
    scroller?.scrollTo({ left: 0, top: scroller.scrollHeight, behavior: 'smooth' })
  }, [rows.length])

  return (
    <div className={styles.container}>
      <table className={styles.table}>
        <thead>
          <tr>
            <th>Timestamp</th>
            <th>ID</th>
            <th>Module</th>
            <th>Action</th>
            <th>Target</th>
            <th>Status</th>
            <th>Duration</th>
            <th>Confidence</th>
            <th>Summary</th>
            <th>Parent</th>
          </tr>
        </thead>
        <tbody ref={tbodyRef}>
          {rows.map((r, i) => (
            <tr key={`${r.run_id}-${r.step_id}-${i}`}>
              <Cell>{r.timestamp}</Cell>
              <Cell>{`${r.run_id} / ${r.step_id}`}</Cell>
              <Cell>{r.module}</Cell>
              <Cell>{r.action}</Cell>
              <Cell>{r.target}</Cell>
              <Cell><span className={`status ${r.status.toLowerCase()}`}>{r.status}</span></Cell>
              <Cell>{r.duration}</Cell>
              <Cell>{r.confidence.toFixed(2)}</Cell>
              <Cell>{r.summary}</Cell>
              <Cell>{r.parent_step_id ?? '-'}</Cell>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default Page