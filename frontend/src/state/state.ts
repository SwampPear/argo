import { create } from 'zustand'
import * as backend from '../../wailsjs/go/app/App'
import { settings } from '../../wailsjs/go/models'
import { EventsOn, LogError } from '../../wailsjs/runtime/runtime'

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

export type AppPage = 'settings' | 'logs' | 'bugs'

export type RemoteState = {
  projectDir: string
  settings: settings.Settings
  logs: LogEntry[]
}

type Store = {
  state: RemoteState
  version: number
  page: AppPage
  setFromServer: (s: RemoteState) => void
  applyOptimistic: (draft: (s: RemoteState) => RemoteState) => Promise<void>
  setProjectDir: (dir: string) => Promise<void>
  setSettings: (cfg: settings.Settings) => Promise<void>
  addLog: (le: LogEntry) => Promise<void>
  selectProjectDirectory: () => Promise<string>
  loadYAMLSettings: (path: string) => Promise<settings.Settings>
  setPage: (page: AppPage) => void
}

const defaultRemote: RemoteState = {
  projectDir: '',
  settings: {} as settings.Settings,
  logs: []
}

/** Convert Wails wire AppState (class) -> plain RemoteState */
function toRemoteState(wire: any): RemoteState {
  if (!wire || typeof wire !== 'object') return defaultRemote
  // Wails respects json tags, but be tolerant of casing just in case
  const s = wire.settings ?? wire.Settings ?? {}
  const dir = wire.projectDir ?? wire.ProjectDir ?? ''
  const logs = wire.logs ?? wire.Logs ?? []
  return {
    projectDir: String(dir || ''),
    settings: s as settings.Settings,
    logs: Array.isArray(logs) ? logs as LogEntry[] : []
  }
}

/** Convert plain RemoteState -> wire shape Wails is happy to JSON-serialize */
function fromRemoteState(s: RemoteState): any {
  // json tags in Go are lowerCamelCase, so keep those keys
  return {
    projectDir: s.projectDir,
    settings: s.settings,
    logs: s.logs
  }
}

export const useAppStore = create<Store>((set, get) => ({
  state: defaultRemote,
  version: 0,
  page: 'settings',

  setFromServer: (s) => set({ state: s }),

  applyOptimistic: async (draft) => {
    const { state, version } = get()
    const next = draft(state)
    set({ state: next }) // optimistic

    // Call backend.SetState with a plain object, then normalize the response
    const resp = await backend.SetState(fromRemoteState(next) as any, version as any)
    const [wireState, serverVersion, ok] = resp as unknown as [any, number, boolean]

    if (!ok) {
      const serverState = toRemoteState(wireState)
      set({ state: serverState, version: serverVersion })
    } else {
      const serverState = toRemoteState(wireState)
      set({ state: serverState, version: serverVersion })
    }
  },

  setProjectDir: async (dir) =>
    get().applyOptimistic(s => ({ ...s, projectDir: dir })),

  setSettings: async (cfg) =>
    get().applyOptimistic(s => ({ ...s, settings: cfg })),

  addLog: async (le) =>
    get().applyOptimistic(s => ({ ...s, logs: [...s.logs, le] })),

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

  setPage: (page) => set({ page })
}))


export async function bootstrapState() {
  const api: any = backend as any
  if (typeof api.GetState !== 'function') {
    LogError('backend.GetState is not a function. Check your import: wailsjs/go/main/App')
    return
  }

  const res = await api.GetState()

  // New canonical shape: a single object with version embedded.
  let wireState: any | null = null
  if (res && typeof res === 'object' && 'version' in res) {
    wireState = res
  } else if (Array.isArray(res) || (res && typeof res === 'object' && '0' in res && '1' in res)) {
    // Back-compat: old tuple/envelope; hydrate but warn.
    LogError('Received legacy GetState() shape; update backend to embedded-version AppState.')
    const s = Array.isArray(res) ? (res as any)[0] : (res as any)[0]
    const v = Array.isArray(res) ? (res as any)[1] : (res as any)[1]
    wireState = { ...(s || {}), version: Number(v ?? 0) }
  } else {
    LogError(`Unexpected GetState() return shape: ${JSON.stringify(res)}`)
    return
  }

  useAppStore.getState().setFromServer(toRemoteState(wireState))

  // Live sync: backend now emits a single AppState payload on "state:update".
  EventsOn('state:update', (nextWire: any) => {
    if (!nextWire || typeof nextWire !== 'object' || !('version' in nextWire)) {
      LogError(`Ignoring unexpected state:update payload: ${JSON.stringify(nextWire)}`)
      return
    }
    useAppStore.getState().setFromServer(toRemoteState(nextWire))
  })
}
