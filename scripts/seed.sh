#!/bin/bash

for i in {1..100}; do
  key="key$i"
  value="value$i"

  # TTL between 3s (3000ms) and 13s (13000ms)
  ttl=$(( (RANDOM % 10000) + 3000 ))

  curl -s -X POST \
    "http://localhost:3000/set?key=$key&ttl=${ttl}ms" \
    -d "$value" > /dev/null

  echo "✅ Sent $key with TTL ${ttl}ms"

  # Optional: sleep to pace the firehose
  sleep 0.03
done

echo "🎉 All 100 keys sent. Watch them get deleted as TTLs expire..."
