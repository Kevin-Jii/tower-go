<template>
  <main class="flex min-h-0 flex-1 flex-col">
    <!-- 灰底 layout-content 之上：统一白底内容区，与门店记账等页层次一致 -->
    <div class="layout-page-surface">
      <router-view v-slot="{ Component, route }">
        <transition name="fade-transform" mode="out-in">
          <!-- 用包裹层接 flex，避免与业务页根节点上的 display:grid 等布局类冲突（如供应商左右分栏） -->
          <div :key="route.fullPath" class="flex min-h-0 min-w-0 flex-1 flex-col">
            <component :is="Component" />
          </div>
        </transition>
      </router-view>
    </div>
  </main>
</template>

<script setup lang="ts"></script>
