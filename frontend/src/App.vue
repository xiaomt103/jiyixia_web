<template>
  <main class="suite-shell">
    <aside class="nav-rail">
      <div class="brand-lockup"><span>记</span><div><b>JiYiXia</b><small>MemberOps Suite</small></div></div>
      <nav class="nav-stack">
        <button :class="{ active: tab === 'admin' }" @click="tab = 'admin'"><i>01</i><span>运营总览</span></button>
        <button :class="{ active: tab === 'user' }" @click="tab = 'user'"><i>02</i><span>会员录入</span></button>
      </nav>
      <div class="rail-insight">
        <p>今日工作流</p>
        <strong>申请 → 审核 → 激活 → 导出</strong>
        <small>Vercel Demo / API Ready</small>
      </div>
    </aside>

    <section class="console">
      <header class="command-bar">
        <div>
          <p class="overline">SaaS Enterprise Console</p>
          <h1>{{ tab === 'admin' ? '会员运营指挥台' : '会员开户注册台' }}</h1>
        </div>
        <div class="command-actions">
          <span class="status-chip">Live Demo</span>
          <span class="shortcut">⌘ K / Search</span>
          <span class="user-chip">YX</span>
        </div>
      </header>

      <section v-if="message" class="toast">{{ message }}</section>
      <AdminPortal v-if="tab === 'admin'" :plans="plans" @refresh-plans="loadPlans" @notice="notice" />
      <UserPortal v-else :plans="plans" @created="onMemberCreated" />
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
  notice(`申请已进入审核队列，会员编号：${member.id}`)
  tab.value = 'admin'
}

function notice(text) {
  message.value = text
  setTimeout(() => (message.value = ''), 3500)
}
</script>
