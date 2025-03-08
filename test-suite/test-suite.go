package testsuite

import (
	"fmt"
	"log"
	"time"

	"github.com/SendHive/Infra-Common/db"
	"github.com/SendHive/Infra-Common/minio"
	"github.com/SendHive/Infra-Common/queue"
)

const (
	bucketName = "mybucket-1"
)

func TestSuite() {
	dbConn, err := db.NewDbRequest()
	if err != nil {
		log.Fatal("the error while creating the database connection: ", err.Error())
	}
	dbConn.InitDB()

	qConn, err := queue.NewQueueRequest()
	if err != nil {
		log.Fatal("the error while creating the queue instance: ", err)
	}
	time.Sleep(3 * time.Second)
	qconn, err := qConn.Connect()
	if err != nil {
		return
	}
	time.Sleep(3 * time.Second)
	queue, err := qConn.DeclareQueue(qconn)
	if err != nil {
		fmt.Println("Error while running the DeclareQueue")
		return
	}
	time.Sleep(3 * time.Second)
	err = qConn.PublishMessage(queue, qconn, "hello world")
	if err != nil {
		fmt.Println("Error while running the PublishMessage")
		return
	}
	time.Sleep(3 * time.Second)
	err = qConn.ConsumeMessage(queue, qconn, true)
	if err != nil {
		fmt.Println("Error while running the ConsumeMessage")
		return
	}

	minioClient, err := minio.NewMinioRequest()
	if err != nil {
		return
	}
	mc, err := minioClient.MinioConnect()
	if err != nil {
		fmt.Println("Error while running the MinioConnect")
		return
	}

	time.Sleep(3 * time.Second)

	err = minioClient.CreateBucket(mc, bucketName)
	if err != nil {
		fmt.Println("Error while running the CreateBucket")
		return
	}

	time.Sleep(3 * time.Second)

	err = minioClient.PutObject(mc, bucketName, "test.txt", "myobject")
	if err != nil {
		fmt.Println("Error while running the PutObject")
		return
	}

	time.Sleep(3 * time.Second)

	err = minioClient.DeleteObject(mc, bucketName)
	if err != nil {
		fmt.Println("Error while running the DeleteObject")
		return
	}

	time.Sleep(3 * time.Second)

	err = minioClient.DeleteBucket(mc, bucketName)
	if err != nil {
		fmt.Println("Error while running the DeleteBucket")
		return
	}

	fmt.Println("All Test Passed.....")
}
