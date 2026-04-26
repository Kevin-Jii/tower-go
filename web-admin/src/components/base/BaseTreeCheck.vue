<template>
  <a-tree
    checkable
    check-strictly
    :selectable="false"
    default-expand-all
    size="small"
    :data="(nodes as unknown) as TreeNodeData[]"
    :field-names="fieldNames"
    :checked-keys="modelValue"
    @update:checked-keys="emit('update:modelValue', $event)"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { TreeNodeData } from '@arco-design/web-vue/es/tree/interface'
import type { BaseTreeNode } from './types'

const props = withDefaults(
  defineProps<{
    nodes: BaseTreeNode[]
    nodeKey?: string
    labelKey?: string
    childrenKey?: string
    modelValue: (string | number)[]
    depth?: number
  }>(),
  {
    nodeKey: 'id',
    labelKey: 'title',
    childrenKey: 'children',
    depth: 0,
  },
)

const emit = defineEmits<{ 'update:modelValue': [(string | number)[]] }>()

const fieldNames = computed(() => ({
  key: props.nodeKey,
  title: props.labelKey,
  children: props.childrenKey,
}))

function setCheckedKeys(keys: (string | number)[]): void {
  emit('update:modelValue', [...keys])
}

function getCheckedKeys(): (string | number)[] {
  return [...props.modelValue]
}

defineExpose({ setCheckedKeys, getCheckedKeys })
</script>
