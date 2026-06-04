<template>
  <section class="user-page">
    <article class="panel intake-panel">
      <div class="panel-title hero-title">
        <div><p class="eyebrow">Customer Intake</p><h2>创建会员申请</h2><small>用户提交后，后台可立即审核与导出。</small></div>
        <span class="panel-icon">✦</span>
      </div>
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
        <button class="primary large" type="submit">提交申请并进入后台审核</button>
      </form>
    </article>

    <article class="panel plan-panel">
      <div class="panel-title"><div><p class="eyebrow">Plan Catalog</p><h2>套餐目录</h2></div></div>
      <div class="plan-grid">
        <button v-for="plan in plans" :key="plan.id" class="plan-card" :class="{ selected: form.plan_id === plan.id }" @click="form.plan_id = plan.id">
          <small>{{ plan.duration_days }} 天权益</small>
          <strong>{{ plan.name }}</strong>
          <span>¥{{ price(plan) }}</span>
          <p>{{ plan.description || '会员权益套餐' }}</p>
        </button>
      </div>
    </article>

    <article class="panel query-panel">
      <div class="panel-title"><div><p class="eyebrow">Lookup</p><h2>会员查询</h2></div></div>
      <form class="inline" @submit.prevent="query">
        <input v-model="queryId" placeholder="输入会员编号" />
        <button class="soft" type="submit">查询</button>
      </form>
      <div v-if="member" class="result-card">
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
