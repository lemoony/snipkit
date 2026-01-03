# SnipKit Documentation

This folder contains the VitePress-based documentation for SnipKit.

## Local Development

```bash
# Install dependencies (from repo root)
npm install

# Start dev server
npm run docs:dev

# Build for production
npm run docs:build

# Preview production build
npm run docs:preview
```

## Structure

```
docs/
├── .vitepress/
│   ├── config.ts              # VitePress configuration
│   └── theme/
│       ├── index.ts           # Theme setup (Catppuccin)
│       ├── Layout.vue         # Custom layout with version switcher
│       └── components/
│           └── VersionSwitcher.vue
├── public/
│   ├── images/                # Static assets
│   └── versions.json          # Mock versions for local dev
├── getting-started/
├── configuration/
├── managers/
├── assistant/
└── index.md                   # Home page
```

## Versioning

Documentation uses directory-based versioning similar to mike:

- Each version is deployed to its own directory on gh-pages (`/dev/`, `/v1.6.1/`, etc.)
- A `versions.json` file at the gh-pages root lists all available versions
- The version switcher reads this file and allows navigation between versions
- Old mkdocs versions remain accessible alongside new VitePress versions

### versions.json Format

```json
[
  {"version": "dev", "title": "dev", "aliases": []},
  {"version": "v1.7.0", "title": "v1.7.0", "aliases": ["latest"]},
  {"version": "v1.6.1", "title": "v1.6.1", "aliases": []}
]
```

## Deployment

Two GitHub Actions workflows handle deployment:

### Dev Deployment (`.github/workflows/vitepress-dev.yml`)
- **Trigger**: Push to `main` branch (when `docs/**` or `package.json` changes)
- **Deploys to**: `gh-pages/dev/`

### Release Deployment (`.github/workflows/vitepress-release.yml`)
- **Trigger**: Push tag `v*` or manual workflow dispatch
- **Deploys to**: `gh-pages/$VERSION/` and updates `gh-pages/latest/`

## Theme

Uses [Catppuccin](https://github.com/catppuccin/vitepress) theme:
- **Dark mode**: Mocha flavor with Mauve accent
- **Light mode**: Latte flavor

## Migration from MkDocs

The docs were migrated from MkDocs. Key syntax conversions:

| MkDocs | VitePress |
|--------|-----------|
| `!!! tip "Title"` | `::: tip Title` |
| `!!! warning` | `::: warning` |
| `` ```yaml title="file"`` | `` ```yaml [file]`` |
| `{ .md-button }` | Plain links |
