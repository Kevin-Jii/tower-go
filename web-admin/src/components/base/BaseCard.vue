<template>
  <a-card class="card-shell min-w-0" :bordered="true">
    <template v-if="$slots.header" #title>
      <slot name="header" />
    </template>
    <template v-else-if="title" #title>{{ title }}</template>
    <slot />
  </a-card>
</template>

<script setup lang="ts">
withDefaults(defineProps<{ title?: string }>(), { title: '' })
</script>
