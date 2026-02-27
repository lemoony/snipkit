import DefaultTheme from 'vitepress/theme'
import '@catppuccin/vitepress/theme/mocha/mauve.css'
import Layout from './Layout.vue'
import type { Theme } from 'vitepress'

export default {
  extends: DefaultTheme,
  Layout
} satisfies Theme
