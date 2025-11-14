<template>
  <div class="transcription-list">
    <div v-if="transcriptions.length === 0" class="empty-state">
      <p>No tienes transcripciones aún</p>
      <router-link to="/upload" class="btn-primary">Subir tu primer archivo</router-link>
    </div>

    <div v-else class="transcription-grid">
      <div
        v-for="transcription in transcriptions"
        :key="transcription.id"
        class="transcription-card"
      >
        <div class="card-header">
          <h3>{{ transcription.file_name }}</h3>
          <span :class="['status-badge', transcription.status]">
            {{ getStatusText(transcription.status) }}
          </span>
        </div>

        <div class="card-info">
          <p class="info-item">
            <strong>Creado:</strong> {{ formatDate(transcription.created_at) }}
          </p>
          <p v-if="transcription.duration" class="info-item">
            <strong>Duración:</strong> {{ formatDuration(transcription.duration) }}
          </p>
          <p v-if="transcription.credits_used" class="info-item">
            <strong>Créditos usados:</strong> {{ transcription.credits_used }} min
          </p>
          <p class="info-item">
            <strong>Idioma:</strong> {{ transcription.language.toUpperCase() }}
          </p>
        </div>

        <div v-if="transcription.status === 'completed'" class="card-actions">
          <button @click="viewTranscript(transcription)" class="btn-secondary">
            Ver Texto
          </button>
          <button @click="downloadFile(transcription.id, 'txt')" class="btn-secondary">
            📄 TXT
          </button>
          <button @click="downloadFile(transcription.id, 'srt')" class="btn-secondary">
            📝 SRT
          </button>
          <button @click="deleteTranscription(transcription.id)" class="btn-danger">
            🗑️
          </button>
        </div>

        <div v-else-if="transcription.status === 'processing' || transcription.status === 'pending'" class="processing-info">
          <div class="spinner"></div>
          <p>Procesando transcripción...</p>
        </div>

        <div v-else-if="transcription.status === 'failed'" class="error-info">
          <p>❌ Error: {{ transcription.error_message || 'Error desconocido' }}</p>
        </div>
      </div>
    </div>

    <!-- Modal for viewing transcript -->
    <div v-if="showModal" class="modal" @click="closeModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>{{ selectedTranscription?.file_name }}</h3>
          <button @click="closeModal" class="btn-close">✕</button>
        </div>
        <div class="modal-body">
          <textarea
            v-model="editableText"
            class="transcript-textarea"
            :readonly="!isEditing"
          ></textarea>
        </div>
        <div class="modal-footer">
          <button v-if="!isEditing" @click="isEditing = true" class="btn-secondary">
            Editar
          </button>
          <button v-if="isEditing" @click="saveEdit" class="btn-primary">
            Guardar
          </button>
          <button @click="closeModal" class="btn-secondary">
            Cerrar
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, defineProps, defineEmits } from 'vue'
import axios from 'axios'

const props = defineProps({
  transcriptions: {
    type: Array,
    required: true
  }
})

const emit = defineEmits(['refresh'])

const showModal = ref(false)
const selectedTranscription = ref(null)
const editableText = ref('')
const isEditing = ref(false)

function getStatusText(status) {
  const statusMap = {
    pending: 'Pendiente',
    processing: 'Procesando',
    completed: 'Completado',
    failed: 'Fallido'
  }
  return statusMap[status] || status
}

function formatDate(dateString) {
  const date = new Date(dateString)
  return date.toLocaleDateString('es-ES', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function formatDuration(seconds) {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

function viewTranscript(transcription) {
  selectedTranscription.value = transcription
  editableText.value = transcription.transcript_text
  isEditing.value = false
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  selectedTranscription.value = null
  isEditing.value = false
}

async function saveEdit() {
  try {
    await axios.put(`/api/transcriptions/${selectedTranscription.value.id}`, {
      transcript_text: editableText.value
    })
    isEditing.value = false
    emit('refresh')
  } catch (error) {
    alert('Error al guardar cambios')
  }
}

async function downloadFile(id, format) {
  try {
    const response = await axios.get(`/api/transcriptions/${id}/download?format=${format}`, {
      responseType: 'blob'
    })

    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `transcription.${format}`)
    document.body.appendChild(link)
    link.click()
    link.remove()
  } catch (error) {
    alert('Error al descargar archivo')
  }
}

async function deleteTranscription(id) {
  if (!confirm('¿Estás seguro de eliminar esta transcripción?')) return

  try {
    await axios.delete(`/api/transcriptions/${id}`)
    emit('refresh')
  } catch (error) {
    alert('Error al eliminar transcripción')
  }
}
</script>

<style scoped>
.transcription-list {
  margin-top: 2rem;
}

.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  background: white;
  border-radius: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.empty-state p {
  color: #666;
  margin-bottom: 1.5rem;
  font-size: 1.125rem;
}

.transcription-grid {
  display: grid;
  gap: 1.5rem;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
}

.transcription-card {
  background: white;
  border-radius: 0.75rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 1.5rem;
  transition: transform 0.2s, box-shadow 0.2s;
}

.transcription-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: start;
  margin-bottom: 1rem;
  gap: 1rem;
}

.card-header h3 {
  font-size: 1rem;
  color: #333;
  word-break: break-word;
}

.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.75rem;
  font-weight: 600;
  white-space: nowrap;
}

.status-badge.completed {
  background: #d4edda;
  color: #155724;
}

.status-badge.processing,
.status-badge.pending {
  background: #fff3cd;
  color: #856404;
}

.status-badge.failed {
  background: #f8d7da;
  color: #721c24;
}

.card-info {
  margin-bottom: 1rem;
}

.info-item {
  font-size: 0.875rem;
  color: #666;
  margin-bottom: 0.5rem;
}

.info-item strong {
  color: #333;
}

.card-actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.btn-primary,
.btn-secondary,
.btn-danger {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 0.375rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary {
  background: #ff6b35;
  color: white;
}

.btn-primary:hover {
  background: #ff5722;
}

.btn-secondary {
  background: #f5f5f5;
  color: #333;
  border: 1px solid #ddd;
}

.btn-secondary:hover {
  background: #e0e0e0;
}

.btn-danger {
  background: #f8d7da;
  color: #721c24;
}

.btn-danger:hover {
  background: #f5c6cb;
}

.processing-info {
  display: flex;
  align-items: center;
  gap: 1rem;
  color: #856404;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #ff6b35;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-info {
  color: #721c24;
  font-size: 0.875rem;
}

/* Modal styles */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 1rem;
}

.modal-content {
  background: white;
  border-radius: 1rem;
  max-width: 800px;
  width: 100%;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}

.modal-header {
  padding: 1.5rem;
  border-bottom: 1px solid #e0e0e0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.modal-header h3 {
  font-size: 1.25rem;
  color: #333;
}

.btn-close {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: #666;
  padding: 0;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-close:hover {
  color: #333;
}

.modal-body {
  padding: 1.5rem;
  flex: 1;
  overflow-y: auto;
}

.transcript-textarea {
  width: 100%;
  min-height: 400px;
  padding: 1rem;
  border: 2px solid #e0e0e0;
  border-radius: 0.5rem;
  font-family: inherit;
  font-size: 0.9375rem;
  line-height: 1.6;
  resize: vertical;
}

.transcript-textarea:focus {
  outline: none;
  border-color: #ff6b35;
}

.transcript-textarea[readonly] {
  background: #f9f9f9;
}

.modal-footer {
  padding: 1.5rem;
  border-top: 1px solid #e0e0e0;
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}
</style>
