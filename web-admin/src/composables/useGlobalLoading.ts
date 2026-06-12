import { useLoadingStore } from '@/store/loading'

/** 全局 math-curve loading（全屏遮罩） */
export function useGlobalLoading() {
  const store = useLoadingStore()
  return {
    visible: store.visible,
    message: store.message,
    show: store.show,
    hide: store.hide,
    wrap: store.wrap,
  }
}
