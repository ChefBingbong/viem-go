# Benchmark Comparison: viem-go vs viem TypeScript

Generated: 2026-02-05T17:16:32.678Z

## Overall Summary

**🏆 Go is 1.52x faster overall**

| Metric | Go | TypeScript |
|--------|----|-----------|
| Avg ns/op | 2,772,653 | 4,216,390 |
| Avg ops/s | 361 | 237 |
| Wins | 12/12 | 0/12 |

## Detailed Results

| Benchmark | Go (ns/op) | TS (ns/op) | Go (ops/s) | TS (ops/s) | Result |
|-----------|------------|------------|------------|------------|--------|
| Multicall_Basic | 221,674 | 438,153 | 4,511 | 2,282 | 🟢 Go 1.98x faster |
| Multicall_WithArgs | 231,625 | 1,018,569 | 4,317 | 982 | 🟢 Go 4.40x faster |
| Multicall_MultiContract | 260,821 | 598,914 | 3,834 | 1,670 | 🟢 Go 2.30x faster |
| Multicall_10Calls | 283,913 | 545,664 | 3,522 | 1,833 | 🟢 Go 1.92x faster |
| Multicall_30Calls | 541,499 | 1,087,441 | 1,847 | 920 | 🟢 Go 2.01x faster |
| Multicall_Deployless | 393,360 | 594,817 | 2,542 | 1,681 | 🟢 Go 1.51x faster |
| Multicall_TokenMetadata | 248,060 | 426,652 | 4,031 | 2,344 | 🟢 Go 1.72x faster |
| Multicall_50Calls | 796,671 | 1,475,536 | 1,255 | 678 | 🟢 Go 1.85x faster |
| Multicall_100Calls | 1,693,048 | 2,581,445 | 591 | 387 | 🟢 Go 1.52x faster |
| Multicall_200Calls | 3,351,026 | 4,961,794 | 298 | 202 | 🟢 Go 1.48x faster |
| Multicall_500Calls | 8,326,050 | 12,061,431 | 120 | 83 | 🟢 Go 1.45x faster |
| Multicall_1000Calls | 16,924,090 | 24,806,263 | 59 | 40 | 🟢 Go 1.47x faster |

## Win Summary

- 🟢 Go wins: 12 (100%)
- 🔵 TS wins: 0 (0%)
- ⚪ Ties: 0 (0%)

## Notes

- Benchmarks run against the same Anvil instance (mainnet fork) for fair comparison
- ns/op = nanoseconds per operation (lower is better)
- ops/s = operations per second (higher is better)
- 🟢 = Go faster, 🔵 = TS faster, ⚪ = Similar (within 5%)
