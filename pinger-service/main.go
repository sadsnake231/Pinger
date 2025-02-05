package main

import (
    "fmt"
    "log"
    "time"
    "encoding/json"
    "net/http"
    "bytes"

    "github.com/go-ping/ping"
)

type PingStats struct {
    Ip           	string              `json:"ip"`
    LastUp          time.Time           `json:"last_up"`
    Min         	time.Duration       `json:"min"`
    Max         	time.Duration       `json:"max"`
    PingTime        time.Time           `json:"time"`
}

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
        stats.LastUp = time.Now()
    }
    stats.Min = pinger.Statistics().MinRtt
    stats.Max = pinger.Statistics().MaxRtt
    stats.PingTime = time.Now()
    return stats
}

func sendPings(stats PingStats) error {
    json, err := json.Marshal(stats)
    if err != nil {
        return err
    }

    resp, err := http.Post("http://localhost:3000", "application/json", bytes.NewBuffer(json))
    if err != nil {
        return err
    }

    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("server did not return 201")
    }

    return nil
}

func main() {
    ips := []string{"192.168.0.1", "87.240.132.67", "127.0.0.1"}

    for _, ip := range ips {
        stats := pingIp(ip)
        fmt.Printf("Результаты пинга для ip" + " " + ip + "\n")
        err := sendPings(stats)
        if err != nil {
            fmt.Println(err.Error())
        }
        
    }
}