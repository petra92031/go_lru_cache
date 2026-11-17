package go_lru_cache

import (
	"testing"
)

func TestBasicGetPut(t *testing.T) {
	c := NewCache[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)

	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("Get(a) = %v, %v; want 1, true", v, ok)
	}
	if v, ok := c.Get("b"); !ok || v != 2 {
		t.Fatalf("Get(b) = %v, %v; want 2, true", v, ok)
	}
	if v, ok := c.Get("c"); ok || v != 0 {
		t.Fatalf("Get(c) = %v, %v; want 0, false", v, ok)
	}
}

func TestEviction(t *testing.T) {
	c := NewCache[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3) // evicts "a"

	if _, ok := c.Get("a"); ok {
		t.Fatal("Get(a) should have failed after eviction")
	}
	if v, ok := c.Get("b"); !ok || v != 2 {
		t.Fatalf("Get(b) = %v, %v; want 2, true", v, ok)
	}
	if v, ok := c.Get("c"); !ok || v != 3 {
		t.Fatalf("Get(c) = %v, %v; want 3, true", v, ok)
	}
}

func TestGetUpdatesRecency(t *testing.T) {
	c := NewCache[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Get("a") // "a" is now most recently used
	c.Put("c", 3) // should evict "b"

	if _, ok := c.Get("b"); ok {
		t.Fatal("Get(b) should have failed after eviction")
	}
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("Get(a) = %v, %v; want 1, true", v, ok)
	}
	if v, ok := c.Get("c"); !ok || v != 3 {
		t.Fatalf("Get(c) = %v, %v; want 3, true", v, ok)
	}
}

func TestUpdateExistingKey(t *testing.T) {
	c := NewCache[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("a", 10) // update "a", should not evict

	if v, ok := c.Get("a"); !ok || v != 10 {
		t.Fatalf("Get(a) = %v, %v; want 10, true", v, ok)
	}
	if v, ok := c.Get("b"); !ok || v != 2 {
		t.Fatalf("Get(b) = %v, %v; want 2, true", v, ok)
	}
	if c.Len() != 2 {
		t.Fatalf("Len() = %d; want 2", c.Len())
	}
}

func TestLen(t *testing.T) {
	c := NewCache[string, int](3)
	if c.Len() != 0 {
		t.Fatalf("Len() = %d; want 0", c.Len())
	}
	c.Put("a", 1)
	if c.Len() != 1 {
		t.Fatalf("Len() = %d; want 1", c.Len())
	}
	c.Put("b", 2)
	c.Put("c", 3)
	if c.Len() != 3 {
		t.Fatalf("Len() = %d; want 3", c.Len())
	}
	c.Put("d", 4) // evicts one
	if c.Len() != 3 {
		t.Fatalf("Len() = %d; want 3", c.Len())
	}
}

func TestCapacityOne(t *testing.T) {
	c := NewCache[string, int](1)
	c.Put("a", 1)
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Fatalf("Get(a) = %v, %v; want 1, true", v, ok)
	}
	c.Put("b", 2) // evicts "a"
	if _, ok := c.Get("a"); ok {
		t.Fatal("Get(a) should have failed after eviction")
	}
	if v, ok := c.Get("b"); !ok || v != 2 {
		t.Fatalf("Get(b) = %v, %v; want 2, true", v, ok)
	}
}

func TestZeroValueNotUsable(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("NewCache(0) did not panic")
		}
	}()
	NewCache[string, int](0)
}

func TestNegativeCapacityPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("NewCache(-1) did not panic")
		}
	}()
	NewCache[string, int](-1)
}

func TestDifferentKeyTypes(t *testing.T) {
	c := NewCache[int, string](2)
	c.Put(1, "one")
	c.Put(2, "two")
	if v, ok := c.Get(1); !ok || v != "one" {
		t.Fatalf("Get(1) = %v, %v; want one, true", v, ok)
	}
	if v, ok := c.Get(2); !ok || v != "two" {
		t.Fatalf("Get(2) = %v, %v; want two, true", v, ok)
	}
}
