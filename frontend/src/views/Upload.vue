<template>
  <div class="upload-container">
    <div class="container">
      <div class="upload-card">
        <h2>Subir Archivo de Audio/Video</h2>
        <p class="subtitle">Sube tu archivo y obtendremos la transcripción automáticamente</p>

        <form @submit.prevent="handleUpload" class="upload-form">
          <div class="file-input-wrapper">
            <input
              type="file"
              id="file"
              ref="fileInput"
              @change="handleFileChange"
              accept="audio/*,video/*,.mp3,.mp4,.wav,.m4a,.flac,.ogg,.webm,.avi,.mov,.mkv"
              required
            />
            <label for="file" class="file-label">
              <span v-if="!file" class="file-placeholder">
                📁 Selecciona un archivo
              </span>
              <span v-else class="file-selected">
                {{ file.name }} ({{ formatFileSize(file.size) }})
              </span>
            </label>
          </div>

          <div class="form-group">
            <label>Idioma</label>
            <select v-model="language">
              <option value="es">Español</option>
              <option value="en">Inglés</option>
              <option value="fr">Francés</option>
              <option value="de">Alemán</option>
              <option value="it">Italiano</option>
              <option value="pt">Portugués</option>
            </select>
          </div>

          <div v-if="error" class="error-message">
            {{ error }}
          </div>

          <div v-if="uploading" class="progress-bar">
            <div class="progress-fill" :style="{ width: uploadProgress + '%' }"></div>
          </div>

          <button type="submit" class="btn-primary" :disabled="uploading || !file">
            {{ uploading ? 'Subiendo...' : 'Subir y Procesar' }}
          </button>
        </form>

        <div class="info-section">
          <h3>Formatos soportados</h3>
          <p>MP3, MP4, WAV, M4A, FLAC, OGG, WebM, AVI, MOV, MKV</p>
          <p class="small">Tamaño máximo: 500MB</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const fileInput = ref(null)
const file = ref(null)
const language = ref('es')
const uploading = ref(false)
const uploadProgress = ref(0)
const error = ref('')

function handleFileChange(event) {
  const selectedFile = event.target.files[0]
  if (selectedFile) {
    file.value = selectedFile
    error.value = ''
  }
}

function formatFileSize(bytes) {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

async function handleUpload() {
  if (!file.value) return

  uploading.value = true
  error.value = ''
  uploadProgress.value = 0

  try {
    const formData = new FormData()
    formData.append('file', file.value)
    formData.append('language', language.value)

    const response = await axios.post('/api/transcriptions/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      },
      onUploadProgress: (progressEvent) => {
        uploadProgress.value = Math.round(
          (progressEvent.loaded * 100) / progressEvent.total
        )
      }
    })

    // Start processing
    await axios.post(`/api/transcriptions/${response.data.transcription.id}/process`)

    // Redirect to dashboard
    router.push('/dashboard')
  } catch (err) {
    error.value = err.response?.data?.error || 'Error al subir el archivo'
  } finally {
    uploading.value = false
  }
}
</script>

<style scoped>
.upload-container {
  padding: 2rem 0;
}

.upload-card {
  background: white;
  border-radius: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  padding: 2rem;
  max-width: 600px;
  margin: 0 auto;
}

.upload-card h2 {
  font-size: 1.75rem;
  color: #333;
  margin-bottom: 0.5rem;
}

.subtitle {
  color: #666;
  margin-bottom: 2rem;
}

.upload-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.file-input-wrapper {
  position: relative;
}

.file-input-wrapper input[type="file"] {
  position: absolute;
  opacity: 0;
  width: 0;
  height: 0;
}

.file-label {
  display: block;
  padding: 3rem 2rem;
  border: 3px dashed #ddd;
  border-radius: 0.5rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
}

.file-label:hover {
  border-color: #ff6b35;
  background: #fff5f2;
}

.file-placeholder {
  color: #999;
  font-size: 1.125rem;
}

.file-selected {
  color: #ff6b35;
  font-weight: 500;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group label {
  font-weight: 500;
  color: #333;
  font-size: 0.875rem;
}

.form-group select {
  padding: 0.75rem 1rem;
  border: 2px solid #e0e0e0;
  border-radius: 0.5rem;
  font-size: 1rem;
  cursor: pointer;
}

.form-group select:focus {
  outline: none;
  border-color: #ff6b35;
}

.error-message {
  background: #fee;
  color: #c33;
  padding: 0.75rem;
  border-radius: 0.5rem;
  font-size: 0.875rem;
}

.progress-bar {
  height: 8px;
  background: #e0e0e0;
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: #ff6b35;
  transition: width 0.3s;
}

.btn-primary {
  background: #ff6b35;
  color: white;
  padding: 1rem;
  border: none;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}

.btn-primary:hover:not(:disabled) {
  background: #ff5722;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(255, 107, 53, 0.3);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.info-section {
  margin-top: 2rem;
  padding-top: 2rem;
  border-top: 1px solid #e0e0e0;
  text-align: center;
}

.info-section h3 {
  font-size: 1rem;
  color: #333;
  margin-bottom: 0.5rem;
}

.info-section p {
  color: #666;
  font-size: 0.875rem;
}

.info-section .small {
  font-size: 0.75rem;
  color: #999;
  margin-top: 0.25rem;
}
</style>
