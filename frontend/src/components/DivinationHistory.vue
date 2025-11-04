<script setup>
import { computed } from 'vue'

const props = defineProps({
  records: {
    type: Array,
    default: () => []
  }
})

const formattedRecords = computed(() =>
  props.records.map((item) => ({
    ...item,
    createdAtText: new Date(item.createdAt).toLocaleString('zh-CN', {
      hour12: false
    })
  }))
)
</script>

<template>
  <div class="history">
    <p v-if="!formattedRecords.length" class="placeholder">暂无记录，快来试一试吧。</p>
    <ul v-else class="list">
      <li v-for="record in formattedRecords" :key="record.id" class="item">
        <div class="item__header">
          <h3>{{ record.subject || '未命名事件' }}</h3>
          <time>{{ record.createdAtText }}</time>
        </div>
        <p class="item__hexagram">卦名：<strong>{{ record.hexagramName }}</strong>（动爻：第{{ record.changingLine }}爻）</p>
        <p class="item__summary">{{ record.summary }}</p>
        <p class="item__numbers">数字：{{ record.number1 }}、{{ record.number2 }}、{{ record.number3 }}</p>
      </li>
    </ul>
  </div>
</template>
