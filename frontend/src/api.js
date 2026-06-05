import { getDemoMember, mergeDemoPlans, saveDemoMember } from './demoStore'

const jsonHeaders = { 'Content-Type': 'application/json' }
const apiBaseUrl = (import.meta.env.VITE_API_BASE_URL || '').replace(/\/$/, '')

export async function request(path, options = {}) {
  const headers = { ...jsonHeaders, ...(options.headers || {}) }
  const response = await fetch(`${apiBaseUrl}${path}`, { ...options, headers })
  const data = await response.json().catch(() => ({}))
  if (!response.ok) throw new Error(data.error || '请求失败')
  return data
}

export const api = {
  plans: async () => mergeDemoPlans(await request('/api/plans')),
  createMember: async (payload) => saveDemoMember(await request('/api/members', requestBody(payload))),
  getMember: async (id) => request(`/api/members/${id}`).catch(error => getDemoMember(id) || Promise.reject(error))
}

function requestBody(payload) {
  return { method: 'POST', body: JSON.stringify(payload) }
}
