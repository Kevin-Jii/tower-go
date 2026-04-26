<template>
  <a-textarea
    :model-value="modelValue"
    :placeholder="placeholder"
    :disabled="disabled"
    :auto-size="{ minRows: rows, maxRows: rows + 12 }"
    allow-clear
    @update:model-value="emit('update:modelValue', $event)"
  />
</template>

<script setup lang="ts">
const emit = defineEmits<{ 'update:modelValue': [string] }>()

withDefaults(
  defineProps<{
    modelValue: string
    rows?: number
    placeholder?: string
    disabled?: boolean
  }>(),
  { rows: 3, disabled: false },
)
</script>
