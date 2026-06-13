<template>
  <span ref="el" class="count-up-number">{{ fallbackText }}</span>
</template>

<script setup lang="ts">
import { CountUp } from 'countup.js'
import { computed, nextTick, onBeforeUnmount, onMounted, ref, useAttrs, watch } from 'vue'

defineOptions({ inheritAttrs: true })

const props = withDefaults(
  defineProps<{
    value?: number | string | null
    decimals?: number
    duration?: number
    prefix?: string
    suffix?: string
    useGrouping?: boolean
  }>(),
  {
    decimals: 0,
    duration: 0.9,
    prefix: '',
    suffix: '',
    useGrouping: false,
  },
)

const el = ref<HTMLElement | null>(null)
const attrs = useAttrs()
let countUp: CountUp | null = null

const numericValue = computed(() => {
  const n = Number(props.value ?? 0)
  return Number.isFinite(n) ? n : 0
})

const fallbackText = computed(() => {
  return `${props.prefix}${numericValue.value.toFixed(props.decimals)}${props.suffix}`
})

function createCounter(): void {
  if (!el.value) return
  countUp = null
  countUp = new CountUp(el.value, numericValue.value, {
    startVal: 0,
    decimalPlaces: props.decimals,
    duration: props.duration,
    useGrouping: props.useGrouping,
    prefix: props.prefix,
    suffix: props.suffix,
  })
  if (!countUp.error) {
    countUp.start()
  }
}

function updateCounter(): void {
  createCounter()
}

onMounted(() => {
  void nextTick(createCounter)
})

watch(
  () => [numericValue.value, props.decimals, props.duration, props.prefix, props.suffix, props.useGrouping, attrs.class, attrs.style] as const,
  () => {
    void nextTick(updateCounter)
  },
)

onBeforeUnmount(() => {
  countUp = null
})
</script>
