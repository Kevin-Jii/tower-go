<template>
  <div class="return-photo-uploader">
    <div class="mb-2 flex items-center justify-between gap-3">
      <span class="text-sm font-medium text-slate-700">返厂照片</span>
      <span class="text-xs tabular-nums text-slate-500">{{ entries.length }}/{{ MAX_PHOTOS }}</span>
    </div>

    <div class="return-photo-uploader__grid">
      <div v-for="entry in entries" :key="entry.id" class="return-photo-uploader__item">
        <a-image :src="entry.previewURL" width="100%" height="100%" fit="cover" :preview="!entry.uploading" />
        <div v-if="entry.uploading" class="return-photo-uploader__progress">
          <a-progress :percent="entry.progress / 100" :show-text="false" size="small" />
          <span>{{ entry.progress }}%</span>
        </div>
        <button
          type="button"
          class="return-photo-uploader__remove"
          :disabled="disabled || entry.uploading"
          aria-label="删除返厂照片"
          title="删除照片"
          @click="removeEntry(entry.id)"
        >
          <IconDelete />
        </button>
      </div>

      <button
        v-if="entries.length < MAX_PHOTOS"
        type="button"
        class="return-photo-uploader__add"
        :disabled="disabled || uploading"
        aria-label="添加返厂照片"
        @click="openFilePicker"
      >
        <IconPlus />
        <span>添加照片</span>
      </button>
    </div>

    <input
      ref="fileInput"
      class="hidden"
      type="file"
      accept=".jpg,.jpeg,.png,.gif,.webp,image/jpeg,image/png,image/gif,image/webp"
      multiple
      @change="onPickFiles"
    />
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, ref, watch } from 'vue'
import { IconDelete, IconPlus } from '@arco-design/web-vue/es/icon'
import { toast } from '@/feedback/toast'
import { uploadGalleryResumable } from '@/utils/galleryUpload'

const MAX_PHOTOS = 3
const MAX_PHOTO_SIZE = 20 * 1024 * 1024
const ALLOWED_EXTENSIONS = new Set(['jpg', 'jpeg', 'png', 'gif', 'webp'])

type PhotoEntry = {
  id: string
  url: string
  previewURL: string
  objectURL?: string
  file?: File
  progress: number
  uploading: boolean
}

const props = withDefaults(defineProps<{
  modelValue?: string[]
  disabled?: boolean
}>(), {
  modelValue: () => [],
  disabled: false,
})

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
}>()

const fileInput = ref<HTMLInputElement | null>(null)
const entries = ref<PhotoEntry[]>([])
const uploading = ref(false)

watch(
  () => props.modelValue,
  (urls) => {
    const normalized = normalizeURLs(urls)
    const currentURLs = entries.value.filter((entry) => entry.url).map((entry) => entry.url)
    if (sameURLs(normalized, currentURLs)) return

    const pending = entries.value.filter((entry) => entry.file && !entry.url)
    const existingByURL = new Map(entries.value.filter((entry) => entry.url).map((entry) => [entry.url, entry]))
    entries.value = [
      ...normalized.map((url) => existingByURL.get(url) ?? createUploadedEntry(url)),
      ...pending,
    ].slice(0, MAX_PHOTOS)
  },
  { immediate: true, deep: true },
)

onBeforeUnmount(() => {
  for (const entry of entries.value) revokeObjectURL(entry)
})

function openFilePicker(): void {
  if (props.disabled || uploading.value) return
  fileInput.value?.click()
}

function onPickFiles(event: Event): void {
  const input = event.target as HTMLInputElement
  const available = MAX_PHOTOS - entries.value.length
  const selected = Array.from(input.files ?? [])
  input.value = ''
  if (available <= 0 || selected.length === 0) return
  if (selected.length > available) toast.warning(`最多上传${MAX_PHOTOS}张返厂照片`)

  for (const file of selected) {
    if (entries.value.length >= MAX_PHOTOS) break
    const extension = file.name.split('.').pop()?.toLowerCase() ?? ''
    if (!ALLOWED_EXTENSIONS.has(extension)) {
      toast.warning(`「${file.name}」不是支持的图片格式`)
      continue
    }
    if (file.size > MAX_PHOTO_SIZE) {
      toast.warning(`「${file.name}」不能超过20MB`)
      continue
    }
    const duplicate = entries.value.some((entry) =>
      entry.file?.name === file.name && entry.file.size === file.size && entry.file.lastModified === file.lastModified,
    )
    if (duplicate) continue

    const objectURL = URL.createObjectURL(file)
    entries.value.push({
      id: createEntryID(),
      url: '',
      previewURL: objectURL,
      objectURL,
      file,
      progress: 0,
      uploading: false,
    })
  }
}

function removeEntry(id: string): void {
  const index = entries.value.findIndex((entry) => entry.id === id)
  if (index < 0) return
  revokeObjectURL(entries.value[index])
  entries.value.splice(index, 1)
  emitUploadedURLs()
}

async function uploadPending(): Promise<string[]> {
  uploading.value = true
  try {
    for (const entry of entries.value) {
      if (!entry.file || entry.url) continue
      entry.uploading = true
      entry.progress = 0
      try {
        const gallery = await uploadGalleryResumable(entry.file, {
          category: 'store-return',
          remark: '返厂记录照片',
          onProgress: (progress) => {
            entry.progress = progress.percent
          },
        })
        entry.url = gallery.url
        entry.previewURL = gallery.url
        entry.file = undefined
        entry.progress = 100
        revokeObjectURL(entry)
        emitUploadedURLs()
      } finally {
        entry.uploading = false
      }
    }
    return entries.value.map((entry) => entry.url).filter(Boolean)
  } finally {
    uploading.value = false
  }
}

function emitUploadedURLs(): void {
  emit('update:modelValue', entries.value.map((entry) => entry.url).filter(Boolean))
}

function normalizeURLs(urls: string[] | undefined): string[] {
  return Array.from(new Set((urls ?? []).map((url) => url.trim()).filter(Boolean))).slice(0, MAX_PHOTOS)
}

function sameURLs(left: string[], right: string[]): boolean {
  return left.length === right.length && left.every((url, index) => url === right[index])
}

function createUploadedEntry(url: string): PhotoEntry {
  return { id: createEntryID(), url, previewURL: url, progress: 100, uploading: false }
}

function createEntryID(): string {
  return globalThis.crypto?.randomUUID?.() ?? `photo_${Date.now()}_${Math.random().toString(16).slice(2)}`
}

function revokeObjectURL(entry: PhotoEntry): void {
  if (!entry.objectURL) return
  URL.revokeObjectURL(entry.objectURL)
  entry.objectURL = undefined
}

defineExpose({ uploadPending })
</script>

<style scoped>
.return-photo-uploader__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 112px));
  gap: 10px;
}

.return-photo-uploader__item,
.return-photo-uploader__add {
  position: relative;
  width: 112px;
  aspect-ratio: 1;
  overflow: hidden;
  border: 1px solid var(--color-border-2);
  border-radius: 6px;
  background: #f8fafc;
}

.return-photo-uploader__add {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border-style: dashed;
  color: #64748b;
  font-size: 12px;
  cursor: pointer;
}

.return-photo-uploader__add:hover:not(:disabled) {
  border-color: #6366f1;
  color: #4f46e5;
  background: #f8faff;
}

.return-photo-uploader__add svg {
  font-size: 22px;
}

.return-photo-uploader__add:disabled,
.return-photo-uploader__remove:disabled {
  cursor: not-allowed;
  opacity: 0.55;
}

.return-photo-uploader__remove {
  position: absolute;
  top: 6px;
  right: 6px;
  display: inline-flex;
  width: 28px;
  height: 28px;
  align-items: center;
  justify-content: center;
  border: 0;
  border-radius: 4px;
  background: rgb(15 23 42 / 72%);
  color: white;
  cursor: pointer;
}

.return-photo-uploader__progress {
  position: absolute;
  inset: auto 0 0;
  display: grid;
  grid-template-columns: 1fr auto;
  align-items: center;
  gap: 6px;
  padding: 8px;
  background: rgb(15 23 42 / 78%);
  color: white;
  font-size: 11px;
}

@media (max-width: 480px) {
  .return-photo-uploader__grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .return-photo-uploader__item,
  .return-photo-uploader__add {
    width: 100%;
  }
}
</style>
