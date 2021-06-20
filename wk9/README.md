redis-cli flushall
redis-benchmark -d 10 -r 1000000 -n 500000 -t set
