package healthcheck_test

import (
	"testing"

	"github.com/yeqown/infrastructure/healthcheck"
)

func Test_CheckerOverTCP(t *testing.T) {
	c := healthcheck.NewTCPChecker("localhost:27017")
	info := c.Check()
	if !info.Healthy {
		t.Error("checking your 27017 port")
		t.FailNow()
	}

	c2 := healthcheck.NewTCPChecker("localhost:1111")
	info2 := c2.Check()
	if info2.Healthy {
		t.Error("checking your 1111 port, we assume there should no service")
		t.FailNow()
	}
}
