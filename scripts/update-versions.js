import fs from 'fs'
import path from 'path'

// Usage: node scripts/update-versions.js <version> <type> [latest]
// Example: node scripts/update-versions.js v1.7.0 vitepress latest
// Example: node scripts/update-versions.js dev vitepress

const args = process.argv.slice(2)
if (args.length < 2) {
  console.error('Usage: node scripts/update-versions.js <version> <type> [latest]')
  console.error('Example: node scripts/update-versions.js v1.7.0 vitepress latest')
  process.exit(1)
}

const [version, type, ...aliases] = args
const versionsPath = path.join(process.cwd(), 'versions.json')

let versions = []
if (fs.existsSync(versionsPath)) {
  versions = JSON.parse(fs.readFileSync(versionsPath, 'utf8'))
}

// Remove 'latest' alias from all existing versions if we're setting a new latest
if (aliases.includes('latest')) {
  versions = versions.map(v => ({
    ...v,
    aliases: (v.aliases || []).filter(a => a !== 'latest')
  }))
}

// Find or create version entry
const existingIndex = versions.findIndex(v => v.version === version)
const newEntry = {
  version,
  title: version,
  aliases: aliases,
  type: type
}

if (existingIndex >= 0) {
  versions[existingIndex] = newEntry
} else {
  // Insert at appropriate position
  // dev always first, then by version number descending
  if (version === 'dev') {
    versions.unshift(newEntry)
  } else {
    // Find insertion point after dev but sorted by version
    let insertIndex = versions.findIndex(v => v.version !== 'dev')
    if (insertIndex === -1) insertIndex = versions.length

    // Insert and sort non-dev versions
    versions.splice(insertIndex, 0, newEntry)
  }
}

fs.writeFileSync(versionsPath, JSON.stringify(versions, null, 2) + '\n')
console.log(`Updated versions.json with ${version} (type: ${type}, aliases: ${aliases.join(', ') || 'none'})`)
