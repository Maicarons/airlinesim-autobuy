import { ref, computed } from 'vue'
import en from './locales/en.json'
import zhCN from './locales/zh-CN.json'
import ko from './locales/ko.json'

export type Locale = 'en' | 'zh-CN' | 'ko'

const locales: Record<Locale, Record<string, string>> = {
  en,
  'zh-CN': zhCN,
  ko,
}

const names: Record<Locale, string> = {
  en: 'English',
  'zh-CN': '简体中文',
  ko: '한국어',
}

const saved = localStorage.getItem('locale') as Locale | null
const currentLocale = ref<Locale>(saved || 'en')

export function useI18n() {
  const t = (key: string, fallback?: string): string => {
    return locales[currentLocale.value]?.[key] || fallback || key
  }

  const setLocale = (locale: Locale) => {
    currentLocale.value = locale
    localStorage.setItem('locale', locale)
    document.documentElement.lang = locale
  }

  const localeNames = computed(() => {
    return Object.entries(names).map(([k, v]) => ({ value: k as Locale, label: v }))
  })

  return { t, setLocale, currentLocale, localeNames }
}