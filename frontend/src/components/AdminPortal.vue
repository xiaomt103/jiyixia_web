<template>
  <section class="data-page">
    <nav class="module-tabs"><button class="active">会员中台</button><button>审核队列</button><button>套餐运营</button><button>导出中心</button></nav>
    <div class="filter-bar"><span>今天</span><span>昨天</span><span>上周</span><b>本月</b><input value="2026-06-01 00:00:00  ~  2026-06-30 23:59:59" readonly /><button class="primary" @click="loadMembers">查询</button></div>
    <div class="metrics"><article v-for="item in metrics" :key="item.label" class="metric"><span>{{ item.label }}</span><strong>{{ item.value }}</strong></article></div>

    <section class="pipeline-grid">
      <button v-for="item in pipeline" :key="item.status" class="pipeline-card" @click="statusFilter = item.status"><b>{{ item.label }}</b><strong>{{ item.count }}</strong><small>{{ item.hint }}</small></button>
    </section>

    <section class="admin-grid">
      <section class="report-layout">
        <aside class="sub-menu"><button class="active">会员列表</button><button>状态分析</button><button>套餐分析</button><button>导出记录</button></aside>
        <article class="report-card">
          <div class="report-title">
            <div><h2>会员数据报表</h2><small>筛选、审核、导出，后续可扩展更多报表页签。</small></div>
            <div class="tools"><select v-model="statusFilter"><option value="all">全部状态</option><option v-for="status in statuses" :key="status" :value="status">{{ status }}</option></select><button @click="exportMembers">导出会员</button><button @click="exportPlans">导出套餐</button></div>
          </div>
          <div class="table-wrap"><table><thead><tr><th>会员</th><th>联系方式</th><th>套餐</th><th>状态</th><th>创建时间</th><th>审核动作</th></tr></thead><tbody>
            <tr v-for="member in filteredMembers" :key="member.id"><td><b>#{{ member.id }}</b> {{ member.name }}</td><td>{{ member.phone }}<small>{{ member.email }}</small></td><td>{{ member.plan_name }}</td><td><span class="status" :class="member.status">{{ member.status }}</span></td><td>{{ formatDate(member.created_at) }}</td><td class="row-actions"><button v-for="status in statuses" :key="status" @click="setStatus(member.id, status)">{{ status }}</button></td></tr>
            <tr v-if="filteredMembers.length === 0"><td class="empty" colspan="6">暂无会员数据，请先在会员录入页创建申请。</td></tr>
          </tbody></table></div>
        </article>
      </section>
      <aside class="admin-aside">
        <article class="panel queue-panel"><div class="section-head"><h3>待审核节点</h3><small>{{ pendingMembers.length }} 个申请待处理</small></div><p v-for="member in pendingMembers.slice(0, 4)" :key="member.id"><b>{{ member.name }}</b><span>{{ member.plan_name || '未分配套餐' }}</span><button @click="setStatus(member.id, 'active')">激活</button></p><div v-if="pendingMembers.length === 0" class="empty-note">当前没有待审核申请。</div></article>
        <article class="panel plan-panel"><div class="section-head"><h3>套餐运营</h3><small>价格、周期、说明统一管理。</small></div><p v-for="plan in props.plans" :key="plan.id"><b>{{ plan.name }}</b><span>¥{{ price(plan) }} / {{ plan.duration_days }} 天</span></p></article>
      </aside>
    </section>

    <section class="admin-forms">
      <article class="panel"><h3>后台登录</h3><form v-if="!token" @submit.prevent="login"><label>用户名<input v-model="credentials.username" required /></label><label>密码<input v-model="credentials.password" required type="password" /></label><button class="primary" type="submit">登录后台</button></form><div v-else class="notice-box"><b>已登录</b><span>可审核会员与导出数据。</span><button @click="loadMembers">同步会员</button></div></article>
      <article class="panel"><h3>新增套餐</h3><form @submit.prevent="createPlan"><label>名称<input v-model="plan.name" required placeholder="例如：年度会员" /></label><label>价格（分）<input v-model.number="plan.price_cents" required type="number" min="0" /></label><label>有效天数<input v-model.number="plan.duration_days" required type="number" min="1" /></label><label>说明<input v-model="plan.description" placeholder="权益摘要" /></label><button class="primary" :disabled="!token" type="submit">创建套餐</button></form></article>
    </section>
  </section>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { api } from '../api'

const props = defineProps({ plans: { type: Array, default: () => [] } })
const emit = defineEmits(['refresh-plans', 'notice'])
const token = ref(localStorage.getItem('admin_token') || '')
const members = ref([])
const statusFilter = ref('all')
const statuses = ['pending', 'active', 'paused', 'cancelled']
const credentials = reactive({ username: 'admin', password: 'admin123' })
const plan = reactive({ name: '', price_cents: 0, duration_days: 30, description: '' })
const filteredMembers = computed(() => statusFilter.value === 'all' ? members.value : members.value.filter(item => item.status === statusFilter.value))
const pendingMembers = computed(() => members.value.filter(item => item.status === 'pending'))
const metrics = computed(() => [{ label: '会员总数', value: members.value.length }, { label: '已激活', value: count('active') }, { label: '待审核', value: count('pending') }, { label: '套餐数', value: props.plans.length }])
const pipeline = computed(() => [{ status: 'pending', label: '待审核', count: count('pending'), hint: '新申请进入队列' }, { status: 'active', label: '有效会员', count: count('active'), hint: '可正常享受权益' }, { status: 'paused', label: '暂停服务', count: count('paused'), hint: '需要人工跟进' }, { status: 'cancelled', label: '已取消', count: count('cancelled'), hint: '留存历史记录' }])

if (token.value) loadMembers()
async function login() { const result = await api.login(credentials); token.value = result.token; localStorage.setItem('admin_token', result.token); emit('notice', '后台登录成功'); await loadMembers() }
async function loadMembers() { members.value = await api.adminMembers(token.value) }
async function setStatus(id, status) { await api.updateStatus(token.value, id, status); emit('notice', `会员 #${id} 状态已更新为 ${status}`); await loadMembers() }
async function createPlan() { await api.createPlan(token.value, plan); Object.assign(plan, { name: '', price_cents: 0, duration_days: 30, description: '' }); emit('refresh-plans'); emit('notice', '套餐已创建') }
function count(status) { return members.value.filter(item => item.status === status).length }
function price(plan) { return (plan.price_cents / 100).toFixed(2) }
function exportMembers() { exportCSV('members.csv', ['编号', '姓名', '手机号', '邮箱', '套餐', '状态', '创建时间'], filteredMembers.value.map(item => [item.id, item.name, item.phone, item.email, item.plan_name, item.status, formatDate(item.created_at)])) }
function exportPlans() { exportCSV('plans.csv', ['编号', '名称', '价格(元)', '有效天数', '说明'], props.plans.map(item => [item.id, item.name, price(item), item.duration_days, item.description])) }
function exportCSV(filename, headers, rows) { const csv = [headers, ...rows].map(row => row.map(csvCell).join(',')).join('\n'); const url = URL.createObjectURL(new Blob([`\uFEFF${csv}`], { type: 'text/csv;charset=utf-8;' })); const link = document.createElement('a'); link.href = url; link.download = filename; link.click(); URL.revokeObjectURL(url) }
function csvCell(value) { return `"${String(value ?? '').replaceAll('"', '""')}"` }
function formatDate(value) { return value ? new Date(value).toLocaleString() : '-' }
</script>
