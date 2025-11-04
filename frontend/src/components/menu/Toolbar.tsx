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
  const scopeFilter = useAppStore(s => s.state.scopeFilter)
  const setScopeFilter = useAppStore(s => s.setScopeFilter)

  const handleSetScopeFilter = async () => {
    try {
      setScopeFilter(!scopeFilter)
    } catch (e: any) {
      LogError(`Project setup error: ${String(e?.message || e)}`)
    }
  }

  const handleStartAnalyzer = async () => {
    try {
      await StartAnalyzer()

      LogInfo(`Analyzer started.`)
    } catch (e: any) {
      LogError(`Analzyer starting error: ${String(e?.message || e)}`)
    }
  }

  return (
    <div className={styles.actionButtonContainer}>
      <button className={styles.actionButton} onClick={handleStartAnalyzer} title="Scan">
        <svg xmlns="http://www.w3.org/2000/svg" className={styles.playIcon} viewBox="0 0 16 16">
          <path d="M10.804 8 5 4.633v6.734zm.792-.696a.802.802 0 0 1 0 1.392l-6.363 3.692C4.713 12.69 4 12.345 4 11.692V4.308c0-.653.713-.998 1.233-.696z"/>
        </svg>
      </button>
      <button className={styles.actionButton} onClick={handleSetScopeFilter} title="Scope">
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
