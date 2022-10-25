package casstore

import (
	"context"
	"encoding/json"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/cas.v2"
)

var _ cas.TicketStore = &etcdTicketStore{}

// NewEtcdTicketStore create a ticket store using etcd.
func NewEtcdTicketStore(config clientv3.Config, ctx context.Context,
	prefix string, maxAge int64) (cas.TicketStore, error) {
	client, err := clientv3.New(config)
	if err != nil {
		return nil, err
	}
	if prefix == "" {
		prefix = "/cas/tickets"
	}
	if maxAge == 0 {
		maxAge = 3600
	}

	return &etcdTicketStore{
		cli:    client,
		ctx:    ctx,
		prefix: prefix,
		maxAge: maxAge,
	}, nil
}

// etcdTicketStore implements the TicketStore interface storing ticket data in etcd.
type etcdTicketStore struct {
	cli    *clientv3.Client
	ctx    context.Context
	prefix string
	maxAge int64
}

// Read returns the AuthenticationResponse for a ticket
func (s *etcdTicketStore) Read(id string) (*cas.AuthenticationResponse, error) {
	key := s.prefix + "/" + id
	resp, err := s.cli.Get(s.ctx, key)
	if err != nil {
		return nil, cas.ErrInvalidTicket
	}
	if resp.Count == 0 {
		return nil, cas.ErrInvalidTicket
	}

	var rsp *cas.AuthenticationResponse
	err = json.Unmarshal(resp.Kvs[0].Value, &rsp)
	if err != nil {
		return nil, cas.ErrInvalidTicket
	}

	return rsp, nil
}

// Write stores the AuthenticationResponse for a ticket
func (s *etcdTicketStore) Write(id string, ticket *cas.AuthenticationResponse) error {
	key := s.prefix + "/" + id
	grant, err := s.cli.Grant(s.ctx, s.maxAge+1)
	if err != nil {
		return err
	}
	data, err := json.Marshal(ticket)
	if err != nil {
		return err
	}
	_, err = s.cli.Put(s.ctx, key, string(data), clientv3.WithLease(grant.ID))
	return err
}

// Delete removes the AuthenticationResponse for a ticket
func (s *etcdTicketStore) Delete(id string) error {
	key := s.prefix + "/" + id
	resp, err := s.cli.Delete(s.ctx, key)
	if err != nil {
		return err
	}

	if resp.Deleted == 0 {
		return fmt.Errorf("key: %s is not found in etcd", key)
	}

	return nil
}

// Clear removes all ticket data
func (s *etcdTicketStore) Clear() error {
	key := s.prefix + "/"
	resp, err := s.cli.Delete(s.ctx, key, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	if resp.Deleted == 0 {
		return fmt.Errorf("key: %s is not found in etcd", key)
	}

	return nil
}
