import { useAppStore } from '@state/state'
import { useMemo, useRef } from 'react'
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

  const logs = useAppStore(s => s.state.logs)
  const scopeFilter = useAppStore(s => s.state.scope_filter)
  const settings = useAppStore(s => s.state.settings)

  const rows = useMemo(() => {
    let filtered = logs.length > MAX_ROWS ? logs.slice(-MAX_ROWS) : logs
    if (scopeFilter) {
      filtered = filtered.filter((log) => {
        if (log.module == 'Analyzer') return true;

        for (let i = 0; i < settings.Assets.InScope.length; i++) {
          let hostname = settings.Assets.InScope[i].Hostname
          hostname += hostname.endsWith('/') ? '' : '/'
  
          if (log.target?.includes(hostname)) {
            return true
          }
        }
  
        return false
      })
    }
    
    return filtered.map((l, i) => {
      const parent = (l as any)?.parent_step_id ?? '-'
      return {
        _key: String(l.step ?? l.id ?? i), // prefer step, then id, then index
        index: i,
        timestamp: l.timestamp ?? '-',
        module: l.module ?? '-',
        action: l.action ?? '-',
        target: l.target ?? '-',
        status: (l.status ?? '-'),
        duration: l.duration ?? '-',
        confidence: typeof l.confidence === 'number' ? l.confidence : Number(l.confidence ?? 0),
        summary: l.summary ?? '-',
        parent_step_id: parent
      }
    })
  }, [logs, scopeFilter])

  return (
    <div className={styles.container}>
      <table className={styles.table}>
        <thead>
          <tr>
            <th style={{ width: '4%' }}>#</th>
            <th style={{ width: '14%' }}>Timestamp</th>
            <th style={{ width: '10%' }}>Module</th>
            <th style={{ width: '8%' }}>Action</th>
            <th style={{ width: '12%' }}>Target</th>
            <th style={{ width: '6%' }}>Status</th>
            <th style={{ width: '6%' }}>Dur</th>
            <th style={{ width: '6%' }}>Conf</th>
            <th style={{ width: '28%' }}>Summary</th>
            <th style={{ width: '6%' }}>Parent</th>
          </tr>
        </thead>
        <tbody ref={tbodyRef}>
          {rows.map((r) => (
            <tr key={r._key}>
              <Cell width={4}>{r.index + 1}</Cell>
              <Cell width={14}>{r.timestamp}</Cell>
              <Cell width={10}>{r.module}</Cell>
              <Cell width={8}>{r.action}</Cell>
              <Cell width={12}>{r.target}</Cell>
              <Cell width={6}>
                <span className={`status ${String(r.status).toLowerCase()}`}>{r.status}</span>
              </Cell>
              <Cell width={6}>{r.duration}</Cell>
              <Cell width={6}>{r.confidence.toFixed(2)}</Cell>
              <Cell width={28}>{r.summary}</Cell>
              <Cell width={6}>{r.parent_step_id}</Cell>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default Page
