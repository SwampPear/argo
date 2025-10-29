import React, { useState } from 'react'
import styles from './Menu.module.css'

interface IButtonProps {
  children: React.ReactNode
}

const Button = ({ children }: IButtonProps) => {
  return (
    <button className={styles.button}>
      {children}
    </button>
  )
}

const Menu = () => {
  const [width, setWidth] = useState<number>(256)

  return (
    <div className={styles.container} style={{width: `${width}px`}}>
      <Button>Settings</Button>
      <Button>Log</Button>
    </div>
  )
}

export default Menu
