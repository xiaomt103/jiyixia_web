<template>
  <main class="shell">
    <header class="hero">
      <div class="hero-copy">
        <p class="eyebrow">Member System</p>
        <h1>记一下会员管理系统</h1>
        <p>面向用户申请、后台审核、套餐管理与数据导出的现代化会员运营工作台。</p>
        <div class="hero-stats">
          <span><b>{{ plans.length }}</b> 套餐</span>
          <span><b>CSV</b> 数据导出</span>
          <span><b>API</b> 可扩展</span>
        </div>
      </div>
      <nav class="tabs">
        <button :class="{ active: tab === 'user' }" @click="tab = 'user'">用户端</button>
        <button :class="{ active: tab === 'admin' }" @click="tab = 'admin'">后台管理端</button>
      </nav>
    </header>

    <section v-if="message" class="message">{{ message }}</section>
    <UserPortal v-if="tab === 'user'" :plans="plans" @created="onMemberCreated" />
    <AdminPortal v-else :plans="plans" @refresh-plans="loadPlans" @notice="notice" />
  </main>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { api } from './api'
import AdminPortal from './components/AdminPortal.vue'
import UserPortal from './components/UserPortal.vue'

const tab = ref('user')
const plans = ref([])
const message = ref('')

onMounted(loadPlans)

async function loadPlans() {
  try {
    plans.value = await api.plans()
  } catch (error) {
    notice(error.message)
  }
}

function onMemberCreated(member) {
  notice(`申请已提交，会员编号：${member.id}`)
}

function notice(text) {
  message.value = text
  setTimeout(() => (message.value = ''), 3500)
}
</script>
