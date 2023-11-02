package cloud

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"vault-uploader/config"

	log "github.com/sirupsen/logrus"
)

func UploadKerberosVault(streamKey string, fileName string, config *config.Config) (bool, bool, error) {

	if config.KStorage.AccessKey == "" ||
		config.KStorage.SecretAccessKey == "" ||
		config.KStorage.Directory == "" ||
		config.KStorage.URI == "" {
		err := "UploadKerberosVault: Kerberos Vault not properly configured"
		log.Info(err)
		return false, false, errors.New(err)
	}

	// KerberosCloud, this means storage is disabled and proxy enabled.
	log.Info("UploadKerberosVault: Uploading to Kerberos Vault (" + config.KStorage.URI + ")")
	log.Info("UploadKerberosVault: Upload started for " + fileName)
	fullname := fileName

	file, err := os.OpenFile(fullname, os.O_RDWR, 0755)
	if file != nil {
		defer file.Close()
	}
	if err != nil {
		err := "UploadKerberosVault: Upload Failed, file doesn't exists anymore"
		log.Info(err)
		return false, false, errors.New(err)
	}

	fileName = filepath.Base(fullname)
	fileName, err = ToKerberosFormat(streamKey, fileName)
	if err != nil {
		log.Warn("Malfunction in parsing kerperos format: " + err.Error())
	}

	publicKey := config.KStorage.CloudKey
	// This is the new way ;)
	if config.HubKey != "" {
		publicKey = config.HubKey
	}

	req, err := http.NewRequest("POST", config.KStorage.URI+"/storage", file)
	if err != nil {
		errorMessage := "UploadKerberosVault: error reading request, " + config.KStorage.URI + "/storage: " + err.Error()
		log.Error(errorMessage)
		return false, true, errors.New(errorMessage)
	}
	req.Header.Set("Content-Type", "video/mp4")
	req.Header.Set("X-Kerberos-Storage-CloudKey", publicKey)
	req.Header.Set("X-Kerberos-Storage-AccessKey", config.KStorage.AccessKey)
	req.Header.Set("X-Kerberos-Storage-SecretAccessKey", config.KStorage.SecretAccessKey)
	req.Header.Set("X-Kerberos-Storage-Provider", config.KStorage.Provider)
	req.Header.Set("X-Kerberos-Storage-FileName", fileName)
	req.Header.Set("X-Kerberos-Storage-Device", streamKey)
	req.Header.Set("X-Kerberos-Storage-Capture", "IPCamera")
	req.Header.Set("X-Kerberos-Storage-Directory", config.KStorage.Directory)

	var client *http.Client
	if os.Getenv("AGENT_TLS_INSECURE") == "true" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	} else {
		client = &http.Client{}
	}

	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err == nil {
		if resp != nil {
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				if resp.StatusCode == 200 {
					log.Info("UploadKerberosVault: Upload Finished, " + resp.Status + ", " + string(body))
					return true, true, nil
				} else {
					log.Info("UploadKerberosVault: Upload Failed, " + resp.Status + ", " + string(body))
					return false, true, nil
				}
			}
		}
	}

	errorMessage := "UploadKerberosVault: Upload Failed, " + err.Error()
	log.Info(errorMessage)
	return false, true, errors.New(errorMessage)
}
