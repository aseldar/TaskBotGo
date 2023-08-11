package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

type Task struct {
	Description string `json:"description"`
}

// Команда для добавления задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Закрываем тело запроса после использования
	defer r.Body.Close()

	// Сохраняем задачу в Redis
	err = redisClient.Set(context.Background(), "task:"+task.Description, task.Description, 0).Err()
	if err != nil {
		http.Error(w, "Failed to save task to Redis", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Команда для просмотра списка задач
func viewTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Здесь можно добавить код для получения списка задач из Redis и отправки пользователю
	keys, err := redisClient.Keys(context.Background(), "task:*").Result()
	if err != nil {
		http.Error(w, "Failed to fetch tasks from Redis", http.StatusInternalServerError)
		return
	}

	// Выводим список ключей в журнал
	log.Println("List of tasks in Redis:")
	for _, key := range keys {
		log.Println("Key:", key)
	}

	// Получаем значения всех задач из Redis
	var tasks []string
	for _, key := range keys {
		value, err := redisClient.Get(context.Background(), key).Result()
		if err != nil {
			http.Error(w, "Failed to fetch tasks from Redis", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, value)
	}

	// Отправляем список задач в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func main() {
	// Инициализация Redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	// Маршрутизация запросов
	http.HandleFunc("/add", addTaskHandler)
	http.HandleFunc("/view", viewTasksHandler)

	// Проверка связи с Redis
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to ping Redis: %v", err)
	}
	log.Printf("Redis ping response: %s", pong)

	// Запуск сервера
	log.Println("Starting API server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
