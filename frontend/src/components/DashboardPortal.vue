<template>
  <section class="console-page simple-page">
    <article class="panel task-panel">
      <div class="section-head"><h3>我的会员控制台</h3><small>这里是用户入口，只保留申请、查询和查看套餐三个常用任务。</small></div>
      <div class="console-grid"><button v-for="item in actions" :key="item.mode" @click="open(item.mode)"><b>{{ item.title }}</b><span>{{ item.text }}</span></button></div>
    </article>

    <article class="report-card">
      <div class="report-title"><div><h2>推荐套餐</h2><small>选择套餐后进入会员申请页。</small></div></div>
      <div class="plans-grid"><button v-for="plan in topPlans" :key="plan.id" class="plan-item" @click="open('apply')"><strong>{{ plan.name }}</strong><span>¥{{ price(plan) }}</span><small>{{ plan.duration_days }} 天 · {{ plan.description || '会员权益套餐' }}</small></button></div>
    </article>
  </section>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({ plans: { type: Array, default: () => [] } })
const emit = defineEmits(['navigate'])
const actions = [
  { mode: 'apply', title: '申请会员', text: '填写资料并选择套餐' },
  { mode: 'query', title: '查询状态', text: '用会员编号查看审核进度' },
  { mode: 'plans', title: '查看套餐', text: '对比价格、周期和权益' }
]
const topPlans = computed(() => props.plans.slice(0, 3))
function open(mode) { emit('navigate', mode) }
function price(plan) { return (plan.price_cents / 100).toFixed(2) }
</script>
