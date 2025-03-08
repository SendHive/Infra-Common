package setup

import (
	"fmt"
	"log"
	"time"

	"github.com/SendHive/Infra-Common/db"
	"github.com/SendHive/Infra-Common/minio"
	"github.com/SendHive/Infra-Common/queue"
)

func Setup() {
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
	fmt.Println(queue)
	time.Sleep(3 * time.Second)

	minioClient, err := minio.NewMinioRequest()
	if err != nil {
		return
	}
	mc, err := minioClient.MinioConnect()
	if err != nil {
		fmt.Println("Error while running the MinioConnect")
		return
	}
	log.Println(mc)
	time.Sleep(3 * time.Second)
	log.Println("Setup done successfully")
}
