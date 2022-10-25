package casstore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/cas.v2"
)

var (
	store        cas.SessionStore
	_defaultEtcd = "http://127.0.0.1:2379"
)

func init() {
	store, _ = NewEtcdSessionStore(clientv3.Config{Endpoints: []string{_defaultEtcd}},
		context.Background(), "/cas/sessions", 3600)
}

func TestSessionStore_Get(t *testing.T) {

	v, ok := store.Get("key1")
	require.False(t, ok)
	require.Equal(t, "", v)

	err := store.Set("key1", "value1")
	require.Nil(t, err)

	v, ok = store.Get("key1")
	require.True(t, ok)
	require.Equal(t, "value1", v)
}

func TestSessionStore_Set(t *testing.T) {
	err := store.Set("key1", "value1")
	require.Nil(t, err)

	err = store.Set("key2", "value2")
	require.Nil(t, err)

	v, ok := store.Get("key1")
	require.True(t, ok)
	require.Equal(t, "value1", v)

	v, ok = store.Get("key2")
	require.True(t, ok)
	require.Equal(t, "value2", v)

	err = store.Set("key2", "value2-new")
	require.Nil(t, err)

	v, ok = store.Get("key2")
	require.True(t, ok)
	require.Equal(t, "value2-new", v)
}

func TestSessionStore_Delete(t *testing.T) {
	err := store.Set("key1", "value1")
	require.Nil(t, err)

	err = store.Set("key2", "value2")
	require.Nil(t, err)

	v, ok := store.Get("key1")
	require.True(t, ok)
	require.Equal(t, "value1", v)

	v, ok = store.Get("key2")
	require.True(t, ok)
	require.Equal(t, "value2", v)

	err = store.Delete("key2")
	require.Nil(t, err)

	v, ok = store.Get("key1")
	require.True(t, ok)
	require.Equal(t, "value1", v)

	v, ok = store.Get("key2")
	require.False(t, ok)
	require.Equal(t, "", v)

	err = store.Delete("key1")
	require.Nil(t, err)

	v, ok = store.Get("key1")
	require.False(t, ok)
	require.Equal(t, "", v)
}
