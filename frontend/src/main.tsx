import React from 'react'
import { createRoot } from 'react-dom/client'
import { LogError } from '../wailsjs/runtime/runtime'
import App from './App'
import { bootstrapState } from './state/state'
import './style.css'

const container = document.getElementById('root')

const root = createRoot(container!)

bootstrapState().catch(err => LogError(err))

root.render(
  <React.StrictMode>
    <App/>
  </React.StrictMode>
)