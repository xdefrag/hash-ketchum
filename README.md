

# Hash Ketchum
[![Build Status](https://travis-ci.org/xdefrag/hash-ketchum.svg?branch=master)](https://travis-ci.org/xdefrag/hash-ketchum)

Gotta catch all hashes with leading zeros!

Start redis, server and client with logs:
```
make dc-up
```

Unit tests:
```
make test-unit
```

Integration tests:
```
make redis // or have redis on :6379
make test-integration
```
