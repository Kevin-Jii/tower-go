import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export const useLoadingStore = defineStore('loading', () => {
  const count = ref(0)
  const message = ref('数据加载中...')

  const visible = computed(() => count.value > 0)

  function show(text = '数据加载中...'): void {
    message.value = text
    count.value += 1
  }

  function hide(): void {
    count.value = Math.max(0, count.value - 1)
  }

  /** 包装异步任务，自动显示/隐藏全局 loading */
  async function wrap<T>(task: () => Promise<T>, text = '数据加载中...'): Promise<T> {
    show(text)
    try {
      return await task()
    } finally {
      hide()
    }
  }

  return { count, message, visible, show, hide, wrap }
})
