package listeners

import (
	"log"
	"os"
	"syscall"

	"github.com/cansulting/elabox-system-tools/foundation/system"
	"github.com/ztrue/shutdown"
)

func ListenToShutdown() {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	shutdown.AddWithParam(func(sig os.Signal) {
		log.Println("shutting down")
		system.SetEnv("ELASHUTDOWNSTATUS", "properly_shutdown")
	})
	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM)
}