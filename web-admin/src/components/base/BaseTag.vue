<template>
  <a-tag :color="arcoColor" size="small" class="font-semibold">
    <slot />
  </a-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    variant?: 'success' | 'warning' | 'info' | 'danger' | 'neutral'
  }>(),
  { variant: 'neutral' },
)

const arcoColor = computed(() => {
  const m: Record<string, string> = {
    success: 'green',
    warning: 'gold',
    info: 'arcoblue',
    danger: 'red',
    neutral: 'gray',
  }
  return m[props.variant] ?? 'gray'
})
</script>
