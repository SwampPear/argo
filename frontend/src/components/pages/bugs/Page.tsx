import styles from './Page.module.css'

const Content = () => {
  /*
  const settings = useAppStore(s => s.state.settings)

  const pretty = useMemo(() => {
    try {
      return JSON.stringify(settings ?? {}, null, 2)
    } catch {
      return 'Unable to render settings.'
    }
  }, [settings])
  */

  return (
    <div className={styles.container}>
      bugs
    </div>
  )
}

export default Content
