import { useAppStore } from '../../../state/useAppStore'
import styles from './Page.module.css'

const Content = () => {
  const { settings } = useAppStore()

  return (
    <div className={styles.container}>
      {JSON.stringify(settings, null, '\t')}
    </div>
  )
}

export default Content