import type { Router } from 'vue-router'
import { useUserStore } from '@/store/user'
import { menusToRoutes } from './transformRoutes'
import { registerDynamicRoutes, registerNotFoundCatchAll } from '@/router'
import { useLoadingStore } from '@/store/loading'

const WHITE = new Set(['/login', '/404'])

export function setupRouterGuard(router: Router): void {
  let routeLoadingActive = false

  router.beforeEach(async (to, _from, next) => {
    const loadingStore = useLoadingStore()
    if (WHITE.has(to.path) || to.meta?.public) {
      next()
      return
    }

    const userStore = useUserStore()
    if (!userStore.token) {
      userStore.hydrateFromStorage()
    }
    if (to.path === '/login' && userStore.token) {
      next({ name: 'Dashboard' })
      return
    }
    if (!userStore.token) {
      next({ path: '/login', query: { redirect: to.fullPath } })
      return
    }

    if (!userStore.dynamicRoutesReady) {
      loadingStore.show('数据加载中...')
      routeLoadingActive = true
      try {
        await userStore.loadMenusAndPermissions()
        registerDynamicRoutes(router, menusToRoutes(userStore.menus))
        registerNotFoundCatchAll(router)
        userStore.markDynamicRoutes(true)
        next({ ...to, replace: true })
      } catch (e) {
        console.error(e)
        userStore.logout()
        next({ path: '/login', query: { redirect: to.fullPath } })
      }
      return
    }

    next()
  })

  router.afterEach(() => {
    if (!routeLoadingActive) return
    const loadingStore = useLoadingStore()
    loadingStore.hide()
    routeLoadingActive = false
  })

  router.onError(() => {
    if (!routeLoadingActive) return
    const loadingStore = useLoadingStore()
    loadingStore.hide()
    routeLoadingActive = false
  })
}
