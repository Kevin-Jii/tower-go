<template>
  <main class="login-screen">
    <section class="login-visual" aria-hidden="true">
      <div class="login-visual__grid"></div>
      <div class="login-radar">
        <span v-for="i in 5" :key="i" :style="{ '--i': i }"></span>
      </div>
      <div class="login-orbit login-orbit--a">
        <span>库存</span>
      </div>
      <div class="login-orbit login-orbit--b">
        <span>记账</span>
      </div>
      <div class="login-orbit login-orbit--c">
        <span>会员</span>
      </div>

      <div class="login-brand-panel">
        <p class="login-kicker">Tower Store OS</p>
        <h1>门店经营后台</h1>
        <p class="login-brand-copy">把记账、库存、会员和经营数据放在同一个工作台里，开店时少一点来回切换。</p>
        <div class="login-metrics">
          <div>
            <strong>实时</strong>
            <span>经营数据</span>
          </div>
          <div>
            <strong>自动</strong>
            <span>库存扣减</span>
          </div>
          <div>
            <strong>清晰</strong>
            <span>门店协作</span>
          </div>
        </div>
      </div>
    </section>

    <section class="login-form-wrap">
      <div class="login-card">
        <div class="login-card__head">
          <p class="login-kicker">Welcome Back</p>
          <h2>登录 Tower</h2>
          <p>使用手机号和密码进入管理后台</p>
        </div>

        <BaseForm class="login-form" @submit="onSubmit">
          <BaseFormItem label="手机号" required>
            <div class="login-field">
              <BaseInput v-model="form.phone" maxlength="11" placeholder="请输入手机号" autocomplete="username" />
            </div>
          </BaseFormItem>

          <BaseFormItem label="密码" required>
            <div class="login-field">
              <BaseInput v-model="form.password" type="password" show-password placeholder="请输入密码"
                autocomplete="current-password" />
            </div>
          </BaseFormItem>

          <div class="login-options">
            <label>
              <input v-model="rememberMe" type="checkbox" />
              <span>记住账号</span>
            </label>
            <span>请勿在公共电脑保存密码</span>
          </div>

          <button class="login-submit" type="submit" :disabled="loading">
            <span v-if="loading" class="login-submit__loading"></span>
            <span>{{ loading ? '正在登录' : '进入后台' }}</span>
          </button>
        </BaseForm>
      </div>
    </section>
  </main>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { BaseForm, BaseFormItem, BaseInput } from '@/components/base'
import { login } from '@/api/auth'
import { useUserStore } from '@/store/user'
import { toast } from '@/feedback/toast'
import { clearRememberedLogin, getRememberedLogin, setRememberedLogin } from '@/utils/storage'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const form = reactive({ phone: '', password: '' })
const loading = ref(false)
const rememberMe = ref(false)

onMounted(() => {
  const saved = getRememberedLogin()
  if (saved) {
    form.phone = saved.phone
    form.password = saved.password
    rememberMe.value = true
  }
})

async function onSubmit(): Promise<void> {
  if (!form.phone.trim() || !form.password.trim()) {
    toast.warning('请填写手机号和密码')
    return
  }
  loading.value = true
  try {
    const data = await login(form.phone, form.password)
    if (rememberMe.value) {
      setRememberedLogin(form.phone, form.password)
    } else {
      clearRememberedLogin()
    }
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

<style scoped>
.login-screen {
  position: relative;
  flex: 1;
  min-height: 0;
  display: grid;
  grid-template-columns: minmax(0, 1.08fr) minmax(420px, 0.92fr);
  overflow: hidden;
  background: #f7f8fb;
  color: #111827;
}

.login-visual {
  position: relative;
  min-height: 100%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 56px;
  background:
    radial-gradient(circle at 25% 20%, rgba(34, 197, 94, 0.18), transparent 30%),
    radial-gradient(circle at 75% 24%, rgba(59, 130, 246, 0.2), transparent 28%),
    linear-gradient(135deg, #101827 0%, #172033 46%, #0f172a 100%);
}

.login-visual__grid {
  position: absolute;
  inset: 0;
  opacity: 0.28;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.08) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.08) 1px, transparent 1px);
  background-size: 42px 42px;
  mask-image: linear-gradient(to bottom, black 30%, transparent 90%);
}

.login-radar {
  position: absolute;
  width: min(54vw, 660px);
  aspect-ratio: 1;
  display: grid;
  place-items: center;
}

.login-radar span {
  --size: calc(150px + var(--i) * 92px);
  position: absolute;
  width: var(--size);
  height: var(--size);
  border-radius: 999px;
  border: 1px solid rgba(148, 163, 184, 0.2);
  animation: login-ripple 4.8s ease-out infinite;
  animation-delay: calc(var(--i) * -0.45s);
}

.login-orbit {
  position: absolute;
  width: 360px;
  height: 360px;
  border: 1px solid rgba(255, 255, 255, 0.14);
  border-radius: 999px;
  animation: login-orbit 22s linear infinite;
}

.login-orbit span {
  position: absolute;
  top: -18px;
  left: 50%;
  transform: translateX(-50%);
  padding: 8px 12px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.92);
  color: #0f172a;
  font-size: 12px;
  font-weight: 700;
  box-shadow: 0 16px 36px rgba(15, 23, 42, 0.28);
}

.login-orbit--b {
  width: 500px;
  height: 500px;
  animation-duration: 30s;
  animation-direction: reverse;
}

.login-orbit--c {
  width: 260px;
  height: 260px;
  animation-duration: 16s;
}

.login-brand-panel {
  position: relative;
  z-index: 2;
  max-width: 520px;
  color: #fff;
  animation: login-rise 0.72s ease both;
}

.login-kicker {
  margin: 0 0 12px;
  color: #38bdf8;
  font-size: 12px;
  font-weight: 800;
  letter-spacing: 0.18em;
  text-transform: uppercase;
}

.login-brand-panel h1 {
  margin: 0;
  font-size: clamp(40px, 5.4vw, 72px);
  line-height: 1;
  font-weight: 850;
  letter-spacing: 0;
}

.login-brand-copy {
  max-width: 440px;
  margin: 22px 0 0;
  color: rgba(226, 232, 240, 0.82);
  font-size: 16px;
  line-height: 1.8;
}

.login-metrics {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin-top: 36px;
}

.login-metrics div {
  padding: 16px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.14);
  backdrop-filter: blur(16px);
}

.login-metrics strong,
.login-metrics span {
  display: block;
}

.login-metrics strong {
  font-size: 20px;
}

.login-metrics span {
  margin-top: 4px;
  color: rgba(226, 232, 240, 0.72);
  font-size: 12px;
}

.login-form-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 36px;
}

.login-card {
  width: min(100%, 420px);
  padding: 34px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.94);
  border: 1px solid rgba(226, 232, 240, 0.95);
  box-shadow: 0 28px 80px rgba(15, 23, 42, 0.12);
  animation: login-rise 0.58s ease 0.08s both;
}

.login-card__head h2 {
  margin: 0;
  font-size: 30px;
  line-height: 1.1;
  font-weight: 820;
  letter-spacing: 0;
}

.login-card__head p:last-child {
  margin: 10px 0 0;
  color: #64748b;
  font-size: 14px;
}

.login-form {
  margin-top: 28px;
}

.login-field {
  position: relative;
  padding: 2px;
  border-radius: 8px;
  width: 100%;
  background: linear-gradient(120deg, rgba(59, 130, 246, 0), rgba(34, 197, 94, 0));
  transition: background 0.25s ease, box-shadow 0.25s ease;
}

.login-field:hover,
.login-field:focus-within {
  background: linear-gradient(120deg, rgba(59, 130, 246, 0.85), rgba(34, 197, 94, 0.75));
  box-shadow: 0 12px 26px rgba(59, 130, 246, 0.12);
}

.login-field :deep(.arco-input-wrapper) {
  height: 44px;
  border: 0;
  border-radius: 6px;
  background: #f8fafc;
  box-shadow: inset 0 0 0 1px #e2e8f0;
}

.login-field :deep(.arco-input-wrapper.arco-input-focus) {
  box-shadow: inset 0 0 0 1px transparent;
}

.login-options {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin: 4px 0 18px;
  color: #64748b;
  font-size: 12px;
}

.login-options label {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: #334155;
  cursor: pointer;
  user-select: none;
}

.login-options input {
  width: 14px;
  height: 14px;
  accent-color: #2563eb;
}

.login-submit {
  position: relative;
  width: 100%;
  height: 46px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border: 0;
  border-radius: 8px;
  color: #fff;
  font-size: 15px;
  font-weight: 800;
  background: linear-gradient(135deg, #111827, #2563eb 58%, #059669);
  box-shadow: 0 18px 34px rgba(37, 99, 235, 0.28);
  cursor: pointer;
  overflow: hidden;
}

.login-submit::after {
  content: '';
  position: absolute;
  inset: auto 12% -1px;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.85), transparent);
  opacity: 0;
  transition: opacity 0.25s ease;
}

.login-submit:hover::after {
  opacity: 1;
}

.login-submit:disabled {
  cursor: not-allowed;
  opacity: 0.78;
}

.login-submit__loading {
  width: 14px;
  height: 14px;
  border-radius: 999px;
  border: 2px solid rgba(255, 255, 255, 0.45);
  border-top-color: #fff;
  animation: login-spin 0.8s linear infinite;
}

@keyframes login-ripple {
  0% {
    opacity: 0.3;
    transform: scale(0.86);
  }

  75% {
    opacity: 0.08;
  }

  100% {
    opacity: 0;
    transform: scale(1.08);
  }
}

@keyframes login-orbit {
  to {
    transform: rotate(360deg);
  }
}

@keyframes login-rise {
  from {
    opacity: 0;
    transform: translateY(18px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes login-spin {
  to {
    transform: rotate(360deg);
  }
}

@media (max-width: 960px) {
  .login-screen {
    grid-template-columns: 1fr;
    overflow-y: auto;
  }

  .login-visual {
    min-height: 250px;
    padding: 24px 22px;
    align-items: flex-end;
    justify-content: flex-start;
  }

  .login-brand-panel h1 {
    font-size: 38px;
  }

  .login-brand-copy,
  .login-metrics {
    display: none;
  }

  .login-form-wrap {
    padding: 18px 22px 22px;
    align-items: flex-start;
  }
}

@media (max-width: 520px) {
  .login-card {
    padding: 24px;
  }

  .login-options {
    align-items: flex-start;
    flex-direction: column;
    gap: 6px;
  }
}
</style>
