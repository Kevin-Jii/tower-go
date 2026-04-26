<template>
  <div
    class="base-table-wrap min-w-0 overflow-hidden rounded-[var(--border-radius-large)] border border-[var(--color-border-2)] bg-[var(--color-bg-2)]"
    :style="wrapStyle"
  >
    <a-table
      :columns="arcoColumns"
      :data="tableData"
      :loading="loading"
      :row-key="rowKey"
      :pagination="false"
      :bordered="{ wrapper: true, cell: true }"
      :scroll="scroll"
      :row-class="rowClassFn"
      :default-expand-all-rows="expandAll"
      size="small"
      table-layout-fixed
      @row-click="onRowClick"
    />
  </div>
</template>

<script setup lang="ts">
import { computed, useSlots } from 'vue'
import type { TableColumnData, TableData } from '@arco-design/web-vue/es/table/interface'
import type { BaseTableColumn } from './types'

const props = withDefaults(
  defineProps<{
    columns: BaseTableColumn[]
    data: Record<string, unknown>[]
    loading?: boolean
    minWidth?: string
    rowKey?: string
    height?: string
    treeChildrenKey?: string
    treeDefaultExpandAll?: boolean
    highlightRowKey?: string | number | null
    rowClickable?: boolean
  }>(),
  {
    loading: false,
    minWidth: '640px',
    rowKey: 'id',
    treeDefaultExpandAll: true,
    rowClickable: false,
  },
)

const emit = defineEmits<{ 'row-click': [Record<string, unknown>] }>()
const slots = useSlots()

const expandAll = computed(() => !!(props.treeChildrenKey && props.treeDefaultExpandAll))

const wrapStyle = computed(() => {
  if (!props.minWidth) return {}
  return { minWidth: props.minWidth }
})

const scroll = computed(() => {
  const s: { y?: string | number; minWidth?: string | number } = {}
  if (props.height) s.y = props.height
  if (props.minWidth) s.minWidth = props.minWidth
  return Object.keys(s).length ? s : undefined
})

function mapNode(row: Record<string, unknown>, ck: string): TableData {
  const out = { ...row } as TableData
  const ch = row[ck]
  if (Array.isArray(ch) && ch.length) {
    out.children = (ch as Record<string, unknown>[]).map((c) => mapNode(c, ck))
  }
  if (ck !== 'children') {
    delete (out as Record<string, unknown>)[ck]
  }
  return out
}

const tableData = computed<TableData[]>(() => {
  const rows = props.data ?? []
  const ck = props.treeChildrenKey
  if (!ck) return rows as TableData[]
  return rows.map((r) => mapNode(r, ck))
})

function rowKeyVal(row: Record<string, unknown>): string | number {
  const v = row[props.rowKey]
  if (typeof v === 'number' || typeof v === 'string') return v
  return String(v)
}

function getByPath(obj: unknown, path?: string): unknown {
  if (!path) return undefined
  return path.split('.').reduce((acc: unknown, key) => {
    if (acc && typeof acc === 'object' && key in (acc as object)) return (acc as Record<string, unknown>)[key]
    return undefined
  }, obj)
}

function cellValue(row: Record<string, unknown>, col: BaseTableColumn): unknown {
  if (col.prop) return getByPath(row, col.prop)
  return row[col.key]
}

function formatCell(v: unknown): string {
  if (v === null || v === undefined) return ''
  if (typeof v === 'object') return ''
  return String(v)
}

function parseSize(w?: string): number | undefined {
  if (!w) return undefined
  const t = w.trim()
  const px = /^([\d.]+)px$/i.exec(t)
  if (px) return Number(px[1])
  const n = Number(t)
  return Number.isFinite(n) ? n : undefined
}

const arcoColumns = computed<TableColumnData[]>(() =>
  props.columns.map((col) => ({
    title: col.label,
    dataIndex: col.prop || col.key,
    width: parseSize(col.width),
    minWidth: parseSize(col.minWidth),
    fixed: col.fixed,
    align: col.align,
    ellipsis: col.ellipsis,
    tooltip: col.ellipsis ? true : undefined,
    render: ({ record }) => {
      const row = record as Record<string, unknown>
      const val = cellValue(row, col)
      const slot = slots[`cell-${col.key}`]
      if (slot) return slot({ row, value: val }) as unknown as string
      return formatCell(val)
    },
  })),
)

function rowClassFn(record: TableData): string | string[] {
  const row = record as Record<string, unknown>
  const k = rowKeyVal(row)
  const active =
    props.highlightRowKey !== undefined && props.highlightRowKey !== null && k === props.highlightRowKey
  if (active) return 'base-table-row--highlight'
  if (props.rowClickable) return 'base-table-row--clickable'
  return ''
}

function onRowClick(record: TableData): void {
  if (!props.rowClickable) return
  emit('row-click', record as Record<string, unknown>)
}
</script>

<style scoped>
:deep(.base-table-row--highlight) > td {
  background-color: color-mix(in srgb, rgb(var(--primary-6)) 12%, var(--color-bg-2));
}
:deep(.base-table-row--clickable) {
  cursor: pointer;
}
</style>
