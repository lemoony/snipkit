package cache

import (
	"path/filepath"

	"github.com/spf13/afero"
)

func (c *cacheImpl) PutData(key DataKey, data []byte) {
	path := c.cacheFilepath(key)
	c.system.CreatePath(path)
	c.system.WriteFile(path, data)
}

func (c *cacheImpl) GetData(key DataKey) ([]byte, bool) {
	bytes, err := afero.ReadFile(c.system.Fs, c.cacheFilepath(key))
	if err != nil {
		return nil, false
	}

	return bytes, true
}

func (c *cacheImpl) cacheFilepath(key DataKey) string {
	return filepath.Join(c.system.HomeDir(), ".cache", string(key))
}
