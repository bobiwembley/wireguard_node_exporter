# This code sets up a Node Exporter for Prometheus in Go to check WireGuard logs.

# The Node Exporter is a Prometheus exporter that collects metrics from the host machine.
# It is written in Go and provides various system-level metrics such as CPU usage, memory usage, disk usage, and network statistics.

# In this code, we are configuring the Node Exporter to specifically monitor WireGuard logs.
# WireGuard is a modern VPN protocol that provides secure and fast communication between devices.

# To use this code, make sure you have Go installed on your machine and the necessary dependencies for building and running the Node Exporter.

# Once the Node Exporter is up and running, you can use Prometheus to scrape the metrics exposed by the Node Exporter and visualize them in a dashboard.

# For more information on how to configure and use the Node Exporter, refer to the official documentation: https://github.com/prometheus/node_exporter
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
