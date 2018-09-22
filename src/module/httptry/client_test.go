package httptry

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	go func() {
		count := DefaultTryTimes
		http.HandleFunc("/engine/http", func(w http.ResponseWriter, r *http.Request) {
			count--
			if count == 0 {
				w.Write([]byte("ok"))
			} else {
				http.NotFound(w, r)
			}
		})
		http.ListenAndServe(":8000", nil)
	}()

	time.Sleep(time.Duration(3 * 1e9))

	resp, err := Get("http://localhost:8000/engine/http")
	if err != nil {
		t.Error(err)
	} else {
		defer resp.Body.Close()
		if resp.Status != "200 OK" {
			t.Error(resp.Status)
		}
	}
}

func TestGetTimout(t *testing.T) {
	println("TestRespGet--------------------")
	client := NewClient(nil, 1, 0, 1)
	_, err := client.Get("http://szcrmsystem.gicp.net")
	if err != nil && !strings.Contains(err.Error(), "timeout") {
		t.Fatal(err)
	}
	if strings.Contains(err.Error(), "timeout") {
		println(err.Error())
	} else {
		t.Fatal("should be time out, but it is not happend")
	}
}
