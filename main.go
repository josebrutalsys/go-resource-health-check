package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/mackerelio/go-osstat/memory"
)

var ALLOW_MEMORY_USAGE float64 = 0.75
var ALLOW_CPU_USAGE float64 = 0.75
var ALLOW_LOAD float64 = 0.75

func healthCheck(w http.ResponseWriter, req *http.Request) {
	memory, err := memory.Get()

	if err != nil {
		log.Println("Error: ", err)
		fmt.Fprintf(w, "ko\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	memoryUsage := 1 - float64(memory.Free)/float64(memory.Total)
	if memoryUsage >= ALLOW_MEMORY_USAGE {
		log.Printf("Memory bigger than allowed %f\n", memoryUsage)
		fmt.Fprintf(w, "ko\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cpu, err := cpu.Get()
	if err != nil {
		log.Println("Error: ", err)
		fmt.Fprintf(w, "ko\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cpuUsage := float64(cpu.System) / float64(cpu.Total)
	if cpuUsage >= ALLOW_CPU_USAGE {
		log.Printf("CPU usage bigger than allowed %f\n", cpuUsage)
		fmt.Fprintf(w, "ko\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	loadAvg, err := loadavg.Get()
	if err != nil {
		log.Println("Error: ", err)
		fmt.Fprintf(w, "ko\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	loadPercentage := float64(loadAvg.Loadavg1) / float64(cpu.CPUCount)

	if loadPercentage >= ALLOW_LOAD {
		log.Printf("CPU load bigger than allowed %f\n", loadPercentage)
		fmt.Fprintf(w, "ko\n")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "ok\n")
}

func main() {
	http.HandleFunc("/health-check", healthCheck)
	http.ListenAndServe(":30001", nil)
}
