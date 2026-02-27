import { defineConfig } from "vitepress";

export default defineConfig({
  title: "SnipKit Documentation",
  description:
    "SnipKit helps you to execute scripts saved in your favorite snippets manager without even leaving the terminal.",

  base: process.env.VITEPRESS_BASE || "/",

  head: [["link", { rel: "icon", href: "/images/logo.png" }]],

  themeConfig: {
    logo: "/images/logo.png",

    nav: [
      { text: "Home", link: "/" },
      { text: "Getting Started", link: "/getting-started/overview" },
      { text: "Configuration", link: "/configuration/overview" },
      { text: "Managers", link: "/managers/overview" },
      { text: "Assistant", link: "/assistant/" },
    ],

    sidebar: {
      "/getting-started/": [
        {
          text: "Getting Started",
          items: [
            { text: "Overview", link: "/getting-started/overview" },
            { text: "Parameters", link: "/getting-started/parameters" },
            { text: "Power Setup", link: "/getting-started/power-setup" },
            { text: "Fzf", link: "/getting-started/fzf" },
          ],
        },
      ],
      "/configuration/": [
        {
          text: "Configuration",
          items: [
            { text: "Overview", link: "/configuration/overview" },
            { text: "Themes", link: "/configuration/themes" },
          ],
        },
      ],
      "/managers/": [
        {
          text: "Managers",
          items: [
            { text: "Overview", link: "/managers/overview" },
            { text: "File System Library", link: "/managers/fslibrary" },
            { text: "GitHub Gist", link: "/managers/githubgist" },
            { text: "SnippetsLab", link: "/managers/snippetslab" },
            { text: "Snip", link: "/managers/pictarinesnip" },
            { text: "Pet", link: "/managers/pet" },
          ],
        },
      ],
      "/assistant/": [
        {
          text: "Assistant",
          items: [
            { text: "Overview", link: "/assistant/" },
            { text: "OpenAI", link: "/assistant/openai" },
            { text: "Anthropic", link: "/assistant/anthropic" },
            { text: "Gemini", link: "/assistant/gemini" },
            { text: "Ollama", link: "/assistant/ollama" },
            { text: "OpenAI-Compatible", link: "/assistant/openai-compatible" },
          ],
        },
      ],
    },

    socialLinks: [
      { icon: "github", link: "https://github.com/lemoony/snipkit" },
    ],

    footer: {
      message: "Released under the Apache License 2.0.",
      copyright: "Copyright © 2026 lemoony",
    },

    search: {
      provider: "local",
    },
  },

  markdown: {
    theme: {
      light: "catppuccin-latte",
      dark: "catppuccin-mocha",
    },
  },
});
