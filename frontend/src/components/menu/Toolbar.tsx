import { StartAnalyzer } from '@go/app/App'
import { useAppStore } from '@state/state'
import { LogError, LogInfo } from '@wails/runtime'
import styles from './Menu.module.css'

const EmptyToolbar = () => {
  return (
    <div className={styles.actionButtonContainer}></div>
  )
}

const LogsToolbar = () => {
  const scopeFilter = useAppStore(s => s.state.scope_filter)
  const setScopeFilter = useAppStore(s => s.setScopeFilter)

  const handleStartAnalyzer = async () => {
    try {
      await StartAnalyzer()
      LogInfo('Analyzer started.')
    } catch (e: any) {
      LogError(`Analzyer starting error: ${String(e?.message || e)}`)
    }
  }

  const handleStopAnalyzer = async () => {

  }

  return (
    <div className={styles.actionButtonContainer}>
      <button className={styles.actionButton} onClick={handleStartAnalyzer} title="Start Analysis">
        <svg xmlns="http://www.w3.org/2000/svg" className={styles.playIcon} viewBox="0 0 16 16">
          <path d="m11.596 8.697-6.363 3.692c-.54.313-1.233-.066-1.233-.697V4.308c0-.63.692-1.01 1.233-.696l6.363 3.692a.802.802 0 0 1 0 1.393"/>
        </svg>
      </button>
      <button className={styles.actionButton} onClick={handleStopAnalyzer} title="Stop Analysis">
        <svg xmlns="http://www.w3.org/2000/svg" className={styles.stopIcon} viewBox="0 0 16 16">
          <path d="M4 0h8a4 4 0 0 1 4 4v8a4 4 0 0 1-4 4H4a4 4 0 0 1-4-4V4a4 4 0 0 1 4-4z" strokeWidth={4}/>
        </svg>
      </button>
      <button className={styles.actionButton} onClick={() => setScopeFilter(!scopeFilter)} title="Scope">
        <svg xmlns="http://www.w3.org/2000/svg" className={scopeFilter ? styles.iconSet : styles.iconUnset} viewBox="0 0 16 16">
          <path d="M1.5 1.5A.5.5 0 0 1 2 1h12a.5.5 0 0 1 .5.5v2a.5.5 0 0 1-.128.334L10 8.692V13.5a.5.5 0 0 1-.342.474l-3 1A.5.5 0 0 1 6 14.5V8.692L1.628 3.834A.5.5 0 0 1 1.5 3.5z"/>
        </svg>
      </button>
    </div>
  )
}

const Toolbar = () => {
  const page = useAppStore(s => s.page)

  switch (page) {
  case 'logs':
    return <LogsToolbar/>
  case 'settings':
    return <EmptyToolbar/>
  default:
    return <EmptyToolbar/>
  }
}

export default Toolbar
