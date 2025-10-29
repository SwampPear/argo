import { useAppStore } from '../../state/useAppStore'
import LogPage from './log/Page'
import styles from './Page.module.css'
import SettingsPage from './settings/Page'


const Page = () => {
  const { page } = useAppStore()

  const RoutePage = () => {
    switch (page) {
      case 'settings':
        return <SettingsPage/>
      case 'log':
        return <LogPage/>
      default:
        return <LogPage/>
    }
  }

  return (
    <div className={styles.container}>
      <RoutePage/>
    </div>
  )
}

export default Page