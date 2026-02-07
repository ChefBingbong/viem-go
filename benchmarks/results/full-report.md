# Full Benchmark Report: viem-go vs viem TypeScript

Generated: 2026-02-07T07:49:09.015Z

---

## Executive Summary

This report compares **59** benchmarks across **9** test suites.

### 🏆 Winner: Go (viem-go)

Go is **4.74x faster** on average across all benchmarks.

### Quick Stats

| Metric | Value |
|--------|-------|
| Total Benchmarks | 59 |
| Test Suites | 9 |
| Go Wins | 58 (98.3%) |
| TypeScript Wins | 0 (0.0%) |
| Ties | 1 (1.7%) |
| Avg Go Latency | 2.23 ms |
| Avg TS Latency | 10.57 ms |
| Go Throughput | 449 ops/s |
| TS Throughput | 95 ops/s |

---

## Suite-by-Suite Analysis

### Abi Suite

**Result:** 🟢 Go 13.97x faster

| Benchmark | Go | TS | Diff | Winner |
|-----------|----|----|------|--------|
| EncodeSimple | 255.1 ns | 6.83 µs | 26.76x | 🟢 |
| EncodeComplex | 323.2 ns | 7.68 µs | 23.75x | 🟢 |
| EncodeMultiArg | 453.6 ns | 8.32 µs | 18.34x | 🟢 |
| DecodeResult | 95.8 ns | 977.4 ns | 10.20x | 🟢 |
| EncodePacked | 327.6 ns | 670.0 ns | 2.05x | 🟢 |
| EncodePackedMulti | 380.7 ns | 1.18 µs | 3.10x | 🟢 |

**Suite Statistics:**
- Benchmarks: 6
- Go wins: 6, TS wins: 0, Ties: 0
- Avg Go: 306.0 ns | Avg TS: 4.27 µs

### Address Suite

**Result:** 🟢 Go 11.23x faster

| Benchmark | Go | TS | Diff | Winner |
|-----------|----|----|------|--------|
| IsAddress | 45.1 ns | 232.9 ns | 5.16x | 🟢 |
| IsAddressLower | 37.3 ns | 239.3 ns | 6.41x | 🟢 |
| Checksum | 75.3 ns | 514.0 ns | 6.83x | 🟢 |
| Create | 738.9 ns | 7.16 µs | 9.69x | 🟢 |
| Create2 | 1.01 µs | 13.23 µs | 13.15x | 🟢 |

**Suite Statistics:**
- Benchmarks: 5
- Go wins: 5, TS wins: 0, Ties: 0
- Avg Go: 380.5 ns | Avg TS: 4.27 µs

### Call Suite

**Result:** 🟢 Go 66.47x faster

| Benchmark | Go | TS | Diff | Winner |
|-----------|----|----|------|--------|
| Basic | 174.80 µs | 18.60 ms | 106.42x | 🟢 |
| WithData | 291.20 µs | 18.74 ms | 64.35x | 🟢 |
| WithAccount | 241.10 µs | 239.67 µs | 1.01x | ⚪ |
| Decimals | 179.53 µs | 17.80 ms | 99.16x | 🟢 |
| Symbol | 191.72 µs | 18.07 ms | 94.27x | 🟢 |
| BalanceOfMultiple | 299.41 µs | 18.12 ms | 60.53x | 🟢 |

**Suite Statistics:**
- Benchmarks: 6
- Go wins: 5, TS wins: 0, Ties: 1
- Avg Go: 229.62 µs | Avg TS: 15.26 ms

### Ens Suite

**Result:** 🟢 Go 11.50x faster

| Benchmark | Go | TS | Diff | Winner |
|-----------|----|----|------|--------|
| Namehash | 1.58 µs | 20.76 µs | 13.15x | 🟢 |
| NamehashDeep | 3.06 µs | 43.59 µs | 14.25x | 🟢 |
| Labelhash | 433.8 ns | 5.40 µs | 12.45x | 🟢 |
| Normalize | 417.1 ns | 966.3 ns | 2.32x | 🟢 |
| NormalizeLong | 882.6 ns | 2.55 µs | 2.89x | 🟢 |

**Suite Statistics:**
- Benchmarks: 5
- Go wins: 5, TS wins: 0, Ties: 0
- Avg Go: 1.27 µs | Avg TS: 14.65 µs

### Event Suite

**Result:** 🟢 Go 24.64x faster

| Benchmark | Go | TS | Diff | Winner |
|-----------|----|----|------|--------|
| DecodeTransfer | 681.5 ns | 10.21 µs | 14.98x | 🟢 |
| DecodeBatch10 | 4.15 µs | 103.64 µs | 24.95x | 🟢 |
| DecodeBatch100 | 40.08 µs | 992.82 µs | 24.77x | 🟢 |

**Suite Statistics:**
- Benchmarks: 3
- Go wins: 3, TS wins: 0, Ties: 0
- Avg Go: 14.97 µs | Avg TS: 368.89 µs

### Hash Suite

**Result:** 🟢 Go 11.66x faster

| Benchmark | Go | TS | Diff | Winner |
|-----------|----|----|------|--------|
| Keccak256Short | 428.7 ns | 20.17 µs | 47.06x | 🟢 |
| Keccak256Long | 2.69 µs | 48.70 µs | 18.14x | 🟢 |
| Keccak256Hex | 468.8 ns | 5.26 µs | 11.23x | 🟢 |
| Sha256Short | 168.4 ns | 1.43 µs | 8.48x | 🟢 |
| Sha256Long | 627.4 ns | 12.86 µs | 20.50x | 🟢 |
| FunctionSelector | 1.94 µs | 6.32 µs | 3.26x | 🟢 |
| EventSelector | 2.37 µs | 6.42 µs | 2.72x | 🟢 |

**Suite Statistics:**
- Benchmarks: 7
- Go wins: 7, TS wins: 0, Ties: 0
- Avg Go: 1.24 µs | Avg TS: 14.45 µs

### Multicall Suite

**Result:** 🟢 Go 4.06x faster

| Benchmark | Go | TS | Diff | Winner |
|-----------|----|----|------|--------|
| Basic | 182.93 µs | 459.46 µs | 2.51x | 🟢 |
| WithArgs | 181.15 µs | 364.05 µs | 2.01x | 🟢 |
| MultiContract | 267.24 µs | 374.59 µs | 1.40x | 🟢 |
| 10Calls | 421.58 µs | 463.72 µs | 1.10x | 🟢 |
| 30Calls | 398.02 µs | 914.26 µs | 2.30x | 🟢 |
| Deployless | 382.37 µs | 558.14 µs | 1.46x | 🟢 |
| TokenMetadata | 215.48 µs | 372.58 µs | 1.73x | 🟢 |
| 50Calls | 543.17 µs | 1.25 ms | 2.31x | 🟢 |
| 100Calls | 864.88 µs | 2.15 ms | 2.49x | 🟢 |
| 200Calls | 1.44 ms | 4.34 ms | 3.01x | 🟢 |
| 500Calls | 2.25 ms | 8.86 ms | 3.94x | 🟢 |
| MixedContracts_100 | 838.69 µs | 2.15 ms | 2.56x | 🟢 |
| 1000Calls | 4.32 ms | 16.61 ms | 3.84x | 🟢 |
| 10000Calls_SingleRPC | 74.43 ms | 162.66 ms | 2.19x | 🟢 |
| 10000Calls_Chunked | 19.56 ms | 164.74 ms | 8.42x | 🟢 |
| 10000Calls_AggressiveChunking | 23.67 ms | 161.16 ms | 6.81x | 🟢 |

**Suite Statistics:**
- Benchmarks: 16
- Go wins: 16, TS wins: 0, Ties: 0
- Avg Go: 8.12 ms | Avg TS: 32.96 ms

### Signature Suite

**Result:** 🟢 Go 57.32x faster

| Benchmark | Go | TS | Diff | Winner |
|-----------|----|----|------|--------|
| HashMessage | 767.8 ns | 6.73 µs | 8.77x | 🟢 |
| HashMessageLong | 1.84 µs | 14.49 µs | 7.88x | 🟢 |
| RecoverAddress | 25.97 µs | 1.57 ms | 60.44x | 🟢 |
| VerifyMessage | 25.95 µs | 1.54 ms | 59.48x | 🟢 |
| ParseSignature | 190.2 ns | 1.74 µs | 9.13x | 🟢 |

**Suite Statistics:**
- Benchmarks: 5
- Go wins: 5, TS wins: 0, Ties: 0
- Avg Go: 10.94 µs | Avg TS: 627.21 µs

### Unit Suite

**Result:** 🟢 Go 2.29x faster

| Benchmark | Go | TS | Diff | Winner |
|-----------|----|----|------|--------|
| ParseEther | 66.0 ns | 315.4 ns | 4.78x | 🟢 |
| ParseEtherLarge | 132.8 ns | 235.2 ns | 1.77x | 🟢 |
| FormatEther | 117.3 ns | 144.6 ns | 1.23x | 🟢 |
| ParseUnits6 | 60.7 ns | 211.1 ns | 3.48x | 🟢 |
| ParseGwei | 66.7 ns | 200.7 ns | 3.01x | 🟢 |
| FormatUnits | 100.4 ns | 140.8 ns | 1.40x | 🟢 |

**Suite Statistics:**
- Benchmarks: 6
- Go wins: 6, TS wins: 0, Ties: 0
- Avg Go: 90.7 ns | Avg TS: 208.0 ns

---

## Category Analysis

### Other

🟢 **Go 4.78x faster**

Benchmarks: 39 | Go wins: 39 | TS wins: 0 | Ties: 0

### Basic Operations

🟢 **Go 53.28x faster**

Benchmarks: 2 | Go wins: 2 | TS wins: 0 | Ties: 0

### With Parameters

🟢 **Go 28.31x faster**

Benchmarks: 3 | Go wins: 3 | TS wins: 0 | Ties: 0

### With Account

⚪ **Similar**

Benchmarks: 1 | Go wins: 0 | TS wins: 0 | Ties: 1

### Simple Reads

🟢 **Go 96.64x faster**

Benchmarks: 2 | Go wins: 2 | TS wins: 0 | Ties: 0

### Batch Operations

🟢 **Go 10.91x faster**

Benchmarks: 4 | Go wins: 4 | TS wins: 0 | Ties: 0

### Multi-Contract

🟢 **Go 2.28x faster**

Benchmarks: 2 | Go wins: 2 | TS wins: 0 | Ties: 0

### Deployless

🟢 **Go 1.46x faster**

Benchmarks: 1 | Go wins: 1 | TS wins: 0 | Ties: 0

### Extreme Stress Tests

🟢 **Go 4.14x faster**

Benchmarks: 5 | Go wins: 5 | TS wins: 0 | Ties: 0

---

## Memory Analysis (Go)

| Benchmark | Bytes/op | Allocs/op |
|-----------|----------|----------|
| Abi_EncodeSimple | 200 | 6 |
| Abi_EncodeComplex | 328 | 10 |
| Abi_EncodeMultiArg | 592 | 15 |
| Abi_DecodeResult | 112 | 3 |
| Abi_EncodePacked | 288 | 8 |
| Abi_EncodePackedMulti | 392 | 11 |
| Address_IsAddress | 48 | 1 |
| Address_IsAddressLower | 48 | 1 |
| Address_Checksum | 96 | 2 |
| Address_Create | 440 | 15 |
| Address_Create2 | 376 | 9 |
| Call_Basic | 28,019 | 236 |
| Call_WithData | 28,244 | 238 |
| Call_WithAccount | 28,491 | 242 |
| Call_Decimals | 27,626 | 236 |
| Call_Symbol | 28,278 | 236 |
| Call_BalanceOfMultiple | 28,503 | 238 |
| Ens_Namehash | 400 | 9 |
| Ens_NamehashDeep | 560 | 13 |
| Ens_Labelhash | 240 | 4 |
| Ens_Normalize | 256 | 11 |
| Ens_NormalizeLong | 696 | 21 |
| Event_DecodeTransfer | 720 | 9 |
| Event_DecodeBatch10 | 7,200 | 91 |
| Event_DecodeBatch100 | 72,001 | 912 |
| Hash_Keccak256Short | 240 | 4 |
| Hash_Keccak256Long | 240 | 4 |
| Hash_Keccak256Hex | 256 | 5 |
| Hash_Sha256Short | 240 | 4 |
| Hash_Sha256Long | 240 | 4 |
| Hash_FunctionSelector | 976 | 98 |
| Hash_EventSelector | 1,376 | 124 |
| Multicall_Basic | 26,404 | 139 |
| Multicall_WithArgs | 25,658 | 140 |
| Multicall_MultiContract | 34,144 | 154 |
| Multicall_10Calls | 68,879 | 212 |
| Multicall_30Calls | 210,995 | 400 |
| Multicall_Deployless | 133,737 | 177 |
| Multicall_TokenMetadata | 32,903 | 144 |
| Multicall_50Calls | 337,804 | 584 |
| Multicall_100Calls | 692,264 | 1,045 |
| Multicall_200Calls | 1,282,941 | 1,953 |
| Multicall_500Calls | 3,307,294 | 4,990 |
| Multicall_MixedContracts_100 | 699,509 | 1,046 |
| Multicall_1000Calls | 6,566,848 | 9,852 |
| Multicall_10000Calls_SingleRPC | 63,701,474 | 96,224 |
| Multicall_10000Calls_Chunked | 62,292,455 | 97,180 |
| Multicall_10000Calls_AggressiveChunking | 56,003,684 | 89,810 |
| Signature_HashMessage | 712 | 12 |
| Signature_HashMessageLong | 2,507 | 15 |
| Signature_RecoverAddress | 3,571 | 54 |
| Signature_VerifyMessage | 3,571 | 54 |
| Signature_ParseSignature | 288 | 7 |
| Unit_ParseEther | 64 | 3 |
| Unit_ParseEtherLarge | 96 | 4 |
| Unit_FormatEther | 56 | 3 |
| Unit_ParseUnits6 | 48 | 3 |
| Unit_ParseGwei | 56 | 3 |
| Unit_FormatUnits | 40 | 3 |

---

## Detailed Raw Data

| Benchmark | Suite | Go ns/op | TS ns/op | Go ops/s | TS ops/s | Ratio | Winner |
|-----------|-------|----------|----------|----------|----------|-------|--------|
| Abi_EncodeSimple | abi | 255 | 6,827 | 3,920,031 | 146,473 | 0.037 | 🟢 |
| Abi_EncodeComplex | abi | 323 | 7,677 | 3,094,059 | 130,262 | 0.042 | 🟢 |
| Abi_EncodeMultiArg | abi | 454 | 8,318 | 2,204,586 | 120,224 | 0.055 | 🟢 |
| Abi_DecodeResult | abi | 96 | 977 | 10,440,593 | 1,023,132 | 0.098 | 🟢 |
| Abi_EncodePacked | abi | 328 | 670 | 3,052,503 | 1,492,562 | 0.489 | 🟢 |
| Abi_EncodePackedMulti | abi | 381 | 1,181 | 2,626,740 | 846,977 | 0.322 | 🟢 |
| Address_IsAddress | address | 45 | 233 | 22,168,034 | 4,293,464 | 0.194 | 🟢 |
| Address_IsAddressLower | address | 37 | 239 | 26,780,932 | 4,178,762 | 0.156 | 🟢 |
| Address_Checksum | address | 75 | 514 | 13,285,506 | 1,945,403 | 0.146 | 🟢 |
| Address_Create | address | 739 | 7,163 | 1,353,363 | 139,611 | 0.103 | 🟢 |
| Address_Create2 | address | 1,006 | 13,226 | 994,036 | 75,610 | 0.076 | 🟢 |
| Call_Basic | call | 174,797 | 18,601,087 | 5,721 | 54 | 0.009 | 🟢 |
| Call_WithData | call | 291,195 | 18,737,540 | 3,434 | 53 | 0.016 | 🟢 |
| Call_WithAccount | call | 241,097 | 239,669 | 4,148 | 4,172 | 1.006 | ⚪ |
| Call_Decimals | call | 179,525 | 17,801,545 | 5,570 | 56 | 0.010 | 🟢 |
| Call_Symbol | call | 191,719 | 18,073,803 | 5,216 | 55 | 0.011 | 🟢 |
| Call_BalanceOfMultiple | call | 299,406 | 18,122,640 | 3,340 | 55 | 0.017 | 🟢 |
| Ens_Namehash | ens | 1,579 | 20,762 | 633,312 | 48,165 | 0.076 | 🟢 |
| Ens_NamehashDeep | ens | 3,059 | 43,586 | 326,904 | 22,943 | 0.070 | 🟢 |
| Ens_Labelhash | ens | 434 | 5,399 | 2,305,210 | 185,214 | 0.080 | 🟢 |
| Ens_Normalize | ens | 417 | 966 | 2,397,507 | 1,034,895 | 0.432 | 🟢 |
| Ens_NormalizeLong | ens | 883 | 2,553 | 1,133,016 | 391,720 | 0.346 | 🟢 |
| Event_DecodeTransfer | event | 682 | 10,210 | 1,467,351 | 97,939 | 0.067 | 🟢 |
| Event_DecodeBatch10 | event | 4,154 | 103,642 | 240,732 | 9,649 | 0.040 | 🟢 |
| Event_DecodeBatch100 | event | 40,076 | 992,822 | 24,953 | 1,007 | 0.040 | 🟢 |
| Hash_Keccak256Short | hash | 429 | 20,174 | 2,332,634 | 49,568 | 0.021 | 🟢 |
| Hash_Keccak256Long | hash | 2,685 | 48,695 | 372,439 | 20,536 | 0.055 | 🟢 |
| Hash_Keccak256Hex | hash | 469 | 5,264 | 2,133,106 | 189,971 | 0.089 | 🟢 |
| Hash_Sha256Short | hash | 168 | 1,428 | 5,938,242 | 700,346 | 0.118 | 🟢 |
| Hash_Sha256Long | hash | 627 | 12,863 | 1,593,880 | 77,740 | 0.049 | 🟢 |
| Hash_FunctionSelector | hash | 1,936 | 6,319 | 516,529 | 158,242 | 0.306 | 🟢 |
| Hash_EventSelector | hash | 2,365 | 6,423 | 422,833 | 155,691 | 0.368 | 🟢 |
| Multicall_Basic | multicall | 182,932 | 459,455 | 5,467 | 2,176 | 0.398 | 🟢 |
| Multicall_WithArgs | multicall | 181,152 | 364,049 | 5,520 | 2,747 | 0.498 | 🟢 |
| Multicall_MultiContract | multicall | 267,241 | 374,585 | 3,742 | 2,670 | 0.713 | 🟢 |
| Multicall_10Calls | multicall | 421,579 | 463,719 | 2,372 | 2,156 | 0.909 | 🟢 |
| Multicall_30Calls | multicall | 398,017 | 914,261 | 2,512 | 1,094 | 0.435 | 🟢 |
| Multicall_Deployless | multicall | 382,374 | 558,138 | 2,615 | 1,792 | 0.685 | 🟢 |
| Multicall_TokenMetadata | multicall | 215,481 | 372,581 | 4,641 | 2,684 | 0.578 | 🟢 |
| Multicall_50Calls | multicall | 543,168 | 1,253,997 | 1,841 | 797 | 0.433 | 🟢 |
| Multicall_100Calls | multicall | 864,877 | 2,150,954 | 1,156 | 465 | 0.402 | 🟢 |
| Multicall_200Calls | multicall | 1,441,860 | 4,344,237 | 694 | 230 | 0.332 | 🟢 |
| Multicall_500Calls | multicall | 2,249,852 | 8,864,462 | 444 | 113 | 0.254 | 🟢 |
| Multicall_MixedContracts_100 | multicall | 838,685 | 2,149,151 | 1,192 | 465 | 0.390 | 🟢 |
| Multicall_1000Calls | multicall | 4,319,635 | 16,605,503 | 232 | 60 | 0.260 | 🟢 |
| Multicall_10000Calls_SingleRPC | multicall | 74,434,656 | 162,662,459 | 13 | 6 | 0.458 | 🟢 |
| Multicall_10000Calls_Chunked | multicall | 19,559,964 | 164,741,932 | 51 | 6 | 0.119 | 🟢 |
| Multicall_10000Calls_AggressiveChunking | multicall | 23,668,209 | 161,155,160 | 42 | 6 | 0.147 | 🟢 |
| Signature_HashMessage | signature | 768 | 6,733 | 1,302,423 | 148,522 | 0.114 | 🟢 |
| Signature_HashMessageLong | signature | 1,839 | 14,494 | 543,774 | 68,996 | 0.127 | 🟢 |
| Signature_RecoverAddress | signature | 25,970 | 1,569,686 | 38,506 | 637 | 0.017 | 🟢 |
| Signature_VerifyMessage | signature | 25,948 | 1,543,377 | 38,539 | 648 | 0.017 | 🟢 |
| Signature_ParseSignature | signature | 190 | 1,737 | 5,257,624 | 575,671 | 0.109 | 🟢 |
| Unit_ParseEther | unit | 66 | 315 | 15,149,220 | 3,170,982 | 0.209 | 🟢 |
| Unit_ParseEtherLarge | unit | 133 | 235 | 7,530,120 | 4,252,017 | 0.565 | 🟢 |
| Unit_FormatEther | unit | 117 | 145 | 8,525,149 | 6,915,772 | 0.811 | 🟢 |
| Unit_ParseUnits6 | unit | 61 | 211 | 16,466,326 | 4,737,299 | 0.288 | 🟢 |
| Unit_ParseGwei | unit | 67 | 201 | 14,997,001 | 4,982,042 | 0.332 | 🟢 |
| Unit_FormatUnits | unit | 100 | 141 | 9,960,159 | 7,104,329 | 0.713 | 🟢 |

---

## Methodology

### Test Environment

- **Network:** Anvil (Mainnet fork)
- **Go Benchmark:** `go test -bench=. -benchmem -benchtime=10s -count=5`
- **TS Benchmark:** `vitest bench` with 10s per benchmark

### Measurement Notes

- **ns/op:** Nanoseconds per operation (lower is better)
- **ops/s:** Operations per second (higher is better)
- **Ratio:** Go time / TS time (>1 means TS is faster)
- **Tie:** Within 5% of each other

### Caveats

- Network latency dominates most benchmarks (RPC calls)
- Results may vary based on network conditions
- CPU-bound operations may show different characteristics
