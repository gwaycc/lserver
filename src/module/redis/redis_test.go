package redis

import (
	"testing"
)

func TestRedis(t *testing.T) {
	r, err := NewRediStore(1, "tcp", "127.0.0.1:6379", "")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	id := 1
	if err := r.Set("testing", id, 60); err != nil {
		t.Fatal(err)
	}
	outId := 0
	if err := r.Scan("testing", &outId); err != nil {
		t.Fatal(err)
	}
	if id != outId {
		t.Fatal(id, outId)
	}
	if err := r.Delete("testing"); err != nil {
		t.Fatal(err)
	}
	if err := r.Scan("testing", &outId); err != nil {
		if err != ErrNil {
			t.Fatal(err)
		}
	} else {
		t.Fatal("found data", outId)
	}
}
