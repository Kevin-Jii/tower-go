<template>
  <div class="supplier-split-layout grid min-h-0 h-[calc(100vh-100px)] gap-4">
    <BaseCard flush-body class="flex min-h-0 h-full min-w-0 flex-col">
      <template #header>
        <div class="flex items-center justify-between gap-2 flex-wrap w-full">
          <div class="flex flex-wrap items-center gap-2">
            <BaseInput v-model="keyword" class="w-44" placeholder="供应商名称" clearable @enter="reload" />
            <BaseButton variant="primary" size="sm" @click="reload">查询</BaseButton>
            <BaseButton v-permission="'supplier:add'" variant="primary" size="sm" @click="openCreate">新增</BaseButton>
          </div>
        </div>
      </template>
      <div class="flex min-h-0 flex-1 flex-col">
        <div ref="supplierTableHost" class="flex min-h-0 min-w-0 flex-1 flex-col overflow-x-auto overflow-y-hidden">
          <BaseTable :columns="supplierColumns" :data="(list as unknown) as Record<string, unknown>[]"
            :loading="loading" :height="supplierTableScrollY" min-width="420px" row-key="id"
            :highlight-row-key="highlightSupplierId" row-clickable class="min-h-0 flex-1" @row-click="onPickSupplierRow"
            @row-dblclick="onSupplierRowDblclick">
            <template #cell-actions="{ row }">
              <div class="flex flex-nowrap items-center justify-end gap-3 whitespace-nowrap shrink-0" @click.stop>
                <BaseButton v-permission="'supplier:edit'" variant="link" size="sm" @click="openEdit(row as Supplier)">
                  编辑
                </BaseButton>
                <BaseButton v-permission="'supplier:delete'" variant="link" size="sm"
                  @click="onDelete(row as Supplier)">删除
                </BaseButton>
              </div>
            </template>
          </BaseTable>
        </div>
        <div class="mt-3 flex shrink-0 justify-end">
          <BasePagination :page="page" :page-size="pageSize" :total="total" @update:page="(p) => (page = p)"
            @update:page-size="(s) => (pageSize = s)" />
        </div>
      </div>
    </BaseCard>

    <BaseCard flush-body class="flex min-h-0 h-full min-w-0 flex-col">
      <template #header>
        <div class="flex flex-col gap-2 w-full">
          <div class="flex gap-2 items-center">
            <BaseInput v-model="productKeywordInput" class="w-56" placeholder="商品名称" clearable
              @enter="applyProductKeyword" />
            <BaseButton variant="primary" size="sm" @click="applyProductKeyword">查询</BaseButton>
            <BaseButton v-permission="'supplier:add'" variant="primary" size="sm" @click="openCategoryCreate">新增供应商分类
            </BaseButton>
            <BaseButton variant="secondary" size="sm" @click="reloadProducts">刷新</BaseButton>
          </div>
        </div>
      </template>
      <div ref="productTableHost" class="flex min-h-0 flex-1 flex-col overflow-hidden">
        <BaseTable :columns="productTreeColumns" :data="(productTreeRows as unknown) as Record<string, unknown>[]"
          :loading="productsLoading || categoriesLoading" :height="productTableScrollY" row-key="id"
          tree-children-key="children" :tree-default-expand-all="true" class="min-h-0 flex-1">
          <template #cell-name="{ row }">
            <span v-if="(row as TreeRow).isCategory" class="font-semibold text-slate-800">{{ (row as TreeRow).name
            }}</span>
            <span v-else>{{ (row as TreeRow).name }}</span>
          </template>
          <template #cell-sale_price="{ row }">
            <span v-if="!(row as TreeRow).isCategory">{{ (row as TreeRow).sale_price ?? '-' }}</span>
            <span v-else class="text-slate-400">-</span>
          </template>
          <template #cell-actions="{ row }">
            <div class="flex flex-nowrap items-center justify-end gap-3 whitespace-nowrap shrink-0" @click.stop>
              <template v-if="(row as TreeRow).isCategory">
                <BaseButton v-permission="'supplier:add'" variant="link" size="sm"
                  @click="openProductCreate((row as TreeRow).categoryId)">新增商品</BaseButton>
              </template>
              <template v-else>
                <BaseButton variant="link" size="sm" @click="openProductDrawer((row as TreeRow).productId!)">查看
                </BaseButton>
                <BaseButton v-permission="'supplier:delete'" variant="link" size="sm"
                  @click="onDeleteProduct((row as TreeRow).raw!)">删除</BaseButton>
              </template>
            </div>
          </template>
        </BaseTable>
      </div>
    </BaseCard>

    <BaseDialog v-model="dlg" :title="isEdit ? '编辑供应商' : '新增供应商'" max-width="min(520px, 96vw)">
      <div class="space-y-4 max-h-[70vh] overflow-y-auto pr-1">
        <BaseFormItem label="名称" required>
          <BaseInput v-model="form.supplier_name" placeholder="供应商名称" />
        </BaseFormItem>
        <BaseFormItem label="联系人">
          <BaseInput v-model="form.contact_person" />
        </BaseFormItem>
        <BaseFormItem label="电话">
          <BaseInput v-model="form.contact_phone" />
        </BaseFormItem>
        <BaseFormItem label="邮箱">
          <BaseInput v-model="form.contact_email" />
        </BaseFormItem>
        <BaseFormItem label="地址">
          <BaseInput v-model="form.supplier_address" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="form.remark" :rows="2" />
        </BaseFormItem>
        <BaseFormItem v-if="isEdit" label="状态">
          <BaseSelect v-model="form.status" :options="[
            { label: '启用', value: 1 },
            { label: '禁用', value: 0 },
          ]" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="dlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="save">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="categoryDlg" title="新增供应商分类" max-width="min(420px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="分类名称" required>
          <BaseInput v-model="categoryForm.name" />
        </BaseFormItem>
        <BaseFormItem label="排序">
          <BaseNumberInput v-model="categoryForm.sort" :min="1" :step="1" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="categoryDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="categorySaving" @click="submitCategory">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="productDlg" title="新增供应商商品（大/小规格）" max-width="min(720px, 96vw)">
      <div class="space-y-4 max-h-[75vh] overflow-y-auto pr-1">
        <p class="m-0 text-sm text-slate-600">
          供应商：<span class="font-medium">{{ currentSupplierName || '-' }}</span>
        </p>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <BaseFormItem label="分类" required>
            <BaseSelect v-model="productForm.category_id" :options="categoryOptions" placeholder="请选择分类" />
          </BaseFormItem>
          <BaseFormItem label="商品名称" required>
            <BaseInput v-model="productForm.name" />
          </BaseFormItem>
        </div>
        <BaseFormItem label="规格描述">
          <BaseInput v-model="productForm.spec" placeholder="如 500ml*24" />
        </BaseFormItem>

        <div class="rounded border border-[var(--color-border-2)] p-3">
          <h4 class="m-0 mb-3 text-sm font-semibold">基础单位（小规格）</h4>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
            <BaseFormItem label="单位编码" required>
              <BaseSelect v-model="smallSpec.unit_code" :options="unitOptions" placeholder="product_unit" />
            </BaseFormItem>
            <BaseFormItem label="换算系数(=1)">
              <BaseNumberInput v-model="smallSpec.factor_to_base" :min="1" :max="1" :step="1" />
            </BaseFormItem>
            <BaseFormItem label="数量精度(0~6)">
              <BaseNumberInput v-model="smallSpec.precision" :min="0" :max="6" :step="1" />
            </BaseFormItem>
            <BaseFormItem label="成本价">
              <BaseNumberInput v-model="smallSpec.cost_price" :min="0" :step="0.01" />
            </BaseFormItem>
            <BaseFormItem label="销售价">
              <BaseNumberInput v-model="smallSpec.sale_price" :min="0" :step="0.01" />
            </BaseFormItem>
          </div>
        </div>

        <div class="rounded border border-[var(--color-border-2)] p-3">
          <div class="flex items-center justify-between mb-3">
            <h4 class="m-0 text-sm font-semibold">大规格单位</h4>
            <BaseSwitch v-model="enableLargeSpec" />
          </div>
          <div v-if="enableLargeSpec" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
            <BaseFormItem label="单位编码" required>
              <BaseSelect v-model="largeSpec.unit_code" :options="unitOptions" placeholder="product_unit" />
            </BaseFormItem>
            <BaseFormItem label="换算系数(>1)" required>
              <BaseNumberInput v-model="largeSpec.factor_to_base" :min="2" :step="1" />
            </BaseFormItem>
            <BaseFormItem label="数量精度(0~6)">
              <BaseNumberInput v-model="largeSpec.precision" :min="0" :max="6" :step="1" />
            </BaseFormItem>
            <BaseFormItem label="成本价">
              <BaseNumberInput v-model="largeSpec.cost_price" :min="0" :step="0.01" />
            </BaseFormItem>
            <BaseFormItem label="销售价">
              <BaseNumberInput v-model="largeSpec.sale_price" :min="0" :step="0.01" />
            </BaseFormItem>
          </div>
        </div>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="productDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="productSaving" @click="submitCreateProduct">保存商品配置</BaseButton>
      </template>
    </BaseDialog>

    <a-drawer :visible="productDrawer" placement="right" :width="560" :mask-closable="true"
      @cancel="productDrawer = false">
      <template #title>商品详情与编辑</template>
      <div class="space-y-4">
        <BaseFormItem label="商品名称">
          <BaseInput v-model="productEdit.name" />
        </BaseFormItem>
        <BaseFormItem label="规格描述">
          <BaseInput v-model="productEdit.spec" />
        </BaseFormItem>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <BaseFormItem label="单价(瓶)">
            <BaseNumberInput v-model="productEdit.bottle_price" :min="0" :step="0.01" />
          </BaseFormItem>
          <BaseFormItem label="单价(箱)">
            <BaseNumberInput v-model="productEdit.case_price" :min="0" :step="0.01" />
          </BaseFormItem>
          <BaseFormItem label="每箱瓶数">
            <BaseNumberInput v-model="productEdit.bottles_per_case" :min="1" :step="1" />
          </BaseFormItem>
          <BaseFormItem label="单位显示">
            <BaseInput v-model="productEdit.unit" />
          </BaseFormItem>
        </div>
        <div class="rounded border border-[var(--color-border-2)] p-3">
          <h4 class="m-0 mb-3 text-sm font-semibold">单位配置</h4>
          <div v-for="(u, idx) in editUnits" :key="idx" class="grid grid-cols-1 md:grid-cols-2 gap-3 mb-3">
            <BaseFormItem label="单位编码">
              <BaseSelect v-model="u.unit_code" :options="unitOptions" />
            </BaseFormItem>
            <BaseFormItem label="换算系数">
              <BaseNumberInput v-model="u.factor_to_base" :min="1" :step="1" />
            </BaseFormItem>
            <BaseFormItem label="成本价">
              <BaseNumberInput v-model="u.cost_price" :min="0" :step="0.01" />
            </BaseFormItem>
            <BaseFormItem label="销售价">
              <BaseNumberInput v-model="u.sale_price" :min="0" :step="0.01" />
            </BaseFormItem>
          </div>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <BaseButton variant="ghost" @click="productDrawer = false">取消</BaseButton>
          <BaseButton variant="primary" :loading="productEditSaving" @click="submitEditProduct">保存</BaseButton>
        </div>
      </template>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useElementSize } from '@vueuse/core'
import { useRouter } from 'vue-router'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import {
  BaseButton,
  BaseCard,
  BaseDialog,
  BaseFormItem,
  BaseInput,
  BaseNumberInput,
  BasePagination,
  BaseSelect,
  BaseSwitch,
  BaseTable,
  BaseTextarea,
} from '@/components/base'
import type { BaseSelectOption, BaseTableColumn } from '@/components/base/types'
import { createSupplier, deleteSupplier, listSuppliers, updateSupplier } from '@/api/supplier'
import {
  batchUpsertProductUnitSpecs,
  createSupplierCategory,
  createSupplierProduct,
  deleteSupplierProduct,
  getSupplierProduct,
  listProductUnitSpecs,
  listSupplierProducts,
  listSupplierCategories,
  updateSupplierProduct,
} from '@/api/supplierProduct'
import { listDictDataByTypeCode } from '@/api/dict'
import type { DictData, ProductUnitSpec, StorePurchasableProduct, Supplier } from '@/api/types'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'

interface TreeRow {
  id: string | number
  isCategory: boolean
  categoryId?: number
  productId?: number
  name: string
  sale_price?: number | string
  children?: TreeRow[]
  raw?: StorePurchasableProduct
}

const qc = useQueryClient()
const router = useRouter()

/** 表格区域撑满卡片剩余高度，用 Arco scroll.y 占满可视区 */
const supplierTableHost = ref<HTMLElement | null>(null)
const productTableHost = ref<HTMLElement | null>(null)
const { height: supplierHostH } = useElementSize(supplierTableHost)
const { height: productHostH } = useElementSize(productTableHost)
const TABLE_HEAD_RESERVE = 48
function scrollYpx(h: number): string {
  const y = Math.floor(h - TABLE_HEAD_RESERVE)
  return `${Math.max(120, y)}px`
}
const supplierTableScrollY = computed(() => scrollYpx(supplierHostH.value))
const productTableScrollY = computed(() => scrollYpx(productHostH.value))
const keyword = ref('')
const page = ref(1)
const pageSize = ref(10)
const queryKey = computed(() => ['suppliers', page.value, pageSize.value, keyword.value.trim()] as const)

const { data: pageData, isLoading: loading } = useQuery({
  queryKey,
  queryFn: () =>
    listSuppliers({
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value.trim() || undefined,
    }),
})

const list = computed(() => pageData.value?.list ?? [])
const total = computed(() => pageData.value?.total ?? 0)
const currentSupplierName = computed(
  () => list.value.find((s) => s.id === activeSupplierId.value)?.supplier_name ?? '',
)

function reload(): void {
  page.value = 1
  void qc.invalidateQueries({ queryKey: ['suppliers'] })
}

watch([page, pageSize], () => {
  void qc.invalidateQueries({ queryKey: ['suppliers'] })
})

const supplierColumns: BaseTableColumn[] = [
  { key: 'supplier_name', label: '供应商名称', prop: 'supplier_name', },
  { key: 'actions', label: '操作', width: '148px', align: 'right' },
]

const activeSupplierId = ref<number | ''>('')
const highlightSupplierId = computed(() => (activeSupplierId.value === '' ? null : (activeSupplierId.value as number)))
const productKeywordInput = ref('')
const productKeyword = ref('')

function onPickSupplierRow(row: Record<string, unknown>): void {
  const s = row as unknown as Supplier
  if (s?.id != null) activeSupplierId.value = s.id
}

/** 双击行打开供应商档案（原「查看」入口） */
function onSupplierRowDblclick(row: Record<string, unknown>): void {
  const s = row as unknown as Supplier
  if (s?.id != null) void router.push(`/public/supplier/${s.id}`)
}

function applyProductKeyword(): void {
  productKeyword.value = productKeywordInput.value.trim()
  void qc.invalidateQueries({ queryKey: ['store-supplier-products'] })
}

const categoriesQueryKey = computed(() => ['supplier-categories', activeSupplierId.value] as const)
const { data: categoriesData, isLoading: categoriesLoading } = useQuery({
  queryKey: categoriesQueryKey,
  queryFn: () => listSupplierCategories(activeSupplierId.value as number),
  enabled: computed(() => activeSupplierId.value !== ''),
})

const productQueryKey = computed(() => ['store-supplier-products', activeSupplierId.value, productKeyword.value] as const)
const { data: productsData, isLoading: productsLoading } = useQuery({
  queryKey: productQueryKey,
  queryFn: () =>
    listSupplierProducts({
      keyword: productKeyword.value || undefined,
      supplier_id: activeSupplierId.value as number,
    }),
  enabled: computed(() => activeSupplierId.value !== ''),
})

const categoryRows = computed(() => categoriesData.value ?? [])
const productRowsRaw = computed(() => productsData.value ?? [])

const productTreeRows = computed<TreeRow[]>(() => {
  if (activeSupplierId.value === '') return []
  const grouped = new Map<number, TreeRow>()
  for (const c of categoryRows.value) {
    grouped.set(c.id, {
      id: `cat-${c.id}`,
      isCategory: true,
      categoryId: c.id,
      name: c.name,
      children: [],
    })
  }
  for (const p of productRowsRaw.value) {
    const cid = p.category_id ?? 0
    if (!grouped.has(cid)) {
      grouped.set(cid, {
        id: `cat-${cid}`,
        isCategory: true,
        categoryId: cid,
        name: p.category?.name ?? '未分类',
        children: [],
      })
    }
    grouped.get(cid)!.children!.push({
      id: p.id,
      isCategory: false,
      productId: p.id,
      categoryId: cid,
      name: p.name,
      sale_price: p.bottle_price ?? p.price ?? 0,
      raw: p,
    })
  }
  return Array.from(grouped.values())
})

function reloadProducts(): void {
  void qc.invalidateQueries({ queryKey: ['store-supplier-products'] })
  void qc.invalidateQueries({ queryKey: ['supplier-categories'] })
}

const productTreeColumns: BaseTableColumn[] = [
  { key: 'name', label: '商品/分类', minWidth: '160px', ellipsis: true },
  { key: 'sale_price', label: '售价', width: '100px' },
  { key: 'actions', label: '操作', width: '168px', minWidth: '168px', align: 'right' },
]

const dlg = ref(false)
const saving = ref(false)
const isEdit = ref(false)
const editId = ref(0)

const form = reactive({
  supplier_name: '',
  contact_person: '',
  contact_phone: '',
  contact_email: '',
  supplier_address: '',
  remark: '',
  status: 1,
})

function openCreate(): void {
  isEdit.value = false
  editId.value = 0
  form.supplier_name = ''
  form.contact_person = ''
  form.contact_phone = ''
  form.contact_email = ''
  form.supplier_address = ''
  form.remark = ''
  form.status = 1
  dlg.value = true
}

function openEdit(row: Supplier): void {
  isEdit.value = true
  editId.value = row.id
  form.supplier_name = row.supplier_name ?? ''
  form.contact_person = row.contact_person ?? ''
  form.contact_phone = row.contact_phone ?? ''
  form.contact_email = row.contact_email ?? ''
  form.supplier_address = row.supplier_address ?? ''
  form.remark = row.remark ?? ''
  form.status = row.status === 1 ? 1 : 0
  dlg.value = true
}

async function save(): Promise<void> {
  if (!form.supplier_name.trim()) {
    toast.warning('请填写供应商名称')
    return
  }
  saving.value = true
  try {
    if (isEdit.value) {
      await updateSupplier(editId.value, {
        supplier_name: form.supplier_name.trim(),
        contact_person: form.contact_person.trim(),
        contact_phone: form.contact_phone.trim(),
        contact_email: form.contact_email.trim() || undefined,
        supplier_address: form.supplier_address.trim(),
        remark: form.remark.trim(),
        status: form.status,
      })
    } else {
      await createSupplier({
        supplier_name: form.supplier_name.trim(),
        contact_person: form.contact_person.trim(),
        contact_phone: form.contact_phone.trim(),
        contact_email: form.contact_email.trim() || undefined,
        supplier_address: form.supplier_address.trim(),
        remark: form.remark.trim(),
      })
    }
    toast.success('已保存')
    dlg.value = false
    await qc.invalidateQueries({ queryKey: ['suppliers'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function onDelete(row: Supplier): Promise<void> {
  const ok = await confirmDialog({ message: `删除供应商「${row.supplier_name}」？` })
  if (!ok) return
  try {
    await deleteSupplier(row.id)
    toast.success('已删除')
    await qc.invalidateQueries({ queryKey: ['suppliers'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}

const categoryDlg = ref(false)
const categorySaving = ref(false)
const categoryForm = reactive({ name: '', sort: 1 })

function openCategoryCreate(): void {
  if (activeSupplierId.value === '') {
    toast.warning('请先在左侧选择供应商')
    return
  }
  categoryForm.name = ''
  categoryForm.sort = 1
  categoryDlg.value = true
}

async function submitCategory(): Promise<void> {
  if (activeSupplierId.value === '' || !categoryForm.name.trim()) {
    toast.warning('请先选择供应商并填写分类名称')
    return
  }
  categorySaving.value = true
  try {
    await createSupplierCategory({
      supplier_id: activeSupplierId.value as number,
      name: categoryForm.name.trim(),
      sort: categoryForm.sort,
    })
    toast.success('分类已创建')
    categoryDlg.value = false
    await qc.invalidateQueries({ queryKey: ['supplier-categories'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '创建分类失败')
  } finally {
    categorySaving.value = false
  }
}

const unitDict = ref<DictData[]>([])
const unitOptions = computed<BaseSelectOption[]>(() => unitDict.value.map((d) => ({ label: d.label, value: d.value })))
const categoryOptions = computed<BaseSelectOption[]>(() =>
  categoryRows.value.map((c) => ({
    label: c.name,
    value: c.id,
  })),
)

const productDlg = ref(false)
const productSaving = ref(false)
const productForm = reactive({
  category_id: undefined as number | undefined,
  name: '',
  spec: '',
})
const enableLargeSpec = ref(false)
const smallSpec = reactive({
  unit_code: '',
  factor_to_base: 1,
  precision: 0,
  cost_price: 0,
  sale_price: 0,
})
const largeSpec = reactive({
  unit_code: '',
  factor_to_base: 2,
  precision: 0,
  cost_price: 0,
  sale_price: 0,
})

function resetProductForm(categoryId?: number): void {
  productForm.category_id = categoryId
  productForm.name = ''
  productForm.spec = ''
  enableLargeSpec.value = false
  smallSpec.unit_code = ''
  smallSpec.factor_to_base = 1
  smallSpec.precision = 0
  smallSpec.cost_price = 0
  smallSpec.sale_price = 0
  largeSpec.unit_code = ''
  largeSpec.factor_to_base = 2
  largeSpec.precision = 0
  largeSpec.cost_price = 0
  largeSpec.sale_price = 0
}

watch(enableLargeSpec, (enabled) => {
  if (enabled) return
  largeSpec.unit_code = ''
  largeSpec.factor_to_base = 2
  largeSpec.precision = 0
  largeSpec.cost_price = 0
  largeSpec.sale_price = 0
})

async function openProductCreate(categoryId?: number): Promise<void> {
  if (activeSupplierId.value === '') {
    toast.warning('请先在左侧选择供应商')
    return
  }
  resetProductForm(categoryId)
  try {
    if (!unitDict.value.length) {
      unitDict.value = await listDictDataByTypeCode('product_unit')
    }
    productDlg.value = true
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '初始化商品配置失败')
  }
}

function unitNameByCode(code: string): string {
  return unitDict.value.find((d) => String(d.value) === code)?.label ?? code
}

async function submitCreateProduct(): Promise<void> {
  if (activeSupplierId.value === '' || !productForm.category_id || !productForm.name.trim()) {
    toast.warning('请先选择供应商并填写分类/商品名称')
    return
  }
  if (!smallSpec.unit_code) {
    toast.warning('请设置基础单位编码')
    return
  }
  if (enableLargeSpec.value && (!largeSpec.unit_code || largeSpec.factor_to_base <= 1)) {
    toast.warning('请设置大规格单位及换算系数(>1)')
    return
  }

  productSaving.value = true
  try {
    await createSupplierProduct({
      supplier_id: activeSupplierId.value as number,
      category_id: productForm.category_id,
      name: productForm.name.trim(),
      unit: unitNameByCode(smallSpec.unit_code),
      bottle_price: Number(smallSpec.sale_price) || 0,
      case_price: enableLargeSpec.value ? Number(largeSpec.sale_price) || 0 : Number(smallSpec.sale_price) || 0,
      bottles_per_case: enableLargeSpec.value ? Math.max(2, Math.round(largeSpec.factor_to_base)) : 1,
      spec: productForm.spec.trim(),
      remark: '',
    })

    const listRes = await listSupplierProducts({
      supplier_id: activeSupplierId.value as number,
      keyword: productForm.name.trim(),
    })
    const created =
      listRes.find((p) => p.name === productForm.name.trim() && p.category_id === productForm.category_id) ?? listRes[0]
    if (!created?.id) throw new Error('商品已创建，但未获取到商品ID')

    const units = [
      {
        unit_code: smallSpec.unit_code,
        unit_name: unitNameByCode(smallSpec.unit_code),
        factor_to_base: 1,
        precision: smallSpec.precision,
        cost_price: Number(smallSpec.cost_price) || 0,
        sale_price: Number(smallSpec.sale_price) || 0,
        is_enabled: true,
      },
    ]
    if (enableLargeSpec.value) {
      units.push({
        unit_code: largeSpec.unit_code,
        unit_name: unitNameByCode(largeSpec.unit_code),
        factor_to_base: Number(largeSpec.factor_to_base),
        precision: largeSpec.precision,
        cost_price: Number(largeSpec.cost_price) || 0,
        sale_price: Number(largeSpec.sale_price) || 0,
        is_enabled: true,
      })
    }
    await batchUpsertProductUnitSpecs({ product_id: created.id, units })

    toast.success('商品与规格配置已保存')
    productDlg.value = false
    await qc.invalidateQueries({ queryKey: ['store-supplier-products'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    productSaving.value = false
  }
}

const productDrawer = ref(false)
const productEditSaving = ref(false)
const productEditingId = ref(0)
const productEdit = reactive({
  name: '',
  spec: '',
  bottle_price: 0,
  case_price: 0,
  bottles_per_case: 1,
  unit: '',
})
const editUnits = ref<
  Array<{
    unit_code: string
    factor_to_base: number
    precision: number
    cost_price: number
    sale_price: number
  }>
>([])

async function openProductDrawer(productId: number): Promise<void> {
  try {
    if (!unitDict.value.length) unitDict.value = await listDictDataByTypeCode('product_unit')
    const p = await getSupplierProduct(productId)
    const units = await listProductUnitSpecs(productId)
    productEditingId.value = productId
    productEdit.name = p.name ?? ''
    productEdit.spec = p.spec ?? ''
    productEdit.bottle_price = Number(p.bottle_price ?? 0)
    productEdit.case_price = Number(p.case_price ?? 0)
    productEdit.bottles_per_case = Number(p.bottles_per_case ?? 1)
    productEdit.unit = p.unit ?? ''
    editUnits.value =
      units.map((u: ProductUnitSpec) => ({
        unit_code: u.unit_code,
        factor_to_base: Number(u.factor_to_base),
        precision: Number(u.precision),
        cost_price: Number(u.cost_price),
        sale_price: Number(u.sale_price),
      })) ?? []
    productDrawer.value = true
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载商品详情失败')
  }
}

async function submitEditProduct(): Promise<void> {
  if (!productEditingId.value) return
  productEditSaving.value = true
  try {
    await updateSupplierProduct(productEditingId.value, {
      name: productEdit.name.trim(),
      spec: productEdit.spec.trim(),
      bottle_price: productEdit.bottle_price,
      case_price: productEdit.case_price,
      bottles_per_case: productEdit.bottles_per_case,
      unit: productEdit.unit.trim(),
    })
    await batchUpsertProductUnitSpecs({
      product_id: productEditingId.value,
      units: editUnits.value.map((u) => ({
        unit_code: u.unit_code,
        unit_name: unitNameByCode(u.unit_code),
        factor_to_base: u.factor_to_base,
        precision: u.precision,
        cost_price: u.cost_price,
        sale_price: u.sale_price,
        is_enabled: true,
      })),
    })
    toast.success('商品已更新')
    productDrawer.value = false
    await qc.invalidateQueries({ queryKey: ['store-supplier-products'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '更新失败')
  } finally {
    productEditSaving.value = false
  }
}

async function onDeleteProduct(row: StorePurchasableProduct): Promise<void> {
  const ok = await confirmDialog({ message: `删除商品「${row.name}」？` })
  if (!ok) return
  try {
    await deleteSupplierProduct(row.id)
    toast.success('商品已删除')
    await qc.invalidateQueries({ queryKey: ['store-supplier-products'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除商品失败')
  }
}
</script>

<style scoped>
/**
 * 不要用 Uno 的 grid-cols-1：它会一直写 grid-template-columns，覆盖掉这里的宽屏两列。
 * 默认单列；≥768px 左右分栏（与常用 md 断点一致）。
 */
.supplier-split-layout {
  grid-template-columns: minmax(0, 1fr);
}

@media (min-width: 768px) {
  .supplier-split-layout {
    grid-template-columns: minmax(0, 0.72fr) minmax(0, 1.28fr);
  }
}
</style>
