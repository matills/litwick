import axios from 'axios'
import { supabase } from './supabase'

const api = axios.create({
  baseURL: '/api',
})

api.interceptors.request.use(async (config) => {
  const { data: { session } } = await supabase.auth.getSession()
  if (session?.access_token) {
    config.headers.Authorization = `Bearer ${session.access_token}`
  }
  return config
})

export default api

export const apiClient = {
  getDashboard: async () => {
    const response = await api.get('/dashboard/')
    return response.data
  },

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

  getTranscriptions: async () => {
    const response = await api.get('/transcriptions/')
    return response.data
  },

  getTranscription: async (id: string) => {
    const response = await api.get(`/transcriptions/${id}`)
    return response.data
  },

  deleteTranscription: async (id: string) => {
    const response = await api.delete(`/transcriptions/${id}`)
    return response.data
  },

  downloadTranscription: async (id: string, format: string) => {
    const response = await api.get(`/transcriptions/${id}/download?format=${format}`, {
      responseType: 'blob',
    })
    return response.data
  },

  getCurrentUser: async () => {
    const { data: { session } } = await supabase.auth.getSession()
    return session?.user || null
  },

  signOut: async () => {
    await supabase.auth.signOut()
  },

  updateSettings: async (settings: {
    default_language?: string
    default_export_format?: string
    include_timestamps?: boolean
    detect_speakers?: boolean
    email_notifications?: boolean
    promotional_emails?: boolean
  }) => {
    const response = await api.put('/auth/settings', settings)
    return response.data
  },

  getCreditPackages: async () => {
    const response = await api.get('/payments/packages')
    return response.data
  },

  createPayment: async (packageId: string) => {
    const response = await api.post('/payments/create', { package_id: packageId })
    return response.data
  },

  getPaymentHistory: async () => {
    const response = await api.get('/payments/history')
    return response.data
  },

  processPaymentSuccess: async (params: URLSearchParams) => {
    const response = await api.get(`/payments/success?${params.toString()}`)
    return response.data
  },
}
