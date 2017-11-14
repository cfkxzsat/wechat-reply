package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/cfkxzsat/wechat-reply/wechat"
)

func init() {
	//读取数据文件，初始化数据库
	file, err := os.Open(wechat.SnippetPath)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	db, er := bolt.Open(wechat.DBName, 0666, nil)
	defer db.Close()
	if er != nil {
		log.Fatal(er)
	}
	count := 0

	for {
		snippet, errr := reader.ReadString('`')
		if errr == io.EOF {
			break
		} else if errr != nil {
			log.Fatal(err)
		}
		snippet = strings.Trim(snippet, "`\r\n")

		db.Update(func(tx *bolt.Tx) error {
			b, _ := tx.CreateBucketIfNotExists(wechat.DBBucket)
			err := b.Put([]byte(strconv.Itoa(count)), []byte(snippet))
			if err != nil {
				log.Fatal(err)
			}
			return nil
		})
		count++
	}

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(wechat.DBBucket)
		status := b.Stats()
		fmt.Println(status.KeyN)
		return nil
	})
	//	fmt.Println("main's init")

}

func main() {
	wechat.Initialize()
	//for test
	// go wechat.PrintMsg()
	// fmt.Print("1") //\r make 1 overwrite some data
	log.Println("Wechat Service: Start!")
	http.HandleFunc("/", wechat.ProcRequest)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("Wechat Service: ListenAndServe failed, ", err)
	}
	log.Println("Wechat Service: Stop!")

}
