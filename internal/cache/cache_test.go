package cache_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/cache"
)

func makeEntry(mod time.Time, size int64) cache.Entry {
	return cache.Entry{
		ModTime:    mod,
		Size:       size,
		Offsets:    []int64{0, 128, 256},
		Timestamps: []time.Time{mod, mod.Add(time.Second), mod.Add(2 * time.Second)},
	}
}

func TestCache_PutAndGet_HitWhenValid(t *testing.T) {
	c := cache.New()
	mod := time.Now().Truncate(time.Second)
	entry := makeEntry(mod, 1024)
	c.Put("/var/log/app.log", entry)

	got, ok := c.Get("/var/log/app.log", mod, 1024)
	if !ok {
		t.Fatal("expected cache hit")
	}
	if len(got.Offsets) != 3 {
		t.Fatalf("expected 3 offsets, got %d", len(got.Offsets))
	}
}

func TestCache_Get_MissOnModTimeChange(t *testing.T) {
	c := cache.New()
	mod := time.Now().Truncate(time.Second)
	c.Put("/var/log/app.log", makeEntry(mod, 1024))

	_, ok := c.Get("/var/log/app.log", mod.Add(time.Minute), 1024)
	if ok {
		t.Fatal("expected cache miss due to modtime change")
	}
}

func TestCache_Get_MissOnSizeChange(t *testing.T) {
	c := cache.New()
	mod := time.Now().Truncate(time.Second)
	c.Put("/var/log/app.log", makeEntry(mod, 1024))

	_, ok := c.Get("/var/log/app.log", mod, 2048)
	if ok {
		t.Fatal("expected cache miss due to size change")
	}
}

func TestCache_Get_MissOnUnknownPath(t *testing.T) {
	c := cache.New()
	_, ok := c.Get("/nonexistent.log", time.Now(), 0)
	if ok {
		t.Fatal("expected cache miss for unknown path")
	}
}

func TestCache_Invalidate_RemovesEntry(t *testing.T) {
	c := cache.New()
	mod := time.Now().Truncate(time.Second)
	c.Put("/var/log/app.log", makeEntry(mod, 512))
	c.Invalidate("/var/log/app.log")

	_, ok := c.Get("/var/log/app.log", mod, 512)
	if ok {
		t.Fatal("expected cache miss after invalidation")
	}
}

func TestCache_Len(t *testing.T) {
	c := cache.New()
	if c.Len() != 0 {
		t.Fatalf("expected 0, got %d", c.Len())
	}
	mod := time.Now()
	c.Put("/a.log", makeEntry(mod, 100))
	c.Put("/b.log", makeEntry(mod, 200))
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
	c.Invalidate("/a.log")
	if c.Len() != 1 {
		t.Fatalf("expected 1, got %d", c.Len())
	}
}
