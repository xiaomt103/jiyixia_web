const memberKey = 'vercel_demo_members'
const planKey = 'vercel_demo_plans'

export function saveDemoMember(member) {
  if (!canUseStorage() || !member?.id) return member
  const members = listDemoMembers()
  const saved = avoidIDConflict(member, members)
  write(memberKey, [...members.filter(item => item.id !== saved.id), saved])
  return saved
}

export function listDemoMembers() {
  return read(memberKey, [])
}

export function getDemoMember(id) {
  return listDemoMembers().find(member => String(member.id) === String(id))
}

export function mergeDemoMembers(remoteMembers) {
  const merged = new Map()
  for (const member of remoteMembers || []) merged.set(String(member.id), member)
  for (const member of listDemoMembers()) merged.set(String(member.id), member)
  return Array.from(merged.values()).sort((left, right) => left.id - right.id)
}

export function updateDemoMemberStatus(id, status) {
  const member = getDemoMember(id)
  if (!member) return null
  const updated = { ...member, status, updated_at: new Date().toISOString() }
  if (status === 'active') setActiveDates(updated)
  saveDemoMember(updated)
  return updated
}

export function saveDemoPlan(plan) {
  if (!canUseStorage() || !plan?.id) return plan
  const plans = listDemoPlans().filter(item => item.id !== plan.id)
  plans.push(plan)
  write(planKey, plans)
  return plan
}

export function mergeDemoPlans(remotePlans) {
  const merged = new Map()
  for (const plan of remotePlans || []) merged.set(String(plan.id), plan)
  for (const plan of listDemoPlans()) merged.set(String(plan.id), plan)
  return Array.from(merged.values()).sort((left, right) => left.id - right.id)
}

function listDemoPlans() {
  return read(planKey, [])
}

function avoidIDConflict(member, members) {
  const existing = members.find(item => item.id === member.id)
  if (!existing || sameMember(existing, member)) return member
  const nextID = Math.max(0, ...members.map(item => Number(item.id) || 0)) + 1
  return { ...member, id: nextID }
}

function sameMember(left, right) {
  return left.name === right.name && left.phone === right.phone && left.email === right.email
}

function setActiveDates(member) {
  const start = new Date()
  const expire = new Date(start)
  expire.setDate(expire.getDate() + 30)
  member.start_at = start.toISOString()
  member.expire_at = expire.toISOString()
}

function read(key, fallback) {
  if (!canUseStorage()) return fallback
  try {
    return JSON.parse(localStorage.getItem(key) || JSON.stringify(fallback))
  } catch (error) {
    return fallback
  }
}

function write(key, value) {
  try {
    localStorage.setItem(key, JSON.stringify(value))
  } catch (error) {
    return undefined
  }
}

function canUseStorage() {
  return typeof localStorage !== 'undefined'
}
