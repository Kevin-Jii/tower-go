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
