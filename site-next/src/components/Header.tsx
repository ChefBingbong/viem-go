import Link from "next/link";
import SearchTrigger from "./SearchTrigger";
import SidebarToggle from "./SidebarToggle";

export default function Header() {
  return (
    <header className="sticky top-0 z-50 w-full border-b border-accent/10 bg-gray-6/80 backdrop-blur-md">
      <div className="max-w-[1400px] mx-auto px-4 sm:px-6 h-16 flex items-center justify-between gap-3">
        {/* Left: hamburger + logo + search */}
        <div className="flex items-center gap-3 sm:gap-5 min-w-0">
          <SidebarToggle />
          <Link href="/" className="flex items-center gap-2 shrink-0">
            <img
              src="/svg/golem-logo-full-light.svg"
              alt="viem-go"
              className="h-[6.25rem] max-sm:h-12"
            />
          </Link>
          <div className="hidden sm:block ml-10">
            <SearchTrigger />
          </div>
        </div>

        {/* Right: nav links */}
        <nav className="flex items-center gap-1 shrink-0">
          <Link
            href="/docs/introduction"
            className="hidden sm:block text-sm font-medium text-gray-2 no-underline px-3 py-2 rounded-md hover:text-white hover:bg-white/5 transition-colors"
          >
            Docs
          </Link>
          <a
            href="https://github.com/ChefBingbong/viem-go"
            target="_blank"
            rel="noopener noreferrer"
            className="hidden sm:block text-sm font-medium text-gray-2 no-underline px-3 py-2 rounded-md hover:text-white hover:bg-white/5 transition-colors"
          >
            GitHub
          </a>
          {/* Mobile search (icon only) */}
          <div className="sm:hidden">
            <SearchTrigger compact />
          </div>
          {/* Version dropdown */}
          <VersionDropdown />
        </nav>
      </div>
    </header>
  );
}

function VersionDropdown() {
  return (
    <div className="relative group">
      <button className="flex items-center gap-1 text-sm font-medium text-gray-2 bg-transparent border border-gray-5 px-3 py-1.5 rounded-md cursor-pointer hover:text-white hover:border-gray-4 hover:bg-white/5 transition-colors">
        v0.1.0
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="12"
          height="12"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          className="transition-transform group-hover:rotate-180"
        >
          <polyline points="6 9 12 15 18 9" />
        </svg>
      </button>
      <div className="absolute top-full right-0 mt-2 min-w-[160px] bg-gray-6 border border-gray-5 rounded-lg p-1 opacity-0 invisible -translate-y-1 transition-all group-hover:opacity-100 group-hover:visible group-hover:translate-y-0 z-50">
        <a
          href="https://github.com/ChefBingbong/viem-go/releases"
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center justify-between gap-2 text-gray-2 no-underline text-sm px-3 py-2 rounded hover:text-white hover:bg-white/[0.08] transition-colors"
        >
          Releases
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="opacity-50"
          >
            <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
            <polyline points="15 3 21 3 21 9" />
            <line x1="10" y1="14" x2="21" y2="3" />
          </svg>
        </a>
        <a
          href="https://github.com/ChefBingbong/viem-go/tree/main/examples"
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center justify-between gap-2 text-gray-2 no-underline text-sm px-3 py-2 rounded hover:text-white hover:bg-white/[0.08] transition-colors"
        >
          Examples
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="opacity-50"
          >
            <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
            <polyline points="15 3 21 3 21 9" />
            <line x1="10" y1="14" x2="21" y2="3" />
          </svg>
        </a>
        <a
          href="https://github.com/ChefBingbong/viem-go/blob/main/.github/CONTRIBUTING.md"
          target="_blank"
          rel="noopener noreferrer"
          className="flex items-center justify-between gap-2 text-gray-2 no-underline text-sm px-3 py-2 rounded hover:text-white hover:bg-white/[0.08] transition-colors"
        >
          Contributing
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="12"
            height="12"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            className="opacity-50"
          >
            <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6" />
            <polyline points="15 3 21 3 21 9" />
            <line x1="10" y1="14" x2="21" y2="3" />
          </svg>
        </a>
      </div>
    </div>
  );
}
