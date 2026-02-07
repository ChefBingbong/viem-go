const cards = [
  {
    title: "Backend Services",
    description:
      "Build high-performance blockchain indexers, APIs, and microservices",
  },
  {
    title: "CLI Tools",
    description:
      "Create command-line utilities for wallet management and contract interactions",
  },
  {
    title: "Infrastructure",
    description:
      "Power validators, relayers, and other blockchain infrastructure",
  },
  {
    title: "Cross-platform",
    description:
      "Compile to a single binary for any OS without dependencies",
  },
];

export default function WhyGoSection() {
  return (
    <section className="w-full py-8 px-8 bg-gradient-to-b from-dark-deep/30 to-transparent">
      <div className="max-w-[1120px] mx-auto mb-10 text-center">
        <h2 className="text-[2rem] font-semibold text-white mb-3">Why Go?</h2>
        <p className="text-lg text-gray-3 max-w-[600px] mx-auto">
          viem-go is designed for teams and projects that need Ethereum tooling
          in the Go ecosystem and have familiarity with the viem typescript
          library.
        </p>
      </div>
      <div className="grid grid-cols-4 gap-5 max-w-[1120px] mx-auto max-lg:grid-cols-2 max-sm:grid-cols-1">
        {cards.map((card) => (
          <div
            key={card.title}
            className="bg-gray-6/50 border border-accent/20 rounded-xl p-6 transition-all duration-200 hover:border-accent/40 hover:bg-gray-6/70 hover:-translate-y-0.5"
          >
            <h3 className="text-lg font-semibold text-white mb-2">
              {card.title}
            </h3>
            <p className="text-[0.9375rem] text-gray-3 leading-normal">
              {card.description}
            </p>
          </div>
        ))}
      </div>
    </section>
  );
}
