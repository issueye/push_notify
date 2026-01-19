import { $get, $post, $put, $delete } from '@/utils/request'

export function getUserList(params) {
  return $get('/users', params)
}

export function getUserDetail(id) {
  return $get(`/users/${id}`)
}

export function createUser(data) {
  return $post('/users', data)
}

export function updateUser(id, data) {
  return $put(`/users/${id}`, data)
}

export function deleteUser(id) {
  return $delete(`/users/${id}`)
}

export function resetUserPassword(id) {
  return $post(`/users/${id}/reset-password`)
}

export function lockUser(id, locked) {
  return $put(`/users/${id}/lock`, { locked })
}
