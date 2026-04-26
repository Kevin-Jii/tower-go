<template>
  <div class="grid grid-cols-1 gap-4 xl:grid-cols-[minmax(0,0.72fr)_minmax(0,1.28fr)]">
    <BaseCard class="min-w-0">
      <template #header>
        <div class="flex items-center justify-between gap-2 flex-wrap w-full">
          <span class="font-semibold text-slate-800">供应商档案</span>
          <div class="flex flex-wrap items-center gap-2">
            <BaseInput v-model="keyword" class="w-44" placeholder="供应商名称" clearable @enter="reload" />
            <BaseButton variant="primary" size="sm" @click="reload">查询</BaseButton>
            <BaseButton v-permission="'supplier:add'" variant="primary" size="sm" @click="openCreate">新增</BaseButton>
          </div>
        </div>
      </template>
      <BaseTable
        :columns="columns"
        :data="(list as unknown) as Record<string, unknown>[]"
        :loading="loading"
        min-width="420px"
        height="calc(100vh - 325px)"
        row-key="id"
        :highlight-row-key="highlightSupplierId"
        row-clickable
        @row-click="onPickSupplierRow"
      >
        <template #cell-actions="{ row }">
          <div class="flex flex-nowrap items-center justify-end gap-3 whitespace-nowrap shrink-0" @click.stop>
            <BaseButton variant="link" size="sm" @click="openView(row as Supplier)">查看</BaseButton>
            <BaseButton v-permission="'supplier:edit'" variant="link" size="sm" @click="openEdit(row as Supplier)">编辑</BaseButton>
            <BaseButton v-permission="'supplier:delete'" variant="link" size="sm" @click="onDelete(row as Supplier)">删除</BaseButton>
          </div>
        </template>
      </BaseTable>
      <div class="flex justify-end mt-3">
        <BasePagination
          :page="page"
          :page-size="pageSize"
          :total="total"
          @update:page="(p) => (page = p)"
          @update:page-size="(s) => (pageSize = s)"
        />
      </div>
    </BaseCard>

    <BaseCard class="min-w-0">
      <template #header>
        <div class="flex flex-col gap-2 w-full">
          <div class="flex items-center justify-between gap-2 flex-wrap">
            <span class="font-semibold text-slate-800">已绑定供应商商品</span>
            <div class="flex items-center gap-2">
              <BaseButton v-permission="'supplier:add'" variant="primary" size="sm" @click="openProductCreate">新增供应商商品</BaseButton>
              <BaseButton variant="secondary" size="sm" @click="reloadProducts">刷新</BaseButton>
            </div>
          </div>
          <div class="flex gap-2 items-center">
            <BaseInput v-model="productKeywordInput" class="w-56" placeholder="商品名称" clearable @enter="applyProductKeyword" />
            <BaseButton variant="primary" size="sm" @click="applyProductKeyword">查询</BaseButton>
          </div>
        </div>
      </template>
      <BaseTable
        :columns="productColumns"
        :data="(productRows as unknown) as Record<string, unknown>[]"
        :loading="productsLoading"
        min-width="860px"
        height="calc(100vh - 290px)"
      >
        <template #cell-supplier_name="{ row }">
          {{ supplierCell(row as StorePurchasableProduct) }}
        </template>
        <template #cell-category="{ row }">
          {{ (row as StorePurchasableProduct).category?.name ?? '-' }}
        </template>
        <template #cell-spec_pair="{ row }">
          基础: {{ (row as StorePurchasableProduct).unit || '-' }} / 大规格x{{ (row as StorePurchasableProduct).bottles_per_case ?? 1 }}
        </template>
      </BaseTable>
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
          <BaseSelect
            v-model="form.status"
            :options="[
              { label: '启用', value: 1 },
              { label: '禁用', value: 0 },
            ]"
          />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="dlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="save">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="productDlg" title="新增供应商商品（大/小规格）" max-width="min(720px, 96vw)">
      <div class="space-y-4 max-h-[75vh] overflow-y-auto pr-1">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <BaseFormItem label="供应商" required>
            <BaseSelect v-model="productForm.supplier_id" :options="supplierSelectOptions" placeholder="请选择供应商" />
          </BaseFormItem>
          <BaseFormItem label="分类" required>
            <BaseSelect v-model="productForm.category_id" :options="categoryOptions" placeholder="请选择分类" />
          </BaseFormItem>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <BaseFormItem label="商品名称" required>
            <BaseInput v-model="productForm.name" />
          </BaseFormItem>
          <BaseFormItem label="规格描述">
            <BaseInput v-model="productForm.spec" placeholder="如 500ml*24" />
          </BaseFormItem>
        </div>

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
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
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
import { listPurchasableProducts, listStoreBoundSuppliers } from '@/api/storeSupplier'
import { batchUpsertProductUnitSpecs, createSupplierProduct, listSupplierCategories, listSupplierProducts } from '@/api/supplierProduct'
import { listDictDataByTypeCode } from '@/api/dict'
import type { DictData, StorePurchasableProduct, Supplier, SupplierCategory } from '@/api/types'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'

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

function reload(): void {
  page.value = 1
  void qc.invalidateQueries({ queryKey: ['suppliers'] })
}

watch([page, pageSize], () => {
  void qc.invalidateQueries({ queryKey: ['suppliers'] })
})

const columns: BaseTableColumn[] = [
  { key: 'supplier_name', label: '供应商名称', prop: 'supplier_name', minWidth: '160px', ellipsis: true },
  { key: 'actions', label: '操作', width: '220px', align: 'right' },
]

const activeSupplierId = ref<number | ''>('')
const highlightSupplierId = computed(() => (activeSupplierId.value === '' ? null : (activeSupplierId.value as number)))
const productKeywordInput = ref('')
const productKeyword = ref('')

function onPickSupplierRow(row: Record<string, unknown>): void {
  const s = row as unknown as Supplier
  if (s?.id != null) activeSupplierId.value = s.id
}

function openView(row: Supplier): void {
  void router.push(`/public/supplier/${row.id}`)
}

const { data: boundRows } = useQuery({
  queryKey: ['store-suppliers', 'bound'],
  queryFn: () => listStoreBoundSuppliers(),
})

function applyProductKeyword(): void {
  productKeyword.value = productKeywordInput.value.trim()
  void qc.invalidateQueries({ queryKey: ['store-supplier-products'] })
}

const productQueryKey = computed(() => ['store-supplier-products', activeSupplierId.value, productKeyword.value] as const)

const { data: productsData, isLoading: productsLoading } = useQuery({
  queryKey: productQueryKey,
  queryFn: () =>
    listPurchasableProducts({
      keyword: productKeyword.value || undefined,
      ...(activeSupplierId.value === '' ? {} : { supplier_id: activeSupplierId.value as number }),
    }),
})

const productRows = computed(() => productsData.value ?? [])

function reloadProducts(): void {
  activeSupplierId.value = ''
  void qc.invalidateQueries({ queryKey: ['store-supplier-products'] })
  void qc.invalidateQueries({ queryKey: ['store-suppliers', 'bound'] })
}

const productColumns: BaseTableColumn[] = [
  { key: 'id', label: 'ID', prop: 'id', width: '72px' },
  { key: 'name', label: '商品名称', prop: 'name', minWidth: '140px', ellipsis: true },
  { key: 'supplier_name', label: '供应商', minWidth: '120px', ellipsis: true },
  { key: 'category', label: '分类', width: '100px', ellipsis: true },
  { key: 'spec_pair', label: '规格', minWidth: '130px', ellipsis: true },
  { key: 'unit', label: '单位', prop: 'unit', width: '64px' },
  { key: 'bottle_price', label: '单价(瓶)', prop: 'bottle_price', width: '88px' },
  { key: 'case_price', label: '单价(箱)', prop: 'case_price', width: '88px' },
  { key: 'bottles_per_case', label: '每箱瓶数', prop: 'bottles_per_case', width: '92px' },
]

function supplierCell(p: StorePurchasableProduct): string {
  return p.supplier?.supplier_name ?? '-'
}

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
    await qc.invalidateQueries({ queryKey: ['store-suppliers', 'bound'] })
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
    await qc.invalidateQueries({ queryKey: ['store-suppliers', 'bound'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}

const supplierSelectOptions = computed<BaseSelectOption[]>(() =>
  (boundRows.value ?? []).map((b) => ({
    label: b.supplier?.supplier_name ?? `供应商 #${b.supplier_id}`,
    value: b.supplier_id,
  })),
)

const productDlg = ref(false)
const productSaving = ref(false)
const productForm = reactive({
  supplier_id: undefined as number | undefined,
  category_id: undefined as number | undefined,
  name: '',
  spec: '',
})

const categoryList = ref<SupplierCategory[]>([])
const categoryOptions = computed<BaseSelectOption[]>(() => categoryList.value.map((c) => ({ label: c.name, value: c.id })))
const unitDict = ref<DictData[]>([])
const unitOptions = computed<BaseSelectOption[]>(() => unitDict.value.map((d) => ({ label: d.label, value: d.value })))

const enableLargeSpec = ref(true)
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

function resetProductForm(): void {
  productForm.supplier_id = activeSupplierId.value === '' ? undefined : (activeSupplierId.value as number)
  productForm.category_id = undefined
  productForm.name = ''
  productForm.spec = ''
  enableLargeSpec.value = true
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

async function openProductCreate(): Promise<void> {
  resetProductForm()
  try {
    if (!unitDict.value.length) {
      unitDict.value = await listDictDataByTypeCode('product_unit')
    }
    if (productForm.supplier_id) {
      categoryList.value = await listSupplierCategories(productForm.supplier_id)
      productForm.category_id = categoryList.value[0]?.id
    } else {
      categoryList.value = []
    }
    productDlg.value = true
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '初始化商品配置失败')
  }
}

watch(
  () => productForm.supplier_id,
  async (sid) => {
    if (!sid) {
      categoryList.value = []
      productForm.category_id = undefined
      return
    }
    try {
      categoryList.value = await listSupplierCategories(sid)
      if (!categoryList.value.find((c) => c.id === productForm.category_id)) {
        productForm.category_id = categoryList.value[0]?.id
      }
    } catch {
      categoryList.value = []
      productForm.category_id = undefined
    }
  },
)

function unitNameByCode(code: string): string {
  return unitDict.value.find((d) => String(d.value) === code)?.label ?? code
}

async function submitCreateProduct(): Promise<void> {
  if (!productForm.supplier_id || !productForm.category_id || !productForm.name.trim()) {
    toast.warning('请完整填写供应商 / 分类 / 商品名称')
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
      supplier_id: productForm.supplier_id,
      category_id: productForm.category_id,
      name: productForm.name.trim(),
      unit: unitNameByCode(smallSpec.unit_code),
      bottle_price: Number(smallSpec.sale_price) || 0,
      case_price: enableLargeSpec.value ? Number(largeSpec.sale_price) || 0 : Number(smallSpec.sale_price) || 0,
      bottles_per_case: enableLargeSpec.value ? Math.max(2, Math.round(largeSpec.factor_to_base)) : 1,
      spec: productForm.spec.trim(),
      remark: '',
    })

    const search = await listSupplierProducts({
      supplier_id: productForm.supplier_id,
      category_id: productForm.category_id,
      keyword: productForm.name.trim(),
      page: 1,
      page_size: 20,
    })
    const created =
      search.list.find((p) => p.name === productForm.name.trim()) ??
      search.list.find((p) => p.supplier_id === productForm.supplier_id) ??
      search.list[0]
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
    await batchUpsertProductUnitSpecs({
      product_id: created.id,
      units,
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
</script>
