# Benchmark Comparison: viem-go vs viem TypeScript

Generated: 2026-02-07T04:58:28.510Z

## Overall Summary

**🏆 Go is 4.41x faster overall**

| Metric | Go | TypeScript |
|--------|----|-----------|
| Avg ns/op | 3,492,003 | 15,401,105 |
| Avg ops/s | 286 | 65 |
| Wins | 55/59 | 3/59 |

## By Suite

| Suite | Benchmarks | Go Wins | TS Wins | Ties | Winner |
|-------|------------|---------|---------|------|--------|
| abi | 6 | 6 | 0 | 0 | 🟢 Go 16.22x faster |
| address | 5 | 2 | 2 | 1 | 🟢 Go 3.96x faster |
| call | 6 | 6 | 0 | 0 | 🟢 Go 73.82x faster |
| ens | 5 | 5 | 0 | 0 | 🟢 Go 14.79x faster |
| event | 3 | 3 | 0 | 0 | 🟢 Go 30.85x faster |
| hash | 7 | 7 | 0 | 0 | 🟢 Go 11.41x faster |
| multicall | 16 | 16 | 0 | 0 | 🟢 Go 3.97x faster |
| signature | 5 | 5 | 0 | 0 | 🟢 Go 59.87x faster |
| unit | 6 | 5 | 1 | 0 | 🟢 Go 1.47x faster |

## Detailed Results

| Benchmark | Go (ns/op) | TS (ns/op) | Go (ops/s) | TS (ops/s) | Result |
|-----------|------------|------------|------------|------------|--------|
| Abi_EncodeSimple | 249 | 8,899 | 4,022,526 | 112,370 | 🟢 Go 35.80x faster |
| Abi_EncodeComplex | 341 | 9,211 | 2,935,995 | 108,571 | 🟢 Go 27.04x faster |
| Abi_EncodeMultiArg | 457 | 9,890 | 2,188,184 | 101,108 | 🟢 Go 21.64x faster |
| Abi_DecodeResult | 99 | 1,061 | 10,108,157 | 942,525 | 🟢 Go 10.72x faster |
| Abi_EncodePacked | 329 | 697 | 3,044,140 | 1,434,746 | 🟢 Go 2.12x faster |
| Abi_EncodePackedMulti | 442 | 1,316 | 2,263,468 | 759,683 | 🟢 Go 2.98x faster |
| Address_IsAddress | 928 | 298 | 1,078,167 | 3,351,033 | 🔵 TS 3.11x faster |
| Address_IsAddressLower | 320 | 302 | 3,128,911 | 3,311,220 | 🔵 TS 1.06x faster |
| Address_Checksum | 840 | 818 | 1,190,476 | 1,223,024 | ⚪ Similar |
| Address_Create | 2,431 | 9,671 | 411,353 | 103,398 | 🟢 Go 3.98x faster |
| Address_Create2 | 2,694 | 17,498 | 371,195 | 57,151 | 🟢 Go 6.50x faster |
| Call_Basic | 196,682 | 18,719,581 | 5,084 | 53 | 🟢 Go 95.18x faster |
| Call_WithData | 193,787 | 18,549,434 | 5,160 | 54 | 🟢 Go 95.72x faster |
| Call_WithAccount | 204,228 | 324,588 | 4,896 | 3,081 | 🟢 Go 1.59x faster |
| Call_Decimals | 201,236 | 17,731,534 | 4,969 | 56 | 🟢 Go 88.11x faster |
| Call_Symbol | 238,621 | 17,735,182 | 4,191 | 56 | 🟢 Go 74.32x faster |
| Call_BalanceOfMultiple | 192,364 | 17,511,172 | 5,198 | 57 | 🟢 Go 91.03x faster |
| Ens_Namehash | 1,603 | 28,984 | 623,830 | 34,502 | 🟢 Go 18.08x faster |
| Ens_NamehashDeep | 3,182 | 55,769 | 314,268 | 17,931 | 🟢 Go 17.53x faster |
| Ens_Labelhash | 442 | 7,205 | 2,263,980 | 138,799 | 🟢 Go 16.31x faster |
| Ens_Normalize | 365 | 979 | 2,740,477 | 1,021,913 | 🟢 Go 2.68x faster |
| Ens_NormalizeLong | 870 | 2,620 | 1,148,897 | 381,724 | 🟢 Go 3.01x faster |
| Event_DecodeTransfer | 389 | 12,163 | 2,572,016 | 82,216 | 🟢 Go 31.28x faster |
| Event_DecodeBatch10 | 3,811 | 117,472 | 262,398 | 8,513 | 🟢 Go 30.82x faster |
| Event_DecodeBatch100 | 38,016 | 1,172,649 | 26,305 | 853 | 🟢 Go 30.85x faster |
| Hash_Keccak256Short | 436 | 7,135 | 2,295,157 | 140,161 | 🟢 Go 16.38x faster |
| Hash_Keccak256Long | 2,718 | 60,568 | 367,918 | 16,510 | 🟢 Go 22.28x faster |
| Hash_Keccak256Hex | 453 | 7,106 | 2,205,558 | 140,718 | 🟢 Go 15.67x faster |
| Hash_Sha256Short | 164 | 1,458 | 6,108,735 | 686,004 | 🟢 Go 8.90x faster |
| Hash_Sha256Long | 639 | 13,967 | 1,564,456 | 71,595 | 🟢 Go 21.85x faster |
| Hash_FunctionSelector | 2,100 | 8,298 | 476,190 | 120,511 | 🟢 Go 3.95x faster |
| Hash_EventSelector | 2,883 | 8,636 | 346,861 | 115,800 | 🟢 Go 3.00x faster |
| Multicall_Basic | 225,243 | 458,142 | 4,440 | 2,183 | 🟢 Go 2.03x faster |
| Multicall_WithArgs | 231,361 | 427,771 | 4,322 | 2,338 | 🟢 Go 1.85x faster |
| Multicall_MultiContract | 283,634 | 466,242 | 3,526 | 2,145 | 🟢 Go 1.64x faster |
| Multicall_10Calls | 303,542 | 542,085 | 3,294 | 1,845 | 🟢 Go 1.79x faster |
| Multicall_30Calls | 576,186 | 1,033,934 | 1,736 | 967 | 🟢 Go 1.79x faster |
| Multicall_Deployless | 392,787 | 624,715 | 2,546 | 1,601 | 🟢 Go 1.59x faster |
| Multicall_TokenMetadata | 274,697 | 433,706 | 3,640 | 2,306 | 🟢 Go 1.58x faster |
| Multicall_50Calls | 888,228 | 1,467,244 | 1,126 | 682 | 🟢 Go 1.65x faster |
| Multicall_100Calls | 1,393,872 | 2,633,866 | 717 | 380 | 🟢 Go 1.89x faster |
| Multicall_200Calls | 2,532,627 | 5,338,743 | 395 | 187 | 🟢 Go 2.11x faster |
| Multicall_500Calls | 3,915,287 | 11,255,849 | 255 | 89 | 🟢 Go 2.87x faster |
| Multicall_MixedContracts_100 | 1,387,717 | 2,613,969 | 721 | 383 | 🟢 Go 1.88x faster |
| Multicall_1000Calls | 6,328,351 | 21,352,328 | 158 | 47 | 🟢 Go 3.37x faster |
| Multicall_10000Calls_SingleRPC | 123,074,232 | 324,348,870 | 8 | 3 | 🟢 Go 2.64x faster |
| Multicall_10000Calls_Chunked | 30,255,180 | 217,004,470 | 33 | 5 | 🟢 Go 7.17x faster |
| Multicall_10000Calls_AggressiveChunking | 32,615,701 | 223,259,137 | 31 | 4 | 🟢 Go 6.85x faster |
| Signature_HashMessage | 752 | 8,456 | 1,330,318 | 118,254 | 🟢 Go 11.25x faster |
| Signature_HashMessageLong | 1,713 | 18,255 | 583,771 | 54,778 | 🟢 Go 10.66x faster |
| Signature_RecoverAddress | 25,976 | 1,594,261 | 38,497 | 627 | 🟢 Go 61.37x faster |
| Signature_VerifyMessage | 25,936 | 1,643,899 | 38,556 | 608 | 🟢 Go 63.38x faster |
| Signature_ParseSignature | 186 | 1,906 | 5,387,931 | 524,661 | 🟢 Go 10.27x faster |
| Unit_ParseEther | 112 | 241 | 8,952,551 | 4,148,035 | 🟢 Go 2.16x faster |
| Unit_ParseEtherLarge | 308 | 247 | 3,246,753 | 4,055,152 | 🔵 TS 1.25x faster |
| Unit_FormatEther | 115 | 153 | 8,665,511 | 6,523,001 | 🟢 Go 1.33x faster |
| Unit_ParseUnits6 | 98 | 212 | 10,249,052 | 4,706,299 | 🟢 Go 2.18x faster |
| Unit_ParseGwei | 100 | 217 | 10,008,006 | 4,598,678 | 🟢 Go 2.18x faster |
| Unit_FormatUnits | 96 | 144 | 10,415,582 | 6,961,428 | 🟢 Go 1.50x faster |

## Win Summary

- 🟢 Go wins: 55 (93%)
- 🔵 TS wins: 3 (5%)
- ⚪ Ties: 1 (2%)

## Notes

- Benchmarks run against the same Anvil instance (mainnet fork) for fair comparison
- ns/op = nanoseconds per operation (lower is better)
- ops/s = operations per second (higher is better)
- 🟢 = Go faster, 🔵 = TS faster, ⚪ = Similar (within 5%)
