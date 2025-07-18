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

	for i := 0; i < len(parts); i++ {
		switch parts[i] {
		case "order":
			if i+1 < len(parts) {
				fmt.Sscanf(parts[i+1], "%d", &t.OrderID)
				i++
			}
		case "cycles":
			if i+1 < len(parts) {
				fmt.Sscanf(parts[i+1], "%d", &t.Cycles)
				i++
			}
		case "link":
			// собираем ссылку из всех частей, пока не встретим другой ключ или конец
			var linkParts []string
			for j := i + 1; j < len(parts); j++ {
				if parts[j] == "order" || parts[j] == "cycles" || parts[j] == "link" {
					break
				}
				linkParts = append(linkParts, parts[j])
			}
			t.Link = strings.Join(linkParts, ":")
			// пропускаем обработанные части
			i += len(linkParts)
		}
	}

	return t
}

func RandomBotID(workerID int) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("rutube-bot-%d-%d", workerID, rand.Intn(10000))
}
