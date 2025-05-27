import { defineConfig } from "vite"
import vue from "@vitejs/plugin-vue"
import { resolve } from "path"
import AutoImport from "unplugin-auto-import/vite"
import Components from "unplugin-vue-components/vite"
import Icons from "unplugin-icons/vite"
import IconsResolver from "unplugin-icons/resolver"
import { ElementPlusResolver } from "unplugin-vue-components/resolvers"
import ElementPlus from "unplugin-element-plus/vite"

// https://vitejs.dev/config/
export default defineConfig({
  //项目根目录(index.ts.html 文件所在的位置,可以是相对路径也可以是绝对路径)，默认为process.cwd()
  root: process.cwd(),
  //开发或生产环境服务的公共基础路径。设置此路径后，项目中的静态资源和路由都会基于这个路径进行生成。
  base: "/",
  mode: "development",
  envDir: process.cwd(),
  //以 envPrefix 开头的环境变量会通过 import.meta.env 暴露在你的客户端源码中。
  envPrefix: "VUE_",
  clearScreen: true,
  cacheDir: "node_modules/.vite",
  define: {
    __VUE_I18N_FULL_INSTALL__: true,
    __VUE_I18N_LEGACY_API__: false,
    __INTLIFY_PROD_DEVTOOLS__: false,
    "process.env": {},
    __APP_VERSION__: JSON.stringify(process.env.npm_package_version)
  },
  resolve: {
    //定义别名，使用绝对路径。
    alias: {
      "@": resolve(__dirname, "./src"),
      "#": resolve(__dirname, "./src/types"),
      "utils": resolve(__dirname, "./src/utils"),
      "services": resolve(__dirname, "./src/services"),
      "models": resolve(__dirname, "./src/models")
    },
    //导入时要省略的扩展名列表
    extensions: [".mjs", ".js", ".ts", ".jsx", ".tsx", ".json"]
  },
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@use "@/styles/theme.scss" as *;`
      }
    }
  },
  plugins: [
    vue(),
    Icons({
      autoInstall: true,
      scale: 1,
      defaultStyle: "",
      defaultClass: "",
      compiler: "vue3"
    }),
    AutoImport({
      //包含匹配模式的正则表达式或字符串数组，指定了哪些文件应该被插件处理<自动导入api>
      include: [],
      //全局引入api并注册
      imports: ["vue", "vue-router", "vue-i18n", "pinia", "@vueuse/core"],
      //为目录下的默认模块导出启用按照文件名自动导入，默认为false
      defaultExportByFilename: false,
      //自动导入目录下的模块导出，默认只扫描目录下一层模块, 使用/**可以嵌套引入所有模块
      dirs: ["./src/stores/**", "./src/utils/**", "./src/services/**", "./src/models/**", "./src/types/**"],
      //生成相应*.d.ts的文件路径, dts的值可为string<文件路径>或者bool，默认值为true,默认路径为'./auto-imports.d.ts'
      dts: true,
      //是否在vue文件的template中可以直接使用api,默认为false
      vueTemplate: true,
      //将自动导入的语句放置在每个<script>标签的末尾，默认为true
      injectAtEnd: true,
      //通过自定义解析器来导入api
      resolvers: [
        ElementPlusResolver({
          importStyle: "sass"
        }),
        IconsResolver({
          prefix: "icon",
          enabledCollections: ["mdi", "noto"]
        })
      ],
      //为eslint检查生成.eslintrc-auto-import.json文件
      eslintrc: {
        enabled: true,
        filepath: "./.eslintrc-auto-import.json",
        globalsPropValue: true
      }
    }),
    Components({
      //注册全局组件类型，默认会将vue-router中的RouterLink和RouterView进行全局注册
      types: [],
      //包含匹配模式的正则表达式或字符串数组，指定了哪些文件应该被插件处理<自动导入并注册组件>
      include: [],
      //组件目录数组，可以将这些目录指定为需要自动导入和注册的组件目录
      dirs: ["src/components"],
      //文件扩展名的数组，指定应该扫描的文件扩展名，默认为['vue']
      extensions: ["vue"],
      //是否查找子目录
      deep: true,
      //是否生成TS声明文件，类型为string<文件路径>或者bool，默认值为true, 路径为'./components.d.ts'，默认为true
      dts: true,
      //默认导入指令，vue3默认值为true, vue2默认值为false
      directives: true,
      //通过自定义解析器来导入并注册组件
      resolvers: [
        ElementPlusResolver({
          importStyle: "sass"
        }),
        IconsResolver({
          prefix: "icon",
          enabledCollections: ["mdi", "noto"]
        })
      ]
    }),
    ElementPlus({
      useSource: true
    })
  ],
  server: {
    host: "localhost",
    port: 5100,
    strictPort: true,
    open: true,
    proxy: {
      "/webapi": "http://localhost:8100"
    }
  },
  preview: {
    host: "localhost",
    port: 5131,
    strictPort: true,
    open: true
  },
  // build: {
  //   outDir: "dist",
  //   assetsDir: "assets",
  //   emptyOutDir: true,
  //   copyPublicDir: true,
  //   reportCompressedSize: false,
  //   chunkSizeWarningLimit: 4096,
  //   rollupOptions: {
  //     //此处可以添加多个入口文件，适用于多页面应用模式
  //     input: {
  //       main: resolve(__dirname, "index.html")
  //     },
  //     output: {
  //       manualChunks(id) {
  //         if (id.includes("node_modules")) {
  //           return id.toString().split("node_modules/")[1].split("/")[0].toString()
  //         }
  //       }
  //     }
  //   }
  // },
  experimental: {
    //可以自定义不同的文件类型的访问路径。
    renderBuiltUrl: (
      filename: string,
      value: {
        type: "asset" | "public"
        hostId: string
        hostType: "js" | "css" | "html"
        ssr: boolean
      }
    ): string | { relative?: boolean; runtime?: string } | undefined => {
      return { relative: false }
    }
  }
})
