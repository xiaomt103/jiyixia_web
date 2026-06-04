<template>
  <section class="grid two">
    <article class="card">
      <h2>后台管理端 · 登录</h2>
      <form v-if="!token" @submit.prevent="login">
        <label>用户名<input v-model="credentials.username" required /></label>
        <label>密码<input v-model="credentials.password" required type="password" /></label>
        <button type="submit">登录后台</button>
      </form>
      <div v-else class="toolbar">
        <span>已登录，可调用后台 API 控制会员状态。</span>
        <button @click="loadMembers">刷新会员</button>
      </div>

      <h3>新增套餐</h3>
      <form @submit.prevent="createPlan">
        <label>名称<input v-model="plan.name" required /></label>
        <label>价格（分）<input v-model.number="plan.price_cents" required type="number" min="0" /></label>
        <label>有效天数<input v-model.number="plan.duration_days" required type="number" min="1" /></label>
        <label>说明<input v-model="plan.description" /></label>
        <button :disabled="!token" type="submit">创建套餐</button>
      </form>
    </article>

    <article class="card wide">
      <h2>会员列表</h2>
      <div class="table">
        <div class="row head"><span>编号</span><span>会员</span><span>套餐</span><span>状态</span><span>操作</span></div>
        <div v-for="member in members" :key="member.id" class="row">
          <span>#{{ member.id }}</span>
          <span>{{ member.name }}<small>{{ member.phone }}</small></span>
          <span>{{ member.plan_name }}</span>
          <span class="badge">{{ member.status }}</span>
          <span class="actions">
            <button v-for="status in statuses" :key="status" @click="setStatus(member.id, status)">{{ status }}</button>
          </span>
        </div>
      </div>
    </article>
  </section>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { api } from '../api'

const emit = defineEmits(['refresh-plans', 'notice'])
const token = ref(localStorage.getItem('admin_token') || '')
const members = ref([])
const statuses = ['active', 'paused', 'cancelled']
const credentials = reactive({ username: 'admin', password: 'admin123' })
const plan = reactive({ name: '', price_cents: 0, duration_days: 30, description: '' })

if (token.value) loadMembers()

async function login() {
  const result = await api.login(credentials)
  token.value = result.token
  localStorage.setItem('admin_token', result.token)
  emit('notice', '后台登录成功')
  await loadMembers()
}

async function loadMembers() {
  members.value = await api.adminMembers(token.value)
}

async function setStatus(id, status) {
  await api.updateStatus(token.value, id, status)
  emit('notice', `会员 #${id} 状态已更新为 ${status}`)
  await loadMembers()
}

async function createPlan() {
  await api.createPlan(token.value, plan)
  Object.assign(plan, { name: '', price_cents: 0, duration_days: 30, description: '' })
  emit('refresh-plans')
  emit('notice', '套餐已创建')
}
</script>
