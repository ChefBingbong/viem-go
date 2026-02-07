export interface NavItem {
  label: string;
  slug?: string;
  items?: NavItem[];
}

export const docsNav: NavItem[] = [
  {
    label: "Introduction",
    items: [
      { label: "Why viem-go", slug: "introduction" },
      { label: "Getting Started", slug: "getting-started" },
      { label: "Examples", slug: "examples" },
    ],
  },
  {
    label: "Clients",
    items: [
      { label: "Introduction", slug: "clients/intro" },
      { label: "Public Client", slug: "clients/public" },
      { label: "Wallet Client", slug: "clients/wallet" },
      {
        label: "Transports",
        items: [
          { label: "HTTP", slug: "clients/transports/http" },
          { label: "WebSocket", slug: "clients/transports/websocket" },
        ],
      },
    ],
  },
  {
    label: "Accounts",
    items: [
      { label: "Overview", slug: "accounts/overview" },
      { label: "Private Key", slug: "accounts/private-key" },
      { label: "Mnemonic", slug: "accounts/mnemonic" },
      { label: "HD Wallet", slug: "accounts/hd-wallet" },
    ],
  },
  {
    label: "Contract",
    items: [
      { label: "Reading Contracts", slug: "contract/read-contract" },
      { label: "Writing Contracts", slug: "contract/write-contract" },
      { label: "ABI Encoding", slug: "contract/abi-encoding" },
    ],
  },
  {
    label: "Utilities",
    items: [
      { label: "Units", slug: "utilities/units" },
      { label: "Hashing", slug: "utilities/hashing" },
      { label: "Signatures", slug: "utilities/signatures" },
      { label: "Addresses", slug: "utilities/addresses" },
    ],
  },
];
