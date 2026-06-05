<template>
  <section class="console-page beauty-console">
    <section class="beauty-hero">
      <article class="member-card">
        <span>JiYiXia Member</span><h2>{{ memberName }}</h2><p>{{ memberLine }}</p><b>{{ memberStatus }}</b>
      </article>
      <article class="quick-panel">
        <h3>我的服务</h3><div class="console-grid"><button v-for="item in actions" :key="item.mode" @click="open(item.mode)"><b>{{ item.title }}</b><span>{{ item.text }}</span></button></div>
      </article>
    </section>

    <section class="beauty-summary">
      <article v-for="item in summary" :key="item.label"><span>{{ item.label }}</span><b>{{ item.value }}</b><small>{{ item.hint }}</small></article>
    </section>

    <section class="beauty-layout">
      <article class="report-card appointment-card"><div class="report-title"><div><h2>下次护理</h2><small>根据会员权益生成服务提醒。</small></div></div><p><b>面部补水护理</b><span>建议 3 天内预约，到店前可联系客服确认时间。</span></p><button class="primary" @click="open('query')">查看会员状态</button></article>
      <article class="report-card service-card"><div class="report-title"><div><h2>推荐项目</h2><small>美业会员常用服务。</small></div></div><div class="service-list"><p v-for="item in services" :key="item"><b>{{ item }}</b><span>适合本月保养计划</span></p></div></article>
    </section>

    <article class="report-card">
      <div class="report-title"><div><h2>会员套餐</h2><small>选择套餐后进入申请页。</small></div></div>
      <div class="plans-grid"><button v-for="plan in topPlans" :key="plan.id" class="plan-item" @click="open('apply')"><strong>{{ plan.name }}</strong><span>¥{{ price(plan) }}</span><small>{{ plan.duration_days }} 天 · {{ plan.description || '会员权益套餐' }}</small></button></div>
    </article>
  </section>
</template>

<script setup>
import { computed } from 'vue'
import { listDemoMembers } from '../demoStore'

const props = defineProps({ plans: { type: Array, default: () => [] } })
const emit = defineEmits(['navigate'])
const latestMember = computed(() => listDemoMembers().at(-1))
const memberName = computed(() => latestMember.value?.name || '欢迎来到会员中心')
const memberStatus = computed(() => latestMember.value?.status || '未申请')
const memberLine = computed(() => latestMember.value ? `${latestMember.value.plan_name || '会员套餐'} · 编号 #${latestMember.value.id}` : '申请会员后可查看专属权益和服务提醒')
const topPlans = computed(() => props.plans.slice(0, 3))
const summary = computed(() => [{ label: '会员状态', value: memberStatus.value, hint: '服务确认后生效' }, { label: '可选套餐', value: props.plans.length, hint: '支持按需申请' }, { label: '护理提醒', value: latestMember.value ? '1项' : '待开启', hint: '根据会员状态生成' }])
const actions = [{ mode: 'apply', title: '申请会员', text: '填写资料并选择套餐' }, { mode: 'query', title: '查询状态', text: '查看会员编号和有效期' }, { mode: 'plans', title: '查看套餐', text: '对比周期和权益' }]
const services = ['面部补水护理', '头皮舒缓护理', '肩颈放松护理']
function open(mode) { emit('navigate', mode) }
function price(plan) { return (plan.price_cents / 100).toFixed(2) }
</script>
