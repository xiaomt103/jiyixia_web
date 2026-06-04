<template>
  <section class="admin-page">
    <div class="kpi-row">
      <article v-for="item in metrics" :key="item.label" class="kpi-card">
        <span>{{ item.icon }}</span><small>{{ item.label }}</small><strong>{{ item.value }}</strong>
      </article>
    </div>

    <section class="admin-grid">
      <article class="panel control-panel">
        <div class="panel-title"><div><p class="eyebrow">Access</p><h2>管理员控制</h2></div><span class="panel-icon">🔐</span></div>
        <form v-if="!token" @submit.prevent="login">
          <label>用户名<input v-model="credentials.username" required /></label>
          <label>密码<input v-model="credentials.password" required type="password" /></label>
          <button class="primary" type="submit">登录后台</button>
        </form>
        <div v-else class="signed-card">
          <span class="live-dot">已登录</span>
          <p>可以刷新会员、审核状态并导出运营数据。</p>
          <button class="soft" @click="loadMembers">刷新会员数据</button>
        </div>

        <div class="divider"></div>
        <div class="panel-title compact"><div><p class="eyebrow">Plans</p><h2>新增套餐</h2></div></div>
        <form @submit.prevent="createPlan">
          <label>套餐名称<input v-model="plan.name" required placeholder="例如：年度会员" /></label>
          <div class="form-row">
            <label>价格（分）<input v-model.number="plan.price_cents" required type="number" min="0" /></label>
            <label>有效天数<input v-model.number="plan.duration_days" required type="number" min="1" /></label>
          </div>
          <label>说明<input v-model="plan.description" placeholder="套餐权益说明" /></label>
          <button class="primary" :disabled="!token" type="submit">创建套餐</button>
        </form>
      </article>

      <article class="panel table-panel">
        <div class="table-head">
          <div><p class="eyebrow">Members</p><h2>会员数据表</h2><small>筛选、审核并导出表格数据</small></div>
          <div class="table-tools">
            <select v-model="statusFilter">
              <option value="all">全部状态</option>
              <option v-for="status in statuses" :key="status" :value="status">{{ status }}</option>
            </select>
            <button class="soft" @click="exportMembers">导出会员</button>
            <button class="soft" @click="exportPlans">导出套餐</button>
          </div>
        </div>
        <div class="data-table">
          <table>
            <thead><tr><th>会员</th><th>联系信息</th><th>套餐</th><th>状态</th><th>创建时间</th><th>操作</th></tr></thead>
            <tbody>
              <tr v-for="member in filteredMembers" :key="member.id">
                <td><div class="member-cell"><b>#{{ member.id }}</b><strong>{{ member.name }}</strong></div></td>
                <td><span>{{ member.phone }}</span><small>{{ member.email }}</small></td>
                <td>{{ member.plan_name }}</td>
                <td><span class="badge" :class="member.status">{{ member.status }}</span></td>
                <td>{{ formatDate(member.created_at) }}</td>
                <td class="actions"><button v-for="status in statuses" :key="status" @click="setStatus(member.id, status)">{{ status }}</button></td>
              </tr>
              <tr v-if="filteredMembers.length === 0"><td class="empty" colspan="6">暂无会员数据，请先在用户门户提交一条申请。</td></tr>
            </tbody>
          </table>
        </div>
      </article>
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
const metrics = computed(() => [
  { icon: '👥', label: '会员总数', value: members.value.length },
  { icon: '✅', label: '已激活', value: members.value.filter(item => item.status === 'active').length },
  { icon: '⏳', label: '待审核', value: members.value.filter(item => item.status === 'pending').length },
  { icon: '💎', label: '套餐数', value: props.plans.length }
])

if (token.value) loadMembers()

async function login() {
  const result = await api.login(credentials)
  token.value = result.token
  localStorage.setItem('admin_token', result.token)
  emit('notice', '后台登录成功')
  await loadMembers()
}

async function loadMembers() {
  members.value = await api.adminMembers(token.value)
}

async function setStatus(id, status) {
  await api.updateStatus(token.value, id, status)
  emit('notice', `会员 #${id} 状态已更新为 ${status}`)
  await loadMembers()
}

async function createPlan() {
  await api.createPlan(token.value, plan)
  Object.assign(plan, { name: '', price_cents: 0, duration_days: 30, description: '' })
  emit('refresh-plans')
  emit('notice', '套餐已创建')
}

function exportMembers() {
  exportCSV('members.csv', ['编号', '姓名', '手机号', '邮箱', '套餐', '状态', '创建时间'], filteredMembers.value.map(item => [
    item.id, item.name, item.phone, item.email, item.plan_name, item.status, formatDate(item.created_at)
  ]))
}

function exportPlans() {
  exportCSV('plans.csv', ['编号', '名称', '价格(元)', '有效天数', '说明'], props.plans.map(item => [
    item.id, item.name, (item.price_cents / 100).toFixed(2), item.duration_days, item.description
  ]))
}

function exportCSV(filename, headers, rows) {
  const csv = [headers, ...rows].map(row => row.map(csvCell).join(',')).join('\n')
  const url = URL.createObjectURL(new Blob([`\uFEFF${csv}`], { type: 'text/csv;charset=utf-8;' }))
  const link = document.createElement('a')
  link.href = url
  link.download = filename
  link.click()
  URL.revokeObjectURL(url)
}

function csvCell(value) {
  return `"${String(value ?? '').replaceAll('"', '""')}"`
}

function formatDate(value) {
  return value ? new Date(value).toLocaleString() : '-'
}
</script>
