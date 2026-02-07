"use client";

import {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from "react";
import { useRouter } from "next/navigation";

interface SearchEntry {
  title: string;
  description: string;
  slug: string;
  section: string;
  content: string;
}

interface SearchResult extends SearchEntry {
  /** matched snippet from content (with query highlighted) */
  snippet: string;
}

/* ------------------------------------------------------------------ */
/*  Fuzzy matching helpers                                             */
/* ------------------------------------------------------------------ */

function fuzzyMatch(text: string, query: string): boolean {
  const lower = text.toLowerCase();
  const q = query.toLowerCase();
  // simple substring match plus word-start matching
  if (lower.includes(q)) return true;
  // check if every word in the query appears somewhere
  const words = q.split(/\s+/).filter(Boolean);
  return words.every((w) => lower.includes(w));
}

function extractSnippet(content: string, query: string): string {
  const lower = content.toLowerCase();
  const q = query.toLowerCase();
  const idx = lower.indexOf(q);
  if (idx === -1) return content.slice(0, 120) + "...";
  const start = Math.max(0, idx - 40);
  const end = Math.min(content.length, idx + query.length + 80);
  let snippet = "";
  if (start > 0) snippet += "...";
  snippet += content.slice(start, end);
  if (end < content.length) snippet += "...";
  return snippet;
}

function highlightMatch(text: string, query: string): React.ReactNode[] {
  if (!query.trim()) return [text];
  const regex = new RegExp(
    `(${query.replace(/[.*+?^${}()|[\]\\]/g, "\\$&")})`,
    "gi"
  );
  const parts = text.split(regex);
  return parts.map((part, i) =>
    regex.test(part) ? (
      <mark key={i} className="bg-accent/25 text-white rounded-sm px-px">
        {part}
      </mark>
    ) : (
      <span key={i}>{part}</span>
    )
  );
}

/* ------------------------------------------------------------------ */
/*  Main component                                                     */
/* ------------------------------------------------------------------ */

export default function SearchModal({
  open,
  onClose,
}: {
  open: boolean;
  onClose: () => void;
}) {
  const router = useRouter();
  const inputRef = useRef<HTMLInputElement>(null);
  const listRef = useRef<HTMLDivElement>(null);
  const [query, setQuery] = useState("");
  const [index, setIndex] = useState<SearchEntry[]>([]);
  const [activeIndex, setActiveIndex] = useState(0);
  const [loaded, setLoaded] = useState(false);

  /* Load search index on first open */
  useEffect(() => {
    if (open && !loaded) {
      fetch("/api/search")
        .then((r) => r.json())
        .then((data: SearchEntry[]) => {
          setIndex(data);
          setLoaded(true);
        })
        .catch(console.error);
    }
  }, [open, loaded]);

  /* Auto-focus input when modal opens */
  useEffect(() => {
    if (open) {
      setQuery("");
      setActiveIndex(0);
      // small delay so the modal finishes rendering
      requestAnimationFrame(() => inputRef.current?.focus());
    }
  }, [open]);

  /* Compute results */
  const results: SearchResult[] = useMemo(() => {
    if (!query.trim()) return [];
    return index
      .filter(
        (entry) =>
          fuzzyMatch(entry.title, query) ||
          fuzzyMatch(entry.description, query) ||
          fuzzyMatch(entry.content, query)
      )
      .map((entry) => ({
        ...entry,
        snippet: extractSnippet(entry.content, query),
      }))
      .slice(0, 20);
  }, [query, index]);

  /* Group results by section */
  const grouped = useMemo(() => {
    const map = new Map<string, SearchResult[]>();
    for (const r of results) {
      const section = r.section || "Other";
      if (!map.has(section)) map.set(section, []);
      map.get(section)!.push(r);
    }
    return map;
  }, [results]);

  /* Flat result list for keyboard navigation */
  const flatResults = results;

  /* Navigate to a result */
  const goTo = useCallback(
    (slug: string) => {
      onClose();
      router.push(`/docs/${slug}`);
    },
    [onClose, router]
  );

  /* Keyboard navigation */
  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (e.key === "ArrowDown") {
        e.preventDefault();
        setActiveIndex((prev) => Math.min(prev + 1, flatResults.length - 1));
      } else if (e.key === "ArrowUp") {
        e.preventDefault();
        setActiveIndex((prev) => Math.max(prev - 1, 0));
      } else if (e.key === "Enter") {
        e.preventDefault();
        if (flatResults[activeIndex]) {
          goTo(flatResults[activeIndex].slug);
        }
      } else if (e.key === "Escape") {
        e.preventDefault();
        onClose();
      }
    },
    [flatResults, activeIndex, goTo, onClose]
  );

  /* Scroll active item into view */
  useEffect(() => {
    const el = listRef.current?.querySelector(`[data-idx="${activeIndex}"]`);
    el?.scrollIntoView({ block: "nearest" });
  }, [activeIndex]);

  /* Reset active index when results change */
  useEffect(() => {
    setActiveIndex(0);
  }, [results]);

  if (!open) return null;

  return (
    <>
      {/* Backdrop + centering wrapper -- single element covers full viewport */}
      <div
        className="fixed top-0 left-0 right-0 bottom-0 z-[100] h-[100vh] bg-black/60 backdrop-blur-sm flex items-start justify-center pt-[min(15vh,120px)] px-4 overflow-y-auto"
        onClick={onClose}
      >
        {/* Modal card -- stop clicks from bubbling to backdrop */}
        <div
          className="w-full max-w-[640px] bg-dark-deep border border-gray-5 rounded-xl shadow-2xl overflow-hidden flex flex-col max-h-[min(70vh,600px)] my-auto mt-0 mb-auto"
          role="dialog"
          aria-label="Search documentation"
          onKeyDown={handleKeyDown}
          onClick={(e) => e.stopPropagation()}
        >
          {/* Search input */}
          <div className="flex items-center gap-3 px-4 border-b border-gray-5">
            {/* Search icon */}
            <svg
              className="w-5 h-5 text-gray-4 shrink-0"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth={2}
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M21 21l-5.197-5.197m0 0A7.5 7.5 0 105.196 5.196a7.5 7.5 0 0010.607 10.607z"
              />
            </svg>
            <input
              ref={inputRef}
              type="text"
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Search docs..."
              className="flex-1 bg-transparent border-none outline-none text-white text-base py-4 placeholder:text-gray-4"
              autoComplete="off"
              spellCheck={false}
            />
            <kbd className="hidden sm:flex items-center gap-0.5 text-[11px] text-gray-4 bg-gray-6 border border-gray-5 rounded px-1.5 py-0.5 font-mono">
              Esc
            </kbd>
          </div>

          {/* Results */}
          <div
            ref={listRef}
            className="flex-1 overflow-y-auto overscroll-contain"
          >
            {query.trim() === "" ? (
              /* Empty state -- no query */
              <div className="px-6 py-12 text-center text-gray-4 text-sm bg-gray-6/80">
                Type to search the documentation
              </div>
            ) : flatResults.length === 0 ? (
              /* No results */
              <div className="px-6 py-12 text-center">
                <p className="text-gray-3 text-sm">
                  No results for &ldquo;
                  <span className="text-white">{query}</span>&rdquo;
                </p>
                <p className="text-gray-4 text-xs mt-1">
                  Try a different search term
                </p>
              </div>
            ) : (
              /* Result groups */
              <div className="py-2">
                {Array.from(grouped.entries()).map(
                  ([section, sectionResults]) => (
                    <div key={section}>
                      <div className="px-4 pt-3 pb-1">
                        <span className="text-[11px] font-semibold text-gray-4 uppercase tracking-wider">
                          {section}
                        </span>
                      </div>
                      {sectionResults.map((result) => {
                        const idx = flatResults.indexOf(result);
                        const isActive = idx === activeIndex;
                        return (
                          <button
                            key={result.slug}
                            data-idx={idx}
                            onClick={() => goTo(result.slug)}
                            onMouseEnter={() => setActiveIndex(idx)}
                            className={`w-full text-left px-4 py-2.5 flex items-start gap-3 cursor-pointer transition-colors ${
                              isActive
                                ? "bg-accent/10"
                                : "hover:bg-white/[0.03]"
                            }`}
                          >
                            {/* Document icon */}
                            <svg
                              className={`w-5 h-5 mt-0.5 shrink-0 ${
                                isActive ? "text-accent" : "text-gray-4"
                              }`}
                              xmlns="http://www.w3.org/2000/svg"
                              fill="none"
                              viewBox="0 0 24 24"
                              strokeWidth={1.5}
                              stroke="currentColor"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m2.25 0H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"
                              />
                            </svg>
                            <div className="min-w-0 flex-1">
                              <div
                                className={`text-sm font-medium truncate ${
                                  isActive ? "text-accent" : "text-white"
                                }`}
                              >
                                {highlightMatch(result.title, query)}
                              </div>
                              {result.description && (
                                <div className="text-xs text-gray-3 truncate mt-0.5">
                                  {highlightMatch(result.description, query)}
                                </div>
                              )}
                              <div className="text-xs text-gray-4 mt-1 line-clamp-1">
                                {highlightMatch(result.snippet, query)}
                              </div>
                            </div>
                            {/* Arrow indicator for active */}
                            {isActive && (
                              <svg
                                className="w-4 h-4 text-accent mt-1 shrink-0"
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                                strokeWidth={2}
                                stroke="currentColor"
                              >
                                <path
                                  strokeLinecap="round"
                                  strokeLinejoin="round"
                                  d="M4.5 12h15m0 0l-6.75-6.75M19.5 12l-6.75 6.75"
                                />
                              </svg>
                            )}
                          </button>
                        );
                      })}
                    </div>
                  )
                )}
              </div>
            )}
          </div>

          {/* Footer with keyboard hints */}
          <div className="flex items-center gap-4 px-4 py-2.5 border-t border-gray-5 text-[11px] text-gray-4">
            <span className="flex items-center gap-1">
              <kbd className="bg-gray-6 border border-gray-5 rounded px-1 py-px font-mono">
                &uarr;
              </kbd>
              <kbd className="bg-gray-6 border border-gray-5 rounded px-1 py-px font-mono">
                &darr;
              </kbd>
              <span className="ml-0.5">to navigate</span>
            </span>
            <span className="flex items-center gap-1">
              <kbd className="bg-gray-6 border border-gray-5 rounded px-1 py-px font-mono">
                &crarr;
              </kbd>
              <span className="ml-0.5">to select</span>
            </span>
            <span className="flex items-center gap-1">
              <kbd className="bg-gray-6 border border-gray-5 rounded px-1 py-px font-mono">
                Esc
              </kbd>
              <span className="ml-0.5">to close</span>
            </span>
          </div>
        </div>
      </div>
    </>
  );
}
