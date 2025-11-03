import { useAppStore } from '../../state/useAppStore'
import styles from './Page.module.css'
import BugsPage from './bugs/Page'
import LogPage from './log/Page'
import SettingsPage from './settings/Page'


const Page = () => {
  const { page } = useAppStore()

  const RoutePage = () => {
    switch (page) {
      case 'settings':
        return <SettingsPage/>
      case 'logs':
        return <LogPage/>
      case 'bugs':
        return <BugsPage/>
      default:
        return <SettingsPage/>
    }
  }

  return (
    <div className={styles.container}>
      <RoutePage/>
    </div>
  )
}

export default Page