import { StartInteractiveBrowser } from '@go/app/App'
import { useAppStore } from '@state/state'
import { LogError, LogInfo } from '@wails/runtime'
import styles from './Menu.module.css'

const PageSelect = () => {
  const setPage = useAppStore(s => s.setPage)

  const handleStartBrowser = async () => {
    try {
      await StartInteractiveBrowser()

      LogInfo(`Browser instance started.`)
    } catch (e: any) {
      LogError(`Browser starting error: ${String(e?.message || e)}`)
    }
  }

  return (
    <div className={styles.pageSelectionContainer}>
      <button className={styles.pageSelectionButton} onClick={handleStartBrowser} title="Browser">
        <svg xmlns="http://www.w3.org/2000/svg" className={styles.pageSelectionButtonIcon} viewBox="0 0 16 16">
          <path fillRule="evenodd" d="M8 0a8 8 0 1 1 0 16A8 8 0 0 1 8 0m0 1a7 7 0 0 0-3.115.73.48.48 0 0 0 .137.292.488.488 0 0 1-.126.78l-.292.145a.7.7 0 0 0-.187.136l-.48.48a1 1 0 0 1-1.023.242l-.02-.007a1 1 0 0 0-.461-.041A6.97 6.97 0 0 0 1 8a6.96 6.96 0 0 0 .883 3.403l.86-.213c.444-.112.757-.512.757-.971v-.184a1 1 0 0 1 .445-.832l.04-.026a1 1 0 0 0 .153-1.54L3.12 6.622a.415.415 0 0 1 .542-.624l1.09.817a.5.5 0 0 0 .523.047A.5.5 0 0 1 6 7.31v.455a.8.8 0 0 0 .13.432l.796 1.193a1 1 0 0 1 .116.238l.73 2.19a1 1 0 0 0 .949.683h.058a1 1 0 0 0 .949-.684l.73-2.189q.042-.127.116-.238l.791-1.187A.45.45 0 0 1 11.743 8c.16 0 .306.083.392.218.557.875 1.63 2.282 2.365 2.282l.04-.003A7 7 0 0 0 8 1"/>
        </svg>
      </button>
      <button className={styles.pageSelectionButton} onClick={() => setPage('bugs')} title="Bugs">
        <svg xmlns="http://www.w3.org/2000/svg" className={styles.pageSelectionButtonIcon} viewBox="0 0 16 16">
          <path d="M4.355.522a.5.5 0 0 1 .623.333l.291.956A5 5 0 0 1 8 1c1.007 0 1.946.298 2.731.811l.29-.956a.5.5 0 1 1 .957.29l-.41 1.352A5 5 0 0 1 13 6h.5a.5.5 0 0 0 .5-.5V5a.5.5 0 0 1 1 0v.5A1.5 1.5 0 0 1 13.5 7H13v1h1.5a.5.5 0 0 1 0 1H13v1h.5a1.5 1.5 0 0 1 1.5 1.5v.5a.5.5 0 1 1-1 0v-.5a.5.5 0 0 0-.5-.5H13a5 5 0 0 1-10 0h-.5a.5.5 0 0 0-.5.5v.5a.5.5 0 1 1-1 0v-.5A1.5 1.5 0 0 1 2.5 10H3V9H1.5a.5.5 0 0 1 0-1H3V7h-.5A1.5 1.5 0 0 1 1 5.5V5a.5.5 0 0 1 1 0v.5a.5.5 0 0 0 .5.5H3c0-1.364.547-2.601 1.432-3.503l-.41-1.352a.5.5 0 0 1 .333-.623M4 7v4a4 4 0 0 0 3.5 3.97V7zm4.5 0v7.97A4 4 0 0 0 12 11V7zM12 6a4 4 0 0 0-1.334-2.982A3.98 3.98 0 0 0 8 2a3.98 3.98 0 0 0-2.667 1.018A4 4 0 0 0 4 6z"/>
        </svg>
      </button>
      <button className={styles.pageSelectionButton} onClick={() => setPage('logs')} title="Logs">
        <svg xmlns="http://www.w3.org/2000/svg" className={styles.pageSelectionButtonIcon} viewBox="0 0 16 16">
          <path d="M2 2a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2zm2-1a1 1 0 0 0-1 1v4h10V2a1 1 0 0 0-1-1zm9 6H6v2h7zm0 3H6v2h7zm0 3H6v2h6a1 1 0 0 0 1-1zm-8 2v-2H3v1a1 1 0 0 0 1 1zm-2-3h2v-2H3zm0-3h2V7H3z"/>
        </svg>
      </button>
      <button className={styles.pageSelectionButton} onClick={() => setPage('settings')} title="Settings">
        <svg xmlns="http://www.w3.org/2000/svg" className={styles.pageSelectionButtonIcon} viewBox="0 0 16 16">
          <path d="M8.932.727c-.243-.97-1.62-.97-1.864 0l-.071.286a.96.96 0 0 1-1.622.434l-.205-.211c-.695-.719-1.888-.03-1.613.931l.08.284a.96.96 0 0 1-1.186 1.187l-.284-.081c-.96-.275-1.65.918-.931 1.613l.211.205a.96.96 0 0 1-.434 1.622l-.286.071c-.97.243-.97 1.62 0 1.864l.286.071a.96.96 0 0 1 .434 1.622l-.211.205c-.719.695-.03 1.888.931 1.613l.284-.08a.96.96 0 0 1 1.187 1.187l-.081.283c-.275.96.918 1.65 1.613.931l.205-.211a.96.96 0 0 1 1.622.434l.071.286c.243.97 1.62.97 1.864 0l.071-.286a.96.96 0 0 1 1.622-.434l.205.211c.695.719 1.888.03 1.613-.931l-.08-.284a.96.96 0 0 1 1.187-1.187l.283.081c.96.275 1.65-.918.931-1.613l-.211-.205a.96.96 0 0 1 .434-1.622l.286-.071c.97-.243.97-1.62 0-1.864l-.286-.071a.96.96 0 0 1-.434-1.622l.211-.205c.719-.695.03-1.888-.931-1.613l-.284.08a.96.96 0 0 1-1.187-1.186l.081-.284c.275-.96-.918-1.65-1.613-.931l-.205.211a.96.96 0 0 1-1.622-.434zM8 12.997a4.998 4.998 0 1 1 0-9.995 4.998 4.998 0 0 1 0 9.996z"/>
        </svg>
      </button>
    </div>
  )
}

export default PageSelect
