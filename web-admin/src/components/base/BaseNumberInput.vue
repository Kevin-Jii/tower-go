<template>
  <a-input-number
    :model-value="modelValue"
    :min="min"
    :max="max"
    :step="step"
    :disabled="disabled"
    :placeholder="placeholder"
    hide-button
    class="w-full"
    @update:model-value="emit('update:modelValue', $event as number | undefined)"
  />
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    modelValue: number | undefined
    min?: number
    max?: number
    step?: number
    disabled?: boolean
    placeholder?: string
  }>(),
  { step: 1, disabled: false },
)

const emit = defineEmits<{ 'update:modelValue': [number | undefined] }>()
</script>
