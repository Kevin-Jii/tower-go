<template>
  <div
    class="math-curve-loader"
    :class="[
      `math-curve-loader--${size}`,
      { 'math-curve-loader--inline': inline, 'math-curve-loader--overlay': overlay },
    ]"
    role="status"
    :aria-label="text || '加载中'"
  >
    <div
      ref="canvasRef"
      class="math-curve-loader__canvas"
      :style="{ width: `${pixelSize}px`, height: `${pixelSize}px` }"
      aria-hidden="true"
    />
    <p v-if="text && !inline" class="math-curve-loader__text">{{ text }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { MATH_CURVE_PRESETS } from './mathCurve/curves'
import { mountMathCurveLoader } from './mathCurve/engine'
import type { MathCurveInstance } from './mathCurve/engine'
import { MATH_CURVE_SIZES } from './mathCurve/types'
import type { MathCurveVariant } from './mathCurve/types'

const props = withDefaults(
  defineProps<{
    size?: keyof typeof MATH_CURVE_SIZES
    variant?: MathCurveVariant
    text?: string
    /** 行内模式（按钮内） */
    inline?: boolean
    /** 全屏遮罩模式 */
    overlay?: boolean
  }>(),
  {
    size: 'md',
    variant: 'thinking',
    text: '',
    inline: false,
    overlay: false,
  },
)

const canvasRef = ref<HTMLElement | null>(null)
let instance: MathCurveInstance | null = null

const pixelSize = computed(() => MATH_CURVE_SIZES[props.size])

function mount(): void {
  instance?.destroy()
  instance = null
  const el = canvasRef.value
  if (!el) return
  instance = mountMathCurveLoader(el, MATH_CURVE_PRESETS[props.variant])
}

onMounted(mount)

watch(() => props.variant, mount)

onBeforeUnmount(() => {
  instance?.destroy()
  instance = null
})
</script>

<style scoped>
.math-curve-loader {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: rgb(var(--primary-6, 99 102 241));
}

.math-curve-loader--inline {
  flex-direction: row;
  gap: 0;
  display: inline-flex;
  vertical-align: middle;
}

.math-curve-loader--overlay {
  position: fixed;
  inset: 0;
  z-index: 10000;
  background: rgba(255, 255, 255, 0.72);
  backdrop-filter: blur(4px);
  color: rgb(var(--primary-6, 99 102 241));
}

.math-curve-loader__canvas {
  display: grid;
  place-items: center;
  flex-shrink: 0;
}

.math-curve-loader__canvas :deep(.math-curve-loader__svg) {
  width: 100%;
  height: 100%;
  overflow: visible;
}

.math-curve-loader__text {
  margin: 0;
  font-size: 14px;
  color: var(--color-text-2, #64748b);
}
</style>
