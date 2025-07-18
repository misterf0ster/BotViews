package main

import (
	"log"
	"sync"
	"time"
	"workers/internal/browser"
	"workers/internal/task"

	"github.com/playwright-community/playwright-go"
)

type DummyQueue struct {
	tasks []string
	mu    sync.Mutex
	index int
}

func (d *DummyQueue) PopTask(queues ...string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.index >= len(d.tasks) {
		time.Sleep(time.Second) // ждать, если задач нет
		return "", nil
	}
	task := d.tasks[d.index]
	d.index++
	return task, nil
}

func main() {
	dummyQueue := &DummyQueue{
		tasks: []string{
			"order:1:link:https://rutube.ru/video/65310d0b0dc633b833e8650265f7895d:cycles:1",
			"order:1:link:https://rutube.ru/video/65310d0b0dc633b833e8650265f7895d:cycles:1",
			"order:1:link:https://rutube.ru/video/65310d0b0dc633b833e8650265f7895d:cycles:1",
			// Добавь сколько хочешь задач
		},
	}

	proxies := []string{""} // или список прокси
	controlURL := "http://localhost/api/report"

	workerCount := 3

	err := playwright.Install()
	if err != nil {
		log.Fatalf("playwright install error: %v", err)
	}

	for i := 0; i < workerCount; i++ {
		go botWorker(dummyQueue, []string{"tasks:order:7"}, proxies, controlURL, i)
	}

	select {} // чтобы программа не завершалась
}

func botWorker(q interface {
	PopTask(...string) (string, error)
}, queues, proxies []string, controlURL string, workerID int) {
	botID := task.RandomBotID(workerID)
	for {
		proxyAddr := browser.RandomProxy(proxies)
		userAgent := browser.RandomUserAgent()

		payload, err := q.PopTask(queues...)
		if err != nil {
			log.Printf("[%s] PopTask error: %v", botID, err)
			time.Sleep(time.Second)
			continue
		}
		if payload == "" {
			continue // если задач нет — ждём
		}
		t := task.Parse(payload)
		log.Printf("[%s] Task: order=%d link=%s cycles=%d proxy=%s", botID, t.OrderID, t.Link, t.Cycles, proxyAddr)

		err = browser.WatchRutubeWithAdLogic(t.Link, userAgent, proxyAddr)
		if err != nil {
			log.Printf("[%s] watchRutube error: %v", botID, err)
			continue
		}

		log.Printf("[%s] Report sent: orderID=%d", botID, t.OrderID)
	}
}
