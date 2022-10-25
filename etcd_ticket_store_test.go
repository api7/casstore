package casstore

import (
	"context"
	"reflect"
	"testing"

	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/cas.v2"
)

func TestEtcdTicketStore(t *testing.T) {
	user1 := &cas.AuthenticationResponse{User: "user1"}
	user2 := &cas.AuthenticationResponse{User: "user2"}
	store, err := NewEtcdTicketStore(clientv3.Config{Endpoints: []string{_defaultEtcd}},
		context.Background(), "/cas/tickets", 3600)

	if err := store.Write("user1", user1); err != nil {
		t.Errorf("Expected store.Write(user1) to succeed, got error: %v", err)
	}

	if err := store.Write("user2", user2); err != nil {
		t.Errorf("Expected store.Write(user2) to succeed, got error: %v", err)
	}

	ar, err := store.Read("user2")
	if err != nil {
		t.Errorf("Expected store.Read(user2) to succeed, got error: %v", err)
	}

	if !reflect.DeepEqual(*ar, *user2) {
		t.Errorf("Expected retrieved AuthenticationResponse to be %v, got %v", *user2, *ar)
	}

	if err := store.Delete("user2"); err != nil {
		t.Errorf("Error while deleting user2, got %v", err)
	}

	if _, err := store.Read("user2"); err != cas.ErrInvalidTicket {
		t.Errorf("Expected store.Read(user2) to fail")
	}

	if err := store.Clear(); err != nil {
		t.Errorf("Expected store.Clear() to succeed, got error: %v", err)
	}

	_, err = store.Read("user1")
	if err == nil {
		t.Errorf("Expected an error from store.Read(user1), got nil")
	}

	if err != cas.ErrInvalidTicket {
		t.Errorf("Expected ErrInvalidTicket from store.Read(user1), got %v", err)
	}
}
