package shardedmap

import (
	"runtime"
	"sync"
	"testing"
)

// TestMap checks if concurrent put call to maps are safe or not
func TestMap(t *testing.T) {
	data := []struct {
		key  string
		val  string
		want string
	}{
		{"key1", "val1", "val1"},
		{"key2", "val2", "val2"},
		{"key3", "val3", "val3"},
		{"key4", "val4", "val4"},
		{"key5", "val5", "val5"},
		{"key6", "val6", "val6"},
		{"key7", "val7", "val7"},
		{"key8", "val8", "val8"},
		{"key9", "val9", "val9"},
		{"key10", "val10", "val10"},
		{"key11", "val11", "val11"},
		{"key12", "val12", "val12"},
	}

	t.Run("concurrent-put", func(t *testing.T) {
		var wg sync.WaitGroup
		mm := NewShardedMap[string, string](64, HashString)
		for _, el := range data {
			wg.Go(func() {
				mm.Set(el.key, el.val)
			})
		}
		wg.Wait()
		for _, d := range data {
			if v, _ := mm.Get(d.key); v != d.want {
				t.Errorf(
					"must be concurrency issue we wanted %s but got %s for %s key",
					d.want,
					v,
					d.key,
				)
			}
		}
	})

	t.Run("concurrent-get", func(t *testing.T) {
		mm := NewShardedMap[string, string](64, HashString)
		for _, el := range data {
			mm.Set(el.key, el.val)
		}
		var wg sync.WaitGroup
		for _, el := range data {
			wg.Go(func() {
				if v, _ := mm.Get(el.key); v == "" {
					t.Error("value should not be empty")
				}
			})
		}
		wg.Wait()
	})

	t.Run("concurrent-put-get", func(t *testing.T) {
		mm := NewShardedMap[string, string](64, HashString)
		var wg sync.WaitGroup
		wg.Go(func() {
			for _, el := range data {
				mm.Set(el.key, el.val)
			}
		})

		wg.Go(func() {
			for _, el := range data {
				if v, _ := mm.Get(el.key); v == "" {
					t.Log("values is empty")
				}
			}
		})
		wg.Wait()
		t.Log(mm.Keys())
	})
}

func runShardedMapBenchmarks(b *testing.B, shardCount uint) {
	m := NewShardedMap[int, string](shardCount, HashInt)

	b.ResetTimer()
	b.ReportAllocs()

	var wg sync.WaitGroup
	wg.Add(8)

	for id := range 8 {
		go func(id int) {
			defer wg.Done()
			for i := id; i < b.N; i += 8 {
				m.Set(i, "v")
			}
		}(id)
	}
	wg.Wait()
}

func Benchmark1Shard(b *testing.B) {
	runShardedMapBenchmarks(b, 1)
}

func Benchmark64Shard(b *testing.B) {
	runShardedMapBenchmarks(b, 64)
}

func BenchmarkMemorySimpleMap(b *testing.B) {
	stats := &runtime.MemStats{}
	m := make(map[int]any)
	for i := range 1_000_000 {
		m[i] = i
	}
	runtime.ReadMemStats(stats)
	b.ReportMetric(float64(stats.Sys/1024/1024), "MB")
}

func BenchmarkMemoryShardedMap(b *testing.B) {
	stats := &runtime.MemStats{}
	m := NewShardedMap[int, any](64, HashInt)
	for i := range 1_000_000 {
		m.Set(i, i)
	}
	runtime.ReadMemStats(stats)
	b.ReportMetric(float64(stats.Sys/1024/1024), "MB")
}
