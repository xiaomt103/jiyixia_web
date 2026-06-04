<template>
  <section class="module-page">
    <nav class="module-tabs"><button class="active">{{ config.title }}</button><button>明细</button><button>导出</button></nav>
    <div class="filter-bar"><span>今天</span><span>本周</span><b>本月</b><input :value="config.range" readonly /><button class="primary" @click="loadRecords">刷新</button></div>
    <div class="metrics"><article v-for="item in metrics" :key="item.label" class="metric"><span>{{ item.label }}</span><strong>{{ item.value }}</strong></article></div>
    <section class="report-layout">
      <aside class="sub-menu"><button class="active">{{ config.menu }}</button><button>状态跟进</button><button>导出记录</button></aside>
      <article class="report-card">
        <div class="report-title"><div><h2>{{ config.title }}台账</h2><small>{{ config.description }}</small></div><div class="tools"><button @click="exportRecords">导出 CSV</button></div></div>
        <div class="table-wrap"><table><thead><tr><th>事项</th><th>关联会员</th><th>金额</th><th>状态</th><th>时间</th><th>操作</th></tr></thead><tbody>
          <tr v-for="item in records" :key="item.id"><td><b>#{{ item.id }}</b> {{ item.title }}</td><td>{{ item.member_id || '-' }}</td><td>¥{{ money(item.amount_cents) }}</td><td><span class="status" :class="item.status">{{ item.status }}</span></td><td>{{ formatDate(item.due_at) }}</td><td class="row-actions"><button v-for="status in statuses" :key="status" @click="setStatus(item.id, status)">{{ status }}</button></td></tr>
          <tr v-if="records.length === 0"><td class="empty" colspan="6">暂无{{ config.title }}数据。</td></tr>
        </tbody></table></div>
      </article>
    </section>
    <article class="panel quick-create"><h3>新增{{ config.title }}记录</h3><form @submit.prevent="createRecord"><label>事项标题<input v-model="form.title" required :placeholder="config.placeholder" /></label><label>金额（分）<input v-model.number="form.amount_cents" type="number" min="0" /></label><button class="primary" type="submit">保存记录</button></form></article>
  </section>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { api } from '../api'

const props = defineProps({ module: { type: String, required: true } })
const emit = defineEmits(['notice'])
const records = ref([])
const statuses = ['pending', 'active', 'done', 'cancelled']
const form = reactive({ title: '', amount_cents: 0 })
const configs = {
  dashboard: ['总览', '综合经营数据概览', '经营摘要', '新增经营摘要'], booking: ['预约', '预约与到店管理', '预约列表', '新增到店预约'],
  order: ['订单', '订单履约管理', '订单列表', '新增订单事项'], marketing: ['营销', '活动与触达管理', '活动列表', '新增营销活动'],
  finance: ['财务', '收款与支出记录', '财务流水', '新增财务记录'], settings: ['设置', '系统配置与记录', '配置记录', '新增配置记录']
}
const config = computed(() => {
  const item = configs[props.module] || configs.dashboard
  return { title: item[0], description: item[1], menu: item[2], placeholder: item[3], range: '2026-06-01 00:00:00  ~  2026-06-30 23:59:59' }
})
const metrics = computed(() => [
  { label: '记录数', value: records.value.length }, { label: '进行中', value: records.value.filter(item => item.status === 'active').length },
  { label: '待处理', value: records.value.filter(item => item.status === 'pending').length }, { label: '总金额', value: `¥${money(records.value.reduce((sum, item) => sum + (item.amount_cents || 0), 0))}` }
])

watch(() => props.module, loadRecords, { immediate: true })

function token() { return localStorage.getItem('admin_token') || 'vercel-demo-admin-token' }
async function loadRecords() { records.value = await api.operations(token(), props.module === 'dashboard' ? '' : props.module) }
async function createRecord() { await api.createOperation(token(), { module: props.module === 'dashboard' ? 'report' : props.module, title: form.title, amount_cents: form.amount_cents, status: 'pending' }); Object.assign(form, { title: '', amount_cents: 0 }); emit('notice', '记录已创建'); await loadRecords() }
async function setStatus(id, status) { await api.updateOperationStatus(token(), id, status); emit('notice', `记录 #${id} 已更新`); await loadRecords() }
function exportRecords() { const rows = records.value.map(item => [item.id, item.module, item.title, item.amount_cents, item.status, formatDate(item.due_at)]); exportCSV(`${props.module}.csv`, ['编号', '模块', '标题', '金额(分)', '状态', '时间'], rows) }
function exportCSV(filename, headers, rows) { const csv = [headers, ...rows].map(row => row.map(cell).join(',')).join('\n'); const link = document.createElement('a'); link.href = URL.createObjectURL(new Blob([`\uFEFF${csv}`], { type: 'text/csv;charset=utf-8;' })); link.download = filename; link.click(); URL.revokeObjectURL(link.href) }
function cell(value) { return `"${String(value ?? '').replaceAll('"', '""')}"` }
function money(cents) { return (Number(cents || 0) / 100).toFixed(2) }
function formatDate(value) { return value ? new Date(value).toLocaleString() : '-' }
</script>
