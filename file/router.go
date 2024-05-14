package file

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"go-web/pkg/global"
	"io"
	"os"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
)

const (
	OSS_ENDPOINT          = "OSS_ENDPOINT"
	OSS_ACCESS_KEY_ID     = "OSS_ACCESS_KEY_ID"
	OSS_ACCESS_KEY_SECRET = "OSS_ACCESS_KEY_SECRET"
	OSS_BUCKET            = "OSS_BUCKET"
)

func InitRouterFile(public *gin.RouterGroup) {
	filerouter := NewFileRouter()
	public.POST("/file", filerouter.createUploadToken)
	public.GET("/file", filerouter.getDownloadUrl)
}

type FileRouter struct {
	ossClient *oss.Client
	ossBucket *oss.Bucket
}

func NewFileRouter() *FileRouter {
	// 从环境变量中获取临时访问凭证。运行本代码示例之前，
	// 请确保已设置环境变量OSS_ACCESS_KEY_ID、OSS_ACCESS_KEY_SECRET、OSS_SESSION_TOKEN。
	// 参考： https://help.aliyun.com/zh/oss/user-guide/authorized-third-party-upload?spm=a2c4g.11186623.0.0.996b6f4f8obfvf#261193c152gdf
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		panic("oss new environment variable credentials provider error")
	}

	ossClient, err := oss.New(os.Getenv(OSS_ENDPOINT), os.Getenv(OSS_ACCESS_KEY_ID), os.Getenv(OSS_ACCESS_KEY_SECRET), oss.SetCredentialsProvider(&provider), oss.UseCname(true))
	if err != nil {
		panic("create oss client error")
	}

	ossBucket, err := ossClient.Bucket(os.Getenv(OSS_BUCKET))
	if err != nil {
		panic("check oss bucket error")
	}

	return &FileRouter{
		ossClient: ossClient,
		ossBucket: ossBucket,
	}
}

var (
	// 指定上传到OSS的文件前缀。
	uploadDir = "user-dir-prefix/"
	// 指定过期时间，单位为秒。
	expireTime = int64(3600)
)

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}
type PolicyToken struct {
	AccessKeyId string `json:"ossAccessKeyId"`
	Host        string `json:"host"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
}

func getGMTISO8601(expireEnd int64) string {
	return time.Unix(expireEnd, 0).UTC().Format("2006-01-02T15:04:05Z")
}

func (cr *FileRouter) createUploadToken(c *gin.Context) {
	now := time.Now().Unix()
	expireEnd := now + expireTime
	tokenExpire := getGMTISO8601(expireEnd)
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, uploadDir)
	config.Conditions = append(config.Conditions, condition)
	result, err := json.Marshal(config)
	if err != nil {
		global.InteralServerErrorWithMsg(c, "json")
		return
	}
	encodedResult := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(sha1.New, []byte(os.Getenv("OSS_ACCESS_KEY_SECRET")))
	io.WriteString(h, encodedResult)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))
	policyToken := PolicyToken{
		AccessKeyId: os.Getenv("OSS_ACCESS_KEY_SECRET"),
		Host:        os.Getenv("OSS_HOST"),
		Signature:   signedStr,
		Policy:      encodedResult,
		Directory:   uploadDir,
	}

	global.SuccessWithData(c, policyToken)
}

func (cr *FileRouter) createPresignedUrl(c *gin.Context) {
	var cmd OssPresignedUrlCmd
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	options := []oss.Option{
		oss.ContentType(cmd.ContentType),
	}

	signedURL, err := cr.ossBucket.SignURL(cmd.FileName, oss.HTTPPut, 600, options...)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, OssPresignedUrlDto{
		Url: signedURL,
	})
}

func (cr *FileRouter) getDownloadUrl(c *gin.Context) {
	filename := c.Query("file_name")
	if len(filename) == 0 {
		c.JSON(400, gin.H{"error": "file_name is empty"})
		return
	}

	signedURL, err := cr.ossBucket.SignURL(filename, oss.HTTPGet, 600000)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, OssDownloadUrlDto{
		Url: signedURL,
	})
}
