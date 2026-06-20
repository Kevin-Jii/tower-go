<template>
  <div class="dashboard-loading-overlay" :class="{ 'dashboard-loading-overlay--fixed': fixed }" role="status" aria-live="polite" aria-busy="true">
    <div class="dashboard-loading-core">
      <span class="dashboard-loading-orbit" />
      <span class="dashboard-loading-orbit dashboard-loading-orbit--second" />
      <span class="dashboard-loading-dot" />
      <strong>{{ title }}</strong>
      <small>{{ subtitle }}</small>
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    title?: string
    subtitle?: string
    fixed?: boolean
  }>(),
  {
    title: '数据加载中...',
    subtitle: '正在准备页面数据',
    fixed: false,
  },
)
</script>

<style scoped>
.dashboard-loading-overlay {
  position: absolute;
  z-index: 30;
  inset: 0;
  display: grid;
  place-items: center;
  min-height: 100%;
  background:
    linear-gradient(180deg, rgba(7, 11, 20, 0.82), rgba(7, 11, 20, 0.94)),
    radial-gradient(circle at 50% 42%, rgba(34, 211, 238, 0.16), transparent 34%);
  backdrop-filter: blur(6px);
}

.dashboard-loading-overlay--fixed {
  position: fixed;
  z-index: 10000;
  min-height: 100vh;
}

.dashboard-loading-core {
  position: relative;
  display: grid;
  place-items: center;
  width: clamp(180px, 18vw, 260px);
  aspect-ratio: 1;
  color: #ecfeff;
}

.dashboard-loading-core::before {
  content: "";
  position: absolute;
  inset: 18%;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.72);
  border: 1px solid rgba(34, 211, 238, 0.26);
  box-shadow:
    inset 0 0 28px rgba(34, 211, 238, 0.08),
    0 0 38px rgba(34, 211, 238, 0.14);
}

.dashboard-loading-core strong,
.dashboard-loading-core small {
  position: absolute;
  left: 50%;
  z-index: 2;
  width: max-content;
  max-width: min(72vw, 360px);
  text-align: center;
  transform: translateX(-50%);
}

.dashboard-loading-core strong {
  top: 35%;
  font-size: clamp(18px, 1.45vw, 26px);
  font-weight: 900;
}

.dashboard-loading-core small {
  top: 63%;
  color: rgba(207, 250, 254, 0.68);
  font-size: 13px;
  font-weight: 700;
}

.dashboard-loading-orbit {
  position: absolute;
  inset: 0;
  border-radius: 999px;
  border: 1px solid rgba(34, 211, 238, 0.18);
  border-top-color: rgba(34, 211, 238, 0.95);
  border-right-color: rgba(167, 139, 250, 0.78);
  animation: dashboard-loading-spin 1.35s linear infinite;
}

.dashboard-loading-orbit--second {
  inset: 14%;
  animation-duration: 1.9s;
  animation-direction: reverse;
  border-top-color: rgba(52, 211, 153, 0.9);
  border-right-color: rgba(34, 211, 238, 0.46);
}

.dashboard-loading-dot {
  position: absolute;
  z-index: 2;
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: #22d3ee;
  box-shadow: 0 0 16px #22d3ee;
  animation: dashboard-loading-pulse 1s ease-in-out infinite alternate;
}

@keyframes dashboard-loading-spin {
  to {
    transform: rotate(360deg);
  }
}

@keyframes dashboard-loading-pulse {
  from {
    opacity: 0.48;
    transform: scale(0.72);
  }

  to {
    opacity: 1;
    transform: scale(1.08);
  }
}
</style>
