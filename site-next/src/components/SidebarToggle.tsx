"use client";

import { useSidebar } from "./SidebarContext";

export default function SidebarToggle() {
  const { open, toggle, hasSidebar } = useSidebar();

  // Don't render on pages without a sidebar (e.g. homepage)
  if (!hasSidebar) return null;

  return (
    <button
      onClick={toggle}
      className="lg:hidden flex items-center justify-center w-9 h-9 rounded-md text-gray-3 hover:text-white hover:bg-white/[0.08] transition-colors shrink-0"
      aria-label={open ? "Close sidebar" : "Open sidebar"}
      aria-expanded={open}
    >
      {open ? (
        /* X icon */
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <line x1="18" y1="6" x2="6" y2="18" />
          <line x1="6" y1="6" x2="18" y2="18" />
        </svg>
      ) : (
        /* Hamburger icon */
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        >
          <line x1="3" y1="6" x2="21" y2="6" />
          <line x1="3" y1="12" x2="21" y2="12" />
          <line x1="3" y1="18" x2="21" y2="18" />
        </svg>
      )}
    </button>
  );
}
