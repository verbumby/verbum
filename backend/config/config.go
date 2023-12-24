package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"gopkg.in/yaml.v3"
)

type configType struct {
	HTTPS struct {
		Addr          string
		CanonicalAddr string `yaml:"canonicalAddr"`
		CertFile      string `yaml:"certFile"`
		KeyFile       string `yaml:"keyFile"`
	}

	HTTP struct {
		Addr              string
		AcmeChallengeRoot string `yaml:"acmeChallengeRoot"`
	}

	Elastic struct {
		Addr string
	}

	Images struct {
		Path string
	}

	Dicts struct {
		Repo struct {
			Path string
		}
	}
}

var config configType

func ReadConfig() error {
	config = configType{}
	config.HTTPS.Addr = ":8443"
	config.HTTPS.CertFile = "cert.pem"
	config.HTTPS.KeyFile = "key.pem"
	config.HTTPS.CanonicalAddr = "https://127.0.0.1:8443"
	config.Elastic.Addr = "http://127.0.0.1:9200"
	config.Images.Path = "./images"
	config.Dicts.Repo.Path = "../slouniki"

	atLeastOneFileRead := false
	files := []string{"config.yaml", "/usr/local/share/verbum/config.yaml"}
	for _, filename := range files {
		f, err := os.Open(filename)
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}
		defer f.Close()
		atLeastOneFileRead = true

		if err := yaml.NewDecoder(f).Decode(&config); err != nil {
			return fmt.Errorf("decode yaml: %w", err)
		}
	}

	if !atLeastOneFileRead {
		return fmt.Errorf("none of config files %v were read", files)
	}

	return nil
}

func HTTPSAddr() string {
	return config.HTTPS.Addr
}

func HTTPSCanonicalAddr() string {
	return config.HTTPS.CanonicalAddr
}

func HTTPSCertFile() string {
	return config.HTTPS.CertFile
}

func HTTPSKeyFile() string {
	return config.HTTPS.KeyFile
}

func HTTPAddr() string {
	return config.HTTP.Addr
}

func HTTPAcmeChallengeRoot() string {
	return config.HTTP.AcmeChallengeRoot
}

func ElasticAddr() string {
	return config.Elastic.Addr
}

func ImagesPath() string {
	return config.Images.Path
}

func DictsRepoPath() string {
	return config.Dicts.Repo.Path
}
