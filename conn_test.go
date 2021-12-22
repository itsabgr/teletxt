package teletxt

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"github.com/itsabgr/go-handy"
	"github.com/itsabgr/teletxt-go/pkg/mockconn"
	"io"
	"testing"
)

func RandHex(n uint) []byte {
	b := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, b)
	handy.Throw(err)
	return []byte(hex.EncodeToString(b))
}

func TestConn(t *testing.T) {
	mock := mockconn.NewConn()
	defer mock.Close()
	a := NewConn(mock.Client)
	b := NewConn(mock.Server)
	for range handy.N(10) {
		k, v := RandHex(10), RandHex(10)
		go func() {
			_, err := a.WriteKV(k, v)
			if err != nil {
				t.Fatal(err)
			}
		}()
		firstByte, err := b.ReadByte()
		if err != nil {
			t.Fatal(err)
		}
		if k[0] != firstByte {
			t.Fatalf("expected first byte  %+v got %+v", k[0], firstByte)
		}
		err = b.UnreadByte()
		if err != nil {
			t.Fatal(err)
		}
		k1, v1, err := b.ReadKV()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(k1, k) {
			t.Fatalf("expected key %+v got %+v", k, k1)
		}
		if !bytes.Equal(v1, v) {
			t.Fatalf("expected value %+v got %+v", v, v1)
		}
	}
	for range handy.N(10) {

		v := RandHex(10)
		go func() {
			_, err := a.WriteValue(v)
			if err != nil {
				t.Fatal(err)
			}
		}()
		firstByte, err := b.ReadByte()
		if err != nil {
			t.Fatal(err)
		}
		if v[0] != firstByte {
			t.Fatalf("expected first byte  %+v got %+v", v[0], firstByte)
		}
		err = b.UnreadByte()
		if err != nil {
			t.Fatal(err)
		}
		v1, empty, err := b.ReadKV()
		if err != nil {
			t.Fatal(err)
		}
		if empty != nil {
			t.Fatalf("expected single got %+v", empty)
		}
		if !bytes.Equal(v1, v) {
			t.Fatalf("expected value %+v got %+v", v, v1)
		}
	}
	for range handy.N(10) {
		go func() {
			_, err := a.WriteEmptyLine()
			if err != nil {
				t.Fatal(err)
			}
		}()
		v, k, err := b.ReadKV()
		if err != nil {
			t.Fatal(err)
		}
		if k != nil {
			t.Fatalf("expected empty key got %+v", k)
		}
		if v != nil {
			t.Fatalf("expected empty value got %+v", v)
		}

	}
}
