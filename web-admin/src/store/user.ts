import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { Menu, User } from '@/api/types'
import { fetchUserMenus, fetchUserPermissions } from '@/api/menu'
import {
  clearStoredTenantId,
  getStoredTenantId,
  setStoredTenantId,
  setToken,
} from '@/utils/storage'

export const useUserStore = defineStore('user', () => {
  function hasStatisticsMenu(list: Menu[] | undefined): boolean {
    if (!list?.length) return false
    for (const m of list) {
      if (m.path === '/store/statistics' || m.component === 'store/statistics/index') return true
      if (m.children?.length && hasStatisticsMenu(m.children)) return true
    }
    return false
  }

  function attachStatisticsMenu(list: Menu[]): Menu[] {
    if (hasStatisticsMenu(list)) return list

    const statsMenu: Menu = {
      id: -10001,
      parent_id: 0,
      name: 'statistics-dash-open',
      title: '数据统计',
      icon: 'apps',
      path: '/store/statistics',
      component: 'store/statistics/index',
      type: 2,
      sort: 999,
      permission: '',
      visible: 1,
      status: 1,
    }

    const cloned = JSON.parse(JSON.stringify(list)) as Menu[]
    const storeNode = cloned.find((m) => m.type === 1 && (m.path === '/store' || m.name === 'store'))
    if (storeNode) {
      storeNode.children = storeNode.children ?? []
      storeNode.children.push({ ...statsMenu, parent_id: storeNode.id })
      return cloned
    }
    cloned.push(statsMenu)
    return cloned
  }

  const token = ref('')
  const userInfo = ref<User | null>(null)
  const permissions = ref<string[]>([])
  const menus = ref<Menu[]>([])
  const tenantId = ref(0)
  const dynamicRoutesReady = ref(false)

  const isLoggedIn = computed(() => !!token.value)

  function hydrateFromStorage(): void {
    const t = localStorage.getItem('tower_token') ?? ''
    token.value = t
    const raw = localStorage.getItem('tower_user')
    if (raw) {
      try {
        userInfo.value = JSON.parse(raw) as User
      } catch {
        userInfo.value = null
      }
    }
    const tid = getStoredTenantId()
    if (tid != null && tid > 0) tenantId.value = tid
    else if (userInfo.value?.store_id) {
      tenantId.value = Number(userInfo.value.store_id)
      setStoredTenantId(tenantId.value)
    }
  }

  function setLogin(payload: { token: string; user: User }): void {
    token.value = payload.token
    userInfo.value = payload.user
    setToken(payload.token)
    localStorage.setItem('tower_user', JSON.stringify(payload.user))
    const sid = Number(payload.user.store_id ?? 0)
    if (getStoredTenantId() == null && sid > 0) {
      tenantId.value = sid
      setStoredTenantId(sid)
    } else if (getStoredTenantId() != null) {
      tenantId.value = getStoredTenantId()!
    } else {
      tenantId.value = sid
    }
  }

  function setPermissions(list: string[]): void {
    permissions.value = list ?? []
  }

  function setMenus(list: Menu[]): void {
    menus.value = attachStatisticsMenu(list ?? [])
  }

  function setTenantId(id: number): void {
    tenantId.value = id
    setStoredTenantId(id)
  }

  function patchUserInfo(partial: User): void {
    if (!userInfo.value) {
      userInfo.value = partial
    } else {
      userInfo.value = { ...userInfo.value, ...partial }
    }
    localStorage.setItem('tower_user', JSON.stringify(userInfo.value))
  }

  async function loadMenusAndPermissions(): Promise<void> {
    const [m, p] = await Promise.all([fetchUserMenus(), fetchUserPermissions()])
    setMenus(m)
    setPermissions(p)
  }

  async function logout(): Promise<void> {
    const { default: router, resetDynamicRoutes } = await import('@/router')
    resetDynamicRoutes(router)
    token.value = ''
    userInfo.value = null
    permissions.value = []
    menus.value = []
    dynamicRoutesReady.value = false
    setToken('')
    localStorage.removeItem('tower_user')
    clearStoredTenantId()
    await router.replace({ name: 'Login' })
  }

  function markDynamicRoutes(ready: boolean): void {
    dynamicRoutesReady.value = ready
  }

  return {
    token,
    userInfo,
    permissions,
    menus,
    tenantId,
    dynamicRoutesReady,
    isLoggedIn,
    hydrateFromStorage,
    setLogin,
    setPermissions,
    setMenus,
    setTenantId,
    patchUserInfo,
    loadMenusAndPermissions,
    logout,
    markDynamicRoutes,
  }
})
