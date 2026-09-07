<script setup lang="ts">
import { ref, inject, onMounted, computed } from 'vue'
import { getRules, createRule, updateRule, deleteRule, toggleRule, getAircraftData, getConfig } from '../api'
import type { Rule, MatchConfig, ActionConfig, AircraftData, AircraftFamily, AircraftType } from '../api/types'

const t = inject('t') as (key: string, fallback?: string) => string

const rules = ref<Rule[]>([])
const showForm = ref(false)
const editingIndex = ref<number | null>(null)
const saving = ref(false)
const formError = ref('')

const aircraftData = ref<AircraftData>({ families: [], types: [] })
const filteredTypes = ref<AircraftType[]>([])
const servers = ref<any[]>([])
const auths = ref<any[]>([])

// Companies for the selected server
const availableCompanies = computed(() => {
  const sv = servers.value[form.value.server_id]
  return sv?.companies || []
})

const form = ref<Rule>(defaultRule())

function defaultRule(): Rule {
  return {
    name: '',
    enabled: true,
    priority: 10,
    server_id: 0,
    auth_id: 0,
    company_name: '',
    match: {
      family_id: '',
      type_id: '',
      types: [],
      price_range: { min: 0, max: 0 },
      max_age: 20,
      max_cycles: 50000,
      condition_min: 50,
      offer_types: ['auction', 'immediate'],
      financing: ['cash', 'credit', 'lease'],
      sort_by: 'price_asc',
    },
    action: {
      auto_buy: false,
      snatch: false,
      max_bid_increment: 100000,
      max_count: 0,
    },
  }
}

const offerTypeOptions = ['auction', 'immediate']
const financingOptions = ['cash', 'credit', 'lease']
const sortOptions = [
  { value: 'price_asc', label: 'Price: Low to High' },
  { value: 'price_desc', label: 'Price: High to Low' },
  { value: 'age_asc', label: 'Age: Low to High' },
  { value: 'age_desc', label: 'Age: High to Low' },
  { value: 'deadline_asc', label: 'Earliest Deadline' },
  { value: 'deadline_desc', label: 'Latest Deadline' },
  { value: 'bid_asc', label: 'Lowest Bid' },
  { value: 'bid_desc', label: 'Highest Bid' },
]

async function fetchRules() {
  try {
    rules.value = await getRules()
  } catch {}
}

async function fetchAircraftData() {
  try {
    aircraftData.value = await getAircraftData()
  } catch {}
}

function onFamilyChange() {
  // Clear type selection when the user manually changes the family
  form.value.match.type_id = ''
}

function openNew() {
  form.value = defaultRule()
  editingIndex.value = null
  formError.value = ''
  showForm.value = true
}

function openEdit(index: number) {
  form.value = JSON.parse(JSON.stringify(rules.value[index]))
  editingIndex.value = index
  formError.value = ''
  showForm.value = true
}

function validate(): boolean {
  if (!form.value.name.trim()) {
    formError.value = 'Rule name is required'
    return false
  }
  if (form.value.server_id < 0 || form.value.server_id >= servers.value.length) {
    formError.value = 'Please select a server'
    return false
  }
  if (form.value.auth_id < 0 || form.value.auth_id >= auths.value.length) {
    formError.value = 'Please select an authentication account'
    return false
  }
  if (form.value.match.price_range.max > 0 && form.value.match.price_range.min > form.value.match.price_range.max) {
    formError.value = 'Min price cannot exceed max price'
    return false
  }
  formError.value = ''
  return true
}

async function save() {
  if (!validate()) return
  saving.value = true
  try {
    if (editingIndex.value !== null) {
      await updateRule(editingIndex.value, form.value)
    } else {
      await createRule(form.value)
    }
    showForm.value = false
    await fetchRules()
  } catch (e: any) {
    formError.value = e.message
  }
  saving.value = false
}

async function remove(index: number) {
  if (!confirm(t('rules.delete_confirm') + ' This action cannot be undone.')) return
  try {
    await deleteRule(index)
    await fetchRules()
  } catch (e: any) {
    alert(e.message)
  }
}

async function toggle(index: number) {
  try {
    await toggleRule(index)
    await fetchRules()
  } catch (e: any) {
    alert(e.message)
  }
}

function toggleOfferType(val: string) {
  const idx = form.value.match.offer_types.indexOf(val)
  if (idx >= 0) {
    form.value.match.offer_types.splice(idx, 1)
  } else {
    form.value.match.offer_types.push(val)
  }
}

function toggleFinancing(val: string) {
  const idx = form.value.match.financing.indexOf(val)
  if (idx >= 0) {
    form.value.match.financing.splice(idx, 1)
  } else {
    form.value.match.financing.push(val)
  }
}

function getFamilyName(id: string): string {
  const f = aircraftData.value.families.find(f => f.id === id)
  return f ? f.name : id
}

function getTypeName(id: string): string {
  const t = aircraftData.value.types.find(t => t.id === id)
  return t ? t.name : id
}

onMounted(() => {
  fetchRules()
  fetchAircraftData()
  getConfig().then(c => { servers.value = c.servers || []; auths.value = c.auths || [] })
})
</script>

<template>
  <div>
    <div class="page-header">
      <div>
        <h1 class="page-title">{{ t('rules.title') }}</h1>
        <p class="page-subtitle">{{ t('rules.subtitle') }}</p>
      </div>
      <button class="btn btn-primary" @click="openNew">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16">
          <path d="M12 4.5v15m7.5-7.5h-15" stroke-linecap="round"/>
        </svg>
        {{ t('rules.new') }}
      </button>
    </div>

    <!-- Empty State -->
    <div v-if="!rules.length" class="card empty-state">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" width="48" height="48">
        <path d="M9 12h3.75M9 15h3.75M9 18h3.75m3 .75H18a2.25 2.25 0 002.25-2.25V6.108c0-1.135-.845-2.098-1.976-2.192a48.424 48.424 0 00-1.123-.08m-5.801 0c-.065.21-.1.433-.1.664 0 .414.336.75.75.75h4.5a.75.75 0 00.75-.75 2.25 2.25 0 00-.1-.664m-5.8 0A2.251 2.251 0 0113.5 2.25H15a2.25 2.25 0 012.15 1.586m-5.8 0c-.376.023-.75.05-1.124.08C9.095 4.01 8.25 4.973 8.25 6.108V8.25m0 0H4.875c-.621 0-1.125.504-1.125 1.125v11.25c0 .621.504 1.125 1.125 1.125h9.75c.621 0 1.125-.504 1.125-1.125V9.375c0-.621-.504-1.125-1.125-1.125H8.25zM6.75 12h.008v.008H6.75V12zm0 3h.008v.008H6.75V15zm0 3h.008v.008H6.75V18z"/>
      </svg>
      <div class="empty-state-title">{{ t('rules.empty_title') }}</div>
      <div class="empty-state-desc">{{ t('rules.empty_desc') }}</div>
    </div>

    <!-- Rule Cards -->
    <div v-for="(rule, index) in rules" :key="index" class="card" style="margin-bottom: 0.75rem;">
      <div class="rule-header">
        <div class="rule-left">
          <button
            :class="['rule-toggle', rule.enabled ? 'enabled' : 'disabled']"
            @click="toggle(index)"
            :title="rule.enabled ? 'Disable rule' : 'Enable rule'"
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
              <path v-if="rule.enabled" d="M4.5 12.75l6 6 9-13.5" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </button>
          <div>
            <div class="rule-name">{{ rule.name || 'Unnamed Rule' }}</div>
            <div class="rule-meta">
              {{ t('rules.priority') }} {{ rule.priority }}
              <span v-if="rule.company_name"> · {{ rule.company_name }}</span>
              <span v-if="rule.match.type_id"> · {{ getTypeName(rule.match.type_id) }}</span>
              <span v-else-if="rule.match.family_id"> · {{ getFamilyName(rule.match.family_id) }}</span>
            </div>
          </div>
        </div>
        <div class="flex gap-2">
          <span :class="['badge', rule.enabled ? 'badge-success' : 'badge-danger']">
            {{ rule.enabled ? t('rules.active') : t('rules.disabled') }}
          </span>
          <button class="btn btn-sm btn-ghost" @click="openEdit(index)" title="Edit rule">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14">
              <path d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </button>
          <button class="btn btn-sm btn-ghost" @click="remove(index)" title="Delete rule">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="14" height="14">
              <path d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </button>
        </div>
      </div>
      <div class="rule-details">
        <div class="rule-detail">
          <span class="detail-label">{{ t('rules.price') }}</span>
          <span class="detail-value">AS$ {{ rule.match.price_range.min.toLocaleString() }} – AS$ {{ rule.match.price_range.max.toLocaleString() }}</span>
        </div>
        <div class="rule-detail">
          <span class="detail-label">{{ t('rules.age') }}</span>
          <span class="detail-value">≤ {{ rule.match.max_age }} years</span>
        </div>
        <div class="rule-detail">
          <span class="detail-label">{{ t('rules.condition') }}</span>
          <span class="detail-value">≥ {{ rule.match.condition_min }}%</span>
        </div>
        <div class="rule-detail">
          <span class="detail-label">{{ t('rules.cycles') }}</span>
          <span class="detail-value">≤ {{ rule.match.max_cycles.toLocaleString() }}</span>
        </div>
        <div class="rule-detail">
          <span class="detail-label">{{ t('rules.auto_buy') }}</span>
          <span :class="['detail-value', rule.action.auto_buy ? 'text-success' : 'text-muted']">
            {{ rule.action.auto_buy ? t('rules.enabled') : t('rules.disabled') }}
          </span>
        </div>
        <div class="rule-detail">
          <span class="detail-label">{{ t('rules.snatch') }}</span>
          <span :class="['detail-value', rule.action.snatch ? 'text-warning' : 'text-muted']">
            {{ rule.action.snatch ? t('rules.enabled') : t('rules.disabled') }}
          </span>
        </div>
      </div>
    </div>

    <!-- Modal -->
    <div v-if="showForm" class="modal-overlay" @click.self="showForm = false">
      <div class="modal">
        <div class="modal-title">{{ editingIndex !== null ? t('rules.edit') : t('rules.create') }}</div>

        <!-- Error -->
        <div v-if="formError" style="background: var(--danger-bg); border: 1px solid rgba(239,68,68,0.3); border-radius: var(--radius-sm); padding: 0.625rem 0.75rem; margin-bottom: 1rem; font-size: 0.8125rem; color: var(--danger);">
          {{ formError }}
        </div>

        <!-- Rule Name -->
        <div class="form-group">
          <label class="form-label">{{ t('rules.name') }}</label>
          <input v-model="form.name" type="text" class="input" :placeholder="t('rules.name_placeholder')" />
        </div>

        <!-- Server Selection -->
        <div class="form-group">
          <label class="form-label">{{ t('rules.server') }}</label>
          <select v-model="form.server_id" class="input">
            <option v-for="(sv, i) in servers" :key="i" :value="i">{{ sv.host }} — {{ sv.base_url }}</option>
          </select>
        </div>

        <!-- Auth Selection -->
        <div class="form-group">
          <label class="form-label">{{ t('rules.account') }}</label>
          <select v-model="form.auth_id" class="input">
            <option v-for="(au, i) in auths" :key="i" :value="i">Account #{{ i+1 }} — {{ au.username || 'unnamed' }}</option>
          </select>
        </div>

        <!-- Company Selection -->
        <div class="form-group">
          <label class="form-label">{{ t('rules.company') }}</label>
          <select v-model="form.company_name" class="input">
            <option value="">{{ t('rules.company_default') }}</option>
            <option v-for="c in availableCompanies" :key="c" :value="c">{{ c }}</option>
          </select>
          <div class="form-hint">{{ t('rules.company_hint') }}</div>
        </div>

        <!-- Aircraft Family -->
        <div class="form-group">
          <label class="form-label">{{ t('rules.family') }}</label>
          <select v-model="form.match.family_id" class="input" @change="onFamilyChange">
            <option value="">{{ t('rules.family_any') }}</option>
            <option v-for="f in aircraftData.families" :key="f.id" :value="f.id">{{ f.name }}</option>
          </select>
          <div class="form-hint">{{ t('rules.family_hint') }}</div>
        </div>

        <!-- Aircraft Type -->
        <div class="form-group">
          <label class="form-label">{{ t('rules.type') }}</label>
          <select v-model="form.match.type_id" class="input">
            <option value="">{{ t('rules.type_any') }}</option>
            <option v-for="t in aircraftData.types" :key="t.id" :value="t.id">{{ t.name }}</option>
          </select>
          <div class="form-hint">{{ t('rules.type_hint') }}</div>
        </div>

        <!-- Sort -->
        <div class="form-group">
          <label class="form-label">{{ t('rules.sort') }}</label>
          <select v-model="form.match.sort_by" class="input">
            <option value="">{{ t('rules.sort_default') }}</option>
            <option v-for="s in sortOptions" :key="s.value" :value="s.value">{{ s.label }}</option>
          </select>
        </div>

        <!-- Price Range -->
        <div class="form-row">
          <div class="form-group">
            <label class="form-label">{{ t('rules.price_min') }}</label>
            <input v-model.number="form.match.price_range.min" type="number" class="input" min="0" step="1000" />
          </div>
          <div class="form-group">
            <label class="form-label">{{ t('rules.price_max') }}</label>
            <input v-model.number="form.match.price_range.max" type="number" class="input" min="0" step="1000" />
          </div>
        </div>

        <!-- Age, Cycles, Condition -->
        <div class="form-row">
          <div class="form-group">
            <label class="form-label">{{ t('rules.max_age') }}</label>
            <input v-model.number="form.match.max_age" type="number" class="input" min="0" />
          </div>
          <div class="form-group">
            <label class="form-label">{{ t('rules.max_cycles') }}</label>
            <input v-model.number="form.match.max_cycles" type="number" class="input" min="0" step="1000" />
          </div>
          <div class="form-group">
            <label class="form-label">{{ t('rules.condition_min') }}</label>
            <input v-model.number="form.match.condition_min" type="number" class="input" min="0" max="100" />
          </div>
        </div>

        <!-- Offer Types -->
        <div class="form-group">
          <label class="form-label">{{ t('rules.offer_types') }}</label>
          <div class="flex gap-3">
            <label v-for="opt in offerTypeOptions" :key="opt" class="checkbox">
              <input type="checkbox" :checked="form.match.offer_types.includes(opt)" @change="toggleOfferType(opt)" />
              <span>{{ opt.charAt(0).toUpperCase() + opt.slice(1) }}</span>
            </label>
          </div>
        </div>

        <!-- Financing -->
        <div class="form-group">
          <label class="form-label">{{ t('rules.financing') }}</label>
          <div class="flex gap-3">
            <label v-for="opt in financingOptions" :key="opt" class="checkbox">
              <input type="checkbox" :checked="form.match.financing.includes(opt)" @change="toggleFinancing(opt)" />
              <span>{{ opt.charAt(0).toUpperCase() + opt.slice(1) }}</span>
            </label>
          </div>
        </div>

        <!-- Priority & Bid Increment & Max Count -->
        <div class="form-row">
          <div class="form-group">
            <label class="form-label">{{ t('rules.priority') }}</label>
            <input v-model.number="form.priority" type="number" class="input" min="0" />
            <div class="form-hint">{{ t('rules.priority_hint') }}</div>
          </div>
          <div class="form-group">
            <label class="form-label">{{ t('rules.bid_increment') }}</label>
            <input v-model.number="form.action.max_bid_increment" type="number" class="input" min="0" step="1000" />
            <div class="form-hint">{{ t('rules.bid_hint') }}</div>
          </div>
          <div class="form-group">
            <label class="form-label">{{ t('rules.max_count') }}</label>
            <input v-model.number="form.action.max_count" type="number" class="input" min="0" />
            <div class="form-hint">{{ t('rules.max_count_hint') }}</div>
          </div>
        </div>

        <!-- Toggles -->
        <div class="form-row" style="margin-top: 0.5rem;">
          <div class="form-group">
            <label class="toggle">
              <input v-model="form.enabled" type="checkbox" />
              <span class="toggle-track"></span>
              <span>{{ t('rules.enabled') }}</span>
            </label>
          </div>
          <div class="form-group">
            <label class="toggle">
              <input v-model="form.action.auto_buy" type="checkbox" />
              <span class="toggle-track"></span>
              <span>{{ t('rules.auto_buy') }}</span>
            </label>
          </div>
          <div class="form-group">
            <label class="toggle">
              <input v-model="form.action.snatch" type="checkbox" />
              <span class="toggle-track"></span>
              <span>{{ t('rules.snatch') }}</span>
            </label>
          </div>
        </div>

        <!-- Footer -->
        <div class="modal-footer">
          <button class="btn" @click="showForm = false">{{ t('rules.cancel') }}</button>
          <button class="btn btn-primary" :disabled="saving" @click="save">
            <svg v-if="saving" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="16" height="16" class="spin">
              <path d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            {{ saving ? t('rules.saving') : t('rules.save') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.rule-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.rule-left {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.rule-toggle {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  border: 1px solid var(--border);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.15s ease;
  background: transparent;
}
.rule-toggle.enabled {
  background: var(--success-bg);
  border-color: rgba(34, 197, 94, 0.3);
  color: var(--success);
}
.rule-toggle.disabled {
  background: var(--bg-input);
  color: var(--text-muted);
}
.rule-toggle:hover {
  opacity: 0.8;
}
.rule-toggle:focus-visible {
  outline: 2px solid var(--accent);
  outline-offset: 2px;
}

.rule-name {
  font-size: 0.9375rem;
  font-weight: 600;
}
.rule-meta {
  font-size: 0.75rem;
  color: var(--text-muted);
  margin-top: 0.125rem;
}

.rule-details {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 0.5rem;
  padding-top: 0.75rem;
  border-top: 1px solid var(--border);
}

.rule-detail {
  font-size: 0.8125rem;
}
.detail-label {
  color: var(--text-muted);
  margin-right: 0.375rem;
}
.detail-value {
  color: var(--text-secondary);
}

.spin {
  animation: spin 1s linear infinite;
}
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>