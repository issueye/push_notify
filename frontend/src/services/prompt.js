import { $get, $post, $put, $delete } from '@/utils/request'

export function getPromptList(params) {
  return $get('/prompts', params)
}

export function getPromptDetail(id) {
  return $get(`/prompts/${id}`)
}

export function createPrompt(data) {
  return $post('/prompts', data)
}

export function updatePrompt(id, data) {
  return $put(`/prompts/${id}`, data)
}

export function deletePrompt(id) {
  return $delete(`/prompts/${id}`)
}

export function testPrompt(id, testData) {
  return $post(`/prompts/${id}/test`, { test_data: testData })
}

export function rollbackPrompt(id, version) {
  return $post(`/prompts/${id}/rollback`, { version })
}

export function exportPrompt(id) {
  return $get(`/prompts/${id}/export`)
}

export function importPrompt(file) {
  const formData = new FormData()
  formData.append('file', file)
  return $post('/prompts/import', formData)
}
