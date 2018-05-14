package platform

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

type registryLocationCache struct {
	url string
	credential  credential
	fs    boshsys.FileSystem
	path string
}

type credential struct {
	username string `json:"username"`
	password string `json:"password"`
}

func NewRegistryLocationCache(fs boshsys.FileSystem, url, username, password string, path string) (*registryLocationCache, error) {
	credential := credential{username: username, password: password}

	registryLocationCache:= registryLocationCache{url:url, credential:credential, fs:fs, path: path}

	return &registryLocationCache, nil
}

func (s *registryLocationCache) SaveRegistryLocationCache() (err error) {
	jsonRegistryLocationCache, err := json.Marshal(*s)
	if err != nil {
		return bosherr.WrapError(err, "Marshalling bootstrap state")
	}

	err = s.fs.WriteFile(s.path, jsonRegistryLocationCache)
	if err != nil {
		return bosherr.WrapError(err, "Writing bootstrap state to file")
	}

	return
}
