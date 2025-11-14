import axios from 'axios'
import { supabase } from './supabase'

const api = axios.create({
  baseURL: '/api',
})

// Add auth token to requests
api.interceptors.request.use(async (config) => {
  const { data: { session } } = await supabase.auth.getSession()
  if (session?.access_token) {
    config.headers.Authorization = `Bearer ${session.access_token}`
  }
  return config
})

export default api

// API methods
export const apiClient = {
  // Dashboard
  getDashboard: async () => {
    const response = await api.get('/dashboard/')
    return response.data
  },

  // Upload
  uploadFile: async (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    const response = await api.post('/upload/', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data
  },

  // Transcriptions
  getTranscriptions: async () => {
    const response = await api.get('/transcriptions/')
    return response.data
  },

  getTranscription: async (id: string) => {
    const response = await api.get(`/transcription/${id}`)
    return response.data
  },

  deleteTranscription: async (id: string) => {
    const response = await api.delete(`/transcription/${id}`)
    return response.data
  },

  downloadTranscription: async (id: string, format: string) => {
    const response = await api.get(`/transcription/${id}/download?format=${format}`, {
      responseType: 'blob',
    })
    return response.data
  },

  // Auth
  getCurrentUser: async () => {
    const { data: { session } } = await supabase.auth.getSession()
    return session?.user || null
  },

  signOut: async () => {
    await supabase.auth.signOut()
  },
}
