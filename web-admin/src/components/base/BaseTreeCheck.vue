<template>
  <a-tree
    checkable
    :check-strictly="checkStrictly"
    :selectable="false"
    default-expand-all
    size="small"
    :data="(treeData as unknown) as TreeNodeData[]"
    :field-names="fieldNames"
    :checked-keys="modelValue"
    @update:checked-keys="onCheckedKeysChange"
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
    checkStrictly?: boolean
  }>(),
  {
    nodeKey: 'id',
    labelKey: 'title',
    childrenKey: 'children',
    depth: 0,
    checkStrictly: true,
  },
)

const emit = defineEmits<{ 'update:modelValue': [(string | number)[]] }>()

type CheckedKeysPayload = (string | number)[] | { checked?: (string | number)[]; halfChecked?: (string | number)[] }

/** Arco Tree 会把节点上的 icon/extra 等当作渲染函数；业务数据里常见字符串 icon（菜单图标名），会触发 RenderFunction 报错 */
const RENDER_FN_KEYS = ['icon', 'extra', 'switcherIcon', 'loadingIcon', 'dragIcon'] as const

function sanitizeTreeNodes(nodes: BaseTreeNode[]): BaseTreeNode[] {
  const ck = props.childrenKey
  return (nodes ?? []).map((node) => {
    const copy = { ...(node as Record<string, unknown>) }
    for (const k of RENDER_FN_KEYS) {
      const v = copy[k]
      if (v !== undefined && typeof v !== 'function') {
        delete copy[k]
      }
    }
    const children = copy[ck]
    if (Array.isArray(children) && children.length) {
      copy[ck] = sanitizeTreeNodes(children as BaseTreeNode[])
    }
    return copy as BaseTreeNode
  })
}

const treeData = computed(() => sanitizeTreeNodes(props.nodes ?? []))

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

function onCheckedKeysChange(payload: CheckedKeysPayload): void {
  if (Array.isArray(payload)) {
    emit('update:modelValue', payload)
    return
  }
  emit('update:modelValue', Array.isArray(payload.checked) ? payload.checked : [])
}

defineExpose({ setCheckedKeys, getCheckedKeys })
</script>
