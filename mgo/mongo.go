package mgo

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"log"
	"time"
)

type (
	Mongo struct {
		host     string
		port     int
		username string
		password string
		database string
		c        *mongo.Client
		db       *mongo.Database
	}

	Col struct {
		c *mongo.Collection
	}
)

func New(host string, port int, username, password string, database ...string) (cli *Mongo, err error) {
	mgo := &Mongo{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}

	if len(database) > 0 {
		db := database[0]
		mgo.database = db
	}

	url := fmt.Sprintf("mongodb://%s:%s@%s:%d", username, password, host, port)
	//fmt.Println(url)
	client, err := mongo.NewClient(url)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return
	}
	mgo.c = client
	mgo.db = client.Database(mgo.database)

	return mgo, nil
}

func (m *Mongo) Collection(s string, opts ...*options.CollectionOptions) *Col {
	return &Col{c: m.db.Collection(s, opts...)}
}

func (c *Col) Indexs(indexs []mongo.IndexModel) (ret []string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	ret, err = c.c.Indexes().CreateMany(ctx, indexs)
	if err != nil {
		log.Println(err)
	}
	return
}

func (c *Col) FindOne(f bson.M) (result bson.M, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := c.c.Find(ctx, f)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer func() {
		_ = cur.Close(ctx)
	}()

	for cur.Next(ctx) {
		err = cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	if err = cur.Err(); err != nil {
		log.Fatal(err)
		return
	}

	return
}

func (c *Col) InsertOne(d interface{}) (id interface{}, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := c.c.InsertOne(ctx, d)
	if err != nil {
		return
	}
	return res.InsertedID, nil
}

func (c *Col) Insert(d []map[string]string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var data []interface{}
	for _, item := range d {
		data = append(data, item)
	}
	res, err := c.c.InsertMany(ctx, data)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("insert many: %v, result: %v", data, res)
	return
}
