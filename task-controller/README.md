task-controller/
├── .env
├── config/
│   └── config.go                  # Загрузка конфига
├── internal/
│   ├── db/
│   │   ├── pg.go                 # Подключение к бд
│   │   └── task.go               # Методы работы с задачами и статусами
│   ├── logger/
│   │   └── logger.go             # Логгер
│   └── queue/
│       ├── queue.go              # Подключение к Redis 
│       └── tasks.go              # Методы добавления/получения задач из Redis
└── cmd/
    └── controller/
        └── main.go               