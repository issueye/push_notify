import { $get, $post, $put, $delete } from '@/utils/request'

export function login(data) {
  return $post('/auth/login', data)
}

export function register(data) {
  return $post('/auth/register', data)
}

export function getUserInfo() {
  return $get('/auth/me')
}

export function changePassword(data) {
  return $put('/auth/password', data)
}

export function refreshToken(refreshToken) {
  return $post('/auth/refresh', { refresh_token: refreshToken })
}

export function logout() {
  return $post('/auth/logout')
}
