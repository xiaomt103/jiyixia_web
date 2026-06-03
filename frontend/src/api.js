const jsonHeaders = { 'Content-Type': 'application/json' }

export async function request(path, options = {}) {
  const headers = { ...jsonHeaders, ...(options.headers || {}) }
  const response = await fetch(path, { ...options, headers })
  const data = await response.json().catch(() => ({}))
  if (!response.ok) throw new Error(data.error || '请求失败')
  return data
}

export const api = {
  plans: () => request('/api/plans'),
  createMember: (payload) => request('/api/members', { method: 'POST', body: JSON.stringify(payload) }),
  getMember: (id) => request(`/api/members/${id}`),
  login: (payload) => request('/api/admin/login', { method: 'POST', body: JSON.stringify(payload) }),
  adminMembers: (token) => request('/api/v1/admin/members', { headers: authHeader(token) }),
  updateStatus: (token, id, status) => request(`/api/v1/admin/members/${id}/status`, {
    method: 'PATCH', headers: authHeader(token), body: JSON.stringify({ status })
  }),
  adminPlans: (token) => request('/api/v1/admin/plans', { headers: authHeader(token) }),
  createPlan: (token, payload) => request('/api/v1/admin/plans', {
    method: 'POST', headers: authHeader(token), body: JSON.stringify(payload)
  })
}

function authHeader(token) {
  return { Authorization: `Bearer ${token}` }
}
