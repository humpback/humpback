import vue from "eslint-plugin-vue"
import typescriptEslint from "@typescript-eslint/eslint-plugin"
import html from "eslint-plugin-html"
import _import from "eslint-plugin-import"
import prettier from "eslint-plugin-prettier"
import vueESLintParser from "vue-eslint-parser"
import { fixupPluginRules } from "@eslint/compat"
import globals from "globals"
import path from "node:path"
import { fileURLToPath } from "node:url"
import { FlatCompat } from "@eslint/eslintrc"
import js from "@eslint/js"

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)
const compat = new FlatCompat({
  baseDirectory: __dirname,
  recommendedConfig: js.configs.recommended,
  allConfig: js.configs.all
})

export default [
  {
    files: ["src/**/*.js", "src/**/*.vue", "src/**/*.ts", "src/**/*.jsx", "src/**/*.tsx"],
    ignores: ["**/vite.config.ts", "**/node_modules", "**/dist"]
  },
  ...compat.extends("eslint:recommended", "plugin:vue/vue3-essential", "plugin:@typescript-eslint/recommended", "prettier", "./.eslintrc-auto-import.json"),
  {
    plugins: {
      vue,
      typescriptEslint,
      html,
      import: fixupPluginRules(_import),
      prettier
    },

    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.node,
        ...globals.es2021
      },

      parser: vueESLintParser,
      ecmaVersion: "latest",
      sourceType: "module",

      parserOptions: {
        parser: "@typescript-eslint/parser"
      }
    },

    rules: {
      "prettier/prettier": "error",
      "no-bitwise": "off",
      "no-tabs": "off",
      "no-console": ["warn", { allow: ["warn", "error"] }],
      "array-element-newline": ["error", "consistent"],
      "object-curly-spacing": ["error", "always"],
      "no-new": "off",
      "linebreak-style": "off",
      "import/extensions": "off",
      "eol-last": "off",
      "no-shadow": "off",
      "no-unused-vars": 0,
      "@typescript-eslint/no-unused-vars": [
        "error",
        {
          args: "none",
          argsIgnorePattern: "^_",
          caughtErrors: "all",
          caughtErrorsIgnorePattern: "^_",
          destructuredArrayIgnorePattern: "^_",
          varsIgnorePattern: "^_"
        }
      ],
      "import/no-cycle": "off",
      "arrow-parens": "off",
      "semi": ["error", "never"],
      "eqeqeq": "off",
      "no-param-reassign": "off",
      "import/prefer-default-export": "off",
      "no-use-before-define": "off",
      "no-continue": "off",
      "prefer-destructuring": "off",
      "no-plusplus": "off",
      "prefer-const": "off",
      "global-require": "off",
      "no-prototype-builtins": "off",
      "consistent-return": "off",
      "vue/require-component-is": "off",
      "prefer-template": "off",
      "one-var-declaration-per-line": "off",
      "one-var": "off",
      "import/named": "off",
      "object-curly-newline": "off",
      "default-case": "off",
      "import/order": "off",
      "no-trailing-spaces": "off",
      "func-names": "off",
      "radix": "off",
      "no-unused-expressions": "off",
      "no-underscore-dangle": "off",
      "no-nested-ternary": "off",
      "no-restricted-syntax": "off",
      "no-mixed-operators": "off",
      "no-await-in-loop": "off",
      "import/no-extraneous-dependencies": "off",
      "import/no-unresolved": "off",
      "no-case-declarations": "off",
      "template-curly-spacing": "off",
      "vue/valid-v-for": "off",
      "@typescript-eslint/no-var-requires": "off",
      "@typescript-eslint/no-empty-function": "off",
      "@typescript-eslint/no-explicit-any": "off",
      "@typescript-eslint/ban-types": "off",
      "@typescript-eslint/ban-ts-comment": "off",
      "no-empty": "off",
      "guard-for-in": "off",
      "class-methods-use-this": "off",
      "no-return-await": "off",
      "vue/html-self-closing": "off",
      "vue/singleline-html-element-content-newline": "off",
      "vue/valid-template-root": "off",
      "vue/multi-word-component-names": "off",
      "vue/return-in-computed-property": "off",
      "vue/no-side-effects-in-computed-properties": "off"
    }
  }
]
