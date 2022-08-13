package goconf

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

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

func (c *Config[T]) IsExist() error {
	if _, err := os.Stat(c.Path); os.IsExist(err) {
		return nil
	} else {
		return err
	}
}

func (c *Config[T]) createIfNotExist() error {
	if err := c.IsExist(); err == nil {
		return nil
	} else if err := c.Save(); err != nil {
		return err
	}
	return nil
}

func (c *Config[T]) bakName() string {
	now := time.Now()
	return fmt.Sprintf("%s.%d-%d-%d-%d-%d-%d", c.Path, now.Local().Year(), now.Local().Month(), now.Local().Day(), now.Local().Hour(), now.Local().Minute(), now.Local().Second())
}

func (c *Config[T]) versionCheck() (bool, error) {
	if c.Content.Version == c.Version {
		return true, nil
	} else if err := c.IsExist(); err != nil {
		return false, err
	} else if err := os.Rename(c.Path, c.bakName()); err != nil {
		return false, err
	} else if err := c.Save(); err != nil {
		return false, err
	}
	return false, nil
}

// bool - is version checked OK, if false, need rerun program
func (c *Config[T]) Load() (bool, error) {
	c.Content.Data = c.Data
	if c.is_load {
		return true, nil
	} else if err := c.createIfNotExist(); err != nil {
		return false, err
	} else if data, err := os.ReadFile(c.Path); err != nil {
		return false, err
	} else if err := json.Unmarshal(data, &c.Content); err != nil {
		return false, err
	} else if update, err := c.versionCheck(); err != nil {
		return false, err
	} else if !update {
		return false, nil
	}
	c.is_load = true
	return true, nil
}
