package logrus_papertrail

import (
	"fmt"
	"testing"

	"github.com/stvp/go-udp-testing"
	"sito/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

func TestWritingToUDP(t *testing.T) {
	port := 16661
	udp.SetAddr(fmt.Sprintf(":%d", port))

	hook, err := NewPapertrailHook("localhost", port, "test")
	if err != nil {
		t.Errorf("Unable to connect to local UDP server.")
	}

	log := logrus.New()
	log.Hooks.Add(hook)

	udp.ShouldReceive(t, "foo", func() {
		log.Info("foo")
	})
}
