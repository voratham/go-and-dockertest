package student_test

import (
	"context"
	"fmt"
	"go-and-dockertest/student"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

var _ = Describe("StudentRepo", Ordered, func() {

	var dbClient *mongo.Client
	var resource *dockertest.Resource
	var _pool *dockertest.Pool
	var _db *mongo.Database
	var repo student.StudentRepo

	BeforeAll(func() {

		pool, err := dockertest.NewPool("")
		if err != nil {
			log.Fatalf("Could not connect to docker: %s", err)
		}
		// pull mongodb docker image for version 5.0
		resource, err = pool.RunWithOptions(&dockertest.RunOptions{
			Repository: "mongo",
			Tag:        "5.0",
			Env: []string{
				// username and password for mongodb superuser
				"MONGO_INITDB_ROOT_USERNAME=root",
				"MONGO_INITDB_ROOT_PASSWORD=password",
			},
		}, func(config *docker.HostConfig) {
			// set AutoRemove to true so that stopped container goes away by itself
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{
				Name: "no",
			}
		})
		if err != nil {
			log.Fatalf("Could not start resource: %s", err)
		}

		// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
		err = pool.Retry(func() error {
			var err error

			dbClient, err = mongo.Connect(
				context.TODO(),
				options.Client().ApplyURI(
					fmt.Sprintf("mongodb://root:password@localhost:%s", resource.GetPort("27017/tcp")),
				),
			)
			if err != nil {
				return err
			}

			return dbClient.Ping(context.TODO(), nil)
		})

		if err != nil {
			log.Fatalf("Could not connect to docker: %s", err)
		}

		// set value
		_pool = pool
		_db = dbClient.Database("test")
		repo = student.NewStudentRepo(_db)

	})

	AfterAll(func() {
		err := _pool.Purge(resource)
		if err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	})

	It("Should be create student correctly", func() {
		err := repo.Create(student.Student{
			Name: "dream",
			Age:  "28",
		})
		Expect(err).To(BeNil())
	})

	It("Should be get student correctly", func() {
		data, err := repo.GetAll()
		Expect(err).To(BeNil())
		Expect(data[0].Name).To(Equal("dream"))
	})

})
