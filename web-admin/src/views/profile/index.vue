<template>
  <div class="mx-auto w-full max-w-6xl flex flex-col md:flex-row gap-6 items-start p-4">

    <BaseCard class="flex-1 min-w-[320px] bg-white shadow-sm rounded-xl">
      <template #header>
        <div class="py-1 border-b border-slate-100">
          <span class="font-semibold text-lg text-slate-800">个人资料</span>
        </div>
      </template>

      <div v-if="loading" class="py-12 flex justify-center">
        <MathCurveLoader size="md" text="加载中…" />
      </div>

      <div v-else class="space-y-5 pt-4">
        <div class="flex items-center gap-4 pb-4 border-b rgb(98 150 201)">
          <a-avatar :size="64" class="!bg-[rgb(var(--primary-6))] shrink-0 shadow-sm font-medium text-lg">
            {{ avatarText }}
          </a-avatar>
          <div class="min-w-0">
            <h3 class="m-0 text-xl font-bold text-slate-800 truncate">
              {{ profile?.nickname || profile?.username }}
            </h3>
            <span
              class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-slate-100 text-slate-600 mt-1.5">
              {{ profile?.role?.name ?? '—' }}
            </span>
          </div>
        </div>

        <div class="space-y-4 pt-2">
          <BaseFormItem label="用户名" class="w-full">
            <BaseInput v-model="form.username" placeholder="登录用户名" />
          </BaseFormItem>

          <BaseFormItem label="昵称" class="w-full">
            <BaseInput v-model="form.nickname" placeholder="显示名称" />
          </BaseFormItem>

          <BaseFormItem label="手机号" class="w-full">
            <BaseInput v-model="form.phone" maxlength="11" placeholder="11 位手机号" />
          </BaseFormItem>

          <BaseFormItem label="邮箱" class="w-full">
            <BaseInput v-model="form.email" type="email" placeholder="选填" />
          </BaseFormItem>

          <BaseFormItem label="性别" class="w-full">
            <BaseSelect v-model="form.gender" :options="[
              { label: '男', value: 1 },
              { label: '女', value: 2 },
            ]" :allow-clear="false" class="w-full" />
          </BaseFormItem>

          <BaseFormItem label="所属门店" class="w-full">
            <BaseInput :model-value="storeLabel" disabled class="bg-slate-50 cursor-not-allowed text-slate-400" />
          </BaseFormItem>
        </div>
      </div>
    </BaseCard>

    <BaseCard class="w-full md:w-[360px] bg-white shadow-sm rounded-xl shrink-0">
      <template #header>
        <div class="py-1 border-b border-slate-100">
          <span class="font-semibold text-lg text-slate-800">修改密码</span>
        </div>
      </template>

      <div class="space-y-5 pt-4">
        <BaseFormItem label="新密码">
          <BaseInput v-model="form.password" type="password" show-password placeholder="至少 6 位，不修改请留空" />
        </BaseFormItem>

        <BaseFormItem label="确认密码">
          <BaseInput v-model="form.passwordConfirm" type="password" show-password placeholder="再次输入新密码" />
        </BaseFormItem>

        <div class="flex justify-end gap-3 pt-4 border-t border-slate-100">
          <BaseButton variant="ghost" class="px-4" @click="loadProfile">重置</BaseButton>
          <BaseButton variant="primary" class="px-5" :loading="saving" @click="save">保存</BaseButton>
        </div>
      </div>
    </BaseCard>

  </div>
</template>
<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { BaseButton, BaseCard, BaseFormItem, BaseInput, BaseSelect } from '@/components/base'
import { MathCurveLoader } from '@/components/loading'
import { getProfile, updateProfile } from '@/api/user'
import type { User } from '@/api/types'
import { useUserStore } from '@/store/user'
import { toast } from '@/feedback/toast'

const userStore = useUserStore()

const loading = ref(true)
const saving = ref(false)
const profile = ref<User | null>(null)

const form = reactive({
  username: '',
  nickname: '',
  phone: '',
  email: '',
  gender: 1 as number,
  password: '',
  passwordConfirm: '',
})

const avatarText = computed(() => {
  const name = profile.value?.nickname || profile.value?.username || '?'
  return name.slice(0, 1).toUpperCase()
})

const storeLabel = computed(() => {
  if (!profile.value?.store_id) return '未绑定门店'
  return profile.value.store?.name ?? `门店 #${profile.value.store_id}`
})

function fillForm(u: User): void {
  form.username = u.username ?? ''
  form.nickname = u.nickname ?? ''
  form.phone = u.phone ?? ''
  form.email = u.email ?? ''
  form.gender = u.gender ?? 1
  form.password = ''
  form.passwordConfirm = ''
}

async function loadProfile(): Promise<void> {
  loading.value = true
  try {
    const u = await getProfile()
    profile.value = u
    fillForm(u)
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载失败')
  } finally {
    loading.value = false
  }
}

async function save(): Promise<void> {
  if (form.password || form.passwordConfirm) {
    if (form.password.length < 6) {
      toast.warning('新密码至少 6 位')
      return
    }
    if (form.password !== form.passwordConfirm) {
      toast.warning('两次输入的密码不一致')
      return
    }
  }
  if (form.phone && form.phone.length !== 11) {
    toast.warning('手机号须为 11 位')
    return
  }

  saving.value = true
  try {
    const body: Record<string, unknown> = {
      username: form.username.trim(),
      nickname: form.nickname.trim(),
      phone: form.phone.trim(),
      email: form.email.trim(),
      gender: form.gender,
    }
    if (form.password) body.password = form.password
    await updateProfile(body)
    toast.success('保存成功')
    await loadProfile()
    if (profile.value) {
      userStore.patchUserInfo(profile.value)
    }
    form.password = ''
    form.passwordConfirm = ''
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  void loadProfile()
})
</script>
