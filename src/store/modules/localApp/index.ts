import type { AppInfo, State } from './helper'
import { defineStore } from 'pinia'
import { defaultSetting } from './helper'

export const useLocalAppStore = defineStore('local-app-store', {
  state: (): State => defaultSetting(),
  actions: {
    setInstalledApps(apps: Record<string, AppInfo>) {
      this.installedApps = apps
    },

  },
})
