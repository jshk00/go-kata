mem_simple=$(go test -bench=BenchmarkMemorySimpleMap | grep BenchmarkMemorySimpleMap | awk '{print int($5)}')
mem_sharded=$(go test -bench=BenchmarkMemoryShardedMap | grep BenchmarkMemoryShardedMap | awk '{print int($5)}')
if [ "$mem_sharded" -gt $((mem_simple+50)) ]; then
  echo "test failed sharded map consumes more memory"
else
  echo "test passed"
fi
echo "simple map memory $mem_simple MB"
echo "sharded map memory $mem_sharded MB"
