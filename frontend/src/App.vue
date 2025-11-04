<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import DivinationHistory from './components/DivinationHistory.vue'

const loading = ref(false)
const errorMessage = ref('')
const records = ref([])
const form = ref({
  subject: '',
  number1: 3,
  number2: 5,
  number3: 7
})

const resetForm = () => {
  form.value = {
    subject: '',
    number1: 3,
    number2: 5,
    number3: 7
  }
}

const loadRecords = async () => {
  try {
    const { data } = await axios.get('/api/divination')
    records.value = data
  } catch (error) {
    console.error(error)
    errorMessage.value = '无法加载历史记录，请稍后再试。'
  }
}

const submit = async () => {
  errorMessage.value = ''
  loading.value = true
  try {
    const payload = {
      subject: form.value.subject.trim(),
      number1: Number(form.value.number1),
      number2: Number(form.value.number2),
      number3: Number(form.value.number3)
    }
    const { data } = await axios.post('/api/divination', payload)
    records.value = [data, ...records.value]
    resetForm()
  } catch (error) {
    if (error.response?.data?.error) {
      errorMessage.value = error.response.data.error
    } else {
      errorMessage.value = '提交失败，请检查网络或稍后再试。'
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadRecords()
})
</script>

<template>
  <main class="layout">
    <header class="hero">
      <h1>梅花易数占卜系统</h1>
      <p>根据数字演算六十四卦，辅助您洞察当下形势。</p>
    </header>

    <section class="card">
      <h2>发起占卜</h2>
      <form class="form" @submit.prevent="submit">
        <label class="field">
          <span>事件/问题</span>
          <input
            v-model="form.subject"
            type="text"
            placeholder="请描述您关心的事情"
            required
          />
        </label>
        <div class="field-group">
          <label class="field">
            <span>数字一</span>
            <input v-model.number="form.number1" type="number" min="0" required />
          </label>
          <label class="field">
            <span>数字二</span>
            <input v-model.number="form.number2" type="number" min="0" required />
          </label>
          <label class="field">
            <span>数字三</span>
            <input v-model.number="form.number3" type="number" min="0" required />
          </label>
        </div>
        <p class="hint">可根据梅花易数的时间、方位或起卦方式任选三个数字。</p>
        <button class="submit" type="submit" :disabled="loading">
          {{ loading ? '计算中...' : '开始占卜' }}
        </button>
        <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
      </form>
    </section>

    <section class="card">
      <h2>占卜历史</h2>
      <DivinationHistory :records="records" />
    </section>
  </main>
</template>
