<template>
  <section class="user-page">
    <article class="card form-card">
      <div class="card-title"><div><p class="caption">Intake</p><h2>创建会员申请</h2><small>提交后自动进入后台审核队列。</small></div></div>
      <form @submit.prevent="submit">
        <div class="form-row">
          <label>姓名<input v-model="form.name" required placeholder="请输入姓名" /></label>
          <label>手机号<input v-model="form.phone" required placeholder="请输入手机号" /></label>
        </div>
        <label>邮箱<input v-model="form.email" required type="email" placeholder="请输入邮箱" /></label>
        <label>套餐
          <select v-model.number="form.plan_id" required>
            <option disabled value="0">请选择套餐</option>
            <option v-for="plan in plans" :key="plan.id" :value="plan.id">{{ plan.name }} · ¥{{ price(plan) }}</option>
          </select>
        </label>
        <button class="primary wide" type="submit">提交申请</button>
      </form>
    </article>

    <article class="card catalog-card">
      <div class="card-title"><div><p class="caption">Plans</p><h2>套餐目录</h2></div></div>
      <div class="plans-grid">
        <button v-for="plan in plans" :key="plan.id" class="plan-item" :class="{ selected: form.plan_id === plan.id }" @click="form.plan_id = plan.id">
          <strong>{{ plan.name }}</strong><span>¥{{ price(plan) }}</span><small>{{ plan.duration_days }} 天 · {{ plan.description || '会员权益套餐' }}</small>
        </button>
      </div>
    </article>

    <article class="card lookup-card">
      <div class="card-title"><div><p class="caption">Lookup</p><h2>会员查询</h2></div></div>
      <form class="inline" @submit.prevent="query">
        <input v-model="queryId" placeholder="输入会员编号" />
        <button class="secondary" type="submit">查询</button>
      </form>
      <div v-if="member" class="result">
        <b>#{{ member.id }} · {{ member.name }}</b>
        <span>{{ member.plan_name }} / {{ member.status }}</span>
        <small>到期：{{ member.expire_at ? new Date(member.expire_at).toLocaleDateString() : '待激活' }}</small>
      </div>
    </article>
  </section>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { api } from '../api'

const props = defineProps({ plans: { type: Array, default: () => [] } })
const emit = defineEmits(['created'])
const form = reactive({ name: '', phone: '', email: '', plan_id: 0 })
const queryId = ref('')
const member = ref(null)

async function submit() {
  const created = await api.createMember(form)
  emit('created', created)
  Object.assign(form, { name: '', phone: '', email: '', plan_id: 0 })
}

async function query() {
  member.value = await api.getMember(queryId.value)
}

function price(plan) {
  return (plan.price_cents / 100).toFixed(2)
}
</script>
