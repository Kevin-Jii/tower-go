<template>
  <a-button
    :type="arcoType"
    :status="arcoStatus"
    :size="arcoSize"
    :html-type="nativeType"
    :disabled="disabled || loading"
    :loading="false"
    @click="emit('click', $event)"
  >
    <span v-if="loading" class="base-btn-loading">
      <MathCurveLoader size="xs" inline />
      <span v-if="$slots.default" class="base-btn-loading__text"><slot /></span>
    </span>
    <slot v-else />
  </a-button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { MathCurveLoader } from '@/components/loading'
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

<style scoped>
.base-btn-loading {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
.base-btn-loading__text {
  line-height: 1;
}
</style>
