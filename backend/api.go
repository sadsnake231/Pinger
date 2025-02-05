package main

import (
	"time"
	"net/http"
	"fmt"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gin-gonic/gin"
)
type PingStats struct {
    Ip           	string              `json:"ip"`
    LastUp          time.Time           `json:"last_up"`
    Min         	time.Duration       `json:"min"`
    Max         	time.Duration       `json:"max"`
    PingTime        time.Time           `json:"time"`
}


var conn *pgxpool.Pool = EstablishConnection()
const timeFormat string = "2006-01-02 15:04:05"



func UpdatePings() gin.HandlerFunc {
	return func (c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		var stats PingStats

		if err := c.ShouldBindJSON(&stats); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}

		query := `
		INSERT INTO results (host, min_time, max_time, last_up, ping_time)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT(host)
		DO UPDATE SET
			min_time = EXCLUDED.min_time,
			max_time = EXCLUDED.max_time,
			last_up = EXCLUDED.last_up,
			ping_time = EXCLUDED.ping_time;
		`

		minConverted := float64(stats.Min) / float64(time.Millisecond) 
		maxConverted := float64(stats.Max) / float64(time.Millisecond) 
		lastUpConverted := stats.LastUp.Format(timeFormat)
		pingTimeConverted := stats.PingTime.Format(timeFormat)

		_, err := conn.Exec(ctx, query, stats.Ip, minConverted, maxConverted, lastUpConverted, pingTimeConverted)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			fmt.Printf(err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{})
	}
}

func PostPings() gin.HandlerFunc {
	return func (c *gin.Context) {
		
	}
}
