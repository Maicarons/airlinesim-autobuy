<script setup lang="ts">
import { ref, inject, onMounted, onUnmounted } from 'vue'
import { getStatus, startEngine, stopEngine } from '../api'

const t = inject('t') as (key: string, fallback?: string) => string

interface Status {
  running: boolean
  start_time: string
  scan_count: number
  found_count: number
  bought_count: number
  failed_count: number
  last_scan_time: string
  last_error: string
}

const status = ref<Status>({
  running: false,
  start_time: '',
  scan_count: 0,
  found_count: 0,
  bought_count: 0,
  failed_count: 0,
  last_scan_time: '',
  last_error: '',
})
const loading = ref(false)
const error = ref('')
let interval: number | undefined

async function fetchStatus() {
  try {
    status.value = await getStatus()
  } catch {
    // server not ready
  }
}

async function handleStart() {
  loading.value = true
  error.value = ''
  try {
    await startEngine()
    await fetchStatus()
  } catch (e: any) {
    error.value = e.message
  }
  loading.value = false
}

async function handleStop() {
  loading.value = true
  error.value = ''
  try {
    await stopEngine()
    await fetchStatus()
  } catch (e: any) {
    error.value = e.message
  }
  loading.value = false
}

function formatTime(t: string): string {
  if (!t) return '—'
  const d = new Date(t)
  return d.toLocaleTimeString()
}

onMounted(() => {
  fetchStatus()
  interval = window.setInterval(fetchStatus, 3000)
})

onUnmounted(() => {
  if (interval) clearInterval(interval)
})
</script>

<template>
  <div>
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('dashboard.title') }}</h1>
        <p class="page-subtitle">{{ t('dashboard.subtitle') }}</p>
      </div>
      <div class="flex gap-2">
        <button
          v-if="!status.running"
          class="btn btn-success"
          :disabled="loading"
          @click="handleStart"
        >
          <svg viewBox="0 0 24 24" fill="currentColor" width="16" height="16">
            <path d="M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.348a1.125 1.125 0 010 1.971l-11.54 6.347a1.125 1.125 0 01-1.667-.985V5.653z"/>
          </svg>
          {{ t('dashboard.start') }}
        </button>
        <button
          v-else
          class="btn btn-danger"
          :disabled="loading"
          @click="handleStop"
        >
          <svg viewBox="0 0 24 24" fill="currentColor" width="16" height="16">
            <path d="M5.25 7.5A2.25 2.25 0 017.5 5.25h9a2.25 2.25 0 012.25 2.25v9a2.25 2.25 0 01-2.25 2.25h-9A2.25 2.25 0 015.25 16.5v-9z"/>
          </svg>
          {{ t('dashboard.stop') }}
        </button>
      </div>
    </div>

    <!-- Engine Status Banner -->
    <div :class="['card', 'engine-banner', status.running ? 'engine-running' : 'engine-stopped']">
      <div class="engine-banner-left">
        <div :class="['engine-indicator', status.running ? 'running' : 'stopped']">
          <svg viewBox="0 0 24 24" fill="currentColor" width="20" height="20">
            <path v-if="status.running" d="M5.25 5.653c0-.856.917-1.398 1.667-.986l11.54 6.348a1.125 1.125 0 010 1.971l-11.54 6.347a1.125 1.125 0 01-1.667-.985V5.653z"/>
            <path v-else d="M5.25 7.5A2.25 2.25 0 017.5 5.25h9a2.25 2.25 0 012.25 2.25v9a2.25 2.25 0 01-2.25 2.25h-9A2.25 2.25 0 015.25 16.5v-9z"/>
          </svg>
        </div>
        <div>
          <div class="engine-status-text">
            {{ t('dashboard.engine') }} <strong>{{ status.running ? t('dashboard.running') : t('dashboard.stopped') }}</strong>
          </div>
          <div class="engine-status-detail">
            {{ status.running ? t('dashboard.monitoring') : t('dashboard.click_start') }}
          </div>
        </div>
      </div>
      <div v-if="status.running && status.start_time" class="engine-time">
        {{ t('dashboard.started') }} {{ formatTime(status.start_time) }}
      </div>
    </div>

    <!-- Stats -->
    <div class="stats-grid">
      <div class="card stat-card">
        <div class="stat-value" style="color: var(--text-accent)">{{ status.scan_count }}</div>
        <div class="stat-label">{{ t('dashboard.scans') }}</div>
      </div>
      <div class="card stat-card">
        <div class="stat-value" style="color: var(--warning)">{{ status.found_count }}</div>
        <div class="stat-label">{{ t('dashboard.found') }}</div>
      </div>
      <div class="card stat-card">
        <div class="stat-value" style="color: var(--success)">{{ status.bought_count }}</div>
        <div class="stat-label">{{ t('dashboard.bought') }}</div>
      </div>
      <div class="card stat-card">
        <div class="stat-value" style="color: var(--danger)">{{ status.failed_count }}</div>
        <div class="stat-label">{{ t('dashboard.failed') }}</div>
      </div>
    </div>

    <!-- Error -->
    <div v-if="error" class="card" style="border-color: rgba(239,68,68,0.3);">
      <div class="flex items-center gap-2" style="color: var(--danger);">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="18" height="18">
          <path d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <span style="font-size: 0.875rem;">{{ error }}</span>
      </div>
    </div>

    <!-- Last Error from Engine -->
    <div v-if="status.last_error" class="card" style="margin-top: 1rem; border-color: rgba(239,68,68,0.2);">
      <div class="card-header">
        <span class="card-title">{{ t('dashboard.last_error') }}</span>
      </div>
      <p style="font-size: 0.875rem; color: var(--danger); font-family: var(--font-mono);">{{ status.last_error }}</p>
    </div>

    <!-- Last Scan Info -->
    <div v-if="status.last_scan_time" class="card" style="margin-top: 1rem;">
      <div class="card-header">
        <span class="card-title">{{ t('dashboard.activity') }}</span>
      </div>
      <div style="font-size: 0.8125rem; color: var(--text-muted);">
        {{ t('dashboard.last_scan') }}: {{ formatTime(status.last_scan_time) }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.engine-banner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding: 1rem 1.25rem;
}
.engine-running {
  border-color: rgba(34, 197, 94, 0.25);
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.05), transparent);
}
.engine-stopped {
  border-color: var(--border);
}

.engine-banner-left {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.engine-indicator {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.engine-indicator.running {
  background: var(--success-bg);
  color: var(--success);
}
.engine-indicator.stopped {
  background: var(--bg-card-hover);
  color: var(--text-muted);
}

.engine-status-text {
  font-size: 0.9375rem;
}
.engine-status-detail {
  font-size: 0.8125rem;
  color: var(--text-muted);
  margin-top: 0.125rem;
}

.engine-time {
  font-size: 0.75rem;
  color: var(--text-muted);
  font-family: var(--font-mono);
}
</style>