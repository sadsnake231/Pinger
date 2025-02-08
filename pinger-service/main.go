package main

import (
    "fmt"
    "log"
    "time"
    "encoding/json"
    "net/http"
    //"os"
    //"bufio"
    "bytes"

    "github.com/go-ping/ping"
)

type PingStats struct {
    Ip           	string              `json:"ip"`
    LastUp          string           	`json:"last_up"`
    Min         	float64      		`json:"min"`
    Max         	float64       		`json:"max"`
    PingTime        string           	`json:"time"`
}

const timeFormat string = "2006-01-02 15:04:05"

func pingIp(ip string) PingStats {
    stats := PingStats{Ip: ip}

    pinger, err := ping.NewPinger(ip)
    if err != nil {
        log.Printf("Ошибка при создании пингера")
    }

    // Пингуем 5 раз, таймаут 5 секунд
    pinger.Count = 5
    pinger.Timeout = 5 * time.Second

    err = pinger.Run()
    if err != nil {
        log.Printf("Ошибка при пинге")
    }

    if pinger.Statistics().PacketsRecv != 0 {
        stats.LastUp = time.Now().Format(timeFormat)
    } else {
        stats.LastUp = ""
    }

    stats.Min = float64(pinger.Statistics().MinRtt) / float64(time.Millisecond)
    stats.Max = float64(pinger.Statistics().MaxRtt) / float64(time.Millisecond)
    stats.PingTime = time.Now().Format(timeFormat)
    return stats
}

//функция отправки напрямую на сервер
func sendPings(stats PingStats) error {
    json, err := json.Marshal(stats)
    if err != nil {
        return err
    }

    resp, err := http.Post("http://backend:5000", "application/json", bytes.NewBuffer(json))
    if err != nil {
        return err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusAccepted {
        return fmt.Errorf("сервер не вернул 202")
    }

    return nil
}

/* отправка через очередь

func sendToQueue(stats PingStats) error {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        return error
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        return err
    }
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "ping_stats",
        false,
        false,
        false,
        false,
        nil
    )
    if err != nil {
        return err
    }

    json, err := json.Marshal(stats)
    if err != nil {
        return err
    }

    err = ch.Publish(
        "",
        q.Name,
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body: body.
        }
    )
    if err != nil {
        return err
    }
    return nil
}
*/
func main() {
    /*
    file, err := os.Open("ips")
    if err != nil {
        log.Fatalf("не смог открыть файл")
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    var ips []string

    for scanner.Scan() {
        ips = append(ips, scanner.Text())
    }
    */
    ips := []string{"test_server1", "test_server2", "test_server3", "87.240.129.133"}

    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        //fmt.Printf("Пингую...\n")
        for _, ip := range ips {
            stats := pingIp(ip)
            err := sendPings(stats)
            if err != nil {
                fmt.Println(err.Error())
            }
        
        }
    }
}