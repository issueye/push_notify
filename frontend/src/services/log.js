import { $get } from '@/utils/request'

export function getSystemLogs(params) {
  return $get('/logs/system', params)
}

export function getOperationLogs(params) {
  return $get('/logs/operations', params)
}

export function getAICallLogs(params) {
  return $get('/logs/ai-calls', params)
}

export function searchLogs(keyword, params) {
  return $get('/logs/search', { keyword, ...params })
}

export function exportLogs(type, startTime, endTime, format = 'csv') {
  return $get('/logs/export', { type, start_time: startTime, end_time: endTime, format })
}

export function getLogStats(startDate, endDate) {
  return $get('/logs/stats', { start_date: startDate, end_date: endDate })
}
