import { useEffect, useRef, useState } from 'react'
import { EventsOn } from '../../../../wailsjs/runtime'
import styles from './Page.module.css'

type LogEvent = {
  step: number
  id: string
  timestamp: string
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

interface ICellProps { 
  children: React.ReactNode
  width: number
}

const Cell = ({ children, width }: ICellProps) => (
  <td style={width ? {width: `${width}%`} : {}}>
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
      step: Number(e?.step ?? 0),
      id: String(e?.id ?? ''),
      timestamp: String(e?.timestamp ?? ''),
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
            <th style={{width: '3%'}}></th>
            <th style={{width: '10%'}}>ID</th>
            <th style={{width: '10%'}}>Timestamp</th>
            <th style={{width: '10%'}}>Module</th>
            <th style={{width: '10%'}}>Action</th>
            <th style={{width: '10%'}}>Target</th>
            <th style={{width: '10%'}}>Status</th>
            <th style={{width: '10%'}}>Duration</th>
            <th style={{width: '10%'}}>Confidence</th>
            <th style={{width: '10%'}}>Summary</th>
            <th style={{width: '10%'}}>Parent</th>
          </tr>
        </thead>
        <tbody ref={tbodyRef}>
          {rows.map((r, i) => (
            <tr key={`${r.id}-${i}`}>
              <Cell width={3}>{r.step}</Cell>
              <Cell width={10}>{r.id}</Cell>
              <Cell width={10}>{r.timestamp}</Cell>
              <Cell width={10}>{r.module}</Cell>
              <Cell width={10}>{r.action}</Cell>
              <Cell width={10}>{r.target}</Cell>
              <Cell width={10}><span className={`status ${r.status.toLowerCase()}`}>{r.status}</span></Cell>
              <Cell width={10}>{r.duration}</Cell>
              <Cell width={10}>{r.confidence.toFixed(2)}</Cell>
              <Cell width={10}>{r.summary}</Cell>
              <Cell width={10}>{r.parent_step_id ?? '-'}</Cell>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default Page