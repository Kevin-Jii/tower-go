<template>
  <AppPageLoading text="正在进入后台..." />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import AppPageLoading from '@/components/AppPageLoading.vue'
import type { Menu } from '@/api/types'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()

function firstPagePath(list: Menu[] | undefined): string {
  if (!list?.length) return '/profile'
  for (const item of list) {
    if (item.visible === 0 || item.status === 0) continue
    if (item.type === 2 && item.path) return item.path
    const childPath = firstPagePath(item.children)
    if (childPath !== '/profile') return childPath
  }
  return '/profile'
}

onMounted(() => {
  void router.replace(firstPagePath(userStore.menus))
})
</script>
