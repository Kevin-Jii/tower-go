<template>
  <a-button
    :type="arcoType"
    :status="arcoStatus"
    :size="arcoSize"
    :html-type="nativeType"
    :disabled="disabled"
    :loading="loading"
    @click="emit('click', $event)"
  >
    <slot />
  </a-button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { BaseBtnSize, BaseBtnVariant } from './types'

const props = withDefaults(
  defineProps<{
    variant?: BaseBtnVariant
    size?: BaseBtnSize
    nativeType?: 'button' | 'submit' | 'reset'
    disabled?: boolean
    loading?: boolean
  }>(),
  {
    variant: 'secondary',
    size: 'md',
    nativeType: 'button',
    disabled: false,
    loading: false,
  },
)

const emit = defineEmits<{ click: [MouseEvent] }>()

const arcoSize = computed(() => {
  const m: Record<BaseBtnSize, 'small' | 'medium' | 'large'> = {
    sm: 'small',
    md: 'medium',
    lg: 'large',
  }
  return m[props.size]
})

const arcoType = computed(() => {
  const v = props.variant
  if (v === 'primary') return 'primary'
  if (v === 'danger') return 'primary'
  if (v === 'link') return 'text'
  if (v === 'ghost') return 'text'
  return 'secondary'
})

const arcoStatus = computed(() => {
  if (props.variant === 'danger') return 'danger'
  return 'normal'
})
</script>
