package goconf

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type ErrorIsUpdated struct {
	Path    string
	BakPath string
}

func (e *ErrorIsUpdated) Error() string {
	return fmt.Sprintf("`%s` updated, backup at `%s`", e.Path, e.BakPath)
}

func IsUpdated(err error) bool {
	_, ok := err.(*ErrorIsUpdated)
	return ok
}

type ErrorIsNewCreated struct {
	Err error
}

func (e *ErrorIsNewCreated) Error() string {
	return e.Err.Error()
}

func IsNewCreated(err error) bool {
	_, ok := err.(*ErrorIsNewCreated)
	return ok
}

// Config definition
type Config[T interface{}] struct {
	is_load bool
	Path    string // json file path (with filename and extension)
	Version int    // config int
	Data    *T     // Data
	Content struct {
		Version int // config version
		Data    *T  // config data payload
	}
}

func (c *Config[T]) Save() error {
	c.Content.Version = c.Version
	c.Content.Data = c.Data
	if data, err := json.Marshal(c.Content); err != nil {
		return err
	} else if err := os.WriteFile(c.Path, data, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (c *Config[T]) IsExist() bool {
	if _, err := os.Stat(c.Path); os.IsExist(err) {
		return true
	} else {
		return false
	}
}

func (c *Config[T]) createIfNotExist() error {
	if c.IsExist() {
		return nil
	} else if err := c.Save(); err != nil {
		return err
	}
	return &ErrorIsNewCreated{
		Err: fmt.Errorf("`%s` is new file", c.Path),
	}
}

func (c *Config[T]) bakName() string {
	now := time.Now()
	return fmt.Sprintf("%s.%d-%d-%d-%d-%d-%d", c.Path, now.Local().Year(), now.Local().Month(), now.Local().Day(), now.Local().Hour(), now.Local().Minute(), now.Local().Second())
}

func (c *Config[T]) versionCheck() error {
	bakPath := c.bakName()
	if c.Content.Version == c.Version {
		return nil
	} else if !c.IsExist() {
		return fmt.Errorf("`%s` not exist", c.Path)
	} else if err := os.Rename(c.Path, bakPath); err != nil {
		return err
	} else if err := c.Save(); err != nil {
		return err
	}
	return &ErrorIsUpdated{
		Path:    c.Path,
		BakPath: bakPath,
	}
}

// bool - is version checked OK, if false, need rerun program
func (c *Config[T]) Load() error {
	if c.Content.Data == nil {
		c.Content.Data = c.Data
	}

	if c.is_load {
		return nil
	} else if err := c.createIfNotExist(); err != nil {
		return err
	} else if data, err := os.ReadFile(c.Path); err != nil {
		return err
	} else if err := json.Unmarshal(data, &c.Content); err != nil {
		return err
	} else if err := c.versionCheck(); err != nil {
		return err
	}

	c.is_load = true
	return nil
}
