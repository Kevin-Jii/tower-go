<template>
  <a-modal
    :visible="modelValue"
    :width="maxWidth"
    :mask-closable="closeOnBackdrop"
    :footer="hasFooter"
    :title="titleForProp"
    unmount-on-close
    @update:visible="onVisibleChange"
    @open="emit('open')"
  >
    <template v-if="slots.title" #title>
      <slot name="title" />
    </template>
    <slot />
    <template v-if="hasFooter" #footer>
      <slot name="footer" />
    </template>
  </a-modal>
</template>

<script setup lang="ts">
import { computed, useSlots } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue: boolean
    title?: string
    maxWidth?: string
    closeOnBackdrop?: boolean
  }>(),
  {
    title: '',
    maxWidth: 'min(520px, 96vw)',
    closeOnBackdrop: true,
  },
)

const emit = defineEmits<{ 'update:modelValue': [boolean]; open: []; close: [] }>()

const slots = useSlots()
const hasFooter = computed(() => !!slots.footer)
const titleForProp = computed(() => (slots.title ? undefined : props.title || undefined))

function onVisibleChange(v: boolean): void {
  emit('update:modelValue', v)
  if (!v) emit('close')
}
</script>
