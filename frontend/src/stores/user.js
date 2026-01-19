import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { $post, $get, $put } from '@/utils/request'
import router from '@/router'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const userInfo = ref(null)
  const roles = ref([])

  const isLoggedIn = computed(() => !!token.value)
  const isAdmin = computed(() => roles.value.includes('admin'))

  async function login(params) {
    const data = await $post('/auth/login', params)
    token.value = data.access_token
    localStorage.setItem('token', data.access_token)
    await getUserProfile()
    return data
  }

  async function register(params) {
    return $post('/auth/register', params)
  }

  async function getUserProfile() {
    const info = await $get('/auth/me')
    userInfo.value = info
    roles.value = [info.role]
    return info
  }

  function logout() {
    token.value = null
    userInfo.value = null
    roles.value = []
    localStorage.removeItem('token')
    router.push('/login')
  }

  async function changePassword(oldPassword, newPassword) {
    return $put('/auth/password', { old_password: oldPassword, new_password: newPassword })
  }

  return {
    token,
    userInfo,
    roles,
    isLoggedIn,
    isAdmin,
    login,
    register,
    getUserProfile,
    logout,
    changePassword
  }
}, {
  persist: true
})
