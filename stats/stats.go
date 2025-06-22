package stats

import (
	"fmt"
	"github.com/Daty26/pomodoro/data"
	"time"
)

func ShowStats(logs []data.Log) {
	totalSec := 0
	now := time.Now()
	y1, m1, d1 := now.Date()
	for _, logEn := range logs {
		if logEn.Working == true {
			y2, m2, d2 := logEn.Timestamp.Date()
			if y1 == y2 && m2 == m1 && d1 == d2 {
				totalSec += logEn.Duration
			}
		}
	}
	fmt.Printf("You worked %.2f minutes today!\n", float64(totalSec)/60)

}
