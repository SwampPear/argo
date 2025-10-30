import { LoadYAMLSettings, Run, SelectProjectDirectory, StartInteractiveBrowser } from '../../../../wailsjs/go/main/App'
import { LogError, LogInfo } from '../../../../wailsjs/runtime'
import { useAppStore } from '../../../state/useAppStore'
import styles from './Menu.module.css'

const Menu = () => {
  const { projectDir, setProjectDir, setSettings, setPage} = useAppStore()

  const handleSetProject = async () => {
    const dir = await SelectProjectDirectory()
    if (dir) {
      setProjectDir(dir)
    }

    try {
      const cfg = await LoadYAMLSettings(dir + '/scope.yaml')
      setSettings(cfg)
    } catch (e) {
      console.error('Load failed:', e)
      LogError(`LoadYAMLSettings error: ${String(e)}`)
    }
  }

  const handleRun = async () => {
    await Run();
  }

  const handleStartBrowser = async () => {
    try {
      const ts = Date.now()
      const artifactsBase = projectDir ? `${projectDir}/artifacts/${ts}` : `/tmp/argo/${ts}`

      await StartInteractiveBrowser('https://www.google.com')

      LogInfo(`Playwright started.`)
    } catch (e: any) {
      console.error('StartInteractiveBrowser failed:', e)
      LogError(`StartInteractiveBrowser error: ${String(e?.message || e)}`)
    }
  }

  /*
  const handleStopBrowser = async () => {
    try {
      //await StopInteractiveBrowser(currentRunId)
    } catch (e: any) {
      console.error('StopInteractiveBrowser failed: ', e)
      LogError('StopInteractiveBrowser error')
    }
  }
  */
  
  return (
    <div className={styles.container}>
      <div className={styles.projectContainer}>
        <svg xmlns="http://www.w3.org/2000/svg" className={styles.projectIcon} viewBox="0 0 16 16">
          <path d="M.54 3.87.5 3a2 2 0 0 1 2-2h3.672a2 2 0 0 1 1.414.586l.828.828A2 2 0 0 0 9.828 3h3.982a2 2 0 0 1 1.992 2.181l-.637 7A2 2 0 0 1 13.174 14H2.826a2 2 0 0 1-1.991-1.819l-.637-7a2 2 0 0 1 .342-1.31zM2.19 4a1 1 0 0 0-.996 1.09l.637 7a1 1 0 0 0 .995.91h10.348a1 1 0 0 0 .995-.91l.637-7A1 1 0 0 0 13.81 4zm4.69-1.707A1 1 0 0 0 6.172 2H2.5a1 1 0 0 0-1 .981l.006.139q.323-.119.684-.12h5.396z"/>
        </svg>
        <button className={styles.projectSelectButton} onClick={handleSetProject}>
          {projectDir ? projectDir : 'select project...'}
        </button>
      </div>
      <div className={styles.actionButtonContainer}>
        <button className={styles.actionButton} onClick={handleRun}>
          <svg xmlns="http://www.w3.org/2000/svg" className={styles.playIcon} viewBox="0 0 16 16">
            <path d="M10.804 8 5 4.633v6.734zm.792-.696a.802.802 0 0 1 0 1.392l-6.363 3.692C4.713 12.69 4 12.345 4 11.692V4.308c0-.653.713-.998 1.233-.696z"/>
          </svg>
        </button>
        <button className={styles.actionButton}>
          <svg xmlns="http://www.w3.org/2000/svg" className={styles.stopIcon} viewBox="0 0 16 16">
            <path d="M14 1a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1zM2 0a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2z"/>
          </svg>
        </button>
      </div>
      <div className={styles.flexGrow}></div>
      <div className={styles.pageSelectionContainer}>
        <button className={styles.pageSelectionButton} onClick={handleStartBrowser}>
          <svg xmlns="http://www.w3.org/2000/svg" className={styles.pageSelectionButtonIcon} viewBox="0 0 16 16">
            <path fill-rule="evenodd" d="M8 0a8 8 0 1 1 0 16A8 8 0 0 1 8 0m0 1a7 7 0 0 0-3.115.73.48.48 0 0 0 .137.292.488.488 0 0 1-.126.78l-.292.145a.7.7 0 0 0-.187.136l-.48.48a1 1 0 0 1-1.023.242l-.02-.007a1 1 0 0 0-.461-.041A6.97 6.97 0 0 0 1 8a6.96 6.96 0 0 0 .883 3.403l.86-.213c.444-.112.757-.512.757-.971v-.184a1 1 0 0 1 .445-.832l.04-.026a1 1 0 0 0 .153-1.54L3.12 6.622a.415.415 0 0 1 .542-.624l1.09.817a.5.5 0 0 0 .523.047A.5.5 0 0 1 6 7.31v.455a.8.8 0 0 0 .13.432l.796 1.193a1 1 0 0 1 .116.238l.73 2.19a1 1 0 0 0 .949.683h.058a1 1 0 0 0 .949-.684l.73-2.189q.042-.127.116-.238l.791-1.187A.45.45 0 0 1 11.743 8c.16 0 .306.083.392.218.557.875 1.63 2.282 2.365 2.282l.04-.003A7 7 0 0 0 8 1"/>
          </svg>
        </button>
        <button className={styles.pageSelectionButton} onClick={() => setPage('log')}>
          <svg xmlns="http://www.w3.org/2000/svg" className={styles.pageSelectionButtonIcon} viewBox="0 0 16 16">
            <path d="M2 2a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2zm2-1a1 1 0 0 0-1 1v4h10V2a1 1 0 0 0-1-1zm9 6H6v2h7zm0 3H6v2h7zm0 3H6v2h6a1 1 0 0 0 1-1zm-8 2v-2H3v1a1 1 0 0 0 1 1zm-2-3h2v-2H3zm0-3h2V7H3z"/>
          </svg>
        </button>
        <button className={styles.pageSelectionButton} onClick={() => setPage('settings')}>
          <svg xmlns="http://www.w3.org/2000/svg" className={styles.pageSelectionButtonIcon} viewBox="0 0 16 16">
            <path d="M8.932.727c-.243-.97-1.62-.97-1.864 0l-.071.286a.96.96 0 0 1-1.622.434l-.205-.211c-.695-.719-1.888-.03-1.613.931l.08.284a.96.96 0 0 1-1.186 1.187l-.284-.081c-.96-.275-1.65.918-.931 1.613l.211.205a.96.96 0 0 1-.434 1.622l-.286.071c-.97.243-.97 1.62 0 1.864l.286.071a.96.96 0 0 1 .434 1.622l-.211.205c-.719.695-.03 1.888.931 1.613l.284-.08a.96.96 0 0 1 1.187 1.187l-.081.283c-.275.96.918 1.65 1.613.931l.205-.211a.96.96 0 0 1 1.622.434l.071.286c.243.97 1.62.97 1.864 0l.071-.286a.96.96 0 0 1 1.622-.434l.205.211c.695.719 1.888.03 1.613-.931l-.08-.284a.96.96 0 0 1 1.187-1.187l.283.081c.96.275 1.65-.918.931-1.613l-.211-.205a.96.96 0 0 1 .434-1.622l.286-.071c.97-.243.97-1.62 0-1.864l-.286-.071a.96.96 0 0 1-.434-1.622l.211-.205c.719-.695.03-1.888-.931-1.613l-.284.08a.96.96 0 0 1-1.187-1.186l.081-.284c.275-.96-.918-1.65-1.613-.931l-.205.211a.96.96 0 0 1-1.622-.434zM8 12.997a4.998 4.998 0 1 1 0-9.995 4.998 4.998 0 0 1 0 9.996z"/>
          </svg>
        </button>
      </div>
    </div>
  )
}

export default Menu
