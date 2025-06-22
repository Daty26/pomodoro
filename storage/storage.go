package storage

import (
	"encoding/json"
	"github.com/Daty26/pomodoro/data"
	"log"
	"os"
)

func LoadLogs() []data.Log {
	var logs []data.Log
	file, err := os.ReadFile("pomodoro_logs.json")
	if err != nil {
		if os.IsNotExist(err) {
			return logs
		}
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &logs)
	if err != nil {
		log.Fatal(err)
	}
	return logs

}

func FileSave(logs []data.Log) error {
	logsEnc, err := json.Marshal(logs)
	if err != nil {
		return err
	}
	if err := os.WriteFile("pomodoro_logs.json", logsEnc, 0644); err != nil {
		return err
	}
	return nil

}
