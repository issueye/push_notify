import { $get, $post, $delete } from '@/utils/request'

export function getPushList(params) {
  return $get('/pushes', params)
}

export function getPushDetail(id) {
  return $get(`/pushes/${id}`)
}

export function retryPush(id) {
  return $post(`/pushes/${id}/retry`)
}

export function deletePush(id) {
  return $delete(`/pushes/${id}`)
}

export function batchRetry(pushIds) {
  return $post('/pushes/batch-retry', { push_ids: pushIds })
}

export function batchDelete(pushIds) {
  return $delete('/pushes/batch-delete', { push_ids: pushIds })
}

export function getPushStats(startDate, endDate) {
  return $get('/pushes/stats', { start_date: startDate, end_date: endDate })
}
