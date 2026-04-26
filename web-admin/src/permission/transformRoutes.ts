import type { RouteRecordRaw } from 'vue-router'
import type { Menu } from '@/api/types'
import { loadView } from './loadView'

function walkMenus(menus: Menu[] | undefined, out: RouteRecordRaw[]): void {
  if (!menus?.length) return
  for (const m of menus) {
    if (m.type === 2 && m.component && m.path) {
      const path = m.path.replace(/^\//, '')
      out.push({
        path,
        name: m.name || `menu-${m.id}`,
        component: loadView(m.component),
        meta: {
          title: m.title,
          icon: m.icon,
          permission: m.permission,
        },
      })
    }
    if (m.children?.length) walkMenus(m.children, out)
  }
}

export function menusToRoutes(menus: Menu[]): RouteRecordRaw[] {
  const routes: RouteRecordRaw[] = []
  walkMenus(menus, routes)
  return routes
}
