package main

import (
    "encoding/json"
    "fmt"
    "log"
    "sync"
    "context"
    "time"

    "github.com/rabbitmq/amqp091-go"
    "github.com/jackc/pgx/v5/pgxpool"

)


var (
    rabbitConn *amqp091.Connection
    rabbitCh   *amqp091.Channel
    queueName  = "ping_queue"
    mu         sync.Mutex
)

// Подключение к RabbitMQ
func SetupRabbitMQ() error {
    mu.Lock()
    defer mu.Unlock()

    var err error
    rabbitConn, err = amqp091.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        return fmt.Errorf("ошибка подключения к RabbitMQ: %w", err)
    }

    rabbitCh, err = rabbitConn.Channel()
    if err != nil {
        return fmt.Errorf("ошибка создания канала: %w", err)
    }

    _, err = rabbitCh.QueueDeclare(
        queueName,
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        return fmt.Errorf("ошибка создания очереди: %w", err)
    }

    log.Println("RabbitMQ подключён")
    return nil
}

// Отправляет данные в очередь
func PublishToQueue(ping PingStats) error {
    mu.Lock()
    defer mu.Unlock()

    body, err := json.Marshal(ping)
    if err != nil {
        return fmt.Errorf("ошибка кодирования JSON: %w", err)
    }

    err = rabbitCh.Publish(
        "",
        queueName,
        false,
        false,
        amqp091.Publishing{
            ContentType: "application/json",
            Body:        body,
        },
    )
    if err != nil {
        return fmt.Errorf("ошибка публикации сообщения: %w", err)
    }

    log.Printf("Отправлено в очередь")
    return nil
}

// Читает очередь и вызывает запись в БД
func StartQueueWorker(db *pgxpool.Pool) {
    msgs, err := rabbitCh.Consume(
        queueName,
        "",
        true,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatalf("Ошибка подписки на очередь: %v", err)
    }
    log.Println("Слушаю очередь...")

    for msg := range msgs {
        var ping PingStats
        err := json.Unmarshal(msg.Body, &ping)
        if err != nil {
            log.Printf("Ошибка декодирования сообщения: %v", err)
            continue
        }

        err = SavePingToDB(db, ping)
        if err != nil {
            log.Printf("Ошибка записи в БД: %v", err)
        }
    }
}

// Запись в БД
func SavePingToDB(db *pgxpool.Pool, stats PingStats) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    flag := stats.LastUp != ""

    query := `
    INSERT INTO results (host, min_time, max_time, last_up, ping_time)
    VALUES ($1, $2, $3, CASE WHEN $6 THEN $4 ELSE 'never' END, $5)
    ON CONFLICT(host) DO UPDATE SET
    min_time = $2,
    max_time = $3,
    last_up = CASE WHEN $6 THEN $4 ELSE results.last_up END,
    ping_time = $5;
    `

    _, err := db.Exec(ctx, query, stats.Ip, stats.Min, stats.Max, stats.LastUp, stats.PingTime, flag)
    if err != nil {
        return fmt.Errorf("ошибка записи в БД: %w", err)
    }

    log.Printf("Записано в БД")
    return nil
}

// Закрытие соединения с RabbitMQ
func CloseRabbitMQ() {
    mu.Lock()
    defer mu.Unlock()

    if rabbitCh != nil {
        rabbitCh.Close()
    }
    if rabbitConn != nil {
        rabbitConn.Close()
    }
    log.Println("Соединение с RabbitMQ закрыто")
}