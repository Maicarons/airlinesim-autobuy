<script setup lang="ts">
import { ref, inject, onMounted, onUnmounted } from 'vue'

const t = inject('t') as (key: string, fallback?: string) => string

const logs = ref<string[]>([])
const filter = ref('')
const autoScroll = ref(true)
let eventSource: EventSource | null = null

onMounted(() => {
  try {
    eventSource = new EventSource('/api/logs')
    eventSource.onmessage = (e) => {
      logs.value.push(e.data)
      if (logs.value.length > 1000) {
        logs.value = logs.value.slice(-500)
      }
    }
    eventSource.onerror = () => {
      // SSE not available
    }
  } catch {}
})

onUnmounted(() => {
  if (eventSource) eventSource.close()
})

function filteredLogs() {
  if (!filter.value) return logs.value
  return logs.value.filter(l => l.toLowerCase().includes(filter.value.toLowerCase()))
}

function scrollToBottom(el: HTMLDivElement) {
  if (autoScroll.value) {
    el.scrollTop = el.scrollHeight
  }
}
</script>

<template>
  <div>
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('logs.title') }}</h1>
        <p class="page-subtitle">{{ t('logs.subtitle') }}</p>
      </div>
      <div class="flex gap-2 items-center">
        <div style="position: relative;">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16" style="position: absolute; left: 0.625rem; top: 50%; transform: translateY(-50%); color: var(--text-muted); pointer-events: none;">
            <path d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <input v-model="filter" type="text" class="input" :placeholder="t('logs.filter')" style="padding-left: 2rem; width: 200px;" />
        </div>
        <label class="toggle" style="font-size: 0.8125rem;">
          <input v-model="autoScroll" type="checkbox" />
          <span class="toggle-track"></span>
          <span>{{ t('logs.auto_scroll') }}</span>
        </label>
        <button class="btn btn-sm btn-ghost" @click="logs = []">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14">
            <path d="M6 18L18 6M6 6l12 12" stroke-linecap="round"/>
          </svg>
          {{ t('logs.clear') }}
        </button>
      </div>
    </div>

    <div class="log-container" ref="logContainer" @scroll="scrollToBottom">
      <div v-if="!filteredLogs().length" class="log-empty">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" width="32" height="32">
          <path d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"/>
        </svg>
        <div style="margin-top: 0.5rem; font-size: 0.875rem;">{{ t('logs.empty') }}</div>
        <div style="font-size: 0.75rem; margin-top: 0.25rem;">{{ t('logs.empty_hint') }}</div>
      </div>
      <div
        v-for="(line, i) in filteredLogs()"
        :key="i"
        class="log-line"
        :class="{ 'log-error': line.includes('error') || line.includes('ERROR'), 'log-warn': line.includes('warn') || line.includes('WARN') }"
      >
        {{ line }}
      </div>
    </div>
  </div>
</template>

<style scoped>
.log-container {
  background: #080e1a;
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  font-family: var(--font-mono);
  font-size: 0.8125rem;
  line-height: 1.6;
  max-height: calc(100vh - 220px);
  overflow-y: auto;
  padding: 0.75rem;
}

.log-empty {
  text-align: center;
  padding: 3rem 1rem;
  color: var(--text-muted);
}

.log-line {
  padding: 0.125rem 0.5rem;
  color: var(--text-secondary);
  white-space: pre-wrap;
  word-break: break-all;
  border-radius: 2px;
  transition: background 0.1s ease;
}
.log-line:hover {
  background: rgba(59, 130, 246, 0.04);
}
.log-line.log-error {
  color: var(--danger);
}
.log-line.log-warn {
  color: var(--warning);
}
</style>