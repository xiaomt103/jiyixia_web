const hopByHopHeaders = new Set(['connection', 'content-length', 'host'])
const demoToken = 'vercel-demo-admin-token'

const state = globalThis.__memberDemoState || seedState()
globalThis.__memberDemoState = state

module.exports = async (req, res) => {
  const baseUrl = normalizedBaseUrl()
  if (!baseUrl || isSelfProxy(baseUrl, req)) {
    await handleDemo(req, res)
    return
  }
  await proxyRequest(req, res, baseUrl)
}

async function proxyRequest(req, res, baseUrl) {
  try {
    const response = await fetch(`${baseUrl}${req.url}`, {
      method: req.method,
      headers: requestHeaders(req.headers),
      body: hasBody(req.method) ? await readBody(req) : undefined
    })
    res.statusCode = response.status
    response.headers.forEach((value, key) => {
      if (!hopByHopHeaders.has(key.toLowerCase())) res.setHeader(key, value)
    })
    res.end(Buffer.from(await response.arrayBuffer()))
  } catch (error) {
    sendJSON(res, 502, { error: 'Failed to reach API service' })
  }
}

async function handleDemo(req, res) {
  const url = new URL(req.url, `https://${req.headers.host || 'localhost'}`)
  const path = url.pathname
  if (req.method === 'GET' && path === '/api/health') {
    sendJSON(res, 200, { status: 'ok', mode: 'vercel-demo' })
    return
  }
  if (req.method === 'GET' && path === '/api/plans') {
    sendJSON(res, 200, state.plans)
    return
  }
  if (req.method === 'POST' && path === '/api/members') {
    const input = await parseJSON(req)
    const plan = state.plans.find(item => item.id === Number(input.plan_id))
    if (!plan || !input.name || !input.phone || !input.email) {
      sendJSON(res, 400, { error: 'invalid member input' })
      return
    }
    const member = createMember(input, plan)
    state.members.push(member)
    sendJSON(res, 201, member)
    return
  }
  if (req.method === 'GET' && path.startsWith('/api/members/')) {
    const member = findByID(state.members, path.replace('/api/members/', ''))
    sendJSON(res, member ? 200 : 404, member || { error: 'member not found' })
    return
  }
  if (req.method === 'POST' && path === '/api/admin/login') {
    const input = await parseJSON(req)
    const username = process.env.ADMIN_USERNAME || 'admin'
    const password = process.env.ADMIN_PASSWORD || 'admin123'
    if (input.username !== username || input.password !== password) {
      sendJSON(res, 401, { error: 'invalid credentials' })
      return
    }
    sendJSON(res, 200, { token: demoToken })
    return
  }
  if (!isAuthorized(req)) {
    sendJSON(res, 401, { error: 'admin token required' })
    return
  }
  if (req.method === 'GET' && path === '/api/v1/admin/members') {
    sendJSON(res, 200, state.members)
    return
  }
  if (req.method === 'PATCH' && path.startsWith('/api/v1/admin/members/') && path.endsWith('/status')) {
    await updateDemoStatus(req, res, path)
    return
  }
  if (req.method === 'GET' && path === '/api/v1/admin/plans') {
    sendJSON(res, 200, state.plans)
    return
  }
  if (req.method === 'GET' && path === '/api/v1/admin/operations') {
    const module = url.searchParams.get('module')
    sendJSON(res, 200, module ? state.operations.filter(item => item.module === module) : state.operations)
    return
  }
  if (req.method === 'POST' && path === '/api/v1/admin/operations') {
    const input = await parseJSON(req)
    if (!input.module || !input.title) {
      sendJSON(res, 400, { error: 'invalid operation input' })
      return
    }
    const record = createOperation(input)
    state.operations.unshift(record)
    sendJSON(res, 201, record)
    return
  }
  if (req.method === 'PATCH' && path.startsWith('/api/v1/admin/operations/') && path.endsWith('/status')) {
    const id = path.replace('/api/v1/admin/operations/', '').replace('/status', '')
    const record = findByID(state.operations, id)
    const input = await parseJSON(req)
    if (!record || !input.status) {
      sendJSON(res, record ? 400 : 404, { error: record ? 'invalid status' : 'operation not found' })
      return
    }
    record.status = input.status
    sendJSON(res, 200, record)
    return
  }
  if (req.method === 'POST' && path === '/api/v1/admin/plans') {
    const input = await parseJSON(req)
    if (!input.name || Number(input.duration_days) <= 0 || Number(input.price_cents) < 0) {
      sendJSON(res, 400, { error: 'invalid plan input' })
      return
    }
    const plan = createPlan(input)
    state.plans.push(plan)
    sendJSON(res, 201, plan)
    return
  }
  sendJSON(res, 404, { error: 'not found' })
}

async function updateDemoStatus(req, res, path) {
  const id = path.replace('/api/v1/admin/members/', '').replace('/status', '')
  const member = findByID(state.members, id)
  const input = await parseJSON(req)
  if (!member) {
    sendJSON(res, 404, { error: 'member not found' })
    return
  }
  if (!['pending', 'active', 'paused', 'cancelled'].includes(input.status)) {
    sendJSON(res, 400, { error: 'invalid status' })
    return
  }
  member.status = input.status
  member.updated_at = new Date().toISOString()
  if (input.status === 'active') activateMember(member)
  sendJSON(res, 200, member)
}

function normalizedBaseUrl() {
  return (process.env.API_BASE_URL || '').replace(/\/$/, '')
}

function isSelfProxy(baseUrl, req) {
  try {
    const target = new URL(baseUrl)
    return target.host === req.headers.host
  } catch (error) {
    return false
  }
}

function requestHeaders(headers) {
  return Object.fromEntries(Object.entries(headers).filter(([key]) => !hopByHopHeaders.has(key.toLowerCase())))
}

function hasBody(method) {
  return !['GET', 'HEAD'].includes(method)
}

function readBody(req) {
  return new Promise((resolve, reject) => {
    const chunks = []
    req.on('data', chunk => chunks.push(chunk))
    req.on('end', () => resolve(Buffer.concat(chunks)))
    req.on('error', reject)
  })
}

async function parseJSON(req) {
  try {
    const body = await readBody(req)
    return body.length ? JSON.parse(body.toString('utf8')) : {}
  } catch (error) {
    return {}
  }
}

function sendJSON(res, status, payload) {
  res.statusCode = status
  res.setHeader('Content-Type', 'application/json')
  res.end(JSON.stringify(payload))
}

function seedState() {
  const createdAt = new Date().toISOString()
  return {
    planSeq: 3,
    memberSeq: 1,
    operationSeq: 6,
    plans: [
      plan(1, '月度会员', 9900, 30, '适合短期体验', createdAt),
      plan(2, '季度会员', 26900, 90, '适合稳定使用', createdAt),
      plan(3, '年度会员', 99900, 365, '全年权益优惠', createdAt)
    ],
    members: [],
    operations: [
      operation(1, 'booking', '明天 10:30 到店预约', 0, 0, 'pending', createdAt),
      operation(2, 'order', '年度会员订单待确认', 0, 99900, 'pending', createdAt),
      operation(3, 'marketing', '六月续费提醒活动', 0, 0, 'active', createdAt),
      operation(4, 'finance', '会员收入日报', 0, 101170, 'done', createdAt),
      operation(5, 'report', '经营分析月报', 0, 0, 'ready', createdAt)
    ]
  }
}

function plan(id, name, price, days, description, createdAt) {
  return { id, name, price_cents: price, duration_days: days, description, created_at: createdAt }
}

function operation(id, module, title, memberID, amount, status, dueAt) {
  return { id, module, title, member_id: memberID, amount_cents: amount, status, due_at: dueAt, created_at: dueAt }
}

function createOperation(input) {
  return operation(++state.operationSeq, input.module, input.title.trim(), Number(input.member_id || 0), Number(input.amount_cents || 0), input.status || 'pending', input.due_at || new Date().toISOString())
}

function createPlan(input) {
  const id = ++state.planSeq
  return plan(id, input.name.trim(), Number(input.price_cents), Number(input.duration_days), input.description || '', new Date().toISOString())
}

function createMember(input, selectedPlan) {
  const now = new Date().toISOString()
  return {
    id: state.memberSeq++,
    name: input.name.trim(),
    phone: input.phone.trim(),
    email: input.email.trim(),
    plan_id: selectedPlan.id,
    plan_name: selectedPlan.name,
    status: 'pending',
    created_at: now,
    updated_at: now
  }
}

function activateMember(member) {
  const start = new Date()
  const selectedPlan = findByID(state.plans, member.plan_id)
  const expire = new Date(start)
  expire.setDate(expire.getDate() + (selectedPlan ? selectedPlan.duration_days : 30))
  member.start_at = start.toISOString()
  member.expire_at = expire.toISOString()
}

function findByID(items, id) {
  const numericID = Number(id)
  return items.find(item => item.id === numericID)
}

function isAuthorized(req) {
  return req.headers.authorization === `Bearer ${demoToken}`
}
