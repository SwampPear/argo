import { useAppStore } from '@state/state'
import { LogError, LogInfo } from '@wails/runtime'
import styles from './Menu.module.css'

const ProjectSelect = () => {
  const projectDir = useAppStore(s => s.state.project_dir)

  const selectProjectDirectory = useAppStore(s => s.selectProjectDirectory)
  const loadYAMLSettings = useAppStore(s => s.loadYAMLSettings)

  const handleSetProject = async () => {
    try {
      const dir = await selectProjectDirectory()
      const path = `${dir}/scope.yaml`

      if (dir) await loadYAMLSettings(path)

      LogInfo(`Settings updated from ${path}.`)
    } catch (e: any) {
      LogError(`Project setup error: ${String(e?.message || e)}`)
    }
  }

  return (
    <div className={styles.projectContainer}>
      <svg xmlns="http://www.w3.org/2000/svg" className={styles.projectIcon} viewBox="0 0 16 16">
        <path d="M.54 3.87.5 3a2 2 0 0 1 2-2h3.672a2 2 0 0 1 1.414.586l.828.828A2 2 0 0 0 9.828 3h3.982a2 2 0 0 1 1.992 2.181l-.637 7A2 2 0 0 1 13.174 14H2.826a2 2 0 0 1-1.991-1.819l-.637-7a2 2 0 0 1 .342-1.31zM2.19 4a1 1 0 0 0-.996 1.09l.637 7a1 1 0 0 0 .995.91h10.348a1 1 0 0 0 .995-.91l.637-7A1 1 0 0 0 13.81 4zm4.69-1.707A1 1 0 0 0 6.172 2H2.5a1 1 0 0 0-1 .981l.006.139q.323-.119.684-.12h5.396z"/>
      </svg>
      <button className={styles.projectSelectButton} onClick={handleSetProject}>
        {projectDir || 'select project...'}
      </button>
    </div>
  )
}

export default ProjectSelect
