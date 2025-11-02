import { create } from 'zustand'
import { settings } from '../../wailsjs/go/models'

export type LogEntry = {
  step: number
  id: string
  timestamp: string
  phase?: string
  module: string
  action: string
  target: string
  status: string
  duration: string
  confidence: number
  summary: string
  parent_step_id?: string | number
}

export type AppPage = 'settings' | 'log'

type AppState = {
  projectDir: string
  settings: settings.Settings | null
  page: AppPage
  logs: LogEntry[]
  setProjectDir: (dir: string) => void
  setSettings: (s: settings.Settings | null) => void
  setPage: (page: AppPage) => void
  addLog: (log: LogEntry) => void
  addLogs: (logs: LogEntry[]) => void
  updateLog: (id: string, patch: Partial<LogEntry>) => void
  removeLog: (id: string) => void
  clearLogs: () => void
  setLogs: (logs: LogEntry[]) => void
}

export const useAppStore = create<AppState>((set) => ({
  projectDir: '',
  settings: null,
  page: 'settings',
  logs: [],
  setProjectDir: (projectDir) => set({ projectDir }),
  setSettings: (settings) => set({ settings }),
  setPage: (page) => set({ page }),
  addLog: (log) =>
    set((state) => ({ logs: [...state.logs, log] })),
  addLogs: (newLogs) =>
    set((state) => ({ logs: [...state.logs, ...newLogs] })),
  updateLog: (id, patch) =>
    set((state) => ({
      logs: state.logs.map((l) => (l.id === id ? { ...l, ...patch } : l))
    })),
  removeLog: (id) =>
    set((state) => ({ logs: state.logs.filter((l) => l.id !== id) })),
  clearLogs: () => set({ logs: [] }),
  setLogs: (logs) => set({ logs })
}))