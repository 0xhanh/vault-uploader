package cloud

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"vault-uploader/config"

	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
)

func UploadS3(streamKey string, fileName string, config *config.Config) (bool, bool, error) {

	// timestamp_microseconds_instanceName_regionCoordinates_numberOfChanges_token
	// 1564859471_6-474162_oprit_577-283-727-375_1153_27.mp4
	// - Timestamp
	// - Size + - + microseconds
	// - device
	// - Region
	// - Number of changes
	// - Token

	if config.S3 == nil {
		errorMessage := "UploadS3: Uploading Failed, as no settings found"
		log.Error(errorMessage)
		return false, false, errors.New(errorMessage)
	}

	// Legacy support, should get rid of it!
	aws_access_key_id := config.S3.Publickey
	aws_secret_access_key := config.S3.Secretkey
	aws_region := config.S3.Region

	// This is the new way ;)
	if config.HubKey != "" {
		aws_access_key_id = config.HubKey
	}
	if config.HubPrivateKey != "" {
		aws_secret_access_key = config.HubPrivateKey
	}

	// Check if we have some credentials otherwise we abort the request.
	if aws_access_key_id == "" || aws_secret_access_key == "" {
		errorMessage := "UploadS3: Uploading Failed, as no credentials found"
		log.Error(errorMessage)
		return false, false, errors.New(errorMessage)
	}

	s3Client, err := minio.NewWithRegion("s3.amazonaws.com", aws_access_key_id, aws_secret_access_key, true, aws_region)
	if err != nil {
		errorMessage := "UploadS3: " + err.Error()
		log.Error(errorMessage)
		return false, true, errors.New(errorMessage)
	}

	// Check if we need to use the proxy.
	if config.S3.ProxyURI != "" {
		var transport http.RoundTripper = &http.Transport{
			Proxy: func(*http.Request) (*url.URL, error) {
				return url.Parse(config.S3.ProxyURI)
			},
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		s3Client.SetCustomTransport(transport)
	}

	fileParts := strings.Split(fileName, "_")
	if len(fileParts) == 1 {
		errorMessage := "UploadS3: " + fileName + " is not a valid name."
		log.Error(errorMessage)
		return false, true, errors.New(errorMessage)
	}

	deviceKey := streamKey
	startRecording, _ := strconv.ParseInt(fileParts[0], 10, 64)
	devicename := fileParts[2]
	coordinates := fileParts[3]
	//numberOfChanges := fileParts[4]
	token, _ := strconv.Atoi(fileParts[5])

	log.Info("UploadS3: Upload started for " + fileName)
	fullname := fileName

	file, err := os.OpenFile(fullname, os.O_RDWR, 0755)
	if file != nil {
		defer file.Close()
	}

	if err != nil {
		errorMessage := "UploadS3: " + err.Error()
		log.Error(errorMessage)
		return false, true, errors.New(errorMessage)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		errorMessage := "UploadS3: " + err.Error()
		log.Error(errorMessage)
		return false, true, errors.New(errorMessage)
	}

	fileName = filepath.Base(fullname)
	fileName, err = ToKerberosFormat(streamKey, fileName)
	if err != nil {
		log.Warn("Malfunction in parsing kerperos format: " + err.Error())
	}

	n, err := s3Client.PutObject(config.S3.Bucket,
		config.S3.Username+"/"+fileName,
		file,
		fileInfo.Size(),
		minio.PutObjectOptions{
			ContentType:  "video/mp4",
			StorageClass: "ONEZONE_IA",
			UserMetadata: map[string]string{
				"event-timestamp":         strconv.FormatInt(startRecording, 10),
				"event-microseconds":      deviceKey,
				"event-instancename":      devicename,
				"event-regioncoordinates": coordinates,
				"event-numberofchanges":   deviceKey,
				"event-token":             strconv.Itoa(token),
				"productid":               deviceKey,
				"publickey":               aws_access_key_id,
				"uploadtime":              "now",
			},
		})

	if err != nil {
		errorMessage := "UploadS3: Uploading Failed, " + err.Error()
		log.Error(errorMessage)
		return false, true, errors.New(errorMessage)
	} else {
		log.Info("UploadS3: Upload Finished, file has been uploaded to bucket: " + strconv.FormatInt(n, 10))
		return true, true, nil
	}
}
