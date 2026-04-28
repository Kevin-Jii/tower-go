<template>
  <div v-if="options.length > 1" class="flex items-center gap-2 max-w-[10rem] md:max-w-xs">
    <span class="text-xs text-[var(--color-text-3)] hidden md:inline whitespace-nowrap">租户</span>
    <a-select
      :model-value="userStore.tenantId"
      :options="selectOptions"
      size="small"
      class="!min-w-[7.5rem]"
      placeholder="门店"
      @change="onChange"
    />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { listAllStores } from '@/api/store'
import { useUserStore } from '@/store/user'
import { applyTenantSwitch } from '@/utils/storage'

const userStore = useUserStore()

const isAdmin = computed(() => {
  const code = userStore.userInfo?.role?.code ?? ''
  return code === 'admin' || code === 'super_admin'
})

const { data: remoteStores } = useQuery({
  queryKey: ['stores', 'all'],
  queryFn: listAllStores,
  enabled: isAdmin,
})

const options = computed(() => {
  if (isAdmin.value && remoteStores.value?.length) return remoteStores.value
  const st = userStore.userInfo?.store
  if (st?.id) return [{ id: st.id, name: st.name || `门店 #${st.id}` }]
  const sid = userStore.userInfo?.store_id
  if (sid) return [{ id: sid, name: `门店 #${sid}` }]
  return []
})

const selectOptions = computed(() =>
  options.value.map((s) => ({
    label: s.name || `门店 #${s.id}`,
    value: s.id,
  })),
)

function onChange(v: string | number | boolean | Record<string, unknown> | unknown[]): void {
  if (v === undefined || v === null || typeof v === 'boolean' || Array.isArray(v) || typeof v === 'object')
    return
  const id = typeof v === 'number' ? v : Number(v)
  if (!Number.isFinite(id)) return
  userStore.setTenantId(id)
  applyTenantSwitch(id, 'reload')
}
</script>
