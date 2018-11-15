package mock

//go:generate mockgen -destination hash_storage_mock.go -package mock github.com/xdefrag/hash-ketchum/pkg/types HashStorager
