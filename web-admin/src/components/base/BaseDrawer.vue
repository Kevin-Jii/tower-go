<template>
  <a-drawer
    :visible="modelValue"
    placement="left"
    :header="false"
    :footer="false"
    :mask-closable="true"
    :closable="false"
    :drawer-style="{ width: size }"
    class="md:!hidden"
    unmount-on-close
    @update:visible="emit('update:modelValue', $event)"
  >
    <slot />
  </a-drawer>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    modelValue: boolean
    size?: string
  }>(),
  { size: 'min(280px, 85vw)' },
)

const emit = defineEmits<{ 'update:modelValue': [boolean] }>()
</script>
