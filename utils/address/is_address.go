package address

import (
	"container/list"
	"strings"
	"sync"
)

// addressCache is a simple LRU cache for IsAddress results.
// Avoids importing utils.LruMap to prevent an import cycle (utils -> utils/address -> utils).
type addressCache struct {
	maxSize int
	cache   map[string]*list.Element
	order   *list.List
	mu      sync.RWMutex
}

type cacheEntry struct {
	key   string
	value bool
}

func newAddressCache(maxSize int) *addressCache {
	return &addressCache{
		maxSize: maxSize,
		cache:   make(map[string]*list.Element),
		order:   list.New(),
	}
}

func (c *addressCache) Get(key string) (bool, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.cache[key]; ok {
		c.order.MoveToFront(elem)
		return elem.Value.(*cacheEntry).value, true
	}
	return false, false
}

func (c *addressCache) Set(key string, value bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.cache[key]; ok {
		elem.Value.(*cacheEntry).value = value
		c.order.MoveToFront(elem)
		return
	}
	entry := &cacheEntry{key: key, value: value}
	elem := c.order.PushFront(entry)
	c.cache[key] = elem
	if c.maxSize > 0 && c.order.Len() > c.maxSize {
		oldest := c.order.Back()
		if oldest != nil {
			c.order.Remove(oldest)
			delete(c.cache, oldest.Value.(*cacheEntry).key)
		}
	}
}

// IsAddressCache is the LRU cache for address validation results.
// Mirrors viem's isAddressCache with 8192 entries.
var IsAddressCache = newAddressCache(8192)

// addressStringCache is a simple LRU cache for string values (used by ChecksumAddress).
type addressStringCache struct {
	maxSize int
	cache   map[string]*list.Element
	order   *list.List
	mu      sync.RWMutex
}

type stringCacheEntry struct {
	key   string
	value string
}

func newAddressStringCache(maxSize int) *addressStringCache {
	return &addressStringCache{
		maxSize: maxSize,
		cache:   make(map[string]*list.Element),
		order:   list.New(),
	}
}

func (c *addressStringCache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.cache[key]; ok {
		c.order.MoveToFront(elem)
		return elem.Value.(*stringCacheEntry).value, true
	}
	return "", false
}

func (c *addressStringCache) Set(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.cache[key]; ok {
		elem.Value.(*stringCacheEntry).value = value
		c.order.MoveToFront(elem)
		return
	}
	entry := &stringCacheEntry{key: key, value: value}
	elem := c.order.PushFront(entry)
	c.cache[key] = elem
	if c.maxSize > 0 && c.order.Len() > c.maxSize {
		oldest := c.order.Back()
		if oldest != nil {
			c.order.Remove(oldest)
			delete(c.cache, oldest.Value.(*stringCacheEntry).key)
		}
	}
}

// IsAddressOptions configures address validation behavior.
type IsAddressOptions struct {
	// Strict enables checksum validation. Default is true.
	Strict bool
}

// IsAddress checks if a string is a valid Ethereum address.
// Results are cached in an LRU cache for repeated lookups.
// By default (strict=true), it validates the checksum if the address contains mixed case.
//
// Example:
//
//	isAddress("0xa5cc3c03994db5b0d9a5eedd10cabab0813678ac") // true (lowercase)
//	isAddress("0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC") // true (valid checksum)
//	isAddress("0xa5cc3c03994DB5b0d9A5eEdD10CabaB0813678AC", IsAddressOptions{Strict: false}) // true
func IsAddress(address string, opts ...IsAddressOptions) bool {
	strict := true
	if len(opts) > 0 {
		strict = opts[0].Strict
	}

	// Build cache key: "address.true" or "address.false"
	cacheKey := address + ".true"
	if !strict {
		cacheKey = address + ".false"
	}

	// Check cache
	if cached, ok := IsAddressCache.Get(cacheKey); ok {
		return cached
	}

	result := isAddressCore(address, strict)
	IsAddressCache.Set(cacheKey, result)
	return result
}

// isAddressCore is the uncached validation logic.
func isAddressCore(address string, strict bool) bool {
	// Check length: "0x" + 40 hex chars = 42
	if len(address) != 42 {
		return false
	}

	// Check prefix
	if address[0] != '0' || address[1] != 'x' {
		return false
	}

	// Check all chars are valid hex
	for i := 2; i < 42; i++ {
		c := address[i]
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') && (c < 'A' || c > 'F') {
			return false
		}
	}

	// If all lowercase, valid (no checksum to verify)
	if strings.ToLower(address) == address {
		return true
	}

	// If strict mode and contains uppercase, verify checksum
	if strict {
		checksummed := ChecksumAddress(address)
		return checksummed == Address(address)
	}

	return true
}
