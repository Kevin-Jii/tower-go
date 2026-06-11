import { fileURLToPath, URL } from "node:url";
import vue from "@vitejs/plugin-vue";
import UnoCSS from "unocss/vite";
import AutoImport from "unplugin-auto-import/vite";
import { defineConfig } from "vite";

export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (!id.includes("node_modules")) return;
          if (id.includes("@arco-design")) return "vendor-arco";
          if (id.includes("echarts") || id.includes("zrender"))
            return "vendor-echarts";
          if (id.includes("handsontable")) return "vendor-handsontable";
          if (id.includes("@tanstack")) return "vendor-vue-query";
          if (id.includes("@vueuse")) return "vendor-vueuse";
          if (id.includes("axios")) return "vendor-axios";
          if (id.includes("vue-router")) return "vendor-vue-runtime";
          if (id.includes("pinia")) return "vendor-vue-runtime";
          if (id.includes("/vue/") || id.includes("\\vue\\"))
            return "vendor-vue-runtime";
        },
      },
    },
    chunkSizeWarningLimit: 900,
  },
  plugins: [
    vue(),
    UnoCSS(),
    AutoImport({
      imports: ["vue", "vue-router", "pinia"],
      dts: "src/auto-imports.d.ts",
    }),
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    port: 5173,
    proxy: {
      "/api": {
        target: "https://tower.usove.online",
        changeOrigin: true,
      },
    },
  },
});
