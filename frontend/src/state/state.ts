import * as backend from '@wails/go/app/App'
import { settings } from '@wails/go/models'
import { EventsOn, LogError } from '@wails/runtime/runtime'
import { create } from 'zustand'

// Log entry for tool calls, general info, and errors.
export type LogEntry = {
  step?: number
  id?: string
  timestamp?: string
  module?: string
  action?: string
  target?: string
  status?: string
  duration?: string
  confidence?: number
  summary?: string
  parent_step_id?: string
}

// Page routing.
export type AppPage = 'settings' | 'logs' | 'bugs'

// Shared remote state.
export type RemoteState = {
  project_dir: string
  settings: settings.Settings
  logs: LogEntry[]
  scope_filter: boolean
}

// Frontend store.
type Store = {
  state: RemoteState
  page: AppPage
  setFromServer: (s: RemoteState) => void
  applyOptimistic: (draft: (s: RemoteState) => RemoteState) => Promise<void>
  addLog: (le: LogEntry) => Promise<void>
  setScopeFilter: (sf: boolean) => Promise<void>
  selectProjectDirectory: () => Promise<string>
  loadYAMLSettings: (path: string) => Promise<settings.Settings>
  setPage: (page: AppPage) => void
}

const defaultRemote: RemoteState = {
  project_dir: '',
  settings: {} as settings.Settings,
  logs: [],
  scope_filter: false
}

// Convert Wails wire AppState (class) -> plain RemoteState.
const toRemoteState = (wire: any): RemoteState => {
  if (!wire || typeof wire !== 'object') return defaultRemote
  
  const s = wire.settings ?? wire.settings ?? {}
  const dir = wire.project_dir ?? wire.project_dir ?? ''
  const logs = wire.logs ?? wire.logs ?? []
  const sf = wire.scope_filter ?? wire.scope_filter ?? false

  return {
    project_dir: String(dir || ''),
    settings: s as settings.Settings,
    logs: Array.isArray(logs) ? logs as LogEntry[] : [],
    scope_filter: sf as boolean
  }
}

// Convert plain RemoteState to JSON-serializable shape.
const fromRemoteState =(s: RemoteState): any => {
  return {
    project_dir: s.project_dir,
    settings: s.settings,
    logs: s.logs,
    scope_filter: s.scope_filter
  }
}

// Store hook.
export const useAppStore = create<Store>((set, get) => ({
  state: defaultRemote,
  page: 'logs',
  setFromServer: (s) => set({ state: s }),
  applyOptimistic: async (draft) => {
    const { state } = get()
    const next = draft(state)
    set({ state: next })

    // call backend.SetState with a plain object, then normalize the response
    const resp = await backend.SetState(fromRemoteState(next) as any)
    const [wireState, ok] = resp as unknown as [any, number, boolean]

    if (!ok) {
      const serverState = toRemoteState(wireState)
      set({ state: serverState })
    } else {
      const serverState = toRemoteState(wireState)
      set({ state: serverState })
    }
  },

  addLog: async (le) => get().applyOptimistic(s => ({ ...s, logs: [...s.logs, le] })),
  setScopeFilter: async (sf) => get().applyOptimistic(s => ({ ...s, scope_filter: sf })),
  selectProjectDirectory: async () => {
    const dir = await backend.SelectProjectDirectory()
    // Backend will emit state:update; no local set needed.
    return dir as unknown as string
  },

  loadYAMLSettings: async (path: string) => {
    const cfg = await backend.LoadYAMLSettings(path)
    // Backend will emit state:update; no local set needed.
    return cfg as unknown as settings.Settings
  },

  setPage: (page) => set({ page }),
}))

export const bootstrapState = async () => {
  const res = await backend.GetState()

  // normalize to an object
  const normalize = (x: any) => {
    let obj: any = null
    if (x && typeof x === 'object' && !Array.isArray(x)) {
      obj = x
    } else if (Array.isArray(x) || (x && typeof x === 'object' && '0' in x)) {
      obj = Array.isArray(x) ? x[0] : x[0]
    } else {
      LogError(`Unexpected GetState() return shape: ${JSON.stringify(x)}`)
      return null
    }
    if (!obj || typeof obj !== 'object') return null

    return obj
  }

  const wireState = normalize(res)
  if (!wireState) return

  useAppStore.getState().setFromServer(toRemoteState(wireState))

  // state update
  EventsOn('state:update', (nextWire: any) => {
    const s = normalize(nextWire)
    if (!s) {
      LogError(`Ignoring unexpected state:update payload: ${JSON.stringify(nextWire)}`)
      return
    }
    
    useAppStore.getState().setFromServer(toRemoteState(s))
  })
}