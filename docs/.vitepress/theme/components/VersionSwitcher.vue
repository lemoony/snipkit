<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'

interface Version {
  version: string
  title: string
  aliases: string[]
  type?: string
}

const versions = ref<Version[]>([])
const currentVersion = ref('dev')
const siteRoot = ref('/')
const isOpen = ref(false)

const sortedVersions = computed(() => {
  return versions.value.slice().sort((a, b) => {
    // dev always first
    if (a.version === 'dev') return -1
    if (b.version === 'dev') return 1
    // latest alias second
    if (a.aliases?.includes('latest')) return -1
    if (b.aliases?.includes('latest')) return 1
    // then by version number descending
    return b.version.localeCompare(a.version, undefined, { numeric: true })
  })
})

onMounted(async () => {
  // Detect site root and current version from URL
  // URL patterns:
  // - GitHub Pages: /snipkit/dev/page.html -> siteRoot=/snipkit, version=dev
  // - GitHub Pages: /snipkit/v1.6.1/page.html -> siteRoot=/snipkit, version=v1.6.1
  // - Local dev: /page.html -> siteRoot='', version=dev (default)

  const pathname = window.location.pathname

  // Pattern: optional repo name, then version (dev, latest, or vX.X.X)
  // e.g., /snipkit/dev/page or /snipkit/v1.6.1/page or /dev/page
  const versionPattern = /^(\/snipkit)?\/(v[\d.]+|dev|latest)(\/|$)/
  const match = pathname.match(versionPattern)

  if (match) {
    // Has a version in the path
    siteRoot.value = match[1] || ''  // e.g., '/snipkit' or ''
    currentVersion.value = match[2]   // e.g., 'dev', 'v1.6.1', 'latest'
  } else {
    // No version in path - local dev mode
    siteRoot.value = ''
    currentVersion.value = 'dev'
  }

  // Fetch versions.json from site root (or current origin for local dev)
  const versionsUrl = siteRoot.value ? `${siteRoot.value}/versions.json` : '/versions.json'
  try {
    const response = await fetch(versionsUrl)
    if (response.ok) {
      versions.value = await response.json()
    }
  } catch (e) {
    console.error('Failed to load versions from', versionsUrl, e)
  }
})

function getVersionLabel(v: Version): string {
  if (v.aliases?.includes('latest')) {
    return `${v.title} (latest)`
  }
  return v.title
}

function switchVersion(version: Version) {
  // Get the current page path without version prefix
  const pathname = window.location.pathname
  const pathAfterVersion = pathname.replace(/^(\/[^/]+)?\/(v[\d.]+|dev|latest)/, '')

  // Build new URL
  const newPath = `${siteRoot.value}/${version.version}${pathAfterVersion || '/'}`
  window.location.href = newPath
}

function handleClickOutside(event: MouseEvent) {
  const target = event.target as HTMLElement
  if (!target.closest('.version-switcher')) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})
</script>

<template>
  <div class="version-switcher" v-if="versions.length > 0">
    <button @click.stop="isOpen = !isOpen" class="version-button">
      <span class="version-text">{{ currentVersion }}</span>
      <svg
        class="arrow"
        :class="{ open: isOpen }"
        width="14"
        height="14"
        viewBox="0 0 24 24"
        fill="currentColor"
      >
        <path d="M7 10l5 5 5-5z"/>
      </svg>
    </button>
    <Transition name="dropdown">
      <ul v-if="isOpen" class="version-dropdown">
        <li v-for="v in sortedVersions" :key="v.version">
          <a
            @click.prevent="switchVersion(v)"
            :class="{ current: v.version === currentVersion }"
          >
            {{ getVersionLabel(v) }}
          </a>
        </li>
      </ul>
    </Transition>
  </div>
</template>

<style scoped>
.version-switcher {
  position: relative;
  margin-left: 16px;
}

.version-button {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 12px;
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  background: var(--vp-c-bg-soft);
  color: var(--vp-c-text-1);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.version-button:hover {
  border-color: var(--vp-c-brand-1);
  background: var(--vp-c-bg-mute);
}

.version-text {
  min-width: 48px;
  text-align: center;
}

.arrow {
  transition: transform 0.2s ease;
}

.arrow.open {
  transform: rotate(180deg);
}

.version-dropdown {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 140px;
  max-height: 320px;
  overflow-y: auto;
  padding: 8px 0;
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  list-style: none;
  z-index: 100;
}

.version-dropdown li a {
  display: block;
  padding: 8px 16px;
  color: var(--vp-c-text-2);
  font-size: 13px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.version-dropdown li a:hover {
  background: var(--vp-c-bg-mute);
  color: var(--vp-c-text-1);
}

.version-dropdown li a.current {
  color: var(--vp-c-brand-1);
  font-weight: 600;
}

/* Dropdown animation */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
