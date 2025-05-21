package main

import (
	"flag"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const DEFAULT_ITEMS_PER_PAGE = 4000
const DEFAULT_CLUSTER_SIZE = 99

type Student struct {
	ID           int
	MongoID      bson.ObjectID `bson:"_id" json:"id"`
	Name         string        `bson:"name"`
	Email        string        `bson:"email"`
	Age          int           `bson:"age"`
	RegisteredAt time.Time     `bson:"registeredAt"`
}

func main() {
	clusterSize := flag.Int("cluster-size", DEFAULT_CLUSTER_SIZE, "")
	itemsPerPage := flag.Int("per-page", DEFAULT_ITEMS_PER_PAGE, "")
	flag.Parse()

	m := NewMongoDB()
	defer m.Close()
	p := NewPostGres()
	defer p.Close()
	total, err := m.Count()
	if err != nil {
		fmt.Printf("failed to count, err: %s\n", err.Error())
		return
	}

	bar := NewProgress(total / *itemsPerPage)

	var wg sync.WaitGroup
	studentsChannel := make(chan []Student)

	for range *clusterSize {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for students := range studentsChannel {
				if err := p.Insert(students); err != nil {
					fmt.Printf("failed to insert, err: %s\n", err.Error())
					break
				}

				bar.Increment()
			}
		}()
	}

	for studentsPage := range m.GetAllPagedData(int64(*itemsPerPage)) {
		studentsChannel <- studentsPage
	}

	close(studentsChannel)
	wg.Wait()

	totalPostgres, err := p.Count()
	if err != nil {
		panic(err)
	}
	fmt.Printf("total on MongoDB %d and total on PostGres %d", total, totalPostgres)
}
