package main

// import (
// 	"fmt"
// 	"time"

// 	"go.uber.org/zap"
// )

// func main() {
// 	log, err := zap.NewProduction()
// 	if err != nil {
// 		panic(err)
// 	}

// 	wal, err := NewWriteAheadLog(&WALOptions{
// 		LogDir:            "data/",
// 		MaxLogSize:        40 * 1024 * 1024, // 40 MB (log rotation size)
// 		MaxSegments:       2,
// 		Log:               log,
// 		MaxWaitBeforeSync: 1 * time.Second,
// 		SyncMaxBytes:      1000,
// 	})
// 	if err != nil {
// 		fmt.Println("Error creating Write-Ahead Log:", err)
// 		return
// 	}
// 	defer wal.Close()

// 	for i := 0; i < 10000000; i++ {
// 		_, err := wal.Write([]byte("Simulate some database changes and write them to the WAL"))
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	// Sync all the logs to the file at the end so that we dont loose any data
// 	err = wal.Sync()
// 	if err != nil {
// 		fmt.Println("Error in final sync", err)
// 	}

// 	startTime := time.Now()
// 	var numLines int
// 	wal.Replay(0, func(b []byte) error {
// 		numLines++
// 		return nil
// 	})

// 	fmt.Println("Total lines replayed:", numLines)
// 	fmt.Println("time taken:", time.Since(startTime))
// }
