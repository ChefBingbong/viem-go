# Benchmark Comparison: viem-go vs viem TypeScript

Generated: 2026-02-07T05:46:12.170Z

## Overall Summary

**🏆 Go is 4.53x faster overall**

| Metric | Go | TypeScript |
|--------|----|-----------|
| Avg ns/op | 3,062,791 | 13,861,286 |
| Avg ops/s | 326 | 72 |
| Wins | 55/59 | 3/59 |

## By Suite

| Suite | Benchmarks | Go Wins | TS Wins | Ties | Winner |
|-------|------------|---------|---------|------|--------|
| abi | 6 | 6 | 0 | 0 | 🟢 Go 17.34x faster |
| address | 5 | 2 | 2 | 1 | 🟢 Go 4.09x faster |
| call | 6 | 6 | 0 | 0 | 🟢 Go 87.21x faster |
| ens | 5 | 5 | 0 | 0 | 🟢 Go 14.77x faster |
| event | 3 | 3 | 0 | 0 | 🟢 Go 28.31x faster |
| hash | 7 | 7 | 0 | 0 | 🟢 Go 12.18x faster |
| multicall | 16 | 16 | 0 | 0 | 🟢 Go 4.01x faster |
| signature | 5 | 5 | 0 | 0 | 🟢 Go 59.31x faster |
| unit | 6 | 5 | 1 | 0 | 🟢 Go 1.37x faster |

## Detailed Results

| Benchmark | Go (ns/op) | TS (ns/op) | Go (ops/s) | TS (ops/s) | Result |
|-----------|------------|------------|------------|------------|--------|
| Abi_EncodeSimple | 215 | 8,644 | 4,649,000 | 115,685 | 🟢 Go 40.19x faster |
| Abi_EncodeComplex | 325 | 9,533 | 3,079,766 | 104,895 | 🟢 Go 29.36x faster |
| Abi_EncodeMultiArg | 455 | 10,030 | 2,199,252 | 99,700 | 🟢 Go 22.06x faster |
| Abi_DecodeResult | 93 | 1,073 | 10,697,475 | 931,629 | 🟢 Go 11.48x faster |
| Abi_EncodePacked | 334 | 703 | 2,994,012 | 1,422,858 | 🟢 Go 2.10x faster |
| Abi_EncodePackedMulti | 382 | 1,281 | 2,621,232 | 780,533 | 🟢 Go 3.36x faster |
| Address_IsAddress | 910 | 293 | 1,099,143 | 3,418,491 | 🔵 TS 3.11x faster |
| Address_IsAddressLower | 304 | 295 | 3,286,231 | 3,395,274 | ⚪ Similar |
| Address_Checksum | 812 | 722 | 1,231,527 | 1,384,784 | 🔵 TS 1.12x faster |
| Address_Create | 2,346 | 9,761 | 426,257 | 102,451 | 🟢 Go 4.16x faster |
| Address_Create2 | 2,591 | 17,410 | 385,951 | 57,437 | 🟢 Go 6.72x faster |
| Call_Basic | 176,772 | 19,322,365 | 5,657 | 52 | 🟢 Go 109.31x faster |
| Call_WithData | 171,302 | 18,632,279 | 5,838 | 54 | 🟢 Go 108.77x faster |
| Call_WithAccount | 175,134 | 289,804 | 5,710 | 3,451 | 🟢 Go 1.65x faster |
| Call_Decimals | 184,363 | 17,701,183 | 5,424 | 56 | 🟢 Go 96.01x faster |
| Call_Symbol | 180,229 | 17,923,138 | 5,548 | 56 | 🟢 Go 99.45x faster |
| Call_BalanceOfMultiple | 169,866 | 18,375,462 | 5,887 | 54 | 🟢 Go 108.18x faster |
| Ens_Namehash | 1,600 | 28,027 | 625,000 | 35,680 | 🟢 Go 17.52x faster |
| Ens_NamehashDeep | 3,067 | 55,594 | 326,052 | 17,987 | 🟢 Go 18.13x faster |
| Ens_Labelhash | 446 | 7,205 | 2,243,662 | 138,784 | 🟢 Go 16.17x faster |
| Ens_Normalize | 354 | 1,000 | 2,823,264 | 1,000,321 | 🟢 Go 2.82x faster |
| Ens_NormalizeLong | 930 | 2,663 | 1,075,731 | 375,460 | 🟢 Go 2.87x faster |
| Event_DecodeTransfer | 417 | 11,837 | 2,398,657 | 84,484 | 🟢 Go 28.39x faster |
| Event_DecodeBatch10 | 4,101 | 123,131 | 243,843 | 8,121 | 🟢 Go 30.02x faster |
| Event_DecodeBatch100 | 42,196 | 1,187,409 | 23,699 | 842 | 🟢 Go 28.14x faster |
| Hash_Keccak256Short | 435 | 7,543 | 2,300,966 | 132,569 | 🟢 Go 17.36x faster |
| Hash_Keccak256Long | 2,688 | 60,971 | 372,024 | 16,401 | 🟢 Go 22.68x faster |
| Hash_Keccak256Hex | 460 | 7,119 | 2,173,913 | 140,464 | 🟢 Go 15.48x faster |
| Hash_Sha256Short | 160 | 1,487 | 6,242,197 | 672,569 | 🟢 Go 9.28x faster |
| Hash_Sha256Long | 663 | 14,108 | 1,507,841 | 70,882 | 🟢 Go 21.27x faster |
| Hash_FunctionSelector | 2,063 | 8,651 | 484,731 | 115,598 | 🟢 Go 4.19x faster |
| Hash_EventSelector | 2,421 | 8,365 | 413,052 | 119,545 | 🟢 Go 3.46x faster |
| Multicall_Basic | 201,757 | 469,942 | 4,956 | 2,128 | 🟢 Go 2.33x faster |
| Multicall_WithArgs | 238,069 | 394,218 | 4,200 | 2,537 | 🟢 Go 1.66x faster |
| Multicall_MultiContract | 229,500 | 466,651 | 4,357 | 2,143 | 🟢 Go 2.03x faster |
| Multicall_10Calls | 237,288 | 522,065 | 4,214 | 1,915 | 🟢 Go 2.20x faster |
| Multicall_30Calls | 482,191 | 1,020,637 | 2,074 | 980 | 🟢 Go 2.12x faster |
| Multicall_Deployless | 354,896 | 687,205 | 2,818 | 1,455 | 🟢 Go 1.94x faster |
| Multicall_TokenMetadata | 210,526 | 430,263 | 4,750 | 2,324 | 🟢 Go 2.04x faster |
| Multicall_50Calls | 628,617 | 1,499,093 | 1,591 | 667 | 🟢 Go 2.38x faster |
| Multicall_100Calls | 1,081,524 | 2,662,194 | 925 | 376 | 🟢 Go 2.46x faster |
| Multicall_200Calls | 1,893,893 | 5,728,689 | 528 | 175 | 🟢 Go 3.02x faster |
| Multicall_500Calls | 3,049,869 | 10,795,604 | 328 | 93 | 🟢 Go 3.54x faster |
| Multicall_MixedContracts_100 | 1,887,659 | 2,553,952 | 530 | 392 | 🟢 Go 1.35x faster |
| Multicall_1000Calls | 6,557,866 | 20,927,595 | 152 | 48 | 🟢 Go 3.19x faster |
| Multicall_10000Calls_SingleRPC | 100,524,853 | 209,200,644 | 10 | 5 | 🟢 Go 2.08x faster |
| Multicall_10000Calls_Chunked | 31,845,494 | 218,966,914 | 31 | 5 | 🟢 Go 6.88x faster |
| Multicall_10000Calls_AggressiveChunking | 30,096,187 | 244,385,249 | 33 | 4 | 🟢 Go 8.12x faster |
| Signature_HashMessage | 805 | 8,647 | 1,241,773 | 115,648 | 🟢 Go 10.74x faster |
| Signature_HashMessageLong | 1,846 | 18,234 | 541,712 | 54,843 | 🟢 Go 9.88x faster |
| Signature_RecoverAddress | 26,312 | 1,673,780 | 38,005 | 597 | 🟢 Go 63.61x faster |
| Signature_VerifyMessage | 26,069 | 1,572,154 | 38,360 | 636 | 🟢 Go 60.31x faster |
| Signature_ParseSignature | 184 | 1,908 | 5,434,783 | 524,037 | 🟢 Go 10.37x faster |
| Unit_ParseEther | 118 | 246 | 8,481,764 | 4,068,068 | 🟢 Go 2.08x faster |
| Unit_ParseEtherLarge | 319 | 233 | 3,136,763 | 4,293,864 | 🔵 TS 1.37x faster |
| Unit_FormatEther | 118 | 143 | 8,453,085 | 6,975,713 | 🟢 Go 1.21x faster |
| Unit_ParseUnits6 | 102 | 218 | 9,832,842 | 4,586,885 | 🟢 Go 2.14x faster |
| Unit_ParseGwei | 107 | 203 | 9,363,296 | 4,927,703 | 🟢 Go 1.90x faster |
| Unit_FormatUnits | 97 | 133 | 10,340,192 | 7,507,014 | 🟢 Go 1.38x faster |

## Win Summary

- 🟢 Go wins: 55 (93%)
- 🔵 TS wins: 3 (5%)
- ⚪ Ties: 1 (2%)

## Notes

- Benchmarks run against the same Anvil instance (mainnet fork) for fair comparison
- ns/op = nanoseconds per operation (lower is better)
- ops/s = operations per second (higher is better)
- 🟢 = Go faster, 🔵 = TS faster, ⚪ = Similar (within 5%)
