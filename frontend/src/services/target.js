import { $get, $post, $put, $delete } from '@/utils/request'

export function getTargetList(params) {
  return $get('/targets', params)
}

export function getTargetDetail(id) {
  return $get(`/targets/${id}`)
}

export function createTarget(data) {
  return $post('/targets', data)
}

export function updateTarget(id, data) {
  return $put(`/targets/${id}`, data)
}

export function deleteTarget(id) {
  return $delete(`/targets/${id}`)
}

export function testTarget(id) {
  return $post(`/targets/${id}/test`)
}

export function addTargetRepo(targetId, repoIds) {
  return $post(`/targets/${targetId}/repos`, { repo_ids: repoIds })
}

export function removeTargetRepo(targetId, repoId) {
  return $delete(`/targets/${targetId}/repos/${repoId}`)
}
