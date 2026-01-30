package hash_test

import (
	"testing"

	"github.com/ChefBingbong/viem-go/utils/hash"
)

func TestKeccak256(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			"hello world bytes",
			[]byte("hello world"),
			"0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad",
		},
		{
			"hello world string",
			"hello world",
			"0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad",
		},
		{
			"hex input",
			"0x68656c6c6f20776f726c64", // "hello world" in hex
			"0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad",
		},
		{
			"empty bytes",
			[]byte{},
			"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hash.Keccak256(tt.input)
			if result != tt.expected {
				t.Errorf("Keccak256(%v) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSha256(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			"hello world bytes",
			[]byte("hello world"),
			"0xb94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
		{
			"hello world string",
			"hello world",
			"0xb94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
		{
			"hex input",
			"0x68656c6c6f20776f726c64", // "hello world" in hex
			"0xb94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hash.Sha256(tt.input)
			if result != tt.expected {
				t.Errorf("Sha256(%v) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRipemd160(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			"hello world bytes",
			[]byte("hello world"),
			"0x98c615784ccb5fe5936fbc0cbe9dfdb408d92f0f",
		},
		{
			"hello world string",
			"hello world",
			"0x98c615784ccb5fe5936fbc0cbe9dfdb408d92f0f",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hash.Ripemd160(tt.input)
			if result != tt.expected {
				t.Errorf("Ripemd160(%v) = %s, want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsHash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid hash", "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad", true},
		{"too short", "0x1234", false},
		{"no prefix", "47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad", false},
		{"too long", "0x47173285a8d7341e5e972fc677286384f802f8ef42a5ec5f03bbfa254cb01fad00", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hash.IsHash(tt.input)
			if result != tt.expected {
				t.Errorf("IsHash(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeSignature(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		{
			"function with names",
			"function transfer(address to, uint256 amount)",
			"transfer(address,uint256)",
			false,
		},
		{
			"event with indexed",
			"event Transfer(address indexed from, address indexed to, uint256 amount)",
			"Transfer(address,address,uint256)",
			false,
		},
		{
			"already normalized",
			"transfer(address,uint256)",
			"transfer(address,uint256)",
			false,
		},
		{
			"function ownerOf",
			"function ownerOf(uint256 tokenId)",
			"ownerOf(uint256)",
			false,
		},
		{
			"complex tuple",
			"function foo((address sender, uint256 amount) data)",
			"foo((address,uint256))",
			false,
		},
		{
			"invalid",
			"not a signature",
			"",
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := hash.NormalizeSignature(tt.input)
			if tt.hasError && err == nil {
				t.Errorf("NormalizeSignature(%q) expected error", tt.input)
			}
			if !tt.hasError && err != nil {
				t.Errorf("NormalizeSignature(%q) unexpected error: %v", tt.input, err)
			}
			if !tt.hasError && result != tt.expected {
				t.Errorf("NormalizeSignature(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToFunctionSelector(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"ownerOf",
			"function ownerOf(uint256 tokenId)",
			"0x6352211e",
		},
		{
			"transfer",
			"function transfer(address to, uint256 amount)",
			"0xa9059cbb",
		},
		{
			"balanceOf",
			"function balanceOf(address owner)",
			"0x70a08231",
		},
		{
			"approve",
			"function approve(address spender, uint256 amount)",
			"0x095ea7b3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := hash.ToFunctionSelector(tt.input)
			if err != nil {
				t.Errorf("ToFunctionSelector(%q) error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("ToFunctionSelector(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToEventSelector(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"Transfer ERC20",
			"event Transfer(address indexed from, address indexed to, uint256 amount)",
			"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		},
		{
			"Approval",
			"event Approval(address indexed owner, address indexed spender, uint256 amount)",
			"0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := hash.ToEventSelector(tt.input)
			if err != nil {
				t.Errorf("ToEventSelector(%q) error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("ToEventSelector(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToSignature(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			"function",
			"function ownerOf(uint256 tokenId)",
			"ownerOf(uint256)",
		},
		{
			"event",
			"event Transfer(address indexed from, address indexed to, uint256 amount)",
			"Transfer(address,address,uint256)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := hash.ToSignature(tt.input)
			if err != nil {
				t.Errorf("ToSignature(%q) error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("ToSignature(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestHashSignature(t *testing.T) {
	// Direct signature hash (no normalization)
	result := hash.HashSignature("transfer(address,uint256)")
	expected := "0xa9059cbb2ab09eb219583f4a59a5d0623ade346d962bcd4e46b11da047c9049b"
	if result != expected {
		t.Errorf("HashSignature = %q, want %q", result, expected)
	}
}

// Benchmark tests

func BenchmarkKeccak256(b *testing.B) {
	data := []byte("hello world")
	for i := 0; i < b.N; i++ {
		hash.Keccak256(data)
	}
}

func BenchmarkToFunctionSelector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash.ToFunctionSelector("function transfer(address to, uint256 amount)")
	}
}

func BenchmarkNormalizeSignature(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash.NormalizeSignature("function transfer(address to, uint256 amount)")
	}
}
