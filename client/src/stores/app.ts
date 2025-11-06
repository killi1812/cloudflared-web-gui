// Utilities
import type { user } from '@/models/user'
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', () => {

  // Discord auth key
  const _authTokenState = ref("")
  const authToken = computed({
    get: () => _authTokenState.value,
    set: (value: string) => _authTokenState.value = value
  })

  // User
  const _userState = ref<user | undefined>(undefined)
  const user = computed({
    get: () => _userState.value,
    set: (value: user) => _userState.value = value
  })

  return {
    authToken,
    user,
  }
})
