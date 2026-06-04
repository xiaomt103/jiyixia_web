<template>
  <section class="grid two">
    <article class="card">
      <h2>用户端 · 会员申请</h2>
      <form @submit.prevent="submit">
        <label>姓名<input v-model="form.name" required placeholder="请输入姓名" /></label>
        <label>手机号<input v-model="form.phone" required placeholder="请输入手机号" /></label>
        <label>邮箱<input v-model="form.email" required type="email" placeholder="请输入邮箱" /></label>
        <label>套餐
          <select v-model.number="form.plan_id" required>
            <option disabled value="0">请选择套餐</option>
            <option v-for="plan in plans" :key="plan.id" :value="plan.id">{{ plan.name }}</option>
          </select>
        </label>
        <button type="submit">提交申请</button>
      </form>
    </article>

    <article class="card">
      <h2>会员查询</h2>
      <form class="inline" @submit.prevent="query">
        <input v-model="queryId" placeholder="会员编号" />
        <button type="submit">查询</button>
      </form>
      <div v-if="member" class="result">
        <strong>{{ member.name }}</strong>
        <span>{{ member.plan_name }} · {{ member.status }}</span>
        <small>到期：{{ member.expire_at ? new Date(member.expire_at).toLocaleDateString() : '待激活' }}</small>
      </div>
      <h3>可选套餐</h3>
      <ul class="plans">
        <li v-for="plan in plans" :key="plan.id">
          <b>{{ plan.name }}</b><span>¥{{ (plan.price_cents / 100).toFixed(2) }} / {{ plan.duration_days }}天</span>
        </li>
      </ul>
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
</script>
