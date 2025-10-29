import { create } from 'zustand'
import { settings } from '../../wailsjs/go/models'

type AppState = {
  projectDir: string
  settings: settings.Settings | null
  page: 'settings' | 'log'
  setProjectDir: (dir: string) => void
  setSettings: (s: settings.Settings | null) => void
  setPage: (page: 'settings' | 'log') => void
}

export const useAppStore = create<AppState>((set) => ({
  projectDir: '',
  settings: null,
  page: 'settings',
  setProjectDir: (projectDir) => set({ projectDir }),
  setSettings: (settings) => set({ settings }),
  setPage: (page) => set({ page })
}))
