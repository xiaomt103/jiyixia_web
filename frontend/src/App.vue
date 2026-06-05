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
        <AdminPortal v-else-if="tab === 'admin'" :plans="plans" @refresh-plans="loadPlans" @notice="notice" />
        <UserPortal v-else-if="tab === 'user'" :plans="plans" :start-mode="userMode" @created="onMemberCreated" />
        <OperationsPortal v-else :module="tab" @notice="notice" />
      </section>
    </div>
  </main>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { api } from './api'
import DashboardPortal from './components/DashboardPortal.vue'
import AdminPortal from './components/AdminPortal.vue'
import UserPortal from './components/UserPortal.vue'
import OperationsPortal from './components/OperationsPortal.vue'

const tab = ref('dashboard')
const plans = ref([])
const userMode = ref('apply')
const message = ref('')
const navItems = [
  { key: 'dashboard', label: '控制台', icon: '◌' },
  { key: 'booking', label: '预约', icon: '◷' },
  { key: 'order', label: '订单', icon: '▤' },
  { key: 'admin', label: '数据', icon: '◕' },
  { key: 'user', label: '会员', icon: '♙' },
  { key: 'marketing', label: '营销', icon: '⌁' },
  { key: 'finance', label: '财务', icon: '￥' },
  { key: 'settings', label: '设置', icon: '⚙' }
]
const titles = { admin: '后台管理', user: '会员服务', dashboard: '用户控制台', booking: '预约', order: '订单', marketing: '营销', finance: '财务', settings: '设置' }
const subtitles = { dashboard: '用户自助入口 · 申请会员、查询状态、查看套餐', user: '当前页面只处理一个会员自助任务', admin: '后台审核和套餐管理入口' }
const currentTitle = computed(() => titles[tab.value] || '业务模块')
const currentSubtitle = computed(() => subtitles[tab.value] || '业务模块管理 · 筛选、列表、创建和导出')

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
