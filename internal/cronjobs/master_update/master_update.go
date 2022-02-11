package master_update

import (
	"log"
	"time"
)

func MasterUpdate() {
	time.Sleep(3 * time.Second)
	times := time.Now().UnixNano()
	log.Printf("master update: %v", times)
}
