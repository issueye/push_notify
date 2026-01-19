import { $get, $post, $put, $delete } from '@/utils/request'

export function getModelList(params) {
  return $get('/models', params)
}

export function getModelDetail(id) {
  return $get(`/models/${id}`)
}

export function createModel(data) {
  return $post('/models', data)
}

export function updateModel(id, data) {
  return $put(`/models/${id}`, data)
}

export function deleteModel(id) {
  return $delete(`/models/${id}`)
}

export function setDefaultModel(id) {
  return $post(`/models/${id}/default`)
}

export function getModelLogs(id, params) {
  return $get(`/models/${id}/logs`, params)
}

export function verifyModel(id) {
  return $post(`/models/${id}/verify`)
}
