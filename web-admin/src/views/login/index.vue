<template>
  <div
    class="flex-1 min-h-0 overflow-y-auto flex flex-col items-center justify-center p-4 bg-gradient-to-br from-slate-100 via-indigo-50/40 to-violet-100/50"
  >
    <div class="w-full max-w-md rounded-2xl bg-white/95 backdrop-blur border border-white/80 p-6 md:p-8 shadow-2xl shadow-indigo-900/10">
      <h1 class="text-xl md:text-2xl font-bold text-center text-slate-900 mb-1 tracking-tight">Tower 管理后台</h1>
      <p class="text-center text-sm text-slate-500 mb-6">使用 UnoCSS 与自研 Base 组件</p>
      <BaseForm @submit="onSubmit">
        <BaseFormItem label="手机号" required>
          <BaseInput v-model="form.phone" maxlength="11" placeholder="11 位手机号" autocomplete="username" />
        </BaseFormItem>
        <BaseFormItem label="密码" required>
          <BaseInput
            v-model="form.password"
            type="password"
            show-password
            placeholder="至少 6 位"
            autocomplete="current-password"
          />
        </BaseFormItem>
        <BaseButton variant="primary" class="w-full mt-2" native-type="submit" size="lg" :loading="loading">
          登录
        </BaseButton>
      </BaseForm>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { BaseButton, BaseForm, BaseFormItem, BaseInput } from '@/components/base'
import { login } from '@/api/auth'
import { useUserStore } from '@/store/user'
import { toast } from '@/feedback/toast'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const form = reactive({ phone: '', password: '' })
const loading = ref(false)

async function onSubmit(): Promise<void> {
  loading.value = true
  try {
    const data = await login(form.phone, form.password)
    userStore.setLogin({ token: data.token, user: data.user_info })
    userStore.markDynamicRoutes(false)
    const redirect = (route.query.redirect as string) || '/'
    await router.replace(redirect)
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '登录失败')
  } finally {
    loading.value = false
  }
}
</script>
