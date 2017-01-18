package ttlmap

import (
	"math"
	"testing"
	"time"
)

func TestNewItemWithExpiration(t *testing.T) {
	ttl := 1 * time.Second
	expiration := time.Now().Add(ttl)
	item := NewItem("foo", WithExpiration(expiration))
	if item.Value() != "foo" {
		t.Fatalf("Invalid value")
	}
	diff := ttl - item.TTL()
	if diff > 10*time.Millisecond {
		t.Fatalf("Invalid TTL")
	}
	if item.Expiration() != expiration {
		t.Fatalf("Invalid expiration")
	}
	if item.Expired() {
		t.Fatalf("Expecting not expired")
	}
	<-time.After(time.Duration(float64(ttl) * 0.8))
	if item.Expired() {
		t.Fatalf("Expecting not expired")
	}
	<-time.After(time.Duration(float64(ttl) * 0.4))
	if !item.Expired() {
		t.Fatalf("Expecting expired")
	}
	if !item.Expires() {
		t.Fatalf("Expecting expires")
	}
}

func TestNewItemWithTTL(t *testing.T) {
	ttl := 10 * time.Second
	item := NewItem("foo", WithTTL(ttl))
	if item.Value() != "foo" {
		t.Fatalf("Invalid value")
	}
	expectedExpiration := time.Now().Add(ttl)
	diff := item.Expiration().Sub(expectedExpiration)
	if diff > 10*time.Millisecond {
		t.Fatalf("Invalid expiration")
	}
	if !item.Expires() {
		t.Fatalf("Expecting expires")
	}
}

func TestNewItemWithoutExpiration(t *testing.T) {
	item := NewItem("foo", nil)
	if item.Value() != "foo" {
		t.Fatalf("Invalid value")
	}
	if item.TTL() != time.Duration(math.MaxInt64) {
		t.Fatalf("Not expecting TTL")
	}
	time.Sleep(1 * time.Second)
	if item.TTL() != time.Duration(math.MaxInt64) {
		t.Fatalf("Not expecting TTL")
	}
	if item.Expired() {
		t.Fatalf("Not expecting expired")
	}
	if item.Expires() {
		t.Fatalf("Not expecting expires")
	}
	if !item.Expiration().IsZero() {
		t.Fatalf("Not expecting expiration")
	}
}
