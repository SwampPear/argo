import { useEffect, useMemo, useRef } from 'react'
import { EventsOn } from '../../../../wailsjs/runtime'
import { useAppStore, type LogEntry } from '../../../state/useAppStore'
import styles from './Page.module.css'

interface ICellProps {
  children: React.ReactNode
  width: number
}

const Cell = ({ children, width }: ICellProps) => (
  <td style={width ? { width: `${width}%`, maxWidth: `${width}%` } : {}}>
    <div className={styles.cellScroll}>{children}</div>
  </td>
)

const MAX_ROWS = 1000

const Page = () => {
  const tbodyRef = useRef<HTMLTableSectionElement | null>(null)

  const logs = useAppStore(s => s.logs)
  const addLog = useAppStore(s => s.addLog)
  const setLogs = useAppStore(s => s.setLogs)

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

  useEffect(() => {
    if (logs.length > MAX_ROWS) setLogs(logs.slice(-MAX_ROWS))
  }, [logs.length])

  // Derive render rows from LogEntry only
  const rows = useMemo(() => {
    return logs.map(l => {
      return {
        _key: l.id,
        timestamp: l.timestamp ?? '-',
        module: l.module ?? '-',
        action: l.action ?? '-',
        target: l.target ?? '-',
        status: l.status ?? '-',
        duration: l.duration ?? '-',
        confidence: Number(l.confidence ?? 0),
        summary: l.summary,
        parent_step_id: l.parent_step_id ?? '-'
      }
    })
  }, [logs])

  return (
    <div className={styles.container}>
      <table className={styles.table}>
        <thead>
          <tr>
            <th style={{ width: '4%' }}>#</th>
            <th style={{ width: '12%' }}>Timestamp</th>
            <th style={{ width: '10%' }}>Module</th>
            <th style={{ width: '10%' }}>Action</th>
            <th style={{ width: '12%' }}>Target</th>
            <th style={{ width: '10%' }}>Status</th>
            <th style={{ width: '10%' }}>Duration</th>
            <th style={{ width: '10%' }}>Confidence</th>
            <th style={{ width: '12%' }}>Summary</th>
            <th style={{ width: '10%' }}>Parent</th>
          </tr>
        </thead>
        <tbody ref={tbodyRef}>
          {rows.map((r, i) => (
            <tr key={`${r._key}-${i}`}>
              <Cell width={4}>{i + 1}</Cell>
              <Cell width={12}>{r.timestamp}</Cell>
              <Cell width={10}>{r.module}</Cell>
              <Cell width={10}>{r.action}</Cell>
              <Cell width={12}>{r.target}</Cell>
              <Cell width={10}>
                <span className={`status ${r.status.toLowerCase()}`}>{r.status}</span>
              </Cell>
              <Cell width={10}>{r.duration}</Cell>
              <Cell width={10}>{r.confidence.toFixed(2)}</Cell>
              <Cell width={12}>{r.summary}</Cell>
              <Cell width={10}>{r.parent_step_id ?? '-'}</Cell>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default Page
