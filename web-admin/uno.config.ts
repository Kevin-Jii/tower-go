import { defineConfig, presetAttributify, presetIcons, presetUno } from 'unocss'

export default defineConfig({
  presets: [
    presetUno(),
    presetAttributify(),
    presetIcons({
      scale: 1.15,
      warn: true,
      extraProperties: {
        display: 'inline-block',
        'vertical-align': 'middle',
      },
    }),
  ],
  shortcuts: {
    'flex-center': 'flex items-center justify-center',
    'transition-sidebar': 'transition-all duration-200 ease-in-out',
    'input-surface':
      'rounded-xl border border-slate-200/90 bg-white text-slate-900 px-3.5 py-2 text-sm outline-none shadow-sm transition focus:border-indigo-400 focus:ring-2 focus:ring-indigo-100',
    'card-shell': 'rounded-2xl border border-slate-200/80 bg-white shadow-sm shadow-slate-900/5 overflow-hidden',
    'page-title': 'text-lg font-semibold text-slate-800 tracking-tight m-0',
  },
})
