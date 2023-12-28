package main

import (
    "bufio"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/exec"
    "strings"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    wireguardLogs = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "wireguard_logs",
        Help: "WireGuard logs",
    })
)

func main() {
    prometheus.MustRegister(wireguardLogs)

    http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
        wireguardLogs.Set(getWireGuardLogs())

        h := promhttp.Handler()
        h.ServeHTTP(w, r)
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}

func getWireGuardLogs() float64 {
    cmd := exec.Command("journalctl", "-u", "wg-quick@wg0.service", "-n", "100", "--no-pager")
    output, err := cmd.Output()
    if err != nil {
        log.Fatal(err)
    }

    scanner := bufio.NewScanner(strings.NewReader(string(output)))
    logCount := 0
    for scanner.Scan() {
        if strings.Contains(scanner.Text(), "wireguard") {
            logCount++
        }
    }

    return float64(logCount)
}
