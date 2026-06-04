<template>
  <main class="app-shell">
    <header class="global-header">
      <div class="brand"><span>记</span><div><b>JiYiXia</b><small>MEMBER.COM</small></div></div>
      <div class="header-tools"><span>🔎 功能搜索</span><span>🔔 消息</span><span>🎧 客服</span><b>YX</b><div><strong>运营员</strong><small>演示门店</small></div></div>
    </header>
    <div class="app-body">
      <aside class="side-nav">
        <button v-for="item in navItems" :key="item.key" :class="{ active: tab === item.key }" :disabled="item.disabled" @click="tab = item.key">
          <span>{{ item.icon }}</span><b>{{ item.label }}</b><small v-if="item.disabled">规划中</small>
        </button>
      </aside>
      <section class="workspace">
        <div class="page-title"><h1>{{ currentTitle }}</h1><p>简洁企业后台 · 支持后续模块扩展</p></div>
        <section v-if="message" class="toast">{{ message }}</section>
        <AdminPortal v-if="tab === 'admin'" :plans="plans" @refresh-plans="loadPlans" @notice="notice" />
        <UserPortal v-else-if="tab === 'user'" :plans="plans" @created="onMemberCreated" />
        <section v-else class="empty-page">该模块已预留，后续可接入更多业务能力。</section>
      </section>
    </div>
  </main>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { api } from './api'
import AdminPortal from './components/AdminPortal.vue'
import UserPortal from './components/UserPortal.vue'

const tab = ref('admin')
const plans = ref([])
const message = ref('')
const navItems = [
  { key: 'dashboard', label: '总览', icon: '◌', disabled: true },
  { key: 'booking', label: '预约', icon: '◷', disabled: true },
  { key: 'order', label: '订单', icon: '▤', disabled: true },
  { key: 'admin', label: '数据', icon: '◕' },
  { key: 'user', label: '会员', icon: '♙' },
  { key: 'marketing', label: '营销', icon: '⌁', disabled: true },
  { key: 'finance', label: '财务', icon: '￥', disabled: true },
  { key: 'settings', label: '设置', icon: '⚙', disabled: true }
]
const currentTitle = computed(() => tab.value === 'admin' ? '数据（会员版）' : '会员录入')

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
