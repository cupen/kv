package kv

type KV [K string | int | int64, T any]struct {
	cache   DB
	db      DB
	creater func(K) (*T, error)
}

func NewKV[ K string | int | int64, T any](db, cache DB, creater func(K) (*T, error)) *KV[K,T] {
	return &KV[K,T]{
		cache:   cache,
		db:      db,
		creater: creater,
	}
}

func (this *KV[K,T]) Get(key K) (*T, error) {
	obj := new(T)
	if err := this.cache.Get(key, obj); err != nil {
		if IsNotFound(err) {
			if this.db == nil {
				return this.creater(key)
			}
			if err := this.cache.Get(key, obj); err != nil {
				if IsNotFound(err) {
					return this.creater(key)
				}
			}
		}
		return nil, err
	}
	return obj, nil
}

func (this *KV[K,T]) Set(key interface{}, value interface{}) (error) {
	this.cache.Set(key, value)
	return this.db.Set(key, value)
}

func (this *KV[K,T]) Has(key interface{}) (bool, error) {
	return false, nil
}