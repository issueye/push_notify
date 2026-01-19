import { $get, $post, $put, $delete } from '@/utils/request'

export function getRepoList(params) {
  return $get('/repos', params)
}

export function getRepoDetail(id) {
  return $get(`/repos/${id}`)
}

export function createRepo(data) {
  return $post('/repos', data)
}

export function updateRepo(id, data) {
  return $put(`/repos/${id}`, data)
}

export function deleteRepo(id) {
  return $delete(`/repos/${id}`)
}

export function testWebhook(id) {
  return $post(`/repos/${id}/test`)
}

export function getRepoTargets(id) {
  return $get(`/repos/${id}/targets`)
}

export function addRepoTarget(repoId, targetId) {
  return $post(`/repos/${repoId}/targets`, { target_id: targetId })
}

export function removeRepoTarget(repoId, targetId) {
  return $delete(`/repos/${repoId}/targets/${targetId}`)
}
