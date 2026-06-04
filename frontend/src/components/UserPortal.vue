<template>
  <section class="intake-page">
    <article class="surface intake-hero">
      <div class="surface-title"><div><p class="overline">Client Intake</p><h2>会员申请录入</h2><small>将线索快速转为可审核会员记录。</small></div><b>↗</b></div>
      <form @submit.prevent="submit">
        <div class="form-pair">
          <label>姓名<input v-model="form.name" required placeholder="请输入姓名" /></label>
          <label>手机号<input v-model="form.phone" required placeholder="请输入手机号" /></label>
        </div>
        <label>邮箱<input v-model="form.email" required type="email" placeholder="请输入邮箱" /></label>
        <label>订阅套餐
          <select v-model.number="form.plan_id" required>
            <option disabled value="0">请选择套餐</option>
            <option v-for="plan in plans" :key="plan.id" :value="plan.id">{{ plan.name }} · ¥{{ price(plan) }}</option>
          </select>
        </label>
        <button class="primary wide" type="submit">创建申请并交给后台审核</button>
      </form>
    </article>

    <article class="surface catalog-card">
      <div class="surface-title"><div><p class="overline">Subscription Menu</p><h2>套餐货架</h2></div></div>
      <div class="catalog-grid">
        <button v-for="plan in plans" :key="plan.id" class="catalog-item" :class="{ selected: form.plan_id === plan.id }" @click="form.plan_id = plan.id">
          <span>{{ plan.duration_days }}D</span><strong>{{ plan.name }}</strong><b>¥{{ price(plan) }}</b><small>{{ plan.description || '企业会员权益套餐' }}</small>
        </button>
      </div>
    </article>

    <article class="surface lookup-card">
      <div class="surface-title"><div><p class="overline">Lookup</p><h2>会员状态查询</h2></div></div>
      <form class="inline" @submit.prevent="query">
        <input v-model="queryId" placeholder="输入会员编号，例如 1" />
        <button class="secondary" type="submit">查询</button>
      </form>
      <div v-if="member" class="lookup-result">
        <strong>#{{ member.id }} / {{ member.name }}</strong>
        <span>{{ member.plan_name }} · {{ member.status }}</span>
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
