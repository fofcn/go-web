package file

import (
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
)

func InitRouterFile(public *gin.RouterGroup) {
	filerouter := NewFileRouter()
	public.POST("/file", filerouter.createPresignedUrl)
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
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	ossClient, err := oss.New(os.Getenv("OSS_ENDPOINT"), os.Getenv("OSS_ACCESS_KEY_ID"), os.Getenv("OSS_ACCESS_KEY_SECRET"), oss.SetCredentialsProvider(&provider), oss.UseCname(true))
	if err != nil {
		panic("create oss client error")
	}

	ossBucket, err := ossClient.Bucket(os.Getenv("OSS_BUCKET"))
	if err != nil {
		panic("check oss bucket error")
	}

	return &FileRouter{
		ossClient: ossClient,
		ossBucket: ossBucket,
	}
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
