package testsuite_test

import (
	"log"
	"testing"
	"time"

	"github.com/SendHive/Infra-Common/minio"
	"github.com/SendHive/Infra-Common/queue"

	"database/sql"

	_ "github.com/lib/pq"
)

func TestSuite(t *testing.T) {
	t.Run("DatabaseConnection", func(t *testing.T) {
		db, err := sql.Open("postgres", "host=localhost port=5432 user=user password=password dbname=mydatabase sslmode=disable")
		if err != nil {
			t.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()
		if err := db.Ping(); err != nil {
			t.Fatalf("Database ping failed: %v", err)
		}
	})

	t.Run("MinioConnection", func(t *testing.T) {
		minioClient, err := minio.NewMinioRequest()
		if err != nil {
			return
		}
		client, err := minioClient.MinioConnect()
		log.Println(client)
		if err != nil {
			t.Fatalf("Failed to connect to Minio: %v", err)
		}
	})

	t.Run("RabbitMQConnection", func(t *testing.T) {
		qConn, err := queue.NewQueueRequest()
		if err != nil {
			log.Fatal("the error while creating the queue instance: ", err)
		}
		time.Sleep(3 * time.Second)
		qconn, err := qConn.Connect()
		if err != nil {
			return
		}
		log.Println(qConn)
		defer qconn.Close()
	})
}
