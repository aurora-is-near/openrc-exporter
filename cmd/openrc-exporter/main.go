package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/aurora-is-near/openrc-exporter/pkg/collector"
	"github.com/aurora-is-near/openrc-exporter/pkg/openrc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddress string

	version     = "dev"
	versionFlag bool
)

func init() {
	flag.StringVar(&listenAddress, "listen-address", ":9816", "Listening address")
	flag.BoolVar(&versionFlag, "version", false, "Print version and exit")
}

func main() {
	flag.Parse()

	if versionFlag {
		fmt.Printf("version %s (%s/%s)\n", version, runtime.GOOS, runtime.GOARCH)
		fmt.Println("golang:", runtime.Version())
		fmt.Println("RC_SVCDIR:", openrc.SvcDir)
		fmt.Println("RC_RUNLEVELDIR:", openrc.RunlevelDir)
		fmt.Println("RC_INITDIR:", openrc.InitDir)
		fmt.Println("RC_CONFDIR:", openrc.ConfDir)
		os.Exit(0)
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)

	prometheus.MustRegister(collector.New(logger))

	http.Handle("/metrics", promhttp.Handler())

	logger.Println("Listening on", listenAddress)
	logger.Fatal(http.ListenAndServe(listenAddress, nil))
}
