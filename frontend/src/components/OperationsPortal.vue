<template>
  <section class="module-page">
    <nav class="module-tabs"><button v-for="tab in config.tabs" :key="tab" :class="{ active: tab === config.tabs[0] }">{{ tab }}</button></nav>
    <div class="filter-bar"><span v-for="item in config.filters" :key="item">{{ item }}</span><b>{{ config.focus }}</b><input :value="config.range" readonly /><button class="primary" @click="loadRecords">刷新</button></div>
    <div class="metrics"><article v-for="item in metrics" :key="item.label" class="metric"><span>{{ item.label }}</span><strong>{{ item.value }}</strong></article></div>

    <section v-if="module === 'dashboard'" class="dashboard-grid">
      <article class="report-card"><div class="report-title"><div><h2>经营摘要</h2><small>聚合预约、订单、营销、财务和报表记录。</small></div></div><div class="summary-list"><p v-for="item in moduleCounts" :key="item.label"><span>{{ item.label }}</span><b>{{ item.value }}</b></p></div></article>
      <article class="report-card"><div class="report-title"><div><h2>待办队列</h2><small>按时间排序展示最近事项。</small></div></div><ul class="timeline"><li v-for="item in records.slice(0, 5)" :key="item.id"><b>{{ item.title }}</b><span>{{ item.module }} · {{ item.status }}</span></li></ul></article>
    </section>

    <section v-else-if="module === 'booking'" class="schedule-board">
      <article v-for="slot in bookingSlots" :key="slot.time" class="slot-card"><time>{{ slot.time }}</time><b>{{ slot.title }}</b><span>{{ slot.status }}</span></article>
    </section>

    <section v-else-if="module === 'order'" class="kanban-board">
      <article v-for="status in statuses" :key="status" class="kanban-col"><h3>{{ status }}</h3><p v-for="item in byStatus(status)" :key="item.id">{{ item.title }}<small>¥{{ money(item.amount_cents) }}</small></p></article>
    </section>

    <section v-else-if="module === 'marketing'" class="campaign-grid">
      <article v-for="item in records" :key="item.id" class="campaign-card"><b>{{ item.title }}</b><span>{{ item.status }}</span><small>预算 ¥{{ money(item.amount_cents) }}</small></article>
    </section>

    <section v-else-if="module === 'finance'" class="finance-strip">
      <article><span>收入</span><b>¥{{ money(totalAmount) }}</b></article><article><span>记录</span><b>{{ records.length }}</b></article><article><span>完成</span><b>{{ byStatus('done').length }}</b></article>
    </section>

    <section v-else-if="module === 'settings'" class="settings-list">
      <article v-for="item in settings" :key="item"><b>{{ item }}</b><span>已预留配置位，后续可接入权限和系统参数。</span></article>
    </section>

    <section class="report-layout">
      <aside class="sub-menu"><button v-for="item in config.menu" :key="item" :class="{ active: item === config.menu[0] }">{{ item }}</button></aside>
      <article class="report-card">
        <div class="report-title"><div><h2>{{ config.title }}明细</h2><small>{{ config.description }}</small></div><div class="tools"><button @click="exportRecords">导出 CSV</button></div></div>
        <div class="table-wrap"><table><thead><tr><th v-for="head in config.columns" :key="head">{{ head }}</th><th>状态</th><th>操作</th></tr></thead><tbody>
          <tr v-for="item in records" :key="item.id"><td v-for="cell in rowCells(item)" :key="cell"><span v-html="cell"></span></td><td><span class="status" :class="item.status">{{ item.status }}</span></td><td class="row-actions"><button v-for="status in statuses" :key="status" @click="setStatus(item.id, status)">{{ status }}</button></td></tr>
          <tr v-if="records.length === 0"><td class="empty" :colspan="config.columns.length + 2">暂无{{ config.title }}数据。</td></tr>
        </tbody></table></div>
      </article>
    </section>

    <article class="panel quick-create"><h3>{{ config.createTitle }}</h3><form @submit.prevent="createRecord"><label>{{ config.inputLabel }}<input v-model="form.title" required :placeholder="config.placeholder" /></label><label>{{ config.amountLabel }}<input v-model.number="form.amount_cents" type="number" min="0" /></label><button class="primary" type="submit">保存记录</button></form></article>
  </section>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { api } from '../api'

const props = defineProps({ module: { type: String, required: true } })
const emit = defineEmits(['notice'])
const records = ref([])
const statuses = ['pending', 'active', 'done', 'cancelled']
const settings = ['门店资料', '角色权限', '通知模板']
const form = reactive({ title: '', amount_cents: 0 })
const views = {
  dashboard: { title: '经营总览', tabs: ['数据总览', '经营分析', '报表中心'], filters: ['今天', '昨天', '上周'], focus: '本月', menu: ['经营摘要', '待办队列', '风险提醒'], description: '跨模块经营摘要。', columns: ['事项', '模块', '金额', '时间'], inputLabel: '摘要标题', amountLabel: '关联金额（分）', placeholder: '新增经营摘要', createTitle: '新增经营摘要' },
  booking: { title: '预约', tabs: ['预约日历', '到店跟进', '爽约记录'], filters: ['今天', '明天', '本周'], focus: '待到店', menu: ['预约列表', '到店确认', '改期记录'], description: '按时间管理预约和到店。', columns: ['预约事项', '会员', '预约时间'], inputLabel: '预约事项', amountLabel: '预约定金（分）', placeholder: '例如：明天 10:30 到店', createTitle: '新增预约' },
  order: { title: '订单', tabs: ['订单看板', '履约中', '售后'], filters: ['待确认', '履约中', '已完成'], focus: '全部订单', menu: ['订单列表', '履约跟进', '退款售后'], description: '跟踪会员订单履约。', columns: ['订单事项', '金额', '时间'], inputLabel: '订单事项', amountLabel: '订单金额（分）', placeholder: '例如：年度会员订单', createTitle: '新增订单事项' },
  marketing: { title: '营销', tabs: ['活动计划', '触达记录', '转化分析'], filters: ['短信', '微信', '优惠券'], focus: '进行中', menu: ['活动列表', '触达任务', '转化记录'], description: '管理续费提醒和活动触达。', columns: ['活动名称', '预算', '时间'], inputLabel: '活动名称', amountLabel: '活动预算（分）', placeholder: '例如：六月续费提醒', createTitle: '新增营销活动' },
  finance: { title: '财务', tabs: ['收款流水', '支出记录', '对账'], filters: ['收入', '支出', '待确认'], focus: '本月', menu: ['财务流水', '收入汇总', '支出记录'], description: '记录收款、支出和对账事项。', columns: ['流水事项', '金额', '时间'], inputLabel: '流水事项', amountLabel: '金额（分）', placeholder: '例如：会员收入日报', createTitle: '新增财务记录' },
  settings: { title: '设置', tabs: ['系统配置', '权限', '通知'], filters: ['基础', '权限', '模板'], focus: '配置项', menu: ['配置记录', '权限变更', '通知模板'], description: '沉淀系统配置和权限变更记录。', columns: ['配置事项', '模块', '时间'], inputLabel: '配置事项', amountLabel: '影响金额（分）', placeholder: '例如：新增店长角色', createTitle: '新增配置记录' }
}
const config = computed(() => ({ ...views[props.module] || views.dashboard, range: '2026-06-01 00:00:00  ~  2026-06-30 23:59:59' }))
const totalAmount = computed(() => records.value.reduce((sum, item) => sum + (item.amount_cents || 0), 0))
const metrics = computed(() => [{ label: '记录数', value: records.value.length }, { label: '进行中', value: byStatus('active').length }, { label: '待处理', value: byStatus('pending').length }, { label: '金额', value: `¥${money(totalAmount.value)}` }])
const moduleCounts = computed(() => ['booking', 'order', 'marketing', 'finance', 'report'].map(label => ({ label, value: records.value.filter(item => item.module === label).length })))
const bookingSlots = computed(() => records.value.slice(0, 6).map((item, index) => ({ time: `${9 + index}:30`, title: item.title, status: item.status })))

watch(() => props.module, loadRecords, { immediate: true })
function token() { return localStorage.getItem('admin_token') || 'vercel-demo-admin-token' }
function apiModule() { return props.module === 'dashboard' ? '' : props.module }
async function loadRecords() { records.value = await api.operations(token(), apiModule()) }
async function createRecord() { await api.createOperation(token(), { module: props.module === 'dashboard' ? 'report' : props.module, title: form.title, amount_cents: form.amount_cents, status: 'pending' }); Object.assign(form, { title: '', amount_cents: 0 }); emit('notice', '记录已创建'); await loadRecords() }
async function setStatus(id, status) { await api.updateOperationStatus(token(), id, status); emit('notice', `记录 #${id} 已更新`); await loadRecords() }
function byStatus(status) { return records.value.filter(item => item.status === status) }
function rowCells(item) {
  const amount = `¥${money(item.amount_cents)}`
  const time = formatDate(item.due_at)
  const member = item.member_id || '-'
  const rows = {
    dashboard: [item.title, item.module, amount, time], booking: [item.title, member, time], order: [item.title, amount, time],
    marketing: [item.title, amount, time], finance: [item.title, amount, time], settings: [item.title, item.module || 'settings', time]
  }
  return rows[props.module] || rows.dashboard
}
function exportRecords() { const rows = records.value.map(item => [item.id, item.module, item.title, item.amount_cents, item.status, formatDate(item.due_at)]); exportCSV(`${props.module}.csv`, ['编号', '模块', '标题', '金额(分)', '状态', '时间'], rows) }
function exportCSV(filename, headers, rows) { const csv = [headers, ...rows].map(row => row.map(cell).join(',')).join('\n'); const link = document.createElement('a'); link.href = URL.createObjectURL(new Blob([`\uFEFF${csv}`], { type: 'text/csv;charset=utf-8;' })); link.download = filename; link.click(); URL.revokeObjectURL(link.href) }
function cell(value) { return `"${String(value ?? '').replaceAll('"', '""')}"` }
function money(cents) { return (Number(cents || 0) / 100).toFixed(2) }
function formatDate(value) { return value ? new Date(value).toLocaleString() : '-' }
</script>
