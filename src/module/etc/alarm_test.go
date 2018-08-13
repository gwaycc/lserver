package etc

import (
	"fmt"
	"testing"
)

func TestRemoveReceiver(t *testing.T) {
	arr := []*Receiver{
		&Receiver{},
		&Receiver{},
		&Receiver{},
	}
	fmt.Println(Remove(arr, 0))
	fmt.Println(Remove(arr, 1))
	fmt.Println(Remove(arr, 2))
}
