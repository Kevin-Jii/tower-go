import ArcoVue from '@arco-design/web-vue'
import '@arco-design/web-vue/dist/arco.css'
import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'
import { createPinia } from 'pinia'
import { createApp, nextTick } from 'vue'
import App from './App.vue'
import { permissionDirective } from './permission/directive'
import router, { resetDynamicRoutes } from './router'
import { setupRouterGuard } from './permission'
import { setLogoutNavigationHandler, useUserStore } from './store/user'
import { setAuthExpiredHandler } from './api/http'
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
const userStore = useUserStore(pinia)
userStore.hydrateFromStorage()
setLogoutNavigationHandler(async (options) => {
  resetDynamicRoutes(router)
  const redirect = options?.redirect
  await router.replace({
    name: 'Login',
    query: redirect && redirect !== '/login' ? { redirect } : undefined,
  })
})
setAuthExpiredHandler(async () => {
  await userStore.logout({
    redirect: router.currentRoute.value.fullPath,
    message: '登录已过期，请重新登录',
  })
})
app.use(router)
app.use(ArcoVue)
app.use(VueQueryPlugin, { queryClient })
app.directive('permission', permissionDirective)

setupRouterGuard(router)

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
