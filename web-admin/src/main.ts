import ArcoVue from '@arco-design/web-vue'
import '@arco-design/web-vue/dist/arco.css'
import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'
import { createPinia } from 'pinia'
import { createApp, nextTick, h } from 'vue'
import App from './App.vue'
import { MathCurveLoader } from './components/loading'
import { permissionDirective } from './permission/directive'
import router from './router'
import { setupRouterGuard } from './permission'
import { useUserStore } from './store/user'
import 'virtual:uno.css'
import './styles/global.css'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 60_000,
      retry: 1,
    },
  },
})

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
useUserStore(pinia).hydrateFromStorage()
app.use(router)
app.use(ArcoVue)
app.use(VueQueryPlugin, { queryClient })
app.directive('permission', permissionDirective)

setupRouterGuard(router)

function mountBootCurveLoader(): void {
  const holder = document.getElementById('app-boot-splash-loader')
  if (!holder) return
  createApp({
    render: () => h(MathCurveLoader, { size: 'lg', inline: true }),
  }).mount(holder)
}

mountBootCurveLoader()

app.mount('#app')

function removeBootSplash(): void {
  const el = document.getElementById('app-boot-splash')
  if (!el) return
  el.classList.add('app-boot-splash--exit')
  const finish = (): void => {
    el.remove()
  }
  el.addEventListener('transitionend', finish, { once: true })
  setTimeout(finish, 500)
}

void Promise.all([
  router.isReady(),
  new Promise<void>((resolve) => {
    requestAnimationFrame(() => requestAnimationFrame(() => resolve()))
  }),
])
  .then(() => nextTick())
  .then(() => {
    removeBootSplash()
  })
  .catch(() => {
    removeBootSplash()
  })
