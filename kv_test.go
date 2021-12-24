package kv

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKV(t *testing.T) {
	assert := assert.New(t)

	rds, err := newRedisTest("hello", "world", 10*time.Second)
	assert.NoError(err)

	// 对象初始化
	objCreater := func(key string) (*TestObject, error) {
		return &TestObject{
			Int: 1,
			Str: key,
		}, nil
	}

	// 创建指定类型的对象仓库
	kv := NewKV(rds, rds, objCreater)

	// 获取对象
	obj, err := kv.Get("123")
	assert.NoError(err)
	assert.Equal(&TestObject{Int: 1, Str: "123"}, obj)

	// 持久对象
	obj.Int = 9527
	err = kv.Set("123", obj)
	assert.NoError(err)

	obj, err = kv.Get("123")
	assert.NoError(err)
	assert.Equal(&TestObject{Int: 9527, Str: "123"}, obj)
}
