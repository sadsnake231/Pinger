package api

import (
	"time"
	"net/http"
	"context"

	"Pinger/queue"
	"Pinger/database"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gin-gonic/gin"

)

var conn *pgxpool.Pool = database.EstablishConnection()
const timeFormat string = "2006-01-02 15:04:05"

type PingStats struct {
    Ip           	string              `json:"ip"`
    LastUp          string           	`json:"last_up"`
    Min         	float64      		`json:"min"`
    Max         	float64       		`json:"max"`
    PingTime        string           	`json:"time"`
}

func UpdatePings() gin.HandlerFunc {
	return func (c *gin.Context) {

		var stats queue.PingStats

		if err := c.ShouldBindJSON(&stats); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//пишем в очередь полученную стату
		err := queue.PublishToQueue(stats)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error:":err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{})
	}
}

func GetPings() gin.HandlerFunc {
	return func (c *gin.Context) {
		
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()
		
		var statsArr []PingStats
		
		query := `SELECT host, min_time, max_time, last_up, ping_time FROM results`

		rows, err := conn.Query(ctx, query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var stats PingStats
			err := rows.Scan(&stats.Ip, &stats.Min, &stats.Max, &stats.LastUp, &stats.PingTime)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			
			statsArr = append(statsArr, stats)
			
		}

		c.JSON(http.StatusOK, statsArr)
		

	}
}
