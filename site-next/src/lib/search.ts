import fs from "fs";
import path from "path";
import matter from "gray-matter";
import { docsNav, type NavItem } from "./docs-nav";

const CONTENT_DIR = path.join(process.cwd(), "src", "content", "docs");

export interface SearchEntry {
  title: string;
  description: string;
  slug: string;
  section: string;
  content: string; // plain text preview (stripped of MDX/HTML)
}

/** Strip MDX/JSX/HTML tags and import lines to get plain text for search */
function stripToPlainText(mdx: string): string {
  return mdx
    // Remove import statements
    .replace(/^import\s+.*$/gm, "")
    // Remove frontmatter (already stripped by gray-matter, but just in case)
    .replace(/^---[\s\S]*?---/m, "")
    // Remove JSX/HTML tags
    .replace(/<[^>]+>/g, " ")
    // Remove JSX expressions { ... }
    .replace(/\{[\s\S]*?\}/g, " ")
    // Remove markdown syntax
    .replace(/#{1,6}\s+/g, "")
    .replace(/\*{1,3}([^*]+)\*{1,3}/g, "$1")
    .replace(/`{1,3}[^`]*`{1,3}/g, "")
    .replace(/\[([^\]]+)\]\([^)]+\)/g, "$1")
    // Collapse whitespace
    .replace(/\s+/g, " ")
    .trim()
    .slice(0, 500); // keep first 500 chars for search
}

/** Resolve which top-level nav section a slug belongs to */
function findSection(slug: string): string {
  for (const group of docsNav) {
    if (containsSlug(group, slug)) {
      return group.label;
    }
  }
  return "";
}

function containsSlug(item: NavItem, slug: string): boolean {
  if (item.slug === slug) return true;
  if (item.items) {
    return item.items.some((child) => containsSlug(child, slug));
  }
  return false;
}

export function buildSearchIndex(): SearchEntry[] {
  const entries: SearchEntry[] = [];

  function walk(dir: string, prefix: string) {
    if (!fs.existsSync(dir)) return;
    const items = fs.readdirSync(dir, { withFileTypes: true });
    for (const item of items) {
      if (item.isDirectory()) {
        walk(path.join(dir, item.name), `${prefix}${item.name}/`);
      } else if (item.name.endsWith(".mdx")) {
        const name = item.name.replace(/\.mdx$/, "");
        const slug =
          name === "index"
            ? prefix.replace(/\/$/, "")
            : `${prefix}${name}`;

        if (!slug) continue; // skip root index (that's the homepage in the old site)

        const raw = fs.readFileSync(path.join(dir, item.name), "utf-8");
        const { data, content } = matter(raw);

        entries.push({
          title: data.title || slug,
          description: data.description || "",
          slug,
          section: findSection(slug),
          content: stripToPlainText(content),
        });
      }
    }
  }

  walk(CONTENT_DIR, "");
  return entries;
}
