<template>
  <main class="app-shell">
    <header class="global-header">
      <div class="brand"><span>记</span><div><b>JiYiXia</b><small>MEMBER.COM</small></div></div>
      <div class="header-tools"><span>🔎 服务搜索</span><span>🔔 通知</span><span>🎧 客服</span><b>YX</b><div><strong>会员用户</strong><small>演示账号</small></div></div>
    </header>
    <div class="app-body">
      <aside class="side-nav">
        <button v-for="item in navItems" :key="item.key" :class="{ active: tab === item.key }" :disabled="item.disabled" @click="tab = item.key">
          <span>{{ item.icon }}</span><b>{{ item.label }}</b><small v-if="item.disabled">规划中</small>
        </button>
      </aside>
      <section class="workspace">
        <div class="page-title"><h1>{{ currentTitle }}</h1><p>{{ currentSubtitle }}</p></div>
        <section v-if="message" class="toast">{{ message }}</section>
        <DashboardPortal v-if="tab === 'dashboard'" :plans="plans" @navigate="openUserTask" />
        <UserPortal v-else :plans="plans" :start-mode="userMode" @created="onMemberCreated" />
      </section>
    </div>
  </main>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { api } from './api'
import DashboardPortal from './components/DashboardPortal.vue'
import UserPortal from './components/UserPortal.vue'

const tab = ref('dashboard')
const plans = ref([])
const userMode = ref('apply')
const message = ref('')
const navItems = [
  { key: 'dashboard', label: '控制台', icon: '◌' },
  { key: 'user', label: '会员服务', icon: '♙' }
]
const titles = { user: '会员服务', dashboard: '用户控制台' }
const subtitles = { dashboard: '用户自助入口 · 申请会员、查询状态、查看套餐', user: '当前页面只处理一个会员自助任务' }
const currentTitle = computed(() => titles[tab.value] || '会员服务')
const currentSubtitle = computed(() => subtitles[tab.value] || '用户自助服务')

onMounted(loadPlans)

async function loadPlans() {
  try {
    plans.value = await api.plans()
  } catch (error) {
    notice(error.message)
  }
}

function openUserTask(mode) {
  userMode.value = mode
  tab.value = 'user'
}

function onMemberCreated(member) {
  notice(`申请已提交，会员编号：${member.id}`)
  tab.value = 'dashboard'
}

function notice(text) {
  message.value = text
  setTimeout(() => (message.value = ''), 3500)
}
</script>
