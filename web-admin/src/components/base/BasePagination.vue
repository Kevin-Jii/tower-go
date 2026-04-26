<template>
  <a-pagination
    class="flex flex-wrap justify-end gap-y-2"
    :total="total"
    :current="page"
    :page-size="pageSize"
    :page-size-options="pageSizes"
    show-total
    show-page-size
    @update:current="emit('update:page', $event)"
    @page-size-change="onPageSizeChange"
  />
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    page: number
    pageSize: number
    total: number
    pageSizes?: number[]
  }>(),
  { pageSizes: () => [10, 20, 50] },
)

const emit = defineEmits<{
  'update:page': [number]
  'update:pageSize': [number]
}>()

function onPageSizeChange(ps: number): void {
  emit('update:pageSize', ps)
  emit('update:page', 1)
}
</script>
