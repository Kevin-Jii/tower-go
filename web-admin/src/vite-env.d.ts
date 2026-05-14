/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE: string
  /** 可选，覆盖默认的 X-Client-Source（默认 web-admin） */
  readonly VITE_CLIENT_SOURCE?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
