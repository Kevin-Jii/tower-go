import type { Router } from 'vue-router'
import { useUserStore } from '@/store/user'
import { menusToRoutes } from './transformRoutes'
import { registerDynamicRoutes, registerNotFoundCatchAll } from '@/router'

const WHITE = new Set(['/login', '/404'])

export function setupRouterGuard(router: Router): void {
  router.beforeEach(async (to, _from, next) => {
    if (WHITE.has(to.path)) {
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
}
