import type { Component } from 'vue'
import {
  IconApps,
  IconArchive,
  IconBook,
  IconCalendar,
  IconDashboard,
  IconFile,
  IconHistory,
  IconImage,
  IconLink,
  IconList,
  IconMenuFold,
  IconMessage,
  IconNotification,
  IconPrinter,
  IconRobot,
  IconSafe,
  IconSend,
  IconSettings,
  IconStorage,
  IconTags,
  IconUser,
  IconUserGroup,
} from '@arco-design/web-vue/es/icon'

export interface MenuIconOption {
  label: string
  value: string
  component: Component
}

const menuIconEntries: Array<[string, Component]> = [
  ['apps', IconApps],
  ['box', IconArchive],
  ['calendar', IconCalendar],
  ['data-board', IconDashboard],
  ['document', IconFile],
  ['history', IconHistory],
  ['link', IconLink],
  ['menu-fold', IconMenuFold],
  ['message', IconMessage],
  ['picture', IconImage],
  ['printer', IconPrinter],
  ['promotion', IconNotification],
  ['read', IconBook],
  ['robot', IconRobot],
  ['setting', IconSettings],
  ['shop', IconStorage],
  ['tickets', IconTags],
  ['user', IconUser],
  ['usergroup', IconUserGroup],
  ['van', IconSend],
  ['view-list', IconList],
  ['wallet', IconSafe],
  ['wallet-cards', IconSafe],
]

const menuIconMap = new Map(menuIconEntries)
menuIconMap.set('data-board', IconDashboard)
menuIconMap.set('user-group', IconUserGroup)
menuIconMap.set('wallet-cards', IconSafe)

export const menuIconOptions: MenuIconOption[] = menuIconEntries.map(([value, component]) => ({
  label: value,
  value,
  component,
}))

export function resolveMenuIcon(name?: string): Component | undefined {
  if (!name) return undefined
  const normalized = name
    .replace(/^Icon/, '')
    .replace(/([a-z0-9])([A-Z])/g, '$1-$2')
    .replace(/([A-Z])([A-Z][a-z])/g, '$1-$2')
    .toLowerCase()
  return menuIconMap.get(normalized) ?? IconFile
}
