package wechat

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

var (
	recordnum int
	channel   chan int
)

func Initialize() {

	db, er := bolt.Open(DBName, 0666, nil)
	defer db.Close()
	if er != nil {
		log.Fatal(er)
	}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(DBBucket)
		status := b.Stats()
		recordnum = status.KeyN
		return nil
	})

	//	channel = make(chan int, recordnum)// maybe
	channel = make(chan int)
	go produceRandomIndex()
	//	fmt.Println("wechat.message's init")

}

func MyMessage() (snippet string) {
	var index int
	index = <-channel

	db, er := bolt.Open(DBName, 0666, nil)
	defer db.Close()
	if er != nil {
		log.Fatal(er)
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(DBBucket)
		data := b.Get([]byte(strconv.Itoa(index)))
		snippet = string(data)
		return nil
	})

	return snippet
}

//for test
func PrintMsg() {

	ticker := time.NewTicker(time.Second * 4)

	for {
		<-ticker.C
		s := MyMessage()
		fmt.Println(s)
	}

}

func produceRandomIndex() {
	for {
		source := rand.NewSource(time.Now().Unix())
		r := rand.New(source)
		a := r.Perm(recordnum)
		for _, index := range a {
			channel <- index
		}
	}

}
