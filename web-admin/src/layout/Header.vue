<template>
  <div class="layout-header">
    <a-button type="text" class="md:!hidden" @click="$emit('toggle-mobile-menu')">
      <template #icon>
        <IconMenu />
      </template>
    </a-button>

    <a-button type="text" class="!hidden md:!inline-flex" @click="layout.toggleSidebar()">
      <template #icon>
        <IconMenuFold v-if="!layout.sidebarCollapsed" />
        <IconMenuUnfold v-else />
      </template>
    </a-button>

    <a-breadcrumb class="min-w-0 flex-1 !hidden sm:!flex">
      <a-breadcrumb-item v-for="(b, i) in crumbs" :key="i">
        <RouterLink
          v-if="b.to && i < crumbs.length - 1"
          :to="b.to"
          class="text-[var(--color-text-2)] hover:text-[rgb(var(--primary-6))]"
        >
          {{ b.label }}
        </RouterLink>
        <span v-else class="text-[var(--color-text-1)] font-medium">{{ b.label }}</span>
      </a-breadcrumb-item>
    </a-breadcrumb>

    <div class="flex items-center gap-2 shrink-0 ml-auto">
      <a-dropdown trigger="click" position="br" @select="onMenuSelect">
        <button type="button" class="user-menu-trigger" aria-label="用户菜单">
          <a-avatar :size="32" class="!bg-[rgb(var(--primary-6))] shrink-0">
            {{ avatarText }}
          </a-avatar>
          <span class="hidden md:inline max-w-[8rem] truncate text-sm text-[var(--color-text-1)]">
            {{ displayName }}
          </span>
          <IconDown class="hidden md:block text-[var(--color-text-3)]" />
        </button>
        <template #content>
          <a-doption value="profile">
            <template #icon><IconUser /></template>
            个人资料
          </a-doption>
          <a-doption value="logout" class="text-[rgb(var(--danger-6))]">
            <template #icon><IconExport /></template>
            退出登录
          </a-doption>
        </template>
      </a-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
import { IconDown, IconExport, IconMenu, IconMenuFold, IconMenuUnfold, IconUser } from '@arco-design/web-vue/es/icon'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useLayoutStore } from '@/store/layout'
import { useUserStore } from '@/store/user'

defineEmits<{ 'toggle-mobile-menu': [] }>()

const route = useRoute()
const router = useRouter()
const layout = useLayoutStore()
const userStore = useUserStore()

const crumbs = computed(() =>
  route.matched
    .filter((r) => r.meta?.title)
    .map((r) => ({
      label: String(r.meta.title),
      to: r.path && r.path !== '/' ? { path: r.path } : undefined,
    })),
)

const displayName = computed(() => userStore.userInfo?.nickname || userStore.userInfo?.username || '用户')

const avatarText = computed(() => displayName.value.slice(0, 1).toUpperCase())

async function onMenuSelect(value: string | number | Record<string, unknown> | undefined): Promise<void> {
  if (value === 'profile') {
    await router.push({ name: 'UserProfile' })
    return
  }
  if (value === 'logout') {
    await userStore.logout()
  }
}
</script>

<style scoped>
.user-menu-trigger {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 4px 8px 4px 4px;
  border: 1px solid var(--color-border-2);
  border-radius: 999px;
  background: var(--color-bg-2);
  cursor: pointer;
  transition: border-color 0.15s ease, background 0.15s ease;
}
.user-menu-trigger:hover {
  border-color: rgb(var(--primary-6));
  background: var(--color-fill-1);
}
</style>
