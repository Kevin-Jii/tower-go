<template>
  <!-- 有 minWidth 时由 Arco scroll.x 在表体内部横向滚动，便于 fixed 列；无 minWidth 时外层 overflow-x-auto 兜底 -->
  <div
    class="base-table-outer w-full min-w-0 max-w-full rounded-[var(--border-radius-large)] border border-[var(--color-border-2)] bg-[var(--color-bg-2)]"
    :class="[minWidth ? 'overflow-x-hidden' : 'overflow-x-auto', height ? 'h-full min-h-0' : '']"
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
      @row-dblclick="onRowDblclick"
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
    minWidth: '',
    rowKey: 'id',
    treeDefaultExpandAll: true,
    rowClickable: false,
  },
)

const emit = defineEmits<{
  'row-click': [Record<string, unknown>]
  'row-dblclick': [Record<string, unknown>]
}>()
const slots = useSlots()

const expandAll = computed(() => !!(props.treeChildrenKey && props.treeDefaultExpandAll))

const scroll = computed(() => {
  const s: { x?: string | number; y?: string | number } = {}
  if (props.minWidth) s.x = props.minWidth
  if (props.height) s.y = props.height
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
  props.columns.map((col) => {
    const w = parseSize(col.width)
    const mw = parseSize(col.minWidth)
    /** 操作列：略给最小宽度；具体按钮由 BaseTableRowActions 收敛为「更多」 */
    const minW =
      col.key === 'actions' ? Math.max(mw ?? 0, w ?? 0, 108) : mw !== undefined ? mw : undefined
    return {
    title: col.label,
    dataIndex: col.prop || col.key,
    width: w,
    minWidth: minW,
    fixed: col.fixed,
    align: col.align,
    className: col.key === 'actions' ? 'base-table-actions-cell' : undefined,
    ellipsis: col.ellipsis,
    tooltip: col.ellipsis ? true : undefined,
    render: ({ record }) => {
      const row = record as Record<string, unknown>
      const val = cellValue(row, col)
      const slot = slots[`cell-${col.key}`]
      if (slot) return slot({ row, value: val }) as unknown as string
      return formatCell(val)
    },
    }
  }
  )
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

function onRowDblclick(record: TableData): void {
  emit('row-dblclick', record as Record<string, unknown>)
}
</script>

<style scoped>
:deep(.base-table-row--highlight) > td {
  background-color: color-mix(in srgb, rgb(var(--primary-6)) 12%, var(--color-bg-2));
}
:deep(.base-table-row--clickable) {
  cursor: pointer;
}

/* 统一兜底：操作列按钮过多时允许自动换行，避免被固定列宽挤压遮挡 */
:deep(td.base-table-actions-cell) {
  white-space: normal !important;
}
:deep(td.base-table-actions-cell > .arco-table-cell) {
  overflow: hidden;
}
:deep(td.base-table-actions-cell > .arco-table-cell > *) {
  max-width: 100%;
}
:deep(td.base-table-actions-cell > .arco-table-cell > div) {
  width: 100%;
  display: flex;
  flex-wrap: wrap !important;
  justify-content: flex-end;
  gap: 4px 8px;
}
:deep(td.base-table-actions-cell .flex-nowrap) {
  flex-wrap: wrap !important;
}
:deep(td.base-table-actions-cell .whitespace-nowrap) {
  white-space: normal !important;
}
:deep(td.base-table-actions-cell .shrink-0) {
  flex-shrink: 1 !important;
}
:deep(td.base-table-actions-cell .arco-select),
:deep(td.base-table-actions-cell .arco-input-wrapper) {
  max-width: 100%;
}
</style>
