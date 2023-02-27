package kv

type (
	Key interface {
		string | int | int32 | int64 | float64
	}
	Value any
)

// objectCreater is a function for create object and initializing it.
type objectCreater[K Key, V Value] func(K) (V, error)

type Store[K Key, V Value] struct {
	cache   DB
	db      DB
	creater func(K) (V, error)
}

func NewStore[K Key, V Value](db, cache DB, creater objectCreater[K, V]) *Store[K, V] {
	return &Store[K, V]{
		cache:   cache,
		db:      db,
		creater: creater,
	}
}

func (this *Store[K, V]) Get(key K) (V, error) {
	obj := new(V)
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
		return *obj, err
	}
	return *obj, nil
}

func (this *Store[K, V]) Set(key K, value *V) error {
	this.cache.Set(key, value)
	return this.db.Set(key, value)
}

func (this *Store[K, V]) Del(key K) error {
	this.cache.Del(key)
	return this.db.Del(key)
}

func (this *Store[K, V]) Has(key K) (bool, error) {
	return false, nil
}

func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}
