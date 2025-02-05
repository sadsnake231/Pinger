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

func sendPings(stats PingStats) error {
    json, err := json.Marshal(stats)
    if err != nil {
        return err
    }

    resp, err := http.Post("http://localhost:5000", "application/json", bytes.NewBuffer(json))
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

    ticker := time.NewTicker(90 * time.Second)
    defer ticker.Stop()

    for range ticker.C {
        fmt.Printf("Пингую...\n")
        for _, ip := range ips {
            stats := pingIp(ip)
            err := sendPings(stats)
            if err != nil {
                fmt.Println(err.Error())
            }
        
        }
    }
}