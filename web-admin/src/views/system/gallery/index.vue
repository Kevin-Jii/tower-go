<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row gap-3 md:items-end justify-between">
      <h2 class="page-title">图库管理</h2>
      <div class="flex flex-col sm:flex-row gap-2 w-full md:w-auto">
        <BaseInput v-model="keyword" class="w-full sm:w-48" placeholder="文件名/备注" clearable @enter="reload" />
        <BaseSelect v-model="category" class="w-full sm:w-36" :options="categoryOptions" placeholder="分类" />
        <div class="flex gap-2">
          <BaseButton variant="primary" @click="reload">查询</BaseButton>
          <BaseButton v-permission="'system:gallery:upload'" variant="primary" @click="openUpload">上传图片</BaseButton>
          <BaseButton v-permission="'system:gallery:delete'" variant="danger" :disabled="selectedIds.length === 0"
            @click="onBatchDelete">
            批量删除
          </BaseButton>
        </div>
      </div>
    </div>

    <BaseTable :columns="columns" :data="(list as unknown) as Record<string, unknown>[]" :loading="loading"
      min-width="1180px">
      <template #cell-select="{ row }">
        <a-checkbox :model-value="selectedIds.includes((row as Gallery).id)"
          @change="toggleSelected((row as Gallery).id)" />
      </template>
      <template #cell-preview="{ row }">
        <a-image :src="(row as Gallery).url" width="56" height="56" fit="cover" :preview="true" />
      </template>
      <template #cell-size="{ row }">
        {{ formatSize((row as Gallery).size) }}
      </template>
      <template #cell-created_at="{ row }">
        {{ formatDateTime((row as Gallery).created_at) }}
      </template>
      <template #cell-actions="{ row }">
        <BaseTableRowActions :actions="galleryRowActions(row as Gallery)" />
      </template>
    </BaseTable>

    <BaseDialog v-model="uploadDlg" title="上传图片" max-width="min(520px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="图片文件" required>
          <input type="file" accept="image/*" @change="onPickFile" />
          <p v-if="uploadForm.fileName" class="m-0 mt-2 text-xs text-slate-500">{{ uploadForm.fileName }}</p>
        </BaseFormItem>
        <BaseFormItem label="分类">
          <BaseSelect v-model="uploadForm.category" :options="categoryOptions.filter((o) => o.value !== '')" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="uploadForm.remark" :rows="2" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="uploadDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="uploading" @click="submitUpload">上传</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="editDlg" title="编辑图片信息" max-width="min(520px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="名称">
          <BaseInput v-model="editForm.name" />
        </BaseFormItem>
        <BaseFormItem label="分类">
          <BaseSelect v-model="editForm.category" :options="categoryOptions.filter((o) => o.value !== '')" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="editForm.remark" :rows="2" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="editDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submitEdit">保存</BaseButton>
      </template>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import {
  BaseButton,
  BaseDialog,
  BaseFormItem,
  BaseInput,
  BaseSelect,
  BaseTable,
  BaseTableRowActions,
  BaseTextarea,
} from '@/components/base'
import type { BaseSelectOption, BaseTableColumn, TableRowAction } from '@/components/base/types'
import { batchDeleteGallery, deleteGallery, listGalleries, updateGallery, uploadGallery } from '@/api/gallery'
import type { Gallery } from '@/api/types'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'

const columns: BaseTableColumn[] = [
  { key: 'select', label: '', width: '56px' },
  { key: 'preview', label: '预览', width: '88px' },
  { key: 'name', label: '文件名', prop: 'name', minWidth: '220px', ellipsis: true },
  { key: 'category', label: '分类', prop: 'category', width: '120px' },
  { key: 'size', label: '大小', width: '96px' },
  { key: 'created_at', label: '上传时间', width: '176px' },
  { key: 'remark', label: '备注', prop: 'remark', minWidth: '160px', ellipsis: true },
  { key: 'actions', label: '操作', width: '140px', align: 'right' },
]

const categoryOptions: BaseSelectOption[] = [
  { label: '全部分类', value: '' },
  { label: '商品', value: 'product' },
  { label: '供应商', value: 'supplier' },
  { label: '头像', value: 'avatar' },
  { label: '采购', value: 'purchase' },
  { label: '其他', value: 'other' },
]

const loading = ref(false)
const list = ref<Gallery[]>([])
const keyword = ref('')
const category = ref<string | number>('')
const selectedIds = ref<number[]>([])

async function load(): Promise<void> {
  loading.value = true
  try {
    const rows = await listGalleries({
      page: 1,
      page_size: 200,
      keyword: keyword.value.trim() || undefined,
      category: String(category.value || '').trim() || undefined,
    })
    list.value = rows ?? []
    selectedIds.value = selectedIds.value.filter((id) => list.value.some((r) => r.id === id))
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载图库失败')
  } finally {
    loading.value = false
  }
}

function reload(): void {
  void load()
}

function toggleSelected(id: number): void {
  if (selectedIds.value.includes(id)) selectedIds.value = selectedIds.value.filter((x) => x !== id)
  else selectedIds.value.push(id)
}

function formatDateTime(v?: string): string {
  if (!v) return '-'
  return v.slice(0, 19).replace('T', ' ')
}

function formatSize(size = 0): string {
  const kb = 1024
  const mb = kb * 1024
  if (size >= mb) return `${(size / mb).toFixed(2)} MB`
  if (size >= kb) return `${(size / kb).toFixed(2)} KB`
  return `${size} B`
}

const uploadDlg = ref(false)
const uploading = ref(false)
const uploadForm = reactive<{ file?: File; fileName: string; category: string | number; remark: string }>({
  file: undefined,
  fileName: '',
  category: 'other',
  remark: '',
})

function openUpload(): void {
  uploadForm.file = undefined
  uploadForm.fileName = ''
  uploadForm.category = 'other'
  uploadForm.remark = ''
  uploadDlg.value = true
}

function onPickFile(e: Event): void {
  const t = e.target as HTMLInputElement
  const f = t.files?.[0]
  uploadForm.file = f
  uploadForm.fileName = f?.name ?? ''
}

async function submitUpload(): Promise<void> {
  if (!uploadForm.file) {
    toast.warning('请选择图片')
    return
  }
  const fd = new FormData()
  fd.append('file', uploadForm.file)
  fd.append('category', String(uploadForm.category || 'other'))
  fd.append('remark', uploadForm.remark.trim())
  uploading.value = true
  try {
    await uploadGallery(fd)
    toast.success('上传成功')
    uploadDlg.value = false
    await load()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '上传失败')
  } finally {
    uploading.value = false
  }
}

const editDlg = ref(false)
const saving = ref(false)
const editId = ref<number | null>(null)
const editForm = reactive<{ name: string; category: string | number; remark: string }>({
  name: '',
  category: 'other',
  remark: '',
})

function openEdit(row: Gallery): void {
  editId.value = row.id
  editForm.name = row.name || ''
  editForm.category = row.category || 'other'
  editForm.remark = row.remark || ''
  editDlg.value = true
}

async function submitEdit(): Promise<void> {
  if (!editId.value) return
  saving.value = true
  try {
    await updateGallery(editId.value, {
      name: editForm.name.trim(),
      category: String(editForm.category || 'other'),
      remark: editForm.remark.trim(),
    })
    toast.success('保存成功')
    editDlg.value = false
    await load()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function onDelete(row: Gallery): Promise<void> {
  const ok = await confirmDialog({ message: `删除图片「${row.name}」？` })
  if (!ok) return
  try {
    await deleteGallery(row.id)
    toast.success('已删除')
    await load()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}

function galleryRowActions(row: Gallery): TableRowAction[] {
  return [
    { label: '编辑', permission: 'system:gallery:edit', onClick: () => openEdit(row) },
    { label: '删除', permission: 'system:gallery:delete', danger: true, onClick: () => void onDelete(row) },
  ]
}

async function onBatchDelete(): Promise<void> {
  if (!selectedIds.value.length) return
  const ok = await confirmDialog({ message: `批量删除 ${selectedIds.value.length} 张图片？` })
  if (!ok) return
  try {
    await batchDeleteGallery(selectedIds.value)
    toast.success('批量删除成功')
    selectedIds.value = []
    await load()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '批量删除失败')
  }
}

void load()
</script>
