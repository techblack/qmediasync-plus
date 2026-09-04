import { defineStore } from 'pinia'
import { ref } from 'vue'
import { get, post } from '../services/http'
import type { User } from '../types'

export const useSessionStore = defineStore('session', () => {
  const user = ref<User | null>(null); const initialized = ref(false)
  window.addEventListener('qms:unauthorized', () => { user.value = null })
  async function restore() { try { user.value = await get<User>('/user/info') } catch { user.value = null } finally { initialized.value = true } }
  async function login(username:string,password:string,rememberMe:boolean) { const data = await post<{user:User}>('/login',{username,password,rememberMe}); user.value = data.user }
  async function logout() { try { await post('/logout',{}) } finally { user.value = null } }
  return { user, initialized, restore, login, logout }
})
