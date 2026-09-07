<script setup lang="ts">
import { ref, inject, onMounted } from 'vue'
import { getConfig, updateConfig } from '../api'

const t = inject('t') as (key: string, fallback?: string) => string

const config = ref<any>({})
const loading = ref(false)
const saved = ref(false)
const error = ref('')

async function fetchConfig() {
  try {
    config.value = await getConfig()
    if (!config.value.servers) config.value.servers = []
    if (!config.value.auths) config.value.auths = []
  } catch {}
}

function addServer() {
  const idx = config.value.servers.length + 1
  config.value.servers.push({ host: 'free' + idx, base_url: 'https://free' + idx + '.airlinesim.aero', companies: [] })
}
function removeServer(i: number) {
  config.value.servers.splice(i, 1)
}

function addAuth() {
  config.value.auths.push({ username: '', password: '', session_file: 'session.json' })
}
function removeAuth(i: number) {
  config.value.auths.splice(i, 1)
}

async function save() {
  loading.value = true
  saved.value = false
  error.value = ''
  try {
    await updateConfig({
      servers: config.value.servers,
      auths: config.value.auths,
      monitor: config.value.monitor,
      notifier: config.value.notifier,
      webui: config.value.webui,
    })
    saved.value = true
  } catch (e: any) {
    error.value = e.message
  }
  loading.value = false
}

onMounted(fetchConfig)
</script>

<template>
  <div>
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('settings.title') }}</h1>
        <p class="page-subtitle">{{ t('settings.subtitle') }}</p>
      </div>
    </div>

    <!-- Server Connections -->
    <div class="card" style="margin-bottom: 1rem;">
      <div class="card-header">
        <span class="card-title">{{ t('settings.servers') }}</span>
        <button class="btn btn-sm" @click="addServer">{{ t('settings.add_server') }}</button>
      </div>
      <div v-for="(sv, i) in config.servers" :key="i" style="padding: 0.75rem 0; border-bottom: 1px solid var(--border);">
        <div class="form-row">
          <div class="form-group">
            <label class="form-label">{{ t('settings.server_host') }}</label>
            <input v-model="sv.host" type="text" class="input" :placeholder="'free' + (i+1)" />
          </div>
          <div class="form-group" style="flex: 2;">
            <label class="form-label">{{ t('settings.server_url') }}</label>
            <input v-model="sv.base_url" type="text" class="input" :placeholder="'https://free' + (i+1) + '.airlinesim.aero'" />
          </div>
          <div class="form-group" style="flex: 1;">
            <label class="form-label">{{ t('settings.server_companies') }}</label>
            <input :value="(sv.companies || []).join(', ')" @input="sv.companies = ($event.target as HTMLInputElement).value.split(',').map((s: string) => s.trim()).filter(Boolean)" type="text" class="input" placeholder="Mai Unm, Mai Civi" />
          </div>
          <div class="form-group" style="flex: 0 0 auto; display: flex; align-items: flex-end;">
            <button v-if="config.servers.length > 1" class="btn btn-sm btn-danger" @click="removeServer(i)" title="Remove server">✕</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Authentication Accounts -->
    <div class="card" style="margin-bottom: 1rem;">
      <div class="card-header">
        <span class="card-title">{{ t('settings.auths') }}</span>
        <button class="btn btn-sm" @click="addAuth">{{ t('settings.add_auth') }}</button>
      </div>
      <div v-for="(au, i) in config.auths" :key="i" style="padding: 0.75rem 0; border-bottom: 1px solid var(--border);">
        <div class="form-row">
          <div class="form-group">
            <label class="form-label">{{ t('settings.auth_username') }}</label>
            <input v-model="au.username" type="text" class="input" placeholder="AirlineSim account email" />
          </div>
          <div class="form-group">
            <label class="form-label">{{ t('settings.auth_password') }}</label>
            <input v-model="au.password" type="password" class="input" placeholder="Password" />
          </div>
          <div class="form-group" style="flex: 0 0 auto; display: flex; align-items: flex-end;">
            <button v-if="config.auths.length > 1" class="btn btn-sm btn-danger" @click="removeAuth(i)" title="Remove account">✕</button>
          </div>
        </div>
      </div>
      <div class="form-hint" style="margin-top: 0.5rem;">Each rule can select which account to use. Cookies are persisted separately.</div>
    </div>

    <!-- Monitoring -->
    <div class="card" style="margin-bottom: 1rem;">
      <div class="card-header">
        <span class="card-title">{{ t('settings.monitoring') }}</span>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">{{ t('settings.interval') }}</label>
          <input v-model.number="config.monitor.interval" type="number" class="input" min="15" />
          <div class="form-hint">{{ t('settings.interval_hint') }}</div>
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('settings.jitter') }}</label>
          <input v-model.number="config.monitor.jitter" type="number" class="input" min="0" />
          <div class="form-hint">{{ t('settings.jitter_hint') }}</div>
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('settings.timeout') }}</label>
          <input v-model.number="config.monitor.request_timeout" type="number" class="input" min="10" />
          <div class="form-hint">{{ t('settings.timeout_hint') }}</div>
        </div>
      </div>
      <div class="form-group" style="margin-top: 0.75rem;">
        <label class="form-label">{{ t('settings.min_balance') }}</label>
        <input v-model.number="config.monitor.min_balance" type="number" class="input" min="0" step="10000" />
        <div class="form-hint">Must retain at least this amount after purchasing aircraft. Default: 1,000,000 AS$</div>
      </div>
    </div>

    <!-- Notifications -->
    <div class="card" style="margin-bottom: 1rem;">
      <div class="card-header">
        <span class="card-title">{{ t('settings.notifications') }}</span>
      </div>
      <div class="form-group">
        <label class="toggle">
          <input v-model="config.notifier.console" type="checkbox" />
          <span class="toggle-track"></span>
          <span>{{ t('settings.console') }}</span>
        </label>
      </div>
      <div class="form-group" style="margin-top: 0.75rem;">
        <label class="form-label">{{ t('settings.discord') }}</label>
        <input v-model="config.notifier.discord_webhook" type="text" class="input" :placeholder="t('settings.discord_placeholder')" />
        <div class="form-hint">{{ t('settings.discord_hint') }}</div>
      </div>
      <div class="form-group" style="margin-top: 0.75rem;">
        <label class="form-label">{{ t('settings.dingtalk') }}</label>
        <input v-model="config.notifier.dingtalk_webhook" type="text" class="input" :placeholder="t('settings.dingtalk_placeholder')" />
        <div class="form-hint">{{ t('settings.dingtalk_hint') }}</div>
      </div>
      <div class="form-group" style="margin-top: 0.75rem;">
        <label class="form-label">{{ t('settings.dingtalk_secret') }}</label>
        <input v-model="config.notifier.dingtalk_secret" type="password" class="input" :placeholder="t('settings.dingtalk_secret_placeholder')" />
        <div class="form-hint">{{ t('settings.dingtalk_secret_hint') }}</div>
      </div>
    </div>

    <!-- Web UI -->
    <div class="card" style="margin-bottom: 1.5rem;">
      <div class="card-header">
        <span class="card-title">{{ t('settings.webui') }}</span>
      </div>
      <div class="form-row">
        <div class="form-group">
          <label class="form-label">{{ t('settings.webui_host') }}</label>
          <input v-model="config.webui.host" type="text" class="input" />
        </div>
        <div class="form-group">
          <label class="form-label">{{ t('settings.webui_port') }}</label>
          <input v-model.number="config.webui.port" type="number" class="input" min="1024" max="65535" />
        </div>
      </div>
    </div>

    <!-- Error -->
    <div v-if="error" style="background: var(--danger-bg); border: 1px solid rgba(239,68,68,0.3); border-radius: var(--radius-sm); padding: 0.625rem 0.75rem; margin-bottom: 1rem; font-size: 0.8125rem; color: var(--danger);">
      {{ error }}
    </div>

    <!-- Actions -->
    <div class="flex items-center gap-4">
      <button class="btn btn-primary" :disabled="loading" @click="save">
        <svg v-if="loading" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16" class="spin">
          <path d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16">
          <path d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        {{ loading ? t('settings.saving') : t('settings.save') }}
      </button>
      <span v-if="saved" style="color: var(--success); font-size: 0.875rem; display: flex; align-items: center; gap: 0.375rem;">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16">
          <path d="M4.5 12.75l6 6 9-13.5" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        {{ t('settings.saved') }}
      </span>
    </div>
  </div>
</template>

<style scoped>
.spin {
  animation: spin 1s linear infinite;
}
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>