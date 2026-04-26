import type { RouteComponent } from 'vue-router'

const viewModules = import.meta.glob('../views/**/*.vue')

export function loadView(component?: string): RouteComponent {
  if (!component) {
    return () => import('../views/error/404.vue')
  }
  const key = `../views/${component}.vue`
  const mod = viewModules[key]
  if (!mod) {
    console.warn(`[loadView] 未找到视图: ${key}`)
    return () => import('../views/error/404.vue')
  }
  return mod as () => Promise<RouteComponent>
}
