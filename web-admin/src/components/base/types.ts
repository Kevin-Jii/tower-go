export type BaseBtnVariant = 'primary' | 'secondary' | 'danger' | 'ghost' | 'link'
export type BaseBtnSize = 'sm' | 'md' | 'lg'

export interface BaseTableColumn {
  /** 列唯一键，用于插槽 `#cell-{key}` */
  key: string
  label: string
  /** 从行对象取值的路径，如 `role.name` */
  prop?: string
  width?: string
  minWidth?: string
  fixed?: 'left' | 'right'
  align?: 'left' | 'center' | 'right'
  /** 单行省略 */
  ellipsis?: boolean
}

export interface BaseTreeNode {
  id: number | string
  [key: string]: unknown
}

export interface BaseSelectOption {
  label: string
  value: string | number
}

/** 表格行操作（配合 BaseTableRowActions：控制外露数量，避免操作列撑破布局） */
export interface TableRowAction {
  label: string
  onClick: () => void
  /** 缺省或空字符串表示不校验权限 */
  permission?: string
  danger?: boolean
  disabled?: boolean
  /**
   * 按钮展示位置：
   * - auto: 按 maxInline 自动分配
   * - inline: 强制外露
   * - more: 强制进入「更多」
   */
  place?: 'auto' | 'inline' | 'more'
}
