import fs from 'fs'
import path from 'path'
import { glob } from 'glob'

// Admonition type mapping from mkdocs to vitepress
const admonitionMap = {
  'tip': 'tip',
  'warning': 'warning',
  'info': 'info',
  'danger': 'danger',
  'note': 'tip',
  'attention': 'warning'
}

function migrateFile(filePath) {
  let content = fs.readFileSync(filePath, 'utf8')
  const originalContent = content

  // Convert admonitions: !!! type "Title" -> ::: type Title
  // Handle multi-line admonitions with 4-space indentation
  content = content.replace(
    /^(!{3}) (\w+)(?: "([^"]*)")?\n((?:    .+(?:\n|$))+)/gm,
    (match, bangs, type, title, body) => {
      const vpType = admonitionMap[type] || 'info'
      const titlePart = title ? ` ${title}` : ''
      // Remove 4-space indentation from body
      const cleanBody = body.replace(/^    /gm, '')
      return `::: ${vpType}${titlePart}\n${cleanBody}:::\n`
    }
  )

  // Convert code block titles: ```lang title="title" -> ```lang [title]
  content = content.replace(
    /```(\w+)(?: linenums="1")? title="([^"]+)"/g,
    (match, lang, title) => {
      const lineNumbers = match.includes('linenums') ? ':line-numbers' : ''
      return `\`\`\`${lang}${lineNumbers} [${title}]`
    }
  )

  // Convert code blocks with just linenums: ```lang linenums="1" -> ```lang:line-numbers
  content = content.replace(
    /```(\w+) linenums="1"(?! title)/g,
    '```$1:line-numbers'
  )

  // Remove button classes: [text](link){ .md-button .md-button--primary } -> [text](link)
  content = content.replace(
    /\[([^\]]+)\]\(([^)]+)\)\{ \.md-button[^}]* \}/g,
    '[$1]($2)'
  )

  // Convert frontmatter hide directive for home page
  content = content.replace(
    /^---\nhide:\n- navigation\n- toc\n---/m,
    '---\nlayout: home\n---'
  )

  // Update image paths: ./images/ and ../images/ -> /images/
  content = content.replace(/\(\.\/images\//g, '(/images/')
  content = content.replace(/\(\.\.\/images\//g, '(/images/')

  // Remove image alignment attributes: { align=left } etc.
  content = content.replace(/\{ align=\w+(?:, width=\d+%)? \}/g, '')

  // Fix relative markdown links by removing .md extension for vitepress
  // e.g., ../configuration/overview.md -> ../configuration/overview
  content = content.replace(/\]\(([^)]+)\.md\)/g, ']($1)')

  // Fix reference-style link definitions that point to .md files
  content = content.replace(/^\[([^\]]+)\]: (.+)\.md$/gm, '[$1]: $2')

  if (content !== originalContent) {
    fs.writeFileSync(filePath, content)
    console.log(`Migrated: ${filePath}`)
    return true
  }
  return false
}

async function main() {
  const docsDir = path.resolve(process.cwd(), 'docs')
  const files = await glob('**/*.md', { cwd: docsDir })

  let migratedCount = 0
  for (const file of files) {
    const filePath = path.join(docsDir, file)
    if (migrateFile(filePath)) {
      migratedCount++
    }
  }

  console.log(`\nMigration complete. ${migratedCount} files updated.`)
}

main().catch(console.error)
