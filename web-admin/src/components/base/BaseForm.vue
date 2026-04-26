<template>
  <a-form
    :model="formModel"
    layout="vertical"
    class="base-form flex flex-col gap-4"
    @submit-success="emit('submit')"
  >
    <slot />
  </a-form>
</template>

<script setup lang="ts">
import { reactive } from 'vue'

/** 业务表单项自行绑定 v-model，此处仅占位以满足 Arco Form 的 model 要求 */
const formModel = reactive<Record<string, unknown>>({})

const emit = defineEmits<{ submit: [] }>()
</script>
