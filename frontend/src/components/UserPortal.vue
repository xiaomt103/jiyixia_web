<template>
  <section class="user-page simple-page">
    <nav class="module-tabs"><button v-for="item in tabs" :key="item.key" :class="{ active: mode === item.key }" @click="mode = item.key">{{ item.label }}</button></nav>

    <article v-if="mode === 'apply'" class="panel task-panel">
      <div class="section-head"><h3>提交会员申请</h3><small>填写资料并选择套餐，提交后等待服务确认。</small></div>
      <form @submit.prevent="submit"><label>姓名<input v-model="form.name" required placeholder="请输入姓名" /></label><label>手机号<input v-model="form.phone" required placeholder="请输入手机号" /></label><label>邮箱<input v-model="form.email" required type="email" placeholder="请输入邮箱" /></label><label>套餐<select v-model.number="form.plan_id" required><option disabled value="0">请选择套餐</option><option v-for="plan in plans" :key="plan.id" :value="plan.id">{{ plan.name }} · ¥{{ price(plan) }}</option></select></label><button class="primary" type="submit">提交申请</button></form>
    </article>

    <article v-else-if="mode === 'query'" class="panel task-panel">
      <div class="section-head"><h3>查询会员状态</h3><small>输入会员编号查看审核状态和有效期。</small></div>
      <form class="inline" @submit.prevent="query"><input v-model="queryId" placeholder="输入会员编号" /><button type="submit">查询</button></form>
      <div v-if="member" class="result"><b>#{{ member.id }} · {{ member.name }}</b><span>{{ member.plan_name }} / {{ member.status }}</span><small>到期：{{ member.expire_at ? new Date(member.expire_at).toLocaleDateString() : '待激活' }}</small></div>
      <div v-else class="empty-note">暂无查询结果。</div>
    </article>

    <article v-else class="panel task-panel">
      <div class="section-head"><h3>套餐列表</h3><small>查看可申请的会员套餐。</small></div>
      <div class="plans-grid"><button v-for="plan in plans" :key="plan.id" class="plan-item" @click="selectPlan(plan)"><strong>{{ plan.name }}</strong><span>¥{{ price(plan) }}</span><small>{{ plan.duration_days }} 天 · {{ plan.description || '会员权益套餐' }}</small></button></div>
    </article>
  </section>
</template>

<script setup>
import { reactive, ref, watch } from 'vue'
import { api } from '../api'

const props = defineProps({ plans: { type: Array, default: () => [] }, startMode: { type: String, default: 'apply' } })
const emit = defineEmits(['created'])
const tabs = [{ key: 'apply', label: '会员申请' }, { key: 'query', label: '状态查询' }, { key: 'plans', label: '套餐列表' }]
const mode = ref(props.startMode)
const form = reactive({ name: '', phone: '', email: '', plan_id: 0 })
const queryId = ref('')
const member = ref(null)

watch(() => props.startMode, value => { mode.value = value })

async function submit() { const created = await api.createMember(form); emit('created', created); Object.assign(form, { name: '', phone: '', email: '', plan_id: 0 }) }
async function query() { member.value = await api.getMember(queryId.value) }
function selectPlan(plan) { form.plan_id = plan.id; mode.value = 'apply' }
function price(plan) { return (plan.price_cents / 100).toFixed(2) }
</script>
