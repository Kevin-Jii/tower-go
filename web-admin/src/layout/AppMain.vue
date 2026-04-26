<template>
  <main class="min-h-0">
    <router-view v-slot="{ Component, route }">
      <transition name="fade-transform" mode="out-in">
        <component :is="Component" :key="route.fullPath" />
      </transition>
    </router-view>
  </main>
</template>

<script setup lang="ts"></script>
