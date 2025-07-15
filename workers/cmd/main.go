package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"workers/internal/browser"
	"workers/internal/queue"
	"workers/internal/report"
	"workers/internal/task"
)

func main() {
	redisAddr := getenv("REDIS_ADDR", "127.0.0.1:6379")
	redisPass := getenv("REDIS_PASS", "")
	queues := strings.Split(getenv("BOT_QUEUES", "tasks:order:7"), ",")
	proxies := strings.Split(getenv("BOT_PROXIES", ""), ",")
	workers := getenvInt("BOT_WORKERS", 3)
	controlURL := getenv("CONTROLLER_URL", "http://controller/api/report")

	log.Printf("RUTUBE BOT: workers=%d queues=%v proxies=%v", workers, queues, proxies)

	q := queue.New(redisAddr, redisPass)

	for i := 0; i < workers; i++ {
		go botWorker(q, queues, proxies, controlURL, i)
	}
	select {}
}

func botWorker(q *queue.Queue, queues, proxies []string, controlURL string, workerID int) {
	botID := task.RandomBotID(workerID)
	for {
		proxyAddr := browser.RandomProxy(proxies)
		userAgent := browser.RandomUserAgent()

		// Взять задачу
		payload, err := q.PopTask(queues...)
		if err != nil {
			log.Printf("[%s] BLPop error: %v", botID, err)
			time.Sleep(time.Second)
			continue
		}
		t := task.Parse(payload)
		log.Printf("[%s] Task: order=%d link=%s cycles=%d proxy=%s", botID, t.OrderID, t.Link, t.Cycles, proxyAddr)

		// Выполнить задачу (эмулировать просмотр на rutube)
		err = browser.WatchRutubeHuman(t.Link, 40*time.Second, userAgent, proxyAddr)
		if err != nil {
			log.Printf("[%s] watchRutube error: %v", botID, err)
			continue
		}

		// Отчитаться в контроллер
		report.Send(controlURL, t.OrderID, 1, botID)
	}
}

func getenv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
func getenvInt(key string, def int) int {
	if val := os.Getenv(key); val != "" {
		var x int
		fmt.Sscanf(val, "%d", &x)
		return x
	}
	return def
}
