package filesystem

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var _ IFilesystem = (*MinioFilesystem)(nil)

type MinioFilesystem struct {
	core   *minio.Core
	config MinioSystemConfig
}

func NewMinioFilesystem(config MinioSystemConfig) *MinioFilesystem {
	client, err := minio.NewCore(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.SecretId, config.SecretKey, ""),
		Secure: config.SSL,
	})

	if err != nil {
		panic(fmt.Sprintf("Unable to initialize minio client, %s", err))
	}

	return &MinioFilesystem{
		core:   client,
		config: config,
	}
}

func (m MinioFilesystem) Driver() string {
	return MinioDriver
}

func (m MinioFilesystem) BucketPublicName() string {
	return m.config.BucketPublic
}

func (m MinioFilesystem) BucketPrivateName() string {
	return m.config.BucketPrivate
}

func (m MinioFilesystem) Stat(bucketName string, objectName string) (*FileStatInfo, error) {
	objInfo, err := m.core.Client.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &FileStatInfo{
		LastModTime: objInfo.LastModified,
		MimeType:    objInfo.ContentType,
		Name:        objInfo.Key,
		Size:        objInfo.Size,
		Ext:         path.Ext(objectName),
	}, nil
}

func (m MinioFilesystem) Write(bucketName string, objectName string, stream []byte) error {
	_, err := m.core.Client.PutObject(context.Background(), bucketName, objectName, strings.NewReader(string(stream)), int64(len(stream)), minio.PutObjectOptions{})
	return err
}

func (m MinioFilesystem) Copy(bucketName string, srcObjectName, objectName string) error {
	return m.CopyObject(bucketName, srcObjectName, bucketName, objectName)
}

func (m MinioFilesystem) CopyObject(srcBucketName string, srcObjectName, dstBucketName string, dstObjectName string) error {
	srcOpts := minio.CopySrcOptions{
		Bucket: srcBucketName,
		Object: srcObjectName,
	}

	dstOpts := minio.CopyDestOptions{
		Bucket: dstBucketName,
		Object: dstObjectName,
	}

	_, err := m.core.Client.CopyObject(context.Background(), dstOpts, srcOpts)
	return err
}

func (m MinioFilesystem) Delete(bucketName string, objectName string) error {
	return m.core.Client.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
}

func (m MinioFilesystem) GetObject(bucketName string, objectName string) ([]byte, error) {
	object, err := m.core.Client.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	defer object.Close()

	return io.ReadAll(object)
}

func (m MinioFilesystem) PublicUrl(bucketName, objectName string) string {
	uri, err := m.core.Client.PresignedGetObject(context.Background(), bucketName, objectName, 30*time.Minute, nil)
	if err != nil {
		panic(err)
	}

	if m.BucketPublicName() == bucketName {
		uri.RawQuery = ""
	}

	return uri.String()
}

func (m MinioFilesystem) PrivateUrl(bucketName, objectName string, filename string, expire time.Duration) string {

	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	uri, err := m.core.Client.PresignedGetObject(context.Background(), bucketName, objectName, expire, reqParams)
	if err != nil {
		panic(err)
	}

	return uri.String()
}

func (m MinioFilesystem) InitiateMultipartUpload(bucketName, objectName string) (string, error) {
	return m.core.NewMultipartUpload(context.Background(), bucketName, objectName, minio.PutObjectOptions{})
}

func (m MinioFilesystem) PutObjectPart(bucketName, objectName string, uploadID string, index int, data io.Reader, size int64) (ObjectPart, error) {
	part, err := m.core.PutObjectPart(context.Background(), bucketName, objectName, uploadID, index, data, size, minio.PutObjectPartOptions{})
	if err != nil {
		return ObjectPart{}, err
	}

	return ObjectPart{
		PartNumber: part.PartNumber,
		ETag:       part.ETag,
	}, nil
}

func (m MinioFilesystem) CompleteMultipartUpload(bucketName, objectName, uploadID string, parts []ObjectPart) error {
	completeParts := make([]minio.CompletePart, 0)

	for _, part := range parts {
		completeParts = append(completeParts, minio.CompletePart{
			PartNumber: part.PartNumber,
			ETag:       part.ETag,
		})
	}

	_, err := m.core.CompleteMultipartUpload(context.Background(), bucketName, objectName, uploadID, completeParts, minio.PutObjectOptions{})
	return err
}

func (m MinioFilesystem) AbortMultipartUpload(bucketName, objectName, uploadID string) error {
	return m.core.AbortMultipartUpload(context.Background(), bucketName, objectName, uploadID)
}
