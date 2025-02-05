package main

import (
 "fmt"
 "log"
 "time"

 "github.com/go-ping/ping"
)

type PingStats struct {
 Ip           	string
 Errors         int
 Min         	time.Duration
 Max         	time.Duration
 Avg         	time.Duration
 TotalRTT       time.Duration
}

func pingIp(ip string) PingStats {
 stats := PingStats{Ip: ip}

 pinger, err := ping.NewPinger(ip)
 if err != nil {
  log.Fatalf("Ошибка при создании пингера")
 }

 // Пингуем 5 раз, таймаут 5 секунд
 pinger.Count = 5
 pinger.Timeout = 5 * time.Second

 err = pinger.Run()
 if err != nil {
  log.Printf("Ошибка при пинге")
  stats.Errors++
 }

 stats.Errors = pinger.Statistics().PacketsSent - pinger.Statistics().PacketsRecv //колво ошибок: отправленные пакеты - полученные
 stats.Min = pinger.Statistics().MinRtt
 stats.Max = pinger.Statistics().MaxRtt
 stats.Avg = pinger.Statistics().AvgRtt

 return stats
}

func main() {
 ips := []string{"192.168.0.1", "87.240.132.67", "127.0.0.1"}

 for _, ip := range ips {
  stats := pingIp(ip)
  fmt.Printf("Результаты пинга для ip" + " " + ip + "\n")
  fmt.Printf("Ошибок: %d\n", stats.Errors)
  fmt.Printf("Min: %v\n", stats.Min)
  fmt.Printf("Max: %v\n", stats.Max)
  fmt.Printf("Avg: %v\n", stats.Avg)
 }
}