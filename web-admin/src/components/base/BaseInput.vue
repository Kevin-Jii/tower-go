<template>
  <a-input
    ref="arcoRef"
    :model-value="displayValue"
    :type="effectiveType"
    :placeholder="placeholder"
    :disabled="disabled"
    :max-length="maxLengthForArco"
    :allow-clear="clearable"
    :autocomplete="autocomplete"
    @update:model-value="onUpdateModel"
    @press-enter="emit('enter', $event)"
  >
    <template v-if="showPasswordToggle" #suffix>
      <a-button type="text" size="mini" class="!px-1" @click="revealed = !revealed">
        {{ revealed ? '隐藏' : '显示' }}
      </a-button>
    </template>
  </a-input>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue: string | number
    type?: string
    placeholder?: string
    disabled?: boolean
    maxlength?: number | string
    autocomplete?: string
    clearable?: boolean
    showPassword?: boolean
  }>(),
  {
    type: 'text',
    disabled: false,
    clearable: false,
    autocomplete: 'off',
    showPassword: false,
  },
)

const revealed = ref(false)
const showPasswordToggle = computed(() => props.type === 'password' && props.showPassword)
const effectiveType = computed((): 'text' | 'password' => {
  if (props.type === 'password' && props.showPassword && revealed.value) return 'text'
  if (props.type === 'password') return 'password'
  return 'text'
})

const displayValue = computed(() =>
  props.modelValue === null || props.modelValue === undefined ? '' : String(props.modelValue),
)

const maxLengthForArco = computed(() => {
  const m = props.maxlength
  if (m === undefined || m === '') return undefined
  const n = typeof m === 'string' ? Number(m) : m
  return Number.isFinite(n) && n > 0 ? n : undefined
})

const emit = defineEmits<{
  'update:modelValue': [string | number]
  enter: [KeyboardEvent]
}>()

const arcoRef = ref<{ focus?: () => void } | null>(null)

function onUpdateModel(v: string): void {
  if (props.type === 'number') {
    if (v === '') {
      emit('update:modelValue', '')
      return
    }
    const n = Number(v)
    emit('update:modelValue', Number.isNaN(n) ? v : n)
    return
  }
  emit('update:modelValue', v)
}

defineExpose({ focus: () => arcoRef.value?.focus?.() })
</script>
