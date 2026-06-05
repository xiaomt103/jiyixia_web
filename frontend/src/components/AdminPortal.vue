<template>
  <section class="data-page simple-page">
    <nav class="module-tabs"><button v-for="item in sections" :key="item.key" :class="{ active: section === item.key }" @click="section = item.key">{{ item.label }}</button></nav>

    <article class="panel login-panel">
      <form v-if="!token" class="inline" @submit.prevent="login"><input v-model="credentials.username" required placeholder="用户名" /><input v-model="credentials.password" required type="password" placeholder="密码" /><button class="primary" type="submit">登录后台</button></form>
      <div v-else class="notice-box"><b>已登录后台</b><span>可审核会员、管理套餐和导出数据。</span><button @click="loadMembers">刷新数据</button></div>
    </article>

    <article v-if="section === 'members'" class="report-card">
      <div class="report-title"><div><h2>会员审核</h2><small>单页只处理会员筛选、审核和导出。</small></div><div class="tools"><select v-model="statusFilter"><option value="all">全部状态</option><option v-for="status in statuses" :key="status" :value="status">{{ status }}</option></select><button @click="exportMembers">导出会员</button></div></div>
      <div class="metrics"><article v-for="item in metrics" :key="item.label" class="metric"><span>{{ item.label }}</span><strong>{{ item.value }}</strong></article></div>
      <div class="table-wrap"><table><thead><tr><th>会员</th><th>联系方式</th><th>套餐</th><th>状态</th><th>操作</th></tr></thead><tbody><tr v-for="member in filteredMembers" :key="member.id"><td><b>#{{ member.id }}</b> {{ member.name }}</td><td>{{ member.phone }}<small>{{ member.email }}</small></td><td>{{ member.plan_name }}</td><td><span class="status" :class="member.status">{{ member.status }}</span></td><td class="row-actions"><button v-for="status in statuses" :key="status" @click="setStatus(member.id, status)">{{ status }}</button></td></tr><tr v-if="filteredMembers.length === 0"><td class="empty" colspan="5">暂无会员数据。</td></tr></tbody></table></div>
    </article>

    <article v-else class="panel task-panel">
      <div class="section-head"><h3>套餐管理</h3><small>新增套餐并导出套餐清单。</small></div>
      <form @submit.prevent="createPlan"><label>名称<input v-model="plan.name" required placeholder="例如：年度会员" /></label><label>价格（分）<input v-model.number="plan.price_cents" required type="number" min="0" /></label><label>有效天数<input v-model.number="plan.duration_days" required type="number" min="1" /></label><label>说明<input v-model="plan.description" placeholder="权益摘要" /></label><button class="primary" :disabled="!token" type="submit">创建套餐</button><button type="button" @click="exportPlans">导出套餐</button></form>
      <div class="plans-grid"><article v-for="plan in props.plans" :key="plan.id" class="plan-item"><strong>{{ plan.name }}</strong><span>¥{{ price(plan) }}</span><small>{{ plan.duration_days }} 天 · {{ plan.description }}</small></article></div>
    </article>
  </section>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { api } from '../api'

const props = defineProps({ plans: { type: Array, default: () => [] } })
const emit = defineEmits(['refresh-plans', 'notice'])
const sections = [{ key: 'members', label: '会员审核' }, { key: 'plans', label: '套餐管理' }]
const section = ref('members')
const token = ref(localStorage.getItem('admin_token') || '')
const members = ref([])
const statusFilter = ref('all')
const statuses = ['pending', 'active', 'paused', 'cancelled']
const credentials = reactive({ username: 'admin', password: 'admin123' })
const plan = reactive({ name: '', price_cents: 0, duration_days: 30, description: '' })
const filteredMembers = computed(() => statusFilter.value === 'all' ? members.value : members.value.filter(item => item.status === statusFilter.value))
const metrics = computed(() => [{ label: '会员总数', value: members.value.length }, { label: '已激活', value: count('active') }, { label: '待审核', value: count('pending') }, { label: '套餐数', value: props.plans.length }])

if (token.value) loadMembers()
async function login() { const result = await api.login(credentials); token.value = result.token; localStorage.setItem('admin_token', result.token); emit('notice', '后台登录成功'); await loadMembers() }
async function loadMembers() { members.value = await api.adminMembers(token.value) }
async function setStatus(id, status) { await api.updateStatus(token.value, id, status); emit('notice', `会员 #${id} 状态已更新为 ${status}`); await loadMembers() }
async function createPlan() { await api.createPlan(token.value, plan); Object.assign(plan, { name: '', price_cents: 0, duration_days: 30, description: '' }); emit('refresh-plans'); emit('notice', '套餐已创建') }
function count(status) { return members.value.filter(item => item.status === status).length }
function price(plan) { return (plan.price_cents / 100).toFixed(2) }
function exportMembers() { exportCSV('members.csv', ['编号', '姓名', '手机号', '邮箱', '套餐', '状态'], filteredMembers.value.map(item => [item.id, item.name, item.phone, item.email, item.plan_name, item.status])) }
function exportPlans() { exportCSV('plans.csv', ['编号', '名称', '价格(元)', '有效天数', '说明'], props.plans.map(item => [item.id, item.name, price(item), item.duration_days, item.description])) }
function exportCSV(filename, headers, rows) { const csv = [headers, ...rows].map(row => row.map(csvCell).join(',')).join('\n'); const url = URL.createObjectURL(new Blob([`\uFEFF${csv}`], { type: 'text/csv;charset=utf-8;' })); const link = document.createElement('a'); link.href = url; link.download = filename; link.click(); URL.revokeObjectURL(url) }
function csvCell(value) { return `"${String(value ?? '').replaceAll('"', '""')}"` }
</script>
