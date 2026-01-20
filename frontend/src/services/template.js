import { $get, $post, $put, $delete } from '@/utils/request'

export function getTemplateList(params) {
  return $get('/templates', params)
}

export function getTemplateDetail(id) {
  return $get(`/templates/${id}`)
}

export function createTemplate(data) {
  return $post('/templates', data)
}

export function updateTemplate(id, data) {
  return $put(`/templates/${id}`, data)
}

export function deleteTemplate(id) {
  return $delete(`/templates/${id}`)
}

export function updateTemplateStatus(id, status) {
  return $put(`/templates/${id}/status`, { status })
}

export function rollbackTemplate(id, version) {
  return $post(`/templates/${id}/rollback`, { version })
}

export function testTemplate(id, targetId) {
  return $post(`/templates/${id}/test`, { target_id: targetId })
}

export function generateTemplate(data) {
  return $post('/templates/generate', data)
}
