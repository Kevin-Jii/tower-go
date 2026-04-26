import type { RouteRecordRaw, Router } from 'vue-router'
import { createRouter, createWebHistory } from 'vue-router'
import LayoutAdmin from '@/layout/index.vue'

const dynamicRouteNames: string[] = []

/** 必须排在所有动态业务路由之后，否则会抢先匹配 /system/... 导致刷新误进 404 */
const NOT_FOUND_CATCH_ALL = 'NotFoundCatchAll'

export function registerDynamicRoutes(router: Router, routes: RouteRecordRaw[]): void {
  for (const r of routes) {
    const name = r.name as string | undefined
    if (name && router.hasRoute(name)) continue
    router.addRoute('Layout', r)
    if (name) dynamicRouteNames.push(name)
  }
}

export function registerNotFoundCatchAll(router: Router): void {
  if (router.hasRoute(NOT_FOUND_CATCH_ALL)) return
  router.addRoute({
    path: '/:pathMatch(.*)*',
    name: NOT_FOUND_CATCH_ALL,
    redirect: '/404',
  })
}

export function resetDynamicRoutes(router: Router): void {
  for (const n of dynamicRouteNames) {
    if (router.hasRoute(n)) router.removeRoute(n)
  }
  dynamicRouteNames.length = 0
  if (router.hasRoute(NOT_FOUND_CATCH_ALL)) router.removeRoute(NOT_FOUND_CATCH_ALL)
}

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { title: '登录', public: true },
  },
  {
    path: '/',
    name: 'Layout',
    component: LayoutAdmin,
    redirect: { name: 'Dashboard' },
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '工作台' },
      },
    ],
  },
  {
    path: '/404',
    name: 'NotFoundPage',
    component: () => import('@/views/error/404.vue'),
    meta: { title: '未找到', public: true },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
