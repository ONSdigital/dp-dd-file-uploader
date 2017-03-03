package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/ONSdigital/dp-dd-file-uploader/config"
	"github.com/ONSdigital/go-ns/handlers/response"
	"net/http"
	"time"
)

type S3Credentials struct {
	EndpointUrl string   `json:"endpoint_url"`
	Params      S3Params `json:"params"`
}

type S3Params struct {
	Key                 string `json:"key"`
	ACL                 string `json:"acl"`
	SuccessActionStatus string `json:"success_action_status"`
	Policy              string `json:"policy"`
	Algorithm           string `json:"x-amz-algorithm"`
	Credential          string `json:"x-amz-credential"`
	Date                string `json:"x-amz-date"`
	Signature           string `json:"x-amz-signature"`
}

type S3Policy struct {
	Expiration string         `json:"expiration"`
	Conditions []interface{}  `json:"conditions"`
}

func GetUploadCredentials(w http.ResponseWriter, req *http.Request) {
	filename := req.URL.Query().Get("filename")
	if filename == "" {
		w.WriteHeader(400)
		w.Write([]byte("Bad request"))
		return
	}
	creds := credentials(filename)
	response.WriteJSON(w, creds, 200)
}

func credentials(filename string) S3Credentials {
	s3Bucket := config.S3URL.Host
	return S3Credentials{
		EndpointUrl: "https://" + s3Bucket + ".s3.amazonaws.com",
		Params:      params(filename),
	}
}

func params(filename string) S3Params {
	creds := credential()
	dateStr := toISO8601(time.Now())
	policyBase64 := policy(filename, creds, dateStr)
	return S3Params{
		Key:                 filename,
		ACL:                 "public-read",
		SuccessActionStatus: "201",
		Policy:              policyBase64,
		Algorithm:           "AWS4-HMAC-SHA256",
		Credential:          creds,
		Date:                dateStr,
		Signature:           signature(dateStr, policyBase64),
	}
}

func policy(filename string, credential string, dateStr string) string {
	policy := S3Policy{
		Expiration: toISO8601_v2(time.Now().Add(5 * time.Minute)),
		Conditions: []interface{}{
			map[string]interface{}{"bucket": config.S3URL.Host},
			map[string]interface{}{"key": filename},
			map[string]interface{}{"acl": "public-read"},
			map[string]interface{}{"success_action_status": "201"},
			//map[string]interface{}{"Content-Type": "text/csv"},
			map[string]interface{}{"x-amz-algorithm": "AWS4-HMAC-SHA256"},
			map[string]interface{}{"x-amz-credential": credential},
			map[string]interface{}{"x-amz-date": dateStr},
		},
	}
	jsonData, err := json.Marshal(policy)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(jsonData)
}

func credential() string {
	return config.AWSAccessKey + "/" + toISO8601(time.Now())[:8] + "/" + config.AWSRegion + "/s3/aws4_request"
}

func toISO8601(t time.Time) string {
	return t.UTC().Format("20060102T150405Z0700")
}

func toISO8601_v2(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05Z0700")
}


func hmacsha256(key []byte, msg []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	return mac.Sum(nil)
}

func signature(dateStr string, policyBase64 string) string {
	// Build up AWS v4 signature from iterated HMAC-SHA-256
	dateKey := hmacsha256([]byte("AWS4"+config.AWSSecretKey), []byte(dateStr[:8]))
	dateRegionKey := hmacsha256(dateKey, []byte(config.AWSRegion))
	dateRegionServiceKey := hmacsha256(dateRegionKey, []byte("s3"))
	signingKey := hmacsha256(dateRegionServiceKey, []byte("aws4_request"))
	return hex.EncodeToString(hmacsha256(signingKey, []byte(policyBase64)))
}
