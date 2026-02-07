"use client";

import { useEffect, useRef, useState } from "react";
import type { TocEntry } from "@/lib/mdx";

export default function TableOfContents({
  headings,
}: {
  headings: TocEntry[];
}) {
  const [activeId, setActiveId] = useState<string>("");
  const observerRef = useRef<IntersectionObserver | null>(null);

  useEffect(() => {
    // Find all heading elements on the page
    const elements = headings
      .map((h) => document.getElementById(h.id))
      .filter(Boolean) as HTMLElement[];

    if (elements.length === 0) return;

    // Use IntersectionObserver to track which heading is in view
    observerRef.current = new IntersectionObserver(
      (entries) => {
        // Find the first heading that is intersecting
        for (const entry of entries) {
          if (entry.isIntersecting) {
            setActiveId(entry.target.id);
            break;
          }
        }
      },
      {
        rootMargin: "-80px 0px -70% 0px",
        threshold: 0,
      }
    );

    for (const el of elements) {
      observerRef.current.observe(el);
    }

    return () => observerRef.current?.disconnect();
  }, [headings]);

  // Also track scroll position for more accurate active state
  useEffect(() => {
    function onScroll() {
      const scrollY = window.scrollY + 120;
      let current = "";
      for (const h of headings) {
        const el = document.getElementById(h.id);
        if (el && el.offsetTop <= scrollY) {
          current = h.id;
        }
      }
      if (current) setActiveId(current);
    }
    window.addEventListener("scroll", onScroll, { passive: true });
    // Set initial
    onScroll();
    return () => window.removeEventListener("scroll", onScroll);
  }, [headings]);

  if (headings.length === 0) return null;

  return (
    <aside className="hidden xl:block w-[220px] shrink-0 h-[calc(100vh-4rem)] sticky top-16 border-l border-accent/10">
      <div className="py-8 pl-4 pr-2 overflow-y-auto h-full">
        <p className="text-xs font-semibold text-gray-3 uppercase tracking-wider mb-3">
          On this page
        </p>
        <nav className="flex flex-col gap-0.5">
          {headings.map((heading) => {
            const isActive = activeId === heading.id;
            return (
              <a
                key={heading.id}
                href={`#${heading.id}`}
                onClick={(e) => {
                  e.preventDefault();
                  const el = document.getElementById(heading.id);
                  if (el) {
                    el.scrollIntoView({ behavior: "smooth" });
                    setActiveId(heading.id);
                    // Update URL hash without jumping
                    history.pushState(null, "", `#${heading.id}`);
                  }
                }}
                className={`block text-[13px] leading-snug no-underline py-1 transition-colors border-l-2 ${
                  heading.depth === 3 ? "pl-5" : "pl-3"
                } ${
                  isActive
                    ? "text-accent border-accent"
                    : "text-gray-4 border-transparent hover:text-gray-2 hover:border-gray-5"
                }`}
              >
                {heading.text}
              </a>
            );
          })}
        </nav>
      </div>
    </aside>
  );
}
