import { useAppStore } from '@state/state'
import styles from './Menu.module.css'

const EmptyToolbar = () => {
  return (
    <div className={styles.actionButtonContainer}></div>
  )
}

const LogsToolbar = () => {
  return (
    <div className={styles.actionButtonContainer}>
      <button className={styles.actionButton} title="Scan">
        <svg xmlns="http://www.w3.org/2000/svg" className={styles.playIcon} viewBox="0 0 16 16">
          <path d="M10.804 8 5 4.633v6.734zm.792-.696a.802.802 0 0 1 0 1.392l-6.363 3.692C4.713 12.69 4 12.345 4 11.692V4.308c0-.653.713-.998 1.233-.696z"/>
        </svg>
      </button>
    </div>
  )
}

const Toolbar = () => {
  const page = useAppStore(s => s.page)

  switch (page) {
  case 'bugs':
    return <LogsToolbar/>
  case 'logs':
    return <LogsToolbar/>
  case 'settings':
    return <EmptyToolbar/>
  default:
    return <EmptyToolbar/>
  }
}

export default Toolbar
