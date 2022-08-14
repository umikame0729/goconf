package goconf

import (
	"fmt"
	"testing"
	"time"
)

type TestConfig struct {
	Name string
}

func TestLoad(t *testing.T) {
	fname := fmt.Sprintf("%d.json", time.Now().UnixMilli())
	conf := &Config[TestConfig]{
		Version: 1,
		Path:    fname,
		Data:    &TestConfig{Name: "Test"},
	}

	if err := conf.Load(); IsNewCreated(err) {
		t.Log(err)
	} else {
		t.Fail()
	}

	conf = &Config[TestConfig]{
		Version: 2,
		Path:    fname,
	}

	if err := conf.Load(); IsUpdated(err) {
		t.Log(err)
	} else {
		t.Fail()
	}

	conf = &Config[TestConfig]{
		Version: 2,
		Path:    fname,
	}

	if err := conf.Load(); err == nil {
		t.Log("PASS")
	} else {
		t.Error(err)
	}
}
