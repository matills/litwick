<template>
  <div class="dashboard-container">
    <div class="container">
      <div class="dashboard-header">
        <h1>Dashboard</h1>
        <router-link to="/upload" class="btn-primary">
          ➕ Nueva Transcripción
        </router-link>
      </div>

      <div v-if="loading" class="loading">Cargando...</div>

      <div v-else>
        <!-- Stats Cards -->
        <div class="stats-grid">
          <div class="stat-card">
            <div class="stat-icon">📊</div>
            <div class="stat-info">
              <p class="stat-label">Total Transcripciones</p>
              <p class="stat-value">{{ stats.total_transcriptions }}</p>
            </div>
          </div>

          <div class="stat-card">
            <div class="stat-icon">✅</div>
            <div class="stat-info">
              <p class="stat-label">Completadas</p>
              <p class="stat-value">{{ stats.completed_count }}</p>
            </div>
          </div>

          <div class="stat-card">
            <div class="stat-icon">⏱️</div>
            <div class="stat-info">
              <p class="stat-label">Procesando</p>
              <p class="stat-value">{{ stats.processing_count }}</p>
            </div>
          </div>

          <div class="stat-card highlight">
            <div class="stat-icon">🔥</div>
            <div class="stat-info">
              <p class="stat-label">Créditos Disponibles</p>
              <p class="stat-value">{{ stats.credits_remaining }} min</p>
            </div>
          </div>
        </div>

        <!-- Transcriptions List -->
        <div class="transcriptions-section">
          <div class="section-header">
            <h2>Mis Transcripciones</h2>
            <button @click="fetchDashboard" class="btn-refresh">
              🔄 Actualizar
            </button>
          </div>

          <TranscriptionList
            :transcriptions="transcriptions"
            @refresh="fetchDashboard"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import TranscriptionList from '@/components/TranscriptionList.vue'

const loading = ref(true)
const transcriptions = ref([])
const stats = ref({
  total_transcriptions: 0,
  completed_count: 0,
  processing_count: 0,
  failed_count: 0,
  total_minutes_used: 0,
  credits_remaining: 0
})

onMounted(() => {
  fetchDashboard()
  // Auto-refresh every 10 seconds if there are processing transcriptions
  setInterval(() => {
    if (stats.value.processing_count > 0) {
      fetchDashboard()
    }
  }, 10000)
})

async function fetchDashboard() {
  try {
    const response = await axios.get('/api/dashboard/')
    transcriptions.value = response.data.transcriptions || []
    stats.value = response.data.stats || stats.value
  } catch (error) {
    console.error('Error fetching dashboard:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.dashboard-container {
  padding: 2rem 0;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.dashboard-header h1 {
  font-size: 2rem;
  color: #333;
}

.btn-primary {
  background: #ff6b35;
  color: white;
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  text-decoration: none;
  font-weight: 600;
  transition: all 0.2s;
  display: inline-block;
  border: none;
  cursor: pointer;
}

.btn-primary:hover {
  background: #ff5722;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(255, 107, 53, 0.3);
}

.loading {
  text-align: center;
  padding: 4rem;
  color: #666;
  font-size: 1.125rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1.5rem;
  margin-bottom: 3rem;
}

.stat-card {
  background: white;
  border-radius: 0.75rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.stat-card.highlight {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.stat-icon {
  font-size: 2.5rem;
  line-height: 1;
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 0.875rem;
  opacity: 0.8;
  margin-bottom: 0.25rem;
}

.stat-card.highlight .stat-label {
  opacity: 0.9;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  line-height: 1;
}

.transcriptions-section {
  margin-top: 3rem;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.section-header h2 {
  font-size: 1.5rem;
  color: #333;
}

.btn-refresh {
  background: white;
  border: 1px solid #ddd;
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  cursor: pointer;
  font-size: 0.875rem;
  transition: all 0.2s;
}

.btn-refresh:hover {
  background: #f5f5f5;
  border-color: #ff6b35;
}
</style>
