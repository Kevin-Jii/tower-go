<template>
  <BaseTag :variant="variant">{{ label }}</BaseTag>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { BaseTag } from '@/components/base'
import { getAllDicts } from '@/api/dict'

const props = defineProps<{
  type: string
  value: string | number
}>()

const { data: dictMap } = useQuery({
  queryKey: ['dict', 'all'],
  queryFn: getAllDicts,
  staleTime: 5 * 60_000,
})

const item = computed(() => {
  const list = dictMap.value?.[props.type] ?? []
  return list.find((d) => d.value === String(props.value))
})

const label = computed(() => item.value?.label ?? String(props.value))

const variant = computed(() => {
  const c = (item.value?.list_class || '').toLowerCase()
  if (c === 'success' || c === 'warning' || c === 'info' || c === 'danger') return c as 'success' | 'warning' | 'info' | 'danger'
  return 'neutral'
})
</script>
