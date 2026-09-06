<script setup lang="ts">
import { ref, onMounted, onUnmounted, provide } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from './i18n'

const { t, setLocale, currentLocale, localeNames } = useI18n()
provide('t', t)

const router = useRouter()
const route = useRoute()
const engineRunning = ref(false)

interface NavItem {
  path: string
  label: string
  icon: string
}

const navItems: NavItem[] = [
  { path: '/', label: t('nav.dashboard'), icon: 'dashboard' },
  { path: '/rules', label: t('nav.rules'), icon: 'rules' },
  { path: '/settings', label: t('nav.settings'), icon: 'settings' },
  { path: '/logs', label: t('nav.logs'), icon: 'logs' },
]

async function fetchStatus() {
  try {
    const res = await fetch('/api/status')
    const data = await res.json()
    engineRunning.value = data.running
  } catch {}
}

let interval: number | undefined

onMounted(() => {
  fetchStatus()
  interval = window.setInterval(fetchStatus, 5000)
})

onUnmounted(() => {
  if (interval) clearInterval(interval)
})
</script>

<template>
  <div class="app-layout">
    <!-- Sidebar -->
    <aside class="sidebar">
      <div class="sidebar-header">
        <svg class="sidebar-logo" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M3.75 15.75l4.5-4.5 3.75 3.75 4.5-4.5 4.5 4.5" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M21 12V5.25H3v6.75" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M3 18.75h18" stroke-linecap="round"/>
        </svg>
        <div>
          <div class="sidebar-title">{{ t('app.title') }}</div>
          <div class="sidebar-version">{{ t('app.version') }}</div>
        </div>
      </div>

      <nav class="sidebar-nav">
        <button
          v-for="item in navItems"
          :key="item.path"
          :class="['nav-item', { active: route.path === item.path }]"
          @click="router.push(item.path)"
        >
          <!-- Dashboard icon -->
          <svg v-if="item.icon === 'dashboard'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="nav-icon">
            <path d="M3.75 6A2.25 2.25 0 016 3.75h2.25A2.25 2.25 0 0110.5 6v2.25a2.25 2.25 0 01-2.25 2.25H6a2.25 2.25 0 01-2.25-2.25V6zM3.75 15.75A2.25 2.25 0 016 13.5h2.25a2.25 2.25 0 012.25 2.25V18a2.25 2.25 0 01-2.25 2.25H6A2.25 2.25 0 013.75 18v-2.25zM13.5 6a2.25 2.25 0 012.25-2.25H18A2.25 2.25 0 0120.25 6v2.25A2.25 2.25 0 0118 10.5h-2.25a2.25 2.25 0 01-2.25-2.25V6zM13.5 15.75a2.25 2.25 0 012.25-2.25H18a2.25 2.25 0 012.25 2.25V18A2.25 2.25 0 0118 20.25h-2.25A2.25 2.25 0 0113.5 18v-2.25z"/>
          </svg>
          <!-- Rules icon -->
          <svg v-else-if="item.icon === 'rules'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="nav-icon">
            <path d="M9 12h3.75M9 15h3.75M9 18h3.75m3 .75H18a2.25 2.25 0 002.25-2.25V6.108c0-1.135-.845-2.098-1.976-2.192a48.424 48.424 0 00-1.123-.08m-5.801 0c-.065.21-.1.433-.1.664 0 .414.336.75.75.75h4.5a.75.75 0 00.75-.75 2.25 2.25 0 00-.1-.664m-5.8 0A2.251 2.251 0 0113.5 2.25H15a2.25 2.25 0 012.15 1.586m-5.8 0c-.376.023-.75.05-1.124.08C9.095 4.01 8.25 4.973 8.25 6.108V8.25m0 0H4.875c-.621 0-1.125.504-1.125 1.125v11.25c0 .621.504 1.125 1.125 1.125h9.75c.621 0 1.125-.504 1.125-1.125V9.375c0-.621-.504-1.125-1.125-1.125H8.25zM6.75 12h.008v.008H6.75V12zm0 3h.008v.008H6.75V15zm0 3h.008v.008H6.75V18z"/>
          </svg>
          <!-- Settings icon -->
          <svg v-else-if="item.icon === 'settings'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="nav-icon">
            <path d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.324.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 011.37.49l1.296 2.247a1.125 1.125 0 01-.26 1.431l-1.003.827c-.293.24-.438.613-.431.992a6.759 6.759 0 010 .255c-.007.378.138.75.43.99l1.005.828c.424.35.534.954.26 1.43l-1.298 2.247a1.125 1.125 0 01-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.57 6.57 0 01-.22.128c-.331.183-.581.495-.644.869l-.213 1.28c-.09.543-.56.941-1.11.941h-2.594c-.55 0-1.02-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 01-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 01-1.369-.49l-1.297-2.247a1.125 1.125 0 01.26-1.431l1.004-.827c.292-.24.437-.613.43-.992a6.932 6.932 0 010-.255c.007-.378-.138-.75-.43-.99l-1.004-.828a1.125 1.125 0 01-.26-1.43l1.297-2.247a1.125 1.125 0 011.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.087.22-.128.332-.183.582-.495.644-.869l.214-1.281z"/>
            <path d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
          </svg>
          <!-- Logs icon -->
          <svg v-else-if="item.icon === 'logs'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" class="nav-icon">
            <path d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"/>
          </svg>
          <span class="nav-label">{{ item.label }}</span>
        </button>
      </nav>

      <!-- Locale Switcher -->
      <div class="locale-switcher">
        <select v-model="currentLocale" @change="setLocale(currentLocale)" class="locale-select">
          <option v-for="loc in localeNames" :key="loc.value" :value="loc.value">{{ loc.label }}</option>
        </select>
      </div>

      <div class="sidebar-footer">
        <div :class="['status-dot', engineRunning ? 'running' : 'stopped']"></div>
        <span class="status-text">{{ engineRunning ? t('nav.engine.running') : t('nav.engine.stopped') }}</span>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<style scoped>
.app-layout {
  display: flex;
  height: 100vh;
  background: var(--bg-primary);
}

.sidebar {
  width: var(--sidebar-width);
  background: var(--bg-sidebar);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  overflow: hidden;
}

.sidebar-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1.25rem 1.25rem 1rem;
  border-bottom: 1px solid var(--border);
}

.sidebar-logo {
  width: 28px;
  height: 28px;
  color: var(--accent);
  flex-shrink: 0;
}

.sidebar-title {
  font-size: 1rem;
  font-weight: 700;
  letter-spacing: -0.01em;
}

.sidebar-version {
  font-size: 0.6875rem;
  color: var(--text-muted);
}

.sidebar-nav {
  flex: 1;
  padding: 0.75rem 0.5rem;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  padding: 0.5625rem 0.75rem;
  border: none;
  border-radius: var(--radius-sm);
  background: transparent;
  color: var(--text-secondary);
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.15s ease;
  text-align: left;
  width: 100%;
  min-height: 38px;
}

.nav-item:hover {
  background: var(--bg-card);
  color: var(--text-primary);
}

.nav-item.active {
  background: var(--accent-light);
  color: var(--accent);
}

.nav-item:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}

.nav-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

.nav-label {
  font-weight: 500;
}

.locale-switcher {
  padding: 0.5rem 0.75rem;
  border-top: 1px solid var(--border);
}

.locale-select {
  width: 100%;
  padding: 0.3rem 0.5rem;
  background: var(--bg-input);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  font-size: 0.75rem;
  cursor: pointer;
  outline: none;
}
.locale-select:focus {
  border-color: var(--accent);
}

.sidebar-footer {
  padding: 0.875rem 1.25rem;
  border-top: 1px solid var(--border);
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}
.status-dot.running {
  background: var(--success);
  box-shadow: 0 0 6px rgba(34, 197, 94, 0.5);
}
.status-dot.stopped {
  background: var(--text-muted);
}

.status-text {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.main-content {
  flex: 1;
  overflow-y: auto;
  padding: 2rem;
  max-width: 960px;
}
</style>