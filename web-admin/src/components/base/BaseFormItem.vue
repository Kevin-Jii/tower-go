<template>
  <a-form-item :label="label" :required="required" :extra="hint || undefined">
    <slot />
  </a-form-item>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    label?: string
    required?: boolean
    hint?: string
  }>(),
  { required: false },
)
</script>
