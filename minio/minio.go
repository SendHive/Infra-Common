package minio

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type IMinioService interface {
	MinioConnect() (*minio.Client, error)
	CreateBucket(minioClient *minio.Client) error
	DeleteBucket(minioClient *minio.Client) error
	ListBucket(minioClient *minio.Client) ([]minio.BucketInfo, error)
	PutObject(minioClient *minio.Client, file string) error
	DeleteObject(minioClient *minio.Client) error
}

type MinioService struct{}

func NewMinioRequest() (IMinioService, error) {
	return &MinioService{}, nil
}

func (m *MinioService) MinioConnect() (*minio.Client, error) {
	endpoint := "localhost:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "minioadmin"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln("Failed to connect to MinIO:", err)
		return nil, err
	}

	log.Println("Successfully connected to MinIO")
	return minioClient, nil
}

func (m *MinioService) CreateBucket(minioClient *minio.Client) error {
	err := minioClient.MakeBucket(context.Background(), "mybucket-1", minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: false})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Successfully created mybucket.")
	return nil
}

func (m *MinioService) ListBucket(minioClient *minio.Client) ([]minio.BucketInfo, error) {
	buckets, err := minioClient.ListBuckets(context.Background())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for _, bucket := range buckets {
		fmt.Println(bucket)
	}
	return buckets, nil
}

func (m *MinioService) DeleteBucket(minioClient *minio.Client) error {
	err := minioClient.RemoveBucket(context.Background(), "mybucket-1")
	if err != nil {
		fmt.Println(err)
		return err
	}
	log.Println("Bucket Deleted successfully !!!!")
	return nil
}

func (m *MinioService) PutObject(minioClient *minio.Client, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return err
	}

	uploadInfo, err := minioClient.PutObject(context.Background(), "mybucket-1", "myobject", file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Successfully uploaded bytes: ", uploadInfo)
	return nil
}

func (m *MinioService) DeleteObject(minioClient *minio.Client) error {
	objectsCh := make(chan minio.ObjectInfo)

	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)
		// List all objects from a bucket-name with a matching prefix.
		for object := range minioClient.ListObjects(context.Background(), "mybucket-1", minio.ListObjectsOptions{WithVersions: true}) {
			if object.Err != nil {
				log.Fatalln(object.Err)
			}
			objectsCh <- object
		}
	}()

	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}

	for rErr := range minioClient.RemoveObjects(context.Background(), "mybucket-1", objectsCh, opts) {
		fmt.Println("Error detected during deletion: ", rErr)
		return &rErr
	}
	return nil
}
