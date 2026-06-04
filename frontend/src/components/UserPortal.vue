<template>
  <section class="user-page">
    <nav class="module-tabs"><button class="active">会员申请</button><button>进度查询</button><button>套餐权益</button></nav>
    <section class="user-hero">
      <div><p class="eyebrow">用户端</p><h2>三步完成会员申请</h2><span>选择套餐、提交资料、等待后台审核，用户侧只保留必要信息。</span></div>
      <div class="hero-metrics"><b>{{ plans.length }}</b><small>可选套餐</small><b>{{ selectedPlan?.name || '未选择' }}</b><small>当前选择</small></div>
    </section>
    <section class="journey"><article v-for="step in steps" :key="step.title"><b>{{ step.title }}</b><span>{{ step.text }}</span></article></section>

    <section class="user-layout">
      <article class="panel form-panel">
        <div class="section-head"><h3>创建会员申请</h3><small>提交后自动进入后台审核队列。</small></div>
        <form @submit.prevent="submit"><label>姓名<input v-model="form.name" required placeholder="请输入姓名" /></label><label>手机号<input v-model="form.phone" required placeholder="请输入手机号" /></label><label>邮箱<input v-model="form.email" required type="email" placeholder="请输入邮箱" /></label><label>套餐<select v-model.number="form.plan_id" required><option disabled value="0">请选择套餐</option><option v-for="plan in plans" :key="plan.id" :value="plan.id">{{ plan.name }} · ¥{{ price(plan) }}</option></select></label><button class="primary" type="submit">提交申请</button></form>
      </article>
      <article class="panel catalog-panel">
        <div class="section-head"><h3>套餐目录</h3><small>突出价格、周期和权益摘要，便于用户选择。</small></div>
        <div class="plans-grid"><button v-for="plan in plans" :key="plan.id" class="plan-item" :class="{ selected: form.plan_id === plan.id }" @click="selectPlan(plan)"><strong>{{ plan.name }}</strong><span>¥{{ price(plan) }}</span><small>{{ plan.duration_days }} 天 · {{ plan.description || '会员权益套餐' }}</small></button></div>
      </article>
      <article class="panel lookup-panel">
        <div class="section-head"><h3>会员查询</h3><small>用于用户自助查看审核和有效期。</small></div>
        <form class="inline" @submit.prevent="query"><input v-model="queryId" placeholder="输入会员编号" /><button type="submit">查询</button></form>
        <div v-if="member" class="result"><b>#{{ member.id }} · {{ member.name }}</b><span>{{ member.plan_name }} / {{ member.status }}</span><small>到期：{{ member.expire_at ? new Date(member.expire_at).toLocaleDateString() : '待激活' }}</small></div>
        <div v-else class="empty-note">提交申请后可用编号查询进度。</div>
      </article>
    </section>
  </section>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { api } from '../api'

const props = defineProps({ plans: { type: Array, default: () => [] } })
const emit = defineEmits(['created'])
const form = reactive({ name: '', phone: '', email: '', plan_id: 0 })
const queryId = ref('')
const member = ref(null)
const steps = [{ title: '选择套餐', text: '按周期和权益选择方案' }, { title: '提交资料', text: '录入基础联系方式' }, { title: '后台审核', text: '管理员激活后生效' }]
const selectedPlan = computed(() => props.plans.find(plan => plan.id === form.plan_id))

async function submit() { const created = await api.createMember(form); emit('created', created); Object.assign(form, { name: '', phone: '', email: '', plan_id: 0 }) }
async function query() { member.value = await api.getMember(queryId.value) }
function selectPlan(plan) { form.plan_id = plan.id }
function price(plan) { return (plan.price_cents / 100).toFixed(2) }
</script>
