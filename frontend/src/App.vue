<template>
  <main class="app-shell">
    <aside class="sidebar">
      <div class="brand"><span>记</span><div><b>JiYiXia</b><small>Member Management</small></div></div>
      <div class="nav-group">
        <p>核心模块</p>
        <button :class="{ active: tab === 'admin' }" @click="tab = 'admin'"><span>概览</span><small>会员 / 套餐 / 导出</small></button>
        <button :class="{ active: tab === 'user' }" @click="tab = 'user'"><span>录入</span><small>申请 / 查询</small></button>
      </div>
      <div class="nav-group muted">
        <p>后续扩展</p>
        <button disabled><span>财务</span><small>收款与发票</small></button>
        <button disabled><span>营销</span><small>活动与触达</small></button>
        <button disabled><span>报表</span><small>经营分析</small></button>
      </div>
      <div class="sidebar-note">Demo 模式已开启；部署后端后可接入真实 PostgreSQL / Redis。</div>
    </aside>

    <section class="content">
      <header class="topbar">
        <div>
          <p class="caption">SaaS 企业后台</p>
          <h1>{{ tab === 'admin' ? '会员运营中心' : '会员录入中心' }}</h1>
        </div>
        <div class="top-actions">
          <span>{{ plans.length }} 个套餐</span>
          <span>CSV 导出</span>
          <b>YX</b>
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
  notice(`申请已提交，会员编号：${member.id}`)
  tab.value = 'admin'
}

function notice(text) {
  message.value = text
  setTimeout(() => (message.value = ''), 3500)
}
</script>
