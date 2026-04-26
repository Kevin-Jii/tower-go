import ArcoVue from '@arco-design/web-vue'
import '@arco-design/web-vue/dist/arco.css'
import { VueQueryPlugin, QueryClient } from '@tanstack/vue-query'
import { createPinia } from 'pinia'
import { createApp } from 'vue'
import App from './App.vue'
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

app.mount('#app')
