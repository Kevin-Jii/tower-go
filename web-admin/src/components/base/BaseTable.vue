<template>
  <!-- 有 minWidth 时由 Arco scroll.x 在表体内部横向滚动，便于 fixed 列；无 minWidth 时外层 overflow-x-auto 兜底 -->
  <div
    ref="tableRoot"
    class="base-table-outer relative w-full min-w-0 max-w-full rounded-[var(--border-radius-large)] border border-[var(--color-border-2)] bg-[var(--color-bg-2)]"
    :class="[minWidth ? 'overflow-x-hidden' : 'overflow-x-auto', height ? 'h-full min-h-0' : '']"
  >
    <div
      v-if="loading"
      class="base-table-loading-mask"
      role="status"
      aria-live="polite"
      aria-label="加载中"
    >
      <MathCurveLoader size="md" />
    </div>
    <a-table
      :columns="arcoColumns"
      :data="tableData"
      :loading="false"
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
import { computed, h, nextTick, onBeforeUnmount, onMounted, onUpdated, ref, useSlots } from 'vue'
import type { TableColumnData, TableData } from '@arco-design/web-vue/es/table/interface'
import { MathCurveLoader } from '@/components/loading'
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
const tableRoot = ref<HTMLElement | null>(null)
const actionColumnWidth = ref(0)
let measureRaf = 0

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
    const autoActionWidth = actionColumnWidth.value || undefined
    const actionWidth = col.key === 'actions' ? Math.max(mw ?? 0, w ?? 0, autoActionWidth ?? 0, 108) : undefined
    const minW = col.key === 'actions' ? actionWidth : mw !== undefined ? mw : undefined
    return {
    title: col.label,
    dataIndex: col.prop || col.key,
    width: col.key === 'actions' ? actionWidth : w,
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
      if (slot) {
        const content = slot({ row, value: val })
        if (col.key === 'actions') {
          return h('div', { class: 'base-table-actions-content' }, content)
        }
        return content as unknown as string
      }
      return formatCell(val)
    },
    }
  }
  )
)

function measureActionRowWidth(row: HTMLElement): number {
  const children = Array.from(row.children).filter((el) => {
    const node = el as HTMLElement
    return node.offsetParent !== null && window.getComputedStyle(node).display !== 'none'
  }) as HTMLElement[]
  if (!children.length) return 0

  const style = window.getComputedStyle(row)
  const gap = Number.parseFloat(style.columnGap || style.gap || '0') || 0
  const contentWidth = children.reduce((sum, child) => sum + Math.ceil(child.getBoundingClientRect().width), 0)
  return contentWidth + gap * Math.max(children.length - 1, 0)
}

function measureActionColumnWidth(): void {
  measureRaf = 0
  const root = tableRoot.value
  if (!root || typeof window === 'undefined') return

  const rows = Array.from(root.querySelectorAll<HTMLElement>('.table-row-actions'))
  const customContents = Array.from(root.querySelectorAll<HTMLElement>('.base-table-actions-content'))
  if (!rows.length && !customContents.length) {
    if (actionColumnWidth.value !== 0) actionColumnWidth.value = 0
    return
  }

  const maxRowWidth = rows.length ? Math.max(...rows.map(measureActionRowWidth)) : 0
  const maxCustomWidth = customContents.length
    ? Math.max(...customContents.map((content) => Math.ceil(content.scrollWidth)))
    : 0
  const maxContentWidth = Math.max(maxRowWidth, maxCustomWidth)
  const nextWidth = Math.ceil(maxContentWidth + 36)
  if (Number.isFinite(nextWidth) && Math.abs(nextWidth - actionColumnWidth.value) > 1) {
    actionColumnWidth.value = nextWidth
  }
}

function scheduleActionColumnMeasure(): void {
  void nextTick(() => {
    if (typeof window === 'undefined') return
    if (measureRaf) window.cancelAnimationFrame(measureRaf)
    measureRaf = window.requestAnimationFrame(measureActionColumnWidth)
  })
}

onMounted(scheduleActionColumnMeasure)
onUpdated(scheduleActionColumnMeasure)
onBeforeUnmount(() => {
  if (measureRaf && typeof window !== 'undefined') window.cancelAnimationFrame(measureRaf)
})

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

/* 操作列由内容测量自动撑宽，按钮保持单行，避免固定右列压住前面的单元格 */
:deep(td.base-table-actions-cell) {
  white-space: nowrap !important;
}
:deep(td.base-table-actions-cell > .arco-table-cell) {
  overflow: visible;
}
:deep(td.base-table-actions-cell > .arco-table-cell > div) {
  width: 100%;
  display: flex;
  flex-wrap: nowrap !important;
  justify-content: flex-end;
  gap: 8px;
}
:deep(.base-table-actions-content) {
  display: inline-flex;
  width: max-content;
  min-width: max-content;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  white-space: nowrap;
}
:deep(td.base-table-actions-cell .flex-nowrap) {
  flex-wrap: nowrap !important;
}
:deep(td.base-table-actions-cell .whitespace-nowrap) {
  white-space: nowrap !important;
}
:deep(td.base-table-actions-cell .shrink-0) {
  flex-shrink: 0 !important;
}
:deep(td.base-table-actions-cell .arco-select),
:deep(td.base-table-actions-cell .arco-input-wrapper) {
  max-width: 100%;
}

.base-table-loading-mask {
  position: absolute;
  inset: 0;
  z-index: 4;
  display: flex;
  align-items: center;
  justify-content: center;
  background: color-mix(in srgb, var(--color-bg-2) 82%, transparent);
  backdrop-filter: blur(1px);
  border-radius: inherit;
}
</style>
