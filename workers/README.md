bot-rutube/
├── cmd/
│   └── bot/                # main.go (точка входа)
├── internal/
│   ├── browser/            # Всё, что связано с chromedp (инициализация, spoof, поведение)
│   │   ├── browser.go
│   │   └── human.go
│   ├── queue/              # Работа с Redis (очередь задач)
│   │   └── queue.go
│   ├── report/             # Отправка отчётов контроллеру
│   │   └── report.go
│   └── task/               # Парсинг и описание задачи
│       └── task.go
├── go.mod
├── go.sum
└── Dockerfile