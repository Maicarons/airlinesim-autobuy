import type { Status, AircraftData } from './types'

const BASE = '/api'

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${url}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error || res.statusText)
  }
  return res.json()
}

export function getStatus(): Promise<Status> {
  return request<Status>('/status')
}

export function startEngine(): Promise<{ status: string }> {
  return request('/control/start', { method: 'POST' })
}

export function stopEngine(): Promise<{ status: string }> {
  return request('/control/stop', { method: 'POST' })
}

export function getConfig(): Promise<any> {
  return request('/config')
}

export function updateConfig(config: any): Promise<any> {
  return request('/config', {
    method: 'PUT',
    body: JSON.stringify(config),
  })
}

export function getRules(): Promise<any[]> {
  return request('/rules')
}

export function createRule(rule: any): Promise<any> {
  return request('/rules', {
    method: 'POST',
    body: JSON.stringify(rule),
  })
}

export function updateRule(id: number, rule: any): Promise<any> {
  return request(`/rules/${id}`, {
    method: 'PUT',
    body: JSON.stringify(rule),
  })
}

export function deleteRule(id: number): Promise<any> {
  return request(`/rules/${id}`, { method: 'DELETE' })
}

export function toggleRule(id: number): Promise<any> {
  return request(`/rules/${id}/toggle`, { method: 'PATCH' })
}

export function reorderRules(order: number[]): Promise<any> {
  return request('/rules/reorder', {
    method: 'PUT',
    body: JSON.stringify(order),
  })
}

export function getAircraftData(): Promise<AircraftData> {
  return request<AircraftData>('/aircraft-data')
}