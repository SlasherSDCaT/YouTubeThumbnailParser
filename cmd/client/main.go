package main

import (
	"YouTubeParser/internal/proto"
	"YouTubeParser/internal/utils"
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
)

var (
	asyncFlag = flag.Bool("async", false, "Асинхронное скачивание")
	srvCheck  = flag.Bool("srvcheck", false, "Проверка доступности сервера")
	fileFlag  = flag.String("file", "", "Файл со списком URL")
)

func printHelp() {
	fmt.Println("Usage: program [FLAGS] [URL1 URL2 ...]")
	fmt.Println("Flags:")
	fmt.Println("  --async        Асинхронное скачивание")
	fmt.Println("  --srvcheck     Проверка доступности сервера")
	fmt.Println("  --file         Путь к файлу со списком URL")
	fmt.Println("Examples:")
	fmt.Println("  program --async https://youtube.com/watch?v=example1")
	fmt.Println("  program --srvcheck")
	fmt.Println("  program --file urls.txt")
}

func main() {
	flag.Usage = printHelp
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 && !*srvCheck && *fileFlag == "" {
		printHelp()
		return
	}

	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Не удалось подключиться к серверу:", err)
	}
	defer conn.Close()
	client := proto.NewThumbnailServiceClient(conn)

	if *srvCheck {
		if err := checkServer(client); err != nil {
			log.Println("Сервер недоступен:", err)
		} else {
			log.Println("Сервер доступен")
		}
	}

	urlList := args
	if *fileFlag != "" {
		urlList, err = utils.ReadURLsFromFile(*fileFlag)
		if err != nil {
			log.Fatal("Ошибка чтения URL из файла:", err)
		}
	}

	if *asyncFlag {
		var wg sync.WaitGroup
		for _, url := range urlList {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				getThumbnail(client, url)
			}(url)
		}
		wg.Wait()
	} else {
		for _, url := range urlList {
			getThumbnail(client, url)
		}
	}
}

func checkServer(client proto.ThumbnailServiceClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	resp, err := client.HealthCheck(ctx, &proto.HealthRequest{})
	if err != nil {
		return err
	}
	log.Println("Состояние сервера:", resp.Status)
	return nil
}

func getThumbnail(client proto.ThumbnailServiceClient, url string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetThumbnail(ctx, &proto.ThumbnailRequest{VideoUrl: url})
	if err != nil {
		log.Println("Ошибка получения превью для", url, ":", err)
		return
	}
	log.Println("Превью для", url, ":", resp.ThumbnailUrl)
}
