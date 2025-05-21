package main

import (
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const ITEMS_PER_PAGE = 4000
const CLUSTER_SIZE = 99

type Student struct {
	ID           int
	MongoID      bson.ObjectID `bson:"_id" json:"id"`
	Name         string        `bson:"name"`
	Email        string        `bson:"email"`
	Age          int           `bson:"age"`
	RegisteredAt time.Time     `bson:"registeredAt"`
}

func main() {
	m := NewMongoDB()
	defer m.Close()
	p := NewPostGres()
	defer p.Close()
	total, err := m.Count()
	if err != nil {
		fmt.Printf("failed to count, err: %s\n", err.Error())
		return
	}

	bar := NewProgress(total / ITEMS_PER_PAGE)

	var wg sync.WaitGroup
	c := make(chan []Student)

	for range CLUSTER_SIZE {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for students := range c {
				if err := p.Insert(students); err != nil {
					fmt.Printf("failed to insert, err: %s\n", err.Error())
					break
				}

				bar.Increment()
			}
		}()
	}

	for ss := range m.GetAllPagedData(ITEMS_PER_PAGE) {
		c <- ss
	}

	close(c)
	wg.Wait()

	totalPostgres, err := p.Count()
	if err != nil {
		panic(err)
	}
	fmt.Printf("total on MongoDB %d and total on PostGres %d", total, totalPostgres)
}
