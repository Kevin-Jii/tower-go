<template>
  <a-select
    :model-value="modelValue"
    :options="arcoOptions"
    :placeholder="placeholder"
    :disabled="disabled"
    :allow-search="searchable"
    allow-clear
    @update:model-value="emit('update:modelValue', $event as string | number | undefined)"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { BaseSelectOption } from './types'

const props = withDefaults(
  defineProps<{
    modelValue: string | number | undefined
    options: BaseSelectOption[]
    placeholder?: string
    disabled?: boolean
    searchable?: boolean
  }>(),
  { disabled: false, searchable: false },
)

const emit = defineEmits<{ 'update:modelValue': [string | number | undefined] }>()

const arcoOptions = computed(() =>
  props.options.map((o) => ({
    label: o.label,
    value: o.value,
  })),
)
</script>
