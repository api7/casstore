package casstore

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/cas.v2"
)

var _ cas.SessionStore = &etcdSessionStore{}

// NewEtcdSessionStore create a session store using etcd.
func NewEtcdSessionStore(config clientv3.Config, ctx context.Context,
	prefix string, maxAge int64) (cas.SessionStore, error) {
	client, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}

	if prefix == "" {
		prefix = "/cas/sessions"
	}

	if maxAge == 0 {
		maxAge = 86400
	}

	return &etcdSessionStore{
		cli:    client,
		ctx:    ctx,
		prefix: prefix,
		maxAge: maxAge,
	}, nil
}

type etcdSessionStore struct {
	cli    *clientv3.Client
	ctx    context.Context
	prefix string
	maxAge int64
}

func (s *etcdSessionStore) Get(sessionID string) (string, bool) {
	key := s.prefix + "/" + sessionID
	resp, err := s.cli.Get(s.ctx, key)
	if err != nil {
		return "", false
	}

	if resp.Count == 0 {
		return "", false
	}

	return string(resp.Kvs[0].Value), true
}

func (s *etcdSessionStore) Set(sessionID, ticket string) error {
	key := s.prefix + "/" + sessionID

	grant, err := s.cli.Grant(s.ctx, s.maxAge+1)
	if err != nil {
		return err
	}

	_, err = s.cli.Put(s.ctx, key, ticket, clientv3.WithLease(grant.ID))

	return err
}

func (s *etcdSessionStore) Delete(sessionID string) error {
	key := s.prefix + "/" + sessionID
	resp, err := s.cli.Delete(s.ctx, key)
	if err != nil {
		return err
	}

	if resp.Deleted == 0 {
		return fmt.Errorf("key: %s is not found in etcd", key)
	}

	return nil
}
