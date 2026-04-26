<template>
  <a-switch
    :model-value="modelValue"
    :disabled="disabled"
    :checked-value="activeValue"
    :unchecked-value="inactiveValue"
    @update:model-value="emit('update:modelValue', $event as number | boolean)"
  />
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    modelValue: number | boolean
    activeValue?: number | boolean
    inactiveValue?: number | boolean
    disabled?: boolean
  }>(),
  {
    activeValue: 1,
    inactiveValue: 0,
    disabled: false,
  },
)

const emit = defineEmits<{ 'update:modelValue': [number | boolean] }>()
</script>
