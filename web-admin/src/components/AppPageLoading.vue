<template>
  <div class="app-page-loading" role="status" aria-live="polite">
    <span class="app-page-loading__ring" aria-hidden="true" />
    <p class="app-page-loading__text">{{ text }}</p>
  </div>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    text?: string
  }>(),
  { text: '加载中…' },
)
</script>

<style scoped>
.app-page-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
  min-height: 200px;
  padding: 32px 16px;
  color: var(--color-text-2, #64748b);
}
.app-page-loading__text {
  margin: 0;
  font-size: 14px;
}
.app-page-loading__ring {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(99, 102, 241, 0.2);
  border-top-color: rgb(99, 102, 241);
  border-radius: 50%;
  animation: app-page-loading-spin 0.75s linear infinite;
}
@keyframes app-page-loading-spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
