<template>
  <a-breadcrumb class="min-w-0 flex-1 !hidden sm:!flex text-sm">
    <a-breadcrumb-item v-for="(b, i) in items" :key="i">
      <RouterLink
        v-if="b.to && i < items.length - 1"
        :to="b.to"
        class="truncate max-w-[10rem] md:max-w-[14rem] text-[var(--color-text-2)] hover:text-[rgb(var(--primary-6))]"
      >
        {{ b.label }}
      </RouterLink>
      <span v-else class="truncate max-w-[12rem] md:max-w-[18rem] font-medium text-[var(--color-text-1)]">
        {{ b.label }}
      </span>
    </a-breadcrumb-item>
  </a-breadcrumb>
</template>

<script setup lang="ts">
import type { RouteLocationRaw } from 'vue-router'

defineProps<{
  items: { label: string; to?: RouteLocationRaw }[]
}>()
</script>
