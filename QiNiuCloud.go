package qiniu

import (
	"context"
	"io"
	"log"
	"mime/multipart"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/viper"
)

var (
	cfg         storage.Config
	AccessKey   string
	SecretKey   string
	Bucket      string
	QiniuServer string
)

func LoadQiniu() {
	AccessKey = viper.GetString("Objectstorage.AccessKey")
	SecretKey = viper.GetString("Objectstorage.SecretKey")
	Bucket = viper.GetString("Objectstorage.Bucket")
	QiniuServer = viper.GetString("Objectstorage.QiniuServer")
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
}
func SignPrivate(key string) string {
	mac := qbox.NewMac(AccessKey, SecretKey)
	deadline := time.Now().Add(time.Second * 86400).Unix() //1小时有效期
	privateAccessURL := storage.MakePrivateURL(mac, QiniuServer, key, deadline)
	return privateAccessURL
}

func SignPublic(key string) string {
	publicAccessURL := storage.MakePublicURL(Bucket, key)
	return publicAccessURL
}

func QiniuUpload(key string, file_size int64, file io.Reader) bool {
	mac := qbox.NewMac(AccessKey, SecretKey)
	putPolicy := storage.PutPolicy{
		Scope: Bucket,
	}
	upToken := putPolicy.UploadToken(mac)
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	log.Printf("upload qiniu cloud key:%s", key)
	err := formUploader.Put(context.Background(), &ret, upToken, key, file, file_size, &putExtra)

	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func QiniuUploadFile(file_path, key string) bool {
	putPolicy := storage.PutPolicy{
		Scope: Bucket,
	}

	mac := qbox.NewMac(AccessKey, SecretKey)
	upToken := putPolicy.UploadToken(mac)

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	log.Printf("%s upload qiniu cloud key:%s", file_path, key)
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, file_path, &putExtra)

	if err != nil {
		log.Fatal(err)
		return false
	}

	return true
}

func Upload(file multipart.File, Key string, file_size int64) string {
	if QiniuUpload(Key, file_size, file) {
		return Key
	} else {
		return ""
	}
}

func UploadFile(file, Key string) string {
	if QiniuUploadFile(file, Key) {
		return Key
	} else {
		return ""
	}
}

func GetUserKey(userid, other string) string {
	return "FileCloud/" + userid + "/" + other
}
