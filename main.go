package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	"github.com/vorkytaka/easyvk-go/easyvk"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type result struct {
	index int
	res   http.Response
	err   error
}

type post struct {
	id string
}

var posts []post
var loaded bool = false

func boundedParallelGet(url string, vk easyvk.VK) []result {
	resultsChan := make(chan *result)

	defer func() {
		close(resultsChan)
	}()

    i := 0

	for {
		go func(i int, url string) {
			_, err := http.Get(url)

			if err != nil {
				fmt.Println(err)
			}

			result := &result{}

			doc, _ := goquery.NewDocument(url)

			doc.Find(".wall_item").Each(func(i int, selection *goquery.Selection) {
				if id, ok := selection.Find(".post__anchor").First().Attr("name"); ok {
					found := false

					for _, post := range posts {
						if post.id == id {
							found = true
							break
						}
					}

					if !found && loaded {
						fmt.Println("New post: " + id + " count posts in database: " + string(len(posts)))

						data := strings.Split(id, "_")
						ownerId, postId := strings.Replace(data[0], "post-", "", 1), data[1]

						params := map[string]string {
							"owner_id": "-" + ownerId,
							"post_id": postId,
							"text": "Мда.",
						}

						resBytes, err := vk.Request("wall.createComment", params); if err != nil {
							fmt.Println("Error: " + err.Error())
						}

						fmt.Println(resBytes)
					}

					posts = append(posts, post{ id: id })
				}
			})

			loaded = true
			resultsChan <- result
			i += 1
		}(i, url)

		if i % 2 == 0 {
			<-time.After(1 * time.Second)
		}
	}
}

func main() {
	benchmark := func(url string) {
		vk, err := easyvk.WithAuth(os.Getenv("VK_EMAIL"), os.Getenv("VK_PASSWORD"), "6278780", "friends,wall,photos"); if err != nil {
			fmt.Println("Error publish post: " + err.Error())
		}

		boundedParallelGet(url, vk)
	}

	err := godotenv.Load(); if err != nil {
		log.Fatal("Error loading .env file")
		panic("Error loading configuration")
	}

	benchmark("https://vk.com/wall-460389?own=1")
}