<template>
  <a-card
    class="card-shell flex h-full min-w-0 flex-col"
    :bordered="true"
    :body-style="bodyStyle"
  >
    <template v-if="$slots.header" #title>
      <slot name="header" />
    </template>
    <template v-else-if="title" #title>{{ title }}</template>
    <slot />
  </a-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    title?: string
    /** 为 true 时去掉 body 默认 padding（Arco `.arco-card-body` 自带 padding 会占高度）；也可用 bodyPadding 自定义 */
    flushBody?: boolean
    /** 覆盖 body 的 padding，如 `12px 16px`；与 flushBody 同时存在时以本项为准 */
    bodyPadding?: string
  }>(),
  { title: '', flushBody: false, bodyPadding: undefined },
)

/** 让卡片主体成为纵向 flex 容器，子元素可用 flex-1 占满剩余高度（如 BaseTable） */
const bodyStyle = computed(() => {
  const pad =
    props.bodyPadding !== undefined && props.bodyPadding !== ''
      ? props.bodyPadding
      : props.flushBody
        ? '0px'
        : undefined
  return {
    flex: 1,
    minHeight: 0,
    display: 'flex',
    flexDirection: 'column' as const,
    overflow: 'hidden',
    ...(pad !== undefined ? { padding: pad } : {}),
  }
})
</script>
