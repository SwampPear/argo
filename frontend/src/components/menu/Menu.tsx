import styles from './Menu.module.css'
import PageSelect from './PageSelect'
import ProjectSelect from './ProjectSelect'
import Toolbar from './Toolbar'

const Menu = () => {
  return (
    <div className={styles.container}>
      <ProjectSelect />
      <Toolbar />
      <div className={styles.flexGrow} />
      <PageSelect />
    </div>
  )
}

export default Menu
