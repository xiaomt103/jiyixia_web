import {
  getDemoMember,
  mergeDemoMembers,
  mergeDemoPlans,
  saveDemoMember,
  saveDemoPlan,
  updateDemoMemberStatus
} from './demoStore'

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
  getMember: async (id) => request(`/api/members/${id}`).catch(error => getDemoMember(id) || Promise.reject(error)),
  login: (payload) => request('/api/admin/login', requestBody(payload)),
  adminMembers: async (token) => mergeDemoMembers(await request('/api/v1/admin/members', { headers: authHeader(token) })),
  updateStatus: async (token, id, status) => {
    try {
      const updated = await request(`/api/v1/admin/members/${id}/status`, {
        method: 'PATCH', headers: authHeader(token), body: JSON.stringify({ status })
      })
      return saveDemoMember(updated)
    } catch (error) {
      return updateDemoMemberStatus(id, status) || Promise.reject(error)
    }
  },
  adminPlans: async (token) => mergeDemoPlans(await request('/api/v1/admin/plans', { headers: authHeader(token) })),
  createPlan: async (token, payload) => saveDemoPlan(await request('/api/v1/admin/plans', {
    method: 'POST', headers: authHeader(token), body: JSON.stringify(payload)
  }))
}

function requestBody(payload) {
  return { method: 'POST', body: JSON.stringify(payload) }
}

function authHeader(token) {
  return { Authorization: `Bearer ${token}` }
}
