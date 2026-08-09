const SCRIPT_BLOCK_RE = /<\s*script\b[^>]*>[\s\S]*?<\s*\/\s*script\s*>/gi
const STYLE_BLOCK_RE = /<\s*style\b[^>]*>[\s\S]*?<\s*\/\s*style\s*>/gi
const LEFTOVER_TAG_RE = /<\s*\/?\s*(script|style)\b[^>]*>/gi
const EVENT_ATTR_RE = /\s+on\w+\s*=\s*("([^"]*)"|'([^']*)'|[^\s>]+)/gi
const DANGEROUS_URL_RE = /\s(href|src)\s*=\s*("|')?(javascript|vbscript|data):/gi

export function sanitizeHtml(html: string): string {
  return html
    .replace(SCRIPT_BLOCK_RE, '')
    .replace(STYLE_BLOCK_RE, '')
    .replace(LEFTOVER_TAG_RE, '')
    .replace(EVENT_ATTR_RE, '')
    .replace(DANGEROUS_URL_RE, '')
}
