package task

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Task struct {
	OrderID int
	Link    string
	Cycles  int
}

func Parse(payload string) Task {
	var t Task
	parts := strings.Split(payload, ":")
	for i, p := range parts {
		if p == "order" && i+1 < len(parts) {
			fmt.Sscanf(parts[i+1], "%d", &t.OrderID)
		}
		if p == "cycles" && i+1 < len(parts) {
			fmt.Sscanf(parts[i+1], "%d", &t.Cycles)
		}
		if p == "link" && i+1 < len(parts) {
			t.Link = parts[i+1]
		}
	}
	return t
}

func RandomBotID(workerID int) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("rutube-bot-%d-%d", workerID, rand.Intn(10000))
}
