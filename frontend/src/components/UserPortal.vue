<template>
  <section class="user-layout">
    <article class="card application-card">
      <div class="section-title">
        <span class="icon">✨</span>
        <div><p class="eyebrow dark">User Portal</p><h2>会员申请</h2></div>
      </div>
      <form @submit.prevent="submit">
        <label>姓名<input v-model="form.name" required placeholder="请输入姓名" /></label>
        <label>手机号<input v-model="form.phone" required placeholder="请输入手机号" /></label>
        <label>邮箱<input v-model="form.email" required type="email" placeholder="请输入邮箱" /></label>
        <label>套餐
          <select v-model.number="form.plan_id" required>
            <option disabled value="0">请选择套餐</option>
            <option v-for="plan in plans" :key="plan.id" :value="plan.id">{{ plan.name }} · ¥{{ price(plan) }}</option>
          </select>
        </label>
        <button class="primary" type="submit">提交申请</button>
      </form>
    </article>

    <article class="card plans-card">
      <div class="section-title compact">
        <span class="icon">💎</span>
        <div><p class="eyebrow dark">Plans</p><h2>可选套餐</h2></div>
      </div>
      <div class="plan-grid">
        <div v-for="plan in plans" :key="plan.id" class="plan-card" :class="{ selected: form.plan_id === plan.id }" @click="form.plan_id = plan.id">
          <div><strong>{{ plan.name }}</strong><small>{{ plan.description || '会员权益套餐' }}</small></div>
          <p><span>¥{{ price(plan) }}</span>/ {{ plan.duration_days }}天</p>
        </div>
      </div>

      <div class="query-box">
        <h3>会员查询</h3>
        <form class="inline" @submit.prevent="query">
          <input v-model="queryId" placeholder="输入会员编号" />
          <button class="ghost" type="submit">查询</button>
        </form>
        <div v-if="member" class="result">
          <strong>{{ member.name }}</strong>
          <span>{{ member.plan_name }} · <b>{{ member.status }}</b></span>
          <small>到期：{{ member.expire_at ? new Date(member.expire_at).toLocaleDateString() : '待激活' }}</small>
        </div>
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
