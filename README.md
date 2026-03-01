# getm

A CLI tool that fetches any URL as clean Markdown, powered by Cloudflare's [markdown.new](https://markdown.new/) API.

## Install

```
go install github.com/bscott/getm@latest
```

## Usage

```
getm <url>
```

The `https://` prefix is added automatically if omitted.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-method` | `auto` | Conversion method: `auto`, `ai`, or `browser` |
| `-images` | `false` | Retain image references in output |
| `-json` | `false` | Output raw JSON response (includes title, token count, timing) |

### Examples

```bash
# Fetch a page as markdown
getm example.com

# Keep images in the output
getm -images https://go.dev/blog

# Force headless browser for JS-heavy sites
getm -method browser https://react.dev

# Get raw JSON with metadata
getm -json https://example.com
```

## API

Uses the [markdown.new](https://markdown.new/) API by Cloudflare. No API key required. Rate limited to 500 requests/day per IP.
