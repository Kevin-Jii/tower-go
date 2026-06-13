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
    redirect: '/dashboard',
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('@/views/dashboard/index.vue'),
    meta: { title: '经营数据大屏' },
  },
  {
    path: '/',
    name: 'Layout',
    component: LayoutAdmin,
    children: [
      {
        path: 'admin',
        name: 'AdminEntry',
        component: () => import('@/views/AdminEntryRedirect.vue'),
        meta: { title: '后台' },
      },
      {
        path: 'profile',
        name: 'UserProfile',
        component: () => import('@/views/profile/index.vue'),
        meta: { title: '个人资料' },
      },
      {
        path: 'system/third-party-account',
        name: 'ThirdPartyAccount',
        component: () => import('@/views/system/third-party-account/index.vue'),
        meta: { title: '第三方账号池' },
      },
      {
        path: 'system/third-party-routes',
        name: 'ThirdPartyRoutes',
        component: () => import('@/views/system/third-party-account/routes.vue'),
        meta: { title: '物流路线图' },
      },
      {
        path: 'system/third-party-routes/import',
        name: 'ThirdPartyRouteImport',
        component: () => import('@/views/system/third-party-account/route-import.vue'),
        meta: { title: '导入第三方订单' },
      },
      {
        path: 'system/third-party-routes/history',
        name: 'ThirdPartyRouteHistory',
        component: () => import('@/views/system/third-party-account/route-history.vue'),
        meta: { title: '历史物流单' },
      },
    ],
  },
  {
    path: '/404',
    name: 'NotFoundPage',
    component: () => import('@/views/error/404.vue'),
    meta: { title: '未找到', public: true },
  },
  {
    path: '/public/supplier/:id',
    name: 'PublicSupplierProfile',
    component: () => import('@/views/public/supplier-profile.vue'),
    meta: { title: '供应商档案', public: true },
  },
  {
    path: '/third-party-account/:id/orders',
    name: 'ThirdPartyAccountOrders',
    component: () => import('@/views/system/third-party-account/orders.vue'),
    meta: { title: '已同步订单' },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
