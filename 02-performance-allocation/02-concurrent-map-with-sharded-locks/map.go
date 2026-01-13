package shardedmap

import (
	"hash/fnv"
	"sync"
	"unsafe"
)

type Map[K comparable, V any] struct {
	m map[K]V
	sync.RWMutex
}

type ShardedMap[K comparable, V any] struct {
	shards []*Map[K, V]
	count  uint
	hash   func(key K) uint64
}

func NewShardedMap[K comparable, V any](count uint, hash func(key K) uint64) *ShardedMap[K, V] {
	shards := make([]*Map[K, V], 0, count)
	for range count {
		shards = append(shards, &Map[K, V]{
			m: make(map[K]V),
		})
	}
	return &ShardedMap[K, V]{shards: shards, count: count, hash: hash}
}

func (sm *ShardedMap[K, V]) GetShard(key K) *Map[K, V] {
	return sm.shards[uint(sm.hash(key))%sm.count]
}

func (sm *ShardedMap[K, V]) Set(k K, v V) {
	s := sm.GetShard(k)
	s.Lock()
	s.m[k] = v
	s.Unlock()
}

func (sm *ShardedMap[K, V]) Get(key K) (V, bool) {
	s := sm.GetShard(key)
	s.RLock()
	v, ok := s.m[key]
	s.RUnlock()
	return v, ok
}

func (sm *ShardedMap[K, V]) Delete(key K) {
	s := sm.GetShard(key)
	s.Lock()
	delete(s.m, key)
	s.Unlock()
}

func (sm *ShardedMap[K, V]) Keys() []K {
	var arr []K
	for _, s := range sm.shards {
		s.RLock()
		for k := range s.m {
			arr = append(arr, k)
		}
		s.RUnlock()
	}
	return arr
}

func HashString(s string) uint64 {
	h := fnv.New64()
	h.Write(unsafe.Slice(unsafe.StringData(s), len(s)))
	return h.Sum64()
}

// Hashes the integer using SplitMix64 Method.
func hasUint64(x uint64) uint64 {
	x += 0x9e3779b97f4a7c15                  // Increment with golden ration costant
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9 // xor the variable with variable right bit shifted 30 and multiply with constant
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb // xor the variable with variable right bit shifted 27 and multiply with constant
	return x ^ (x >> 31)                     // xor the variable with variable right bit shifted 31
}

func HashUint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64](x T) uint64 {
	return hasUint64(uint64(x))
}

func HashInt[T ~int | ~int8 | ~int16 | ~int32 | ~int64](x T) uint64 {
	return hasUint64(uint64(int64(x)))
}
