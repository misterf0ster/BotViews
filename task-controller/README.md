task-controller/
├── .env
├── config/
│   └── config.go                  # Загрузка конфига
├── internal/
│   ├── db/
│   │   ├── pg.go                 # Подключение и базовые методы
│   │   └── task.go               # Методы работы с задачами и статусами
│   ├── logger/
│   │   └── logger.go             # Логгер с zap (или zerolog)
│   └── queue/
│       ├── queue.go              # Подключение к Redis и базовые методы
│       └── tasks.go              # Методы добавления/получения задач из Redis
└── cmd/
    └── controller/
        └── main.go               # Точка входа, инициализация, цикл