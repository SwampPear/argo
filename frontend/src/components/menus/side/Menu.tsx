import { useState } from 'react'
import styles from './Menu.module.css'

const Menu = () => {
  const [width, setWidth] = useState<number>(256)

  return (
    <div className={styles.container} style={{width: `${width}px`}}>
      asdf
    </div>
  )
}

export default Menu
