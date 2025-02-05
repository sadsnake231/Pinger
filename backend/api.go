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
    LastUp          string           	`json:"last_up"`
    Min         	float64      		`json:"min"`
    Max         	float64       		`json:"max"`
    PingTime        string           	`json:"time"`
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
		
		var flag bool
		if stats.LastUp == "" {
			flag = false
		} else {
			flag = true
		}

		query := `
		INSERT INTO results (host, min_time, max_time, last_up, ping_time)
		VALUES ($1, $2, $3, CASE WHEN $6 THEN $4 ELSE 'never' END, $5)
		ON CONFLICT(host)
		DO UPDATE SET
			min_time = $2,
			max_time = $3,
			last_up = CASE WHEN $6 THEN $4 ELSE results.last_up END,
			ping_time = $5;
		`


		_, err := conn.Exec(ctx, query, stats.Ip, stats.Min, stats.Max, stats.LastUp, stats.PingTime, flag)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{})
			fmt.Printf(err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{})
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
