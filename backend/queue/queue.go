package queue

import (
    "encoding/json"
    "log"
    "sync"
    "context"
    "time"

    "github.com/rabbitmq/amqp091-go"
    "github.com/jackc/pgx/v5/pgxpool"

)

type PingStats struct {
    Ip           	string              `json:"ip"`
    LastUp          string           	`json:"last_up"`
    Min         	float64      		`json:"min"`
    Max         	float64       		`json:"max"`
    PingTime        string           	`json:"time"`
}


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
    rabbitConn, err = amqp091.Dial("amqp://guest:guest@rabbitmq:5672/")
    if err != nil {
        return err
    }

    rabbitCh, err = rabbitConn.Channel()
    if err != nil {
        return err
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
        return err
    }

    log.Println("RabbitMQ подключён")
    return nil
}

// Отправляет данные в очередь
func PublishToQueue(stats PingStats) error {
    mu.Lock()
    defer mu.Unlock()

    body, err := json.Marshal(stats)
    if err != nil {
        return err
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
        return err
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
        log.Fatalf("Ошибка подписки на очередь")
    }
    log.Println("Слушаю очередь...")

    for msg := range msgs {
        var stats PingStats
        err := json.Unmarshal(msg.Body, &stats)
        if err != nil {
            log.Printf("Ошибка декодирования сообщения")
            continue
        }

        err = SavePingToDB(db, stats)
        if err != nil {
            log.Printf("Ошибка записи в БД")
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
        return err
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