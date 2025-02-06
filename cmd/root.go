package cmd

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go-mem-thief/mem"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	filePathSource      string
	filePathDestination string
	fileSize            int
	intervalSeconds     int
	runs                int

	// Metriken für die Anzahl der Schreiboperationen und das Ziel
	fileWriteOperationsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "file_write_operations_total",
			Help: "Total number of file write operations",
		},
		[]string{"destination", "status", "size"},
	)
)

func init() {
	// Registrieren der Metriken
	prometheus.MustRegister(fileWriteOperationsTotal)

	rootCmd.Flags().StringVarP(&filePathSource, "sourcePath", "p", "/tmp/source.bin", "Path to source File")
	rootCmd.Flags().StringVarP(&filePathDestination, "destinationPath", "d", "/tmp/target.bin", "Path to destination File")
	rootCmd.Flags().IntVarP(&fileSize, "size", "s", 200, "Size of File in MB")
	rootCmd.Flags().IntVarP(&intervalSeconds, "interval", "t", 5, "Interval in seconds")
	rootCmd.Flags().IntVarP(&runs, "runs", "n", 2, "Number of runs")
}

var rootCmd = &cobra.Command{
	Use:   "go-mem-thief",
	Short: "Erstellt und liest eine große Datei",
	Long:  `Dieses Programm erstellt eine große Datei und liest sie anschließend wieder ein.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("===Mem Thief===")
		fmt.Println("==Config==")
		fmt.Println("filePathSource: ", filePathSource)
		fmt.Println("fileSize: ", fileSize)
		fmt.Println("intervalSeconds: ", intervalSeconds)
		fmt.Println("runs: ", runs)

		go func() {
			http.Handle("/metrics", promhttp.Handler())
			err := http.ListenAndServe(":2112", nil)
			if err != nil {
				log.Fatal(err)
				return
			}
		}()

		ticker := time.NewTicker(time.Second * time.Duration(intervalSeconds))
		defer ticker.Stop()
		ioOps := mem.IoOps{}

		randomData := make([]byte, int64(fileSize)*1024*1024)

		// Erstellen der Datei
		err := ioOps.WriteFile(filePathSource, randomData)
		if err != nil {
			fmt.Printf("Err Creating the Sourcefile: %v\n", err)
			return
		}

		fmt.Printf("Source File %s (%d MB) Created\n", filePathSource, fileSize)

		i := 0
		for {
			select {
			case <-ticker.C:
				if i <= runs {
					fmt.Println("Read source file", filePathSource)

					// Datei lesen
					content, err := ioOps.ReadLargeFile(filePathSource)
					if err != nil {
						fmt.Println(err)
						return
					}

					fmt.Println("File read, Size:", len(content), "Bytes")

					// Datei schreiben und Metriken zählen
					err = ioOps.WriteFile(filePathDestination, content)
					if err != nil {
						fmt.Println(err)
					} else {

						fileWriteOperationsTotal.WithLabelValues(filePathDestination, "success", strconv.Itoa(fileSize)).Inc()

					}

					fmt.Println("file written")

				} else {
					fmt.Printf("%d read writes done... idling\n", runs)
				}
				i++
			}
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
