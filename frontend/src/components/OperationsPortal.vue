<template>
  <section class="module-page">
    <nav class="module-tabs"><button v-for="tab in config.tabs" :key="tab" :class="{ active: tab === activeTab }" @click="activeTab = tab">{{ tab }}</button></nav>
    <div class="filter-bar"><button v-for="item in periodOptions" :key="item" :class="{ active: item === activePeriod }" @click="activePeriod = item">{{ item }}</button><input :value="rangeText" readonly /><button class="primary" @click="loadRecords">刷新</button></div>
    <div class="metrics"><article v-for="item in metrics" :key="item.label" class="metric"><span>{{ item.label }}</span><strong>{{ item.value }}</strong></article></div>
    <section class="node-strip"><article v-for="node in config.nodes" :key="node.title"><b>{{ node.title }}</b><span>{{ node.text }}</span></article></section>

    <template v-if="module === 'dashboard'">
      <section v-if="activeTab === '数据总览'" class="dashboard-grid">
        <article class="report-card"><div class="report-title"><div><h2>经营摘要</h2><small>按当前时间范围聚合各模块记录。</small></div></div><div class="summary-list"><p v-for="item in moduleCounts" :key="item.label"><span>{{ item.label }}</span><b>{{ item.value }}</b></p></div></article>
        <article class="report-card"><div class="report-title"><div><h2>待办队列</h2><small>最近待处理事项优先展示。</small></div></div><ul class="timeline"><li v-for="item in visibleRecords.slice(0, 5)" :key="item.id"><b>{{ item.title }}</b><span>{{ item.module }} · {{ item.status }}</span></li></ul></article>
      </section>
      <section v-else-if="activeTab === '经营分析'" class="analysis-grid"><article v-for="item in analysisCards" :key="item.label" class="analysis-card"><span>{{ item.label }}</span><b>{{ item.value }}</b><small>{{ item.hint }}</small></article></section>
      <section v-else class="report-grid"><article v-for="item in reportCards" :key="item.title" class="report-card"><div class="report-title"><div><h2>{{ item.title }}</h2><small>{{ item.desc }}</small></div><button @click="exportReport(item.type)">生成报表</button></div></article></section>
    </template>

    <section v-else-if="module === 'booking'" class="schedule-board"><article v-for="slot in bookingSlots" :key="slot.time" class="slot-card"><time>{{ slot.time }}</time><b>{{ slot.title }}</b><span>{{ slot.status }}</span></article></section>
    <section v-else-if="module === 'order'" class="kanban-board"><article v-for="status in statuses" :key="status" class="kanban-col"><h3>{{ status }}</h3><p v-for="item in byStatus(status)" :key="item.id">{{ item.title }}<small>¥{{ money(item.amount_cents) }}</small></p></article></section>
    <section v-else-if="module === 'marketing'" class="campaign-grid"><article v-for="item in visibleRecords" :key="item.id" class="campaign-card"><b>{{ item.title }}</b><span>{{ item.status }}</span><small>预算 ¥{{ money(item.amount_cents) }}</small></article></section>
    <section v-else-if="module === 'finance'" class="finance-strip"><article><span>收入</span><b>¥{{ money(totalAmount) }}</b></article><article><span>记录</span><b>{{ visibleRecords.length }}</b></article><article><span>完成</span><b>{{ byStatus('done').length }}</b></article></section>
    <section v-else-if="module === 'settings'" class="settings-list"><article v-for="item in settings" :key="item"><b>{{ item }}</b><span>已预留配置位，后续可接入权限和系统参数。</span></article></section>

    <section class="report-layout">
      <aside class="sub-menu"><button v-for="item in config.menu" :key="item" :class="{ active: item === config.menu[0] }">{{ item }}</button></aside>
      <article class="report-card">
        <div class="report-title"><div><h2>{{ config.title }}明细</h2><small>{{ config.description }}</small></div><div class="tools"><button @click="exportRecords">导出 CSV</button></div></div>
        <div class="table-wrap"><table><thead><tr><th v-for="head in config.columns" :key="head">{{ head }}</th><th>状态</th><th>操作</th></tr></thead><tbody>
          <tr v-for="item in visibleRecords" :key="item.id"><td v-for="cell in rowCells(item)" :key="cell">{{ cell }}</td><td><span class="status" :class="item.status">{{ item.status }}</span></td><td class="row-actions"><button v-for="status in statuses" :key="status" @click="setStatus(item.id, status)">{{ status }}</button></td></tr>
          <tr v-if="visibleRecords.length === 0"><td class="empty" :colspan="config.columns.length + 2">当前时间范围暂无{{ config.title }}数据。</td></tr>
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
const activeTab = ref('')
const activePeriod = ref('本月')
const statuses = ['pending', 'active', 'done', 'cancelled']
const settings = ['门店资料', '角色权限', '通知模板']
const form = reactive({ title: '', amount_cents: 0 })
const views = {
  dashboard: view('经营总览', ['数据总览', '经营分析', '报表中心'], ['今天', '昨天', '上周'], '本月', ['经营摘要', '待办队列', '风险提醒'], ['采集:汇总各模块数据', '诊断:识别待办和风险', '复盘:沉淀经营报表'], '跨模块经营摘要。', ['事项', '模块', '金额', '时间'], '摘要标题', '关联金额（分）', '新增经营摘要'),
  booking: view('预约', ['预约日历', '到店跟进', '爽约记录'], ['今天', '明天', '本周'], '今天', ['预约列表', '到店确认', '改期记录'], ['预约:锁定服务时间', '提醒:到店前触达', '核销:确认服务完成'], '按时间管理预约和到店。', ['预约事项', '会员', '预约时间'], '预约事项', '预约定金（分）', '例如：明天 10:30 到店'),
  order: view('订单', ['订单看板', '履约中', '售后'], ['今天', '昨天', '上周'], '本月', ['订单列表', '履约跟进', '退款售后'], ['下单:生成交易事项', '履约:跟踪服务交付', '售后:处理退款和异常'], '跟踪会员订单履约。', ['订单事项', '金额', '时间'], '订单事项', '订单金额（分）', '例如：年度会员订单'),
  marketing: view('营销', ['活动计划', '触达记录', '转化分析'], ['今天', '昨天', '上周'], '本月', ['活动列表', '触达任务', '转化记录'], ['圈选:筛选目标会员', '触达:执行活动消息', '转化:统计预算产出'], '管理续费提醒和活动触达。', ['活动名称', '预算', '时间'], '活动名称', '活动预算（分）', '例如：六月续费提醒'),
  finance: view('财务', ['收款流水', '支出记录', '对账'], ['今天', '昨天', '上周'], '本月', ['财务流水', '收入汇总', '支出记录'], ['记账:登记收入支出', '对账:核对订单流水', '报表:输出经营数据'], '记录收款、支出和对账事项。', ['流水事项', '金额', '时间'], '流水事项', '金额（分）', '例如：会员收入日报'),
  settings: view('设置', ['系统配置', '权限', '通知'], ['今天', '昨天', '上周'], '本月', ['配置记录', '权限变更', '通知模板'], ['配置:维护基础参数', '权限:控制后台角色', '通知:管理消息模板'], '沉淀系统配置和权限变更记录。', ['配置事项', '模块', '时间'], '配置事项', '影响金额（分）', '例如：新增店长角色')
}
const config = computed(() => ({ ...(views[props.module] || views.dashboard), range: rangeText.value }))
const periodOptions = computed(() => [...config.value.filters, config.value.focus].filter((item, index, list) => list.indexOf(item) === index))
const range = computed(() => periodRange(activePeriod.value))
const rangeText = computed(() => `${formatRange(range.value.start)}  ~  ${formatRange(range.value.end)}`)
const visibleRecords = computed(() => records.value.filter(inPeriod))
const totalAmount = computed(() => visibleRecords.value.reduce((sum, item) => sum + (item.amount_cents || 0), 0))
const metrics = computed(() => [{ label: '记录数', value: visibleRecords.value.length }, { label: '进行中', value: byStatus('active').length }, { label: '待处理', value: byStatus('pending').length }, { label: '金额', value: `¥${money(totalAmount.value)}` }])
const moduleCounts = computed(() => ['booking', 'order', 'marketing', 'finance', 'report'].map(label => ({ label, value: visibleRecords.value.filter(item => item.module === label).length })))
const bookingSlots = computed(() => visibleRecords.value.slice(0, 6).map((item, index) => ({ time: `${9 + index}:30`, title: item.title, status: item.status })))
const analysisCards = computed(() => [{ label: '完成率', value: percent(byStatus('done').length, visibleRecords.value.length), hint: '已完成 / 当前范围记录' }, { label: '待办压力', value: byStatus('pending').length, hint: '需要优先处理的事项' }, { label: '平均金额', value: `¥${money(Math.round(totalAmount.value / Math.max(visibleRecords.value.length, 1)))}`, hint: '当前范围记录均值' }])
const reportCards = [{ title: '经营日报', desc: '导出当前时间范围的明细台账。', type: 'records' }, { title: '模块汇总', desc: '导出各模块记录数和占比。', type: 'summary' }, { title: '经营分析', desc: '导出完成率、待办和金额分析。', type: 'analysis' }]

watch(() => props.module, async () => { activeTab.value = config.value.tabs[0]; activePeriod.value = config.value.focus; await loadRecords() }, { immediate: true })
function view(title, tabs, filters, focus, menu, nodeText, description, columns, inputLabel, amountLabel, placeholder) { return { title, tabs, filters, focus, menu, nodes: nodeText.map(text => { const [title, body] = text.split(':'); return { title, text: body } }), description, columns, inputLabel, amountLabel, placeholder, createTitle: `新增${title}记录` } }
function token() { return localStorage.getItem('admin_token') || 'vercel-demo-admin-token' }
function apiModule() { return props.module === 'dashboard' ? '' : props.module }
async function loadRecords() { records.value = await api.operations(token(), apiModule()) }
async function createRecord() { await api.createOperation(token(), { module: props.module === 'dashboard' ? 'report' : props.module, title: form.title, amount_cents: form.amount_cents, status: 'pending' }); Object.assign(form, { title: '', amount_cents: 0 }); emit('notice', '记录已创建'); await loadRecords() }
async function setStatus(id, status) { await api.updateOperationStatus(token(), id, status); emit('notice', `记录 #${id} 已更新`); await loadRecords() }
function byStatus(status) { return visibleRecords.value.filter(item => item.status === status) }
function inPeriod(item) { const date = new Date(item.due_at || item.created_at || Date.now()); return date >= range.value.start && date <= range.value.end }
function periodRange(label) { const now = new Date(); const start = new Date(now); start.setHours(0, 0, 0, 0); const end = new Date(start); end.setHours(23, 59, 59, 999); if (label === '昨天') { start.setDate(start.getDate() - 1); end.setDate(end.getDate() - 1) } if (label === '上周') start.setDate(start.getDate() - 7); if (label === '本周') start.setDate(start.getDate() - start.getDay()); if (label === '明天') { start.setDate(start.getDate() + 1); end.setDate(end.getDate() + 1) } if (label === '本月') { start.setDate(1); end.setMonth(start.getMonth() + 1, 0) } return { start, end } }
function rowCells(item) { const amount = `¥${money(item.amount_cents)}`; const time = formatDate(item.due_at); const member = item.member_id || '-'; return ({ dashboard: [item.title, item.module, amount, time], booking: [item.title, member, time], order: [item.title, amount, time], marketing: [item.title, amount, time], finance: [item.title, amount, time], settings: [item.title, item.module || 'settings', time] }[props.module] || [item.title, item.module, amount, time]) }
function exportRecords() { const rows = visibleRecords.value.map(item => [item.id, item.module, item.title, item.amount_cents, item.status, formatDate(item.due_at)]); exportCSV(`${props.module}-${activePeriod.value}.csv`, ['编号', '模块', '标题', '金额(分)', '状态', '时间'], rows) }
function exportReport(type) { if (type === 'records') return exportRecords(); const rows = type === 'summary' ? moduleCounts.value.map(item => [item.label, item.value, percent(item.value, visibleRecords.value.length)]) : analysisCards.value.map(item => [item.label, item.value, item.hint]); exportCSV(`dashboard-${type}.csv`, type === 'summary' ? ['模块', '记录数', '占比'] : ['指标', '数值', '说明'], rows) }
function exportCSV(filename, headers, rows) { const csv = [headers, ...rows].map(row => row.map(cell).join(',')).join('\n'); const link = document.createElement('a'); link.href = URL.createObjectURL(new Blob([`\uFEFF${csv}`], { type: 'text/csv;charset=utf-8;' })); link.download = filename; link.click(); URL.revokeObjectURL(link.href) }
function cell(value) { return `"${String(value ?? '').replaceAll('"', '""')}"` }
function percent(count, total) { return `${Math.round((count / Math.max(total, 1)) * 100)}%` }
function money(cents) { return (Number(cents || 0) / 100).toFixed(2) }
function formatDate(value) { return value ? new Date(value).toLocaleString() : '-' }
function formatRange(value) { return value.toISOString().slice(0, 19).replace('T', ' ') }
</script>
