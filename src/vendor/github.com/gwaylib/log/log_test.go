package log

import (
	"log"
	"testing"
)

func TestApi(t *testing.T) {
	Debug("echo debug")
	Debugf("echo debug of int:%d", 1)

	Info("echo info")
	Infof("echo info of int:%d", 1)

	Warn("echo warn")
	Warnf("echo warn of int:%d", 1)

	Error("echo error")
	Errorf("echo error of int:%d", 1)

	/*
	   Fatal("echo fatal")
	   Fatalf("echo fatal of int:%d", 1)
	*/
	log.Println("testing syslog")
	log.Printf("testing syslog:%d", 1)
	log.Printf("testing syslog:%d", 2)
}
