<template>
  <main class="app-frame">
    <aside class="sidebar">
      <div class="brand-mark"><span>记</span><div><strong>JiYiXia</strong><small>Membership OS</small></div></div>
      <nav class="side-nav">
        <button :class="{ active: tab === 'user' }" @click="tab = 'user'"><span>✦</span>用户门户</button>
        <button :class="{ active: tab === 'admin' }" @click="tab = 'admin'"><span>▦</span>运营后台</button>
      </nav>
      <div class="sidebar-card">
        <small>当前能力</small>
        <strong>{{ plans.length }} 个套餐 · CSV 导出</strong>
        <p>支持用户申请、后台审核、状态管理和开放 API。</p>
      </div>
    </aside>

    <section class="workspace">
      <header class="topbar">
        <div>
          <p class="eyebrow">SaaS Enterprise Console</p>
          <h1>{{ tab === 'user' ? '会员申请工作台' : '会员运营管理' }}</h1>
        </div>
        <div class="topbar-actions">
          <span class="env-pill">Vercel Demo</span>
          <span class="search-box">⌘K 搜索会员 / 套餐</span>
          <span class="avatar">YX</span>
        </div>
      </header>

      <section v-if="message" class="message">{{ message }}</section>
      <UserPortal v-if="tab === 'user'" :plans="plans" @created="onMemberCreated" />
      <AdminPortal v-else :plans="plans" @refresh-plans="loadPlans" @notice="notice" />
    </section>
  </main>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { api } from './api'
import AdminPortal from './components/AdminPortal.vue'
import UserPortal from './components/UserPortal.vue'

const tab = ref('admin')
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
  tab.value = 'admin'
}

function notice(text) {
  message.value = text
  setTimeout(() => (message.value = ''), 3500)
}
</script>
