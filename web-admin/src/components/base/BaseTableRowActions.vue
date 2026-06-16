<template>
  <div class="table-row-actions inline-flex min-w-max items-center justify-end gap-2 whitespace-nowrap" @click.stop>
    <BaseButton
      v-for="(a, i) in inlineActions"
      :key="'i-' + i"
      variant="link"
      size="sm"
      :class="a.danger ? '!text-rose-600' : undefined"
      :disabled="a.disabled"
      @click="run(a)"
    >
      {{ a.label }}
    </BaseButton>
    <a-dropdown v-if="overflowActions.length > 0" trigger="click" position="br">
      <span class="inline-flex shrink-0 cursor-pointer">
        <BaseButton variant="link" size="sm">{{ moreLabel }}</BaseButton>
      </span>
      <template #content>
        <a-doption
          v-for="(a, i) in overflowActions"
          :key="'m-' + i"
          :disabled="a.disabled"
          @click="run(a)"
        >
          <span :class="a.danger ? 'text-rose-600' : ''">{{ a.label }}</span>
        </a-doption>
      </template>
    </a-dropdown>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import BaseButton from './BaseButton.vue'
import type { TableRowAction } from './types'
import { usePermission } from '@/hooks/usePermission'

const props = withDefaults(
  defineProps<{
    actions: TableRowAction[]
    /** 前几条直接展示为链接，其余收入「更多」 */
    maxInline?: number
    moreLabel?: string
  }>(),
  {
    maxInline: 2,
    moreLabel: '更多',
  },
)

const { hasPerm } = usePermission()

function allowed(a: TableRowAction): boolean {
  return hasPerm(a.permission ?? '')
}

const permittedActions = computed(() => props.actions.filter(allowed))

const inlineActions = computed(() => {
  const inline: TableRowAction[] = []
  for (const action of permittedActions.value) {
    if (action.place === 'more') continue
    if (action.place === 'inline') {
      inline.push(action)
      continue
    }
    if (inline.length < props.maxInline) inline.push(action)
  }
  return inline
})

const overflowActions = computed(() => {
  const overflow: TableRowAction[] = []
  const inlineSet = new Set(inlineActions.value)
  for (const action of permittedActions.value) {
    if (action.place === 'more') {
      overflow.push(action)
      continue
    }
    if (!inlineSet.has(action)) overflow.push(action)
  }
  return overflow
})

function run(a: TableRowAction): void {
  if (a.disabled) return
  a.onClick()
}
</script>
