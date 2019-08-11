package alarm

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
	fmt.Println(RemoveReceiver(arr, 0))
	fmt.Println(RemoveReceiver(arr, 1))
	fmt.Println(RemoveReceiver(arr, 2))
}
