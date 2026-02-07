# Benchmark Comparison: viem-go vs viem TypeScript

Generated: 2026-02-07T08:02:28.983Z

## Overall Summary

**🏆 Go is 7.12x faster overall** (geometric mean)

| Metric | Go | TypeScript |
|--------|----|-----------|
| Geometric mean speedup | 7.12x | - |
| Avg ns/op | 2,020,982 | 10,567,395 |
| Avg ops/s | 495 | 95 |
| Wins | 59/59 | 0/59 |

## By Suite

| Suite | Benchmarks | Go Wins | TS Wins | Ties | Winner |
|-------|------------|---------|---------|------|--------|
| abi | 6 | 6 | 0 | 0 | 🟢 Go 7.25x faster |
| address | 5 | 5 | 0 | 0 | 🟢 Go 6.27x faster |
| call | 6 | 6 | 0 | 0 | 🟢 Go 45.87x faster |
| ens | 5 | 5 | 0 | 0 | 🟢 Go 7.09x faster |
| event | 3 | 3 | 0 | 0 | 🟢 Go 21.82x faster |
| hash | 7 | 7 | 0 | 0 | 🟢 Go 10.57x faster |
| multicall | 16 | 16 | 0 | 0 | 🟢 Go 2.87x faster |
| signature | 5 | 5 | 0 | 0 | 🟢 Go 18.74x faster |
| unit | 6 | 6 | 0 | 0 | 🟢 Go 2.20x faster |

## Detailed Results

| Benchmark | Go (ns/op) | TS (ns/op) | Go (ops/s) | TS (ops/s) | Result |
|-----------|------------|------------|------------|------------|--------|
| Abi_EncodeSimple | 215 | 6,827 | 4,659,832 | 146,473 | 🟢 Go 31.81x faster |
| Abi_EncodeComplex | 544 | 7,677 | 1,838,235 | 130,262 | 🟢 Go 14.11x faster |
| Abi_EncodeMultiArg | 685 | 8,318 | 1,459,002 | 120,224 | 🟢 Go 12.14x faster |
| Abi_DecodeResult | 115 | 977 | 8,733,624 | 1,023,132 | 🟢 Go 8.54x faster |
| Abi_EncodePacked | 523 | 670 | 1,912,046 | 1,492,562 | 🟢 Go 1.28x faster |
| Abi_EncodePackedMulti | 483 | 1,181 | 2,072,109 | 846,977 | 🟢 Go 2.45x faster |
| Address_IsAddress | 38 | 233 | 26,116,479 | 4,293,464 | 🟢 Go 6.08x faster |
| Address_IsAddressLower | 59 | 239 | 17,044,486 | 4,178,762 | 🟢 Go 4.08x faster |
| Address_Checksum | 99 | 514 | 10,112,246 | 1,945,403 | 🟢 Go 5.20x faster |
| Address_Create | 1,167 | 7,163 | 856,898 | 139,611 | 🟢 Go 6.14x faster |
| Address_Create2 | 1,079 | 13,226 | 926,784 | 75,610 | 🟢 Go 12.26x faster |
| Call_Basic | 274,381 | 18,601,087 | 3,645 | 54 | 🟢 Go 67.79x faster |
| Call_WithData | 187,542 | 18,737,540 | 5,332 | 53 | 🟢 Go 99.91x faster |
| Call_WithAccount | 175,643 | 239,669 | 5,693 | 4,172 | 🟢 Go 1.36x faster |
| Call_Decimals | 174,224 | 17,801,545 | 5,740 | 56 | 🟢 Go 102.18x faster |
| Call_Symbol | 180,635 | 18,073,803 | 5,536 | 55 | 🟢 Go 100.06x faster |
| Call_BalanceOfMultiple | 183,921 | 18,122,640 | 5,437 | 55 | 🟢 Go 98.53x faster |
| Ens_Namehash | 1,661 | 20,762 | 602,047 | 48,165 | 🟢 Go 12.50x faster |
| Ens_NamehashDeep | 3,029 | 43,586 | 330,142 | 22,943 | 🟢 Go 14.39x faster |
| Ens_Labelhash | 434 | 5,399 | 2,304,678 | 185,214 | 🟢 Go 12.44x faster |
| Ens_Normalize | 349 | 966 | 2,864,509 | 1,034,895 | 🟢 Go 2.77x faster |
| Ens_NormalizeLong | 883 | 2,553 | 1,132,888 | 391,720 | 🟢 Go 2.89x faster |
| Event_DecodeTransfer | 639 | 10,210 | 1,564,456 | 97,939 | 🟢 Go 15.97x faster |
| Event_DecodeBatch10 | 3,982 | 103,642 | 251,130 | 9,649 | 🟢 Go 26.03x faster |
| Event_DecodeBatch100 | 39,738 | 992,822 | 25,165 | 1,007 | 🟢 Go 24.98x faster |
| Hash_Keccak256Short | 435 | 20,174 | 2,301,496 | 49,568 | 🟢 Go 46.43x faster |
| Hash_Keccak256Long | 2,812 | 48,695 | 355,619 | 20,536 | 🟢 Go 17.32x faster |
| Hash_Keccak256Hex | 460 | 5,264 | 2,175,332 | 189,971 | 🟢 Go 11.45x faster |
| Hash_Sha256Short | 161 | 1,428 | 6,211,180 | 700,346 | 🟢 Go 8.87x faster |
| Hash_Sha256Long | 632 | 12,863 | 1,582,028 | 77,740 | 🟢 Go 20.35x faster |
| Hash_FunctionSelector | 1,920 | 6,319 | 520,833 | 158,242 | 🟢 Go 3.29x faster |
| Hash_EventSelector | 2,377 | 6,423 | 420,698 | 155,691 | 🟢 Go 2.70x faster |
| Multicall_Basic | 182,793 | 459,455 | 5,471 | 2,176 | 🟢 Go 2.51x faster |
| Multicall_WithArgs | 189,193 | 364,049 | 5,286 | 2,747 | 🟢 Go 1.92x faster |
| Multicall_MultiContract | 220,815 | 374,585 | 4,529 | 2,670 | 🟢 Go 1.70x faster |
| Multicall_10Calls | 227,875 | 463,719 | 4,388 | 2,156 | 🟢 Go 2.03x faster |
| Multicall_30Calls | 377,187 | 914,261 | 2,651 | 1,094 | 🟢 Go 2.42x faster |
| Multicall_Deployless | 352,429 | 558,138 | 2,837 | 1,792 | 🟢 Go 1.58x faster |
| Multicall_TokenMetadata | 202,770 | 372,581 | 4,932 | 2,684 | 🟢 Go 1.84x faster |
| Multicall_50Calls | 484,157 | 1,253,997 | 2,065 | 797 | 🟢 Go 2.59x faster |
| Multicall_100Calls | 869,667 | 2,150,954 | 1,150 | 465 | 🟢 Go 2.47x faster |
| Multicall_200Calls | 1,410,831 | 4,344,237 | 709 | 230 | 🟢 Go 3.08x faster |
| Multicall_500Calls | 2,172,946 | 8,864,462 | 460 | 113 | 🟢 Go 4.08x faster |
| Multicall_MixedContracts_100 | 843,830 | 2,149,151 | 1,185 | 465 | 🟢 Go 2.55x faster |
| Multicall_1000Calls | 2,909,357 | 16,605,503 | 344 | 60 | 🟢 Go 5.71x faster |
| Multicall_10000Calls_SingleRPC | 66,612,371 | 162,662,459 | 15 | 6 | 🟢 Go 2.44x faster |
| Multicall_10000Calls_Chunked | 18,517,578 | 164,741,932 | 54 | 6 | 🟢 Go 8.90x faster |
| Multicall_10000Calls_AggressiveChunking | 22,368,029 | 161,155,160 | 45 | 6 | 🟢 Go 7.20x faster |
| Signature_HashMessage | 772 | 6,733 | 1,295,672 | 148,522 | 🟢 Go 8.72x faster |
| Signature_HashMessageLong | 1,809 | 14,494 | 552,792 | 68,996 | 🟢 Go 8.01x faster |
| Signature_RecoverAddress | 25,877 | 1,569,686 | 38,644 | 637 | 🟢 Go 60.66x faster |
| Signature_VerifyMessage | 26,062 | 1,543,377 | 38,370 | 648 | 🟢 Go 59.22x faster |
| Signature_ParseSignature | 189 | 1,737 | 5,293,806 | 575,671 | 🟢 Go 9.20x faster |
| Unit_ParseEther | 65 | 315 | 15,396,459 | 3,170,982 | 🟢 Go 4.86x faster |
| Unit_ParseEtherLarge | 130 | 235 | 7,680,492 | 4,252,017 | 🟢 Go 1.81x faster |
| Unit_FormatEther | 116 | 145 | 8,605,852 | 6,915,772 | 🟢 Go 1.24x faster |
| Unit_ParseUnits6 | 89 | 211 | 11,196,954 | 4,737,299 | 🟢 Go 2.36x faster |
| Unit_ParseGwei | 67 | 201 | 14,992,504 | 4,982,042 | 🟢 Go 3.01x faster |
| Unit_FormatUnits | 97 | 141 | 10,351,967 | 7,104,329 | 🟢 Go 1.46x faster |

## Win Summary

- 🟢 Go wins: 59 (100%)
- 🔵 TS wins: 0 (0%)
- ⚪ Ties: 0 (0%)

## Notes

- Benchmarks run against the same Anvil instance (mainnet fork) for fair comparison
- ns/op = nanoseconds per operation (lower is better)
- ops/s = operations per second (higher is better)
- 🟢 = Go faster, 🔵 = TS faster, ⚪ = Similar (within 5%)
