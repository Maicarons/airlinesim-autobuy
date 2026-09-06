import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import Dashboard from './views/Dashboard.vue'
import Rules from './views/Rules.vue'
import Settings from './views/Settings.vue'
import Logs from './views/Logs.vue'
import './assets/main.css'

const routes = [
  { path: '/', name: 'Dashboard', component: Dashboard },
  { path: '/rules', name: 'Rules', component: Rules },
  { path: '/settings', name: 'Settings', component: Settings },
  { path: '/logs', name: 'Logs', component: Logs },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

const app = createApp(App)
app.use(router)
app.mount('#app')