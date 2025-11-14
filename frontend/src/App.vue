<template>
  <div id="app">
    <nav v-if="authStore.isAuthenticated" class="navbar">
      <div class="container">
        <div class="nav-brand">
          <h1>🔥 Litwick</h1>
        </div>
        <div class="nav-menu">
          <router-link to="/dashboard">Dashboard</router-link>
          <router-link to="/upload">Subir Archivo</router-link>
          <div class="user-info">
            <span>{{ authStore.user?.email }}</span>
            <span class="credits">{{ authStore.user?.credits_remaining || 0 }} min</span>
            <button @click="handleLogout" class="btn-logout">Salir</button>
          </div>
        </div>
      </div>
    </nav>

    <main class="main-content">
      <router-view v-if="!authStore.loading" />
      <div v-else class="loading">Cargando...</div>
    </main>
  </div>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const router = useRouter()

async function handleLogout() {
  await authStore.signOut()
  router.push('/login')
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
  background: #f5f5f5;
  color: #333;
}

#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.navbar {
  background: white;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 1rem 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.nav-brand h1 {
  font-size: 1.5rem;
  color: #ff6b35;
}

.nav-menu {
  display: flex;
  gap: 1.5rem;
  align-items: center;
}

.nav-menu a {
  text-decoration: none;
  color: #666;
  font-weight: 500;
  transition: color 0.2s;
}

.nav-menu a:hover,
.nav-menu a.router-link-active {
  color: #ff6b35;
}

.user-info {
  display: flex;
  gap: 1rem;
  align-items: center;
  padding-left: 1rem;
  border-left: 1px solid #ddd;
}

.credits {
  background: #ff6b35;
  color: white;
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.875rem;
  font-weight: 600;
}

.btn-logout {
  background: transparent;
  border: 1px solid #ddd;
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  cursor: pointer;
  font-size: 0.875rem;
  transition: all 0.2s;
}

.btn-logout:hover {
  background: #f5f5f5;
  border-color: #ff6b35;
  color: #ff6b35;
}

.main-content {
  flex: 1;
  padding: 2rem 0;
}

.loading {
  text-align: center;
  padding: 4rem;
  font-size: 1.25rem;
  color: #666;
}
</style>
