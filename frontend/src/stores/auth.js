import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { supabase } from '@/config/supabase'
import axios from 'axios'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const session = ref(null)
  const loading = ref(true)

  const isAuthenticated = computed(() => !!session.value)

  // Initialize auth state
  async function initialize() {
    loading.value = true
    try {
      const { data: { session: currentSession } } = await supabase.auth.getSession()
      session.value = currentSession

      if (currentSession) {
        await fetchUserData()
      }
    } catch (error) {
      console.error('Error initializing auth:', error)
    } finally {
      loading.value = false
    }

    // Listen for auth changes
    supabase.auth.onAuthStateChange(async (event, newSession) => {
      session.value = newSession
      if (newSession) {
        await fetchUserData()
      } else {
        user.value = null
      }
    })
  }

  // Fetch user data from our backend
  async function fetchUserData() {
    try {
      const token = session.value?.access_token
      if (!token) return

      const response = await axios.get('/api/auth/me', {
        headers: {
          Authorization: `Bearer ${token}`
        }
      })
      user.value = response.data.user
    } catch (error) {
      console.error('Error fetching user data:', error)
    }
  }

  // Sign up
  async function signUp(email, password) {
    const { data, error } = await supabase.auth.signUp({
      email,
      password
    })
    if (error) throw error
    return data
  }

  // Sign in
  async function signIn(email, password) {
    const { data, error } = await supabase.auth.signInWithPassword({
      email,
      password
    })
    if (error) throw error
    session.value = data.session
    await fetchUserData()
    return data
  }

  // Sign out
  async function signOut() {
    const { error } = await supabase.auth.signOut()
    if (error) throw error
    session.value = null
    user.value = null
  }

  // Get access token for API requests
  function getAccessToken() {
    return session.value?.access_token
  }

  return {
    user,
    session,
    loading,
    isAuthenticated,
    initialize,
    fetchUserData,
    signUp,
    signIn,
    signOut,
    getAccessToken
  }
})
