package database

import (
	"fmt"
	"testing"
	"time"
)

// for id insert
type ReflectTestStruct1 struct {
	Id    int64     `db:"id"`
	A     int       `db:"a"`
	B     int       `db:"-"`
	T     time.Time `db:"time"`
	Slice []byte    `db:"data"`
	Byte  byte      `db:"byte"`
	DBStr DBData    `db:"dbdata"`
	C     string
}

// for autoincrement
type ReflectTestStruct2 struct {
	Id int64 `db:"id,autoincrement"`
	A  int   `db:"a"`
	B  int   `db:"-"`
	C  string
}

func (r *ReflectTestStruct2) SetLastInsertId(id int64, err error) {
	if err != nil {
		panic(err)
	}
	r.Id = id
}

type ReflectTestStruct3 struct {
	ReflectTestStruct1
	D string `db:"d"`
}
type ReflectTestStruct4 struct {
	ReflectTestStruct3
	*ReflectTestStruct2
	E string `db:"e"`
}

func TestReflect(t *testing.T) {
	s1 := &ReflectTestStruct1{
		Id:    1,
		A:     100,
		B:     200,
		C:     "testing",
		Slice: []byte("abc"),
	}
	names, inputs, vals, err := reflectInsertStruct(s1, "mysql")
	if err != nil {
		t.Fatal(err)
	}
	if names != "id,a,time,data,byte,dbdata,C" {
		t.Fatal(names)
	}
	if inputs != "?,?,?,?,?,?,?" {
		t.Fatal(inputs)
	}
	if len(vals) != 7 {
		t.Fatal(fmt.Printf("%+v\n", vals))
	}

	s2 := &ReflectTestStruct2{
		A: 101,
		B: 201,
		C: "testing1",
	}
	names, inputs, vals, err = reflectInsertStruct(s2, "mysql")
	if err != nil {
		t.Fatal(err)
	}
	if names != "a,C" {
		t.Fatal(names)
	}
	if inputs != "?,?" {
		t.Fatal(inputs)
	}
	if len(vals) != 2 {
		t.Fatal(fmt.Printf("%+v\n", vals))
	}

	s4 := &ReflectTestStruct4{
		ReflectTestStruct2: s2,
		ReflectTestStruct3: ReflectTestStruct3{
			ReflectTestStruct1: *s1,
			D:                  "d",
		},
		E: "e",
	}
	names, inputs, vals, err = reflectInsertStruct(s4, "oracle")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", names)
	fmt.Printf("%+v\n", inputs)
	fmt.Printf("%+v\n", vals)
}
