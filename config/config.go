package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	HubURI            string `yaml:"hub_uri" json:"hub_uri" bson:"hub_uri"`
	HubKey            string `yaml:"hub_key" json:"hub_key" bson:"hub_key"`
	HubPrivateKey     string `yaml:"hub_private_key" json:"hub_private_key" bson:"hub_private_key"`
	HubSite           string `yaml:"hub_site" json:"hub_site" bson:"hub_site"`
	Key               string `yaml:"key" json:"key"`
	RemoveAfterUpload bool   `yaml:"remove_after_upload" json:"remove_after_upload"`

	Cloud    string    `yaml:"cloud" json:"cloud" bson:"cloud"`
	S3       *S3       `yaml:"s3,omitempty" json:"s3,omitempty" bson:"s3,omitempty"`
	KStorage *KStorage `yaml:"kstorage,omitempty" json:"kstorage,omitempty" bson:"kstorage,omitempty"`
	Dropbox  *Dropbox  `yaml:"dropbox,omitempty" json:"dropbox,omitempty" bson:"dropbox,omitempty"`
}

// S3 integration
type S3 struct {
	Proxy     string `yaml:"proxy,omitempty" json:"proxy,omitempty" bson:"proxy,omitempty"`
	ProxyURI  string `yaml:"proxyuri,omitempty" json:"proxyuri,omitempty" bson:"proxyuri,omitempty"`
	Bucket    string `yaml:"bucket,omitempty" json:"bucket,omitempty" bson:"bucket,omitempty"`
	Region    string `yaml:"region,omitempty" json:"region,omitempty" bson:"region,omitempty"`
	Username  string `yaml:"username,omitempty" json:"username,omitempty" bson:"username,omitempty"`
	Publickey string `yaml:"publickey,omitempty" json:"publickey,omitempty" bson:"publickey,omitempty"`
	Secretkey string `yaml:"secretkey,omitempty" json:"secretkey,omitempty" bson:"secretkey,omitempty"`
}

// KStorage contains the credentials of the Kerberos Storage/Kerberos Cloud instance.
// By defining KStorage you can make your recordings available in the cloud, at a centrel place.
type KStorage struct {
	URI             string `yaml:"uri,omitempty" json:"uri,omitempty" bson:"uri,omitempty"`
	CloudKey        string `yaml:"cloud_key,omitempty" json:"cloud_key,omitempty" bson:"cloud_key,omitempty"` /* old way, remove this */
	AccessKey       string `yaml:"access_key,omitempty" json:"access_key,omitempty" bson:"access_key,omitempty"`
	SecretAccessKey string `yaml:"secret_access_key,omitempty" json:"secret_access_key,omitempty" bson:"secret_access_key,omitempty"`
	Provider        string `yaml:"provider,omitempty" json:"provider,omitempty" bson:"provider,omitempty"`
	Directory       string `yaml:"directory,omitempty" json:"directory,omitempty" bson:"directory,omitempty"`
}

// Dropbox integration
type Dropbox struct {
	AccessToken string `yaml:"access_token,omitempty" json:"access_token,omitempty" bson:"access_token,omitempty"`
	Directory   string `yaml:"directory,omitempty" json:"directory,omitempty" bson:"directory,omitempty"`
}

func ReadConf(filename string) (*Config, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil {
		return nil, err
	}
	buf, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	return c, err
}
