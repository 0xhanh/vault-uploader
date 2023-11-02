package cloud

import (
	"os"
	"vault-uploader/config"

	log "github.com/sirupsen/logrus"
)

func HandleUpload(streamKey string, filename string, config *config.Config) {

	log.Debug("HandleUpload: started")

	// Half a second delay between two uploads
	uploaded := false
	configured := false
	var err error
	if config.Cloud == "s3" || config.Cloud == "kerberoshub" {
		uploaded, configured, err = UploadKerberosHub(streamKey, filename, config)
	} else if config.Cloud == "kstorage" || config.Cloud == "kerberosvault" {
		uploaded, configured, err = UploadKerberosVault(streamKey, filename, config)
	} else if config.Cloud == "dropbox" {
		// Todo: implement dropbox upload
		// uploaded, configured, err = UploadDropbox(config, filename)
	} else if config.Cloud == "gdrive" {
		// Todo: implement gdrive upload
	} else if config.Cloud == "onedrive" {
		// Todo: implement onedrive upload
	} else if config.Cloud == "minio" {
		// Todo: implement minio upload
	} else if config.Cloud == "webdav" {
		// Todo: implement webdav upload
	} else if config.Cloud == "ftp" {
		// Todo: implement ftp upload
	} else if config.Cloud == "sftp" {
		// Todo: implement sftp upload
	} else if config.Cloud == "aws" {
		// Todo: need to be updated, was previously used for hub.
		uploaded, configured, err = UploadS3(streamKey, filename, config)
	} else if config.Cloud == "azure" {
		// Todo: implement azure upload
	} else if config.Cloud == "google" {
		// Todo: implement google upload
	}

	// Check if the file is uploaded, if so, remove it.
	if uploaded {
		// Check if we need to remove the original recording
		// removeAfterUpload is set to false by default
		if config.RemoveAfterUpload {
			err := os.Remove(filename)
			if err != nil {
				log.Error("error remove: " + err.Error())
			} else {
				log.Info("removed: " + filename)
			}
		}
	} else if !configured {
		// err := os.Remove(filename)
		// if err != nil {
		// 	log.Error("HandleUpload: " + err.Error())
		// }
	} else {
		if err != nil {
			log.Error("HandleUpload: " + err.Error())
		}
	}

	log.Debug("HandleUpload: finished")
}
