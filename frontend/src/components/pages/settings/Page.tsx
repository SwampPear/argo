import { useMemo } from 'react'
import { useAppStore } from '../../../state/useAppStore'
import styles from './Page.module.css'

const Content = () => {
  const settings = useAppStore(s => s.state.settings)

  const pretty = useMemo(() => {
    try {
      return JSON.stringify(settings ?? {}, null, 2)
    } catch {
      return 'Unable to render settings.'
    }
  }, [settings])

  return (
    <div className={styles.container}>
      <pre className={styles.code}>{pretty}</pre>
    </div>
  )
}

export default Content
