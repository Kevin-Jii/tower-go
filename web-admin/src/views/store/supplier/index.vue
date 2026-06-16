<template>
  <div class="supplier-workspace">
    <aside class="supplier-panel">
      <div class="supplier-panel__head">
        <div>
          <h2 class="supplier-title">供应商</h2>
          <p class="supplier-subtitle">{{ list.length }} / {{ total }}</p>
        </div>
        <BaseButton v-permission="'supplier:add'" variant="primary" size="sm" @click="openCreate">新增</BaseButton>
      </div>

      <div class="supplier-search">
        <BaseInput v-model="keyword" placeholder="搜索供应商名称" clearable @enter="reload" />
        <BaseButton variant="primary" size="sm" @click="reload">查询</BaseButton>
      </div>

      <div class="supplier-list">
        <button
          v-for="item in list"
          :key="item.id"
          type="button"
          class="supplier-item"
          :class="{ 'is-active': activeSupplierId === item.id }"
          @click="pickSupplier(item)"
          @dblclick="openSupplierProfile(item)"
        >
          <div class="supplier-item__body">
            <div class="supplier-item__top">
              <span class="supplier-item__name">{{ item.supplier_name }}</span>
              <span class="supplier-status" :class="item.status === 1 ? 'is-on' : 'is-off'">{{ supplierStatusLabel(item.status) }}</span>
            </div>
            <div class="supplier-item__meta">
              {{ item.contact_person || '未填联系人' }}<span v-if="item.contact_phone"> · {{ item.contact_phone }}</span>
            </div>
          </div>
          <div class="supplier-item__actions" @click.stop>
            <BaseButton v-permission="'supplier:edit'" variant="link" size="sm" @click="openEdit(item)">编辑</BaseButton>
            <BaseButton v-permission="'supplier:delete'" variant="link" size="sm" class="!text-rose-600" @click="onDelete(item)">删除</BaseButton>
          </div>
        </button>

        <div v-if="!loading && list.length === 0" class="supplier-empty">暂无供应商</div>
      </div>

      <div class="supplier-pagination">
        <BasePagination
          :page="page"
          :page-size="pageSize"
          :total="total"
          @update:page="(p) => (page = p)"
          @update:page-size="(s) => (pageSize = s)"
        />
      </div>
    </aside>

    <section class="product-panel">
      <div class="product-panel__head">
        <div class="min-w-0">
          <h2 class="supplier-title truncate">{{ currentSupplierName || '商品管理' }}</h2>
          <p class="supplier-subtitle">
            <template v-if="activeSupplierId !== ''">{{ categoryRows.length }} 个分类 · {{ productRowsRaw.length }} 个商品</template>
            <template v-else>请选择左侧供应商</template>
          </p>
        </div>
        <div class="product-actions">
          <BaseButton v-permission="'supplier:add'" variant="secondary" size="sm" :disabled="activeSupplierId === ''" @click="openCategoryCreate">新增分类</BaseButton>
          <BaseButton v-permission="'supplier:add'" variant="primary" size="sm" :disabled="activeSupplierId === ''" @click="openProductCreate()">新增商品</BaseButton>
        </div>
      </div>

      <div class="product-toolbar">
        <BaseInput v-model="productKeywordInput" class="product-search" placeholder="搜索商品名称" clearable @enter="applyProductKeyword" />
        <BaseButton variant="primary" size="sm" :disabled="activeSupplierId === ''" @click="applyProductKeyword">查询</BaseButton>
        <BaseButton variant="secondary" size="sm" :disabled="activeSupplierId === ''" @click="reloadProducts">刷新</BaseButton>
      </div>

      <div v-if="activeSupplierId !== ''" class="product-table-wrap">
        <BaseTable :columns="productTreeColumns" :data="(productTreeRows as unknown) as Record<string, unknown>[]"
          :loading="productsLoading || categoriesLoading" height="calc(100vh - 285px)" row-key="id"
          tree-children-key="children" :tree-default-expand-all="true" class="min-h-0 flex-1">
          <template #cell-name="{ row }">
            <span v-if="(row as TreeRow).isCategory" class="font-semibold text-slate-800">{{ (row as TreeRow).name }}</span>
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
                <BaseButton variant="link" size="sm" @click="openProductDrawer((row as TreeRow).productId!)">查看</BaseButton>
                <BaseButton v-permission="'supplier:delete'" variant="link" size="sm"
                  @click="onDeleteProduct((row as TreeRow).raw!)">删除</BaseButton>
              </template>
            </div>
          </template>
        </BaseTable>
      </div>
      <div v-else class="supplier-empty product-empty">请选择供应商后维护商品</div>
    </section>

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

    <BaseDialog v-model="productDlg" title="新增供应商商品（多规格）" max-width="min(1100px, 96vw)">
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
          <div class="flex items-center justify-between mb-3">
            <h4 class="m-0 text-sm font-semibold">规格价格</h4>
            <BaseButton variant="secondary" size="sm" @click="addCreateUnit">添加规格</BaseButton>
          </div>
          <div class="unit-spec-editor">
            <div class="unit-spec-editor__head">
              <span>单位编码</span>
              <span>规格名称</span>
              <span>换算基础量</span>
              <span>成本价</span>
              <span>销售价</span>
              <span>操作</span>
            </div>
            <div v-for="(u, idx) in createUnits" :key="idx" class="unit-spec-editor__row">
              <BaseSelect v-model="u.unit_code" class="unit-spec-editor__control" :options="unitOptions" />
              <BaseInput v-model="u.unit_name" class="unit-spec-editor__control" placeholder="如 2L桶" />
              <BaseNumberInput v-model="u.factor_to_base" class="unit-spec-editor__control" :min="0.000001" :step="0.01" />
              <BaseNumberInput v-model="u.cost_price" class="unit-spec-editor__control" :min="0" :step="0.01" />
              <BaseNumberInput v-model="u.sale_price" class="unit-spec-editor__control" :min="0" :step="0.01" />
              <BaseButton variant="ghost" size="sm" :disabled="createUnits.length <= 1" @click="removeCreateUnit(idx)">移除</BaseButton>
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="productDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="productSaving" @click="submitCreateProduct">保存商品配置</BaseButton>
      </template>
    </BaseDialog>

    <a-drawer :visible="productDrawer" placement="right" width="min(1080px, 96vw)" :mask-closable="true"
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
          <div class="mb-3 flex justify-end">
            <BaseButton variant="secondary" size="sm" @click="addEditUnit">添加规格</BaseButton>
          </div>
          <div class="unit-spec-editor">
            <div class="unit-spec-editor__head">
              <span>单位编码</span>
              <span>规格名称</span>
              <span>换算基础量</span>
              <span>成本价</span>
              <span>销售价</span>
              <span>操作</span>
            </div>
            <div v-for="(u, idx) in editUnits" :key="idx" class="unit-spec-editor__row">
              <BaseSelect v-model="u.unit_code" class="unit-spec-editor__control" :options="unitOptions" />
              <BaseInput v-model="u.unit_name" class="unit-spec-editor__control" placeholder="如 2L桶" />
              <BaseNumberInput v-model="u.factor_to_base" class="unit-spec-editor__control" :min="0.000001" :step="0.01" />
              <BaseNumberInput v-model="u.cost_price" class="unit-spec-editor__control" :min="0" :step="0.01" />
              <BaseNumberInput v-model="u.sale_price" class="unit-spec-editor__control" :min="0" :step="0.01" />
              <BaseButton variant="ghost" size="sm" :disabled="editUnits.length <= 1" @click="removeEditUnit(idx)">移除</BaseButton>
            </div>
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
import { useRouter } from 'vue-router'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import {
  BaseButton,
  BaseDialog,
  BaseFormItem,
  BaseInput,
  BaseNumberInput,
  BasePagination,
  BaseSelect,
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

const activeSupplierId = ref<number | ''>('')
const productKeywordInput = ref('')
const productKeyword = ref('')

function pickSupplier(row: Supplier): void {
  activeSupplierId.value = row.id
}

function openSupplierProfile(row: Supplier): void {
  void router.push(`/public/supplier/${row.id}`)
}

function supplierStatusLabel(status?: number): string {
  return Number(status) === 0 ? '禁用' : '启用'
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
type UnitFormLine = {
  unit_code: string
  unit_name: string
  factor_to_base: number
  precision: number
  cost_price: number
  sale_price: number
}

function makeUnitLine(overrides: Partial<UnitFormLine> = {}): UnitFormLine {
  return {
    unit_code: '',
    unit_name: '',
    factor_to_base: 1,
    precision: 0,
    cost_price: 0,
    sale_price: 0,
    ...overrides,
  }
}

const createUnits = ref<UnitFormLine[]>([makeUnitLine()])

function addCreateUnit(): void {
  createUnits.value.push(makeUnitLine())
}

function removeCreateUnit(idx: number): void {
  createUnits.value = createUnits.value.filter((_, i) => i !== idx)
  if (!createUnits.value.length) createUnits.value.push(makeUnitLine())
}

function addEditUnit(): void {
  editUnits.value.push(makeUnitLine())
}

function removeEditUnit(idx: number): void {
  editUnits.value = editUnits.value.filter((_, i) => i !== idx)
  if (!editUnits.value.length) editUnits.value.push(makeUnitLine())
}

function resetProductForm(categoryId?: number): void {
  productForm.category_id = categoryId
  productForm.name = ''
  productForm.spec = ''
  createUnits.value = [makeUnitLine()]
}

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

function normalizeUnitLines(lines: UnitFormLine[]): UnitFormLine[] {
  return lines
    .map((u) => ({
      unit_code: String(u.unit_code || '').trim(),
      unit_name: String(u.unit_name || '').trim(),
      factor_to_base: Number(u.factor_to_base || 0),
      precision: Number(u.precision || 0),
      cost_price: Number(u.cost_price || 0),
      sale_price: Number(u.sale_price || 0),
    }))
    .filter((u) => u.unit_code && u.factor_to_base > 0)
}

function duplicatedUnitLine(units: UnitFormLine[]): UnitFormLine | null {
  const seen = new Set<string>()
  for (const u of units) {
    const name = u.unit_name || unitNameByCode(u.unit_code)
    const key = `${u.unit_code.trim().toLowerCase()}\u0000${name.trim().toLowerCase()}`
    if (seen.has(key)) return u
    seen.add(key)
  }
  return null
}

async function submitCreateProduct(): Promise<void> {
  if (activeSupplierId.value === '' || !productForm.category_id || !productForm.name.trim()) {
    toast.warning('请先选择供应商并填写分类/商品名称')
    return
  }
  const units = normalizeUnitLines(createUnits.value)
  if (!units.length) {
    toast.warning('请至少配置一个规格')
    return
  }
  if (duplicatedUnitLine(units)) {
    toast.warning('同一商品下规格编码和规格名称不能重复，请填写如 1L桶、2L桶')
    return
  }

  productSaving.value = true
  try {
    const defaultUnit = units[0]
    await createSupplierProduct({
      supplier_id: activeSupplierId.value as number,
      category_id: productForm.category_id,
      name: productForm.name.trim(),
      unit: defaultUnit.unit_name || unitNameByCode(defaultUnit.unit_code),
      bottle_price: Number(defaultUnit.sale_price) || 0,
      case_price: Number(defaultUnit.sale_price) || 0,
      bottles_per_case: Math.max(1, Math.round(defaultUnit.factor_to_base || 1)),
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

    await batchUpsertProductUnitSpecs({
      product_id: created.id,
      units: units.map((u) => ({
        unit_code: u.unit_code,
        unit_name: u.unit_name || unitNameByCode(u.unit_code),
        factor_to_base: u.factor_to_base,
        precision: u.precision,
        cost_price: u.cost_price,
        sale_price: u.sale_price,
        is_enabled: true,
      })),
    })

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
  UnitFormLine[]
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
        unit_name: u.unit_name || unitNameByCode(u.unit_code),
        factor_to_base: Number(u.factor_to_base),
        precision: Number(u.precision),
        cost_price: Number(u.cost_price),
        sale_price: Number(u.sale_price),
      })) ?? []
    if (!editUnits.value.length) {
      editUnits.value = [
        makeUnitLine({
          unit_code: unitOptions.value[0]?.value != null ? String(unitOptions.value[0].value) : '',
          unit_name: p.unit ?? '',
          sale_price: Number(p.bottle_price ?? p.price ?? 0),
        }),
      ]
    }
    productDrawer.value = true
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载商品详情失败')
  }
}

async function submitEditProduct(): Promise<void> {
  if (!productEditingId.value) return
  const units = normalizeUnitLines(editUnits.value)
  if (!units.length) {
    toast.warning('请至少配置一个规格')
    return
  }
  if (duplicatedUnitLine(units)) {
    toast.warning('同一商品下规格编码和规格名称不能重复，请填写如 1L桶、2L桶')
    return
  }
  const defaultUnit = units[0]
  productEditSaving.value = true
  try {
    await updateSupplierProduct(productEditingId.value, {
      name: productEdit.name.trim(),
      spec: productEdit.spec.trim(),
      bottle_price: Number(defaultUnit.sale_price) || productEdit.bottle_price,
      case_price: Number(defaultUnit.sale_price) || productEdit.case_price,
      bottles_per_case: Math.max(1, Math.round(defaultUnit.factor_to_base || productEdit.bottles_per_case || 1)),
      unit: defaultUnit.unit_name || unitNameByCode(defaultUnit.unit_code) || productEdit.unit.trim(),
    })
    await batchUpsertProductUnitSpecs({
      product_id: productEditingId.value,
      units: units.map((u) => ({
        unit_code: u.unit_code,
        unit_name: u.unit_name || unitNameByCode(u.unit_code),
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
.supplier-workspace {
  display: grid;
  grid-template-columns: minmax(280px, 360px) minmax(0, 1fr);
  gap: 16px;
  height: calc(100vh - 100px);
  min-height: 560px;
  overflow: hidden;
}

.supplier-panel,
.product-panel {
  min-width: 0;
  min-height: 0;
  border: 1px solid var(--color-border-2);
  border-radius: var(--border-radius-large);
  background: var(--color-bg-2);
}

.supplier-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
}

.product-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
}

.supplier-panel__head,
.product-panel__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-height: 36px;
}

.supplier-title {
  margin: 0;
  color: #0f172a;
  font-size: 18px;
  font-weight: 650;
  line-height: 24px;
}

.supplier-subtitle {
  margin: 4px 0 0;
  color: #64748b;
  font-size: 12px;
  line-height: 18px;
}

.supplier-search,
.product-toolbar,
.product-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.supplier-search {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
}

.product-toolbar {
  flex-wrap: wrap;
}

.product-search {
  width: min(280px, 100%);
}

.supplier-list {
  display: flex;
  min-height: 0;
  flex: 1;
  flex-direction: column;
  gap: 8px;
  overflow-y: auto;
  overflow-x: hidden;
  padding-right: 2px;
}

.supplier-item {
  display: flex;
  width: 100%;
  align-items: center;
  gap: 10px;
  padding: 11px 12px;
  border: 1px solid transparent;
  border-radius: 8px;
  background: transparent;
  color: inherit;
  cursor: pointer;
  text-align: left;
  transition: background 0.16s ease, border-color 0.16s ease;
}

.supplier-item:hover {
  background: #f8fafc;
}

.supplier-item.is-active {
  border-color: #93c5fd;
  background: #eff6ff;
}

.supplier-item__body {
  min-width: 0;
  flex: 1;
}

.supplier-item__top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  min-width: 0;
}

.supplier-item__name {
  overflow: hidden;
  color: #0f172a;
  font-size: 14px;
  font-weight: 600;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.supplier-item__meta {
  margin-top: 4px;
  overflow: hidden;
  color: #64748b;
  font-size: 12px;
  line-height: 18px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.supplier-item__actions {
  display: flex;
  flex-shrink: 0;
  gap: 4px;
}

.supplier-status {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  height: 22px;
  padding: 0 8px;
  border-radius: 999px;
  font-size: 12px;
  line-height: 22px;
}

.supplier-status.is-on {
  color: #047857;
  background: #d1fae5;
}

.supplier-status.is-off {
  color: #64748b;
  background: #f1f5f9;
}

.supplier-pagination {
  display: flex;
  flex-shrink: 0;
  justify-content: flex-end;
}

.product-table-wrap {
  min-height: 0;
  flex: 1;
}

.supplier-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 160px;
  color: #94a3b8;
  font-size: 14px;
}

.product-empty {
  flex: 1;
  border: 1px dashed var(--color-border-2);
  border-radius: var(--border-radius-large);
}

.unit-spec-editor {
  width: 100%;
  overflow-x: auto;
}

.unit-spec-editor__head,
.unit-spec-editor__row {
  display: grid;
  grid-template-columns: minmax(150px, 1.1fr) minmax(150px, 1.1fr) minmax(130px, 0.9fr) minmax(120px, 0.8fr) minmax(120px, 0.8fr) 72px;
  gap: 12px;
  align-items: center;
  min-width: 840px;
}

.unit-spec-editor__head {
  margin-bottom: 8px;
  color: var(--color-text-2);
  font-size: 13px;
  font-weight: 600;
}

.unit-spec-editor__row {
  margin-bottom: 10px;
}

.unit-spec-editor__control {
  width: 100%;
  min-width: 0;
}

@media (max-width: 900px) {
  .supplier-workspace {
    grid-template-columns: 1fr;
    height: auto;
    min-height: 0;
    overflow: visible;
  }

  .supplier-panel {
    max-height: min(520px, 70vh);
  }

  .product-panel__head {
    align-items: flex-start;
    flex-direction: column;
  }
}
</style>
