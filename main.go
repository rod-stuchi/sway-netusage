// Public domain.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/nettest"
)

var (
	iface        = flag.String("interface", "", "What interface to use")
	showDownload = flag.Bool("down", false, "show download only, upload next line (tooltip)")
	showUpload   = flag.Bool("up", false, "show upload only, download next line (tooltip)")
)

type Output struct {
	Text    string `json:"text"`
	Alt     string `json:"alt"`
	Tooltip string `json:"tooltip"`
	Class   string `json:"class"`
}

// stats fetches the cumulative rx/tx bytes for network interface
// iface.
func stats() (rx, tx uint64) {
	b, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return 0, 0
	}
	buff := bytes.NewBuffer(b)
	for l, err := buff.ReadString('\n'); err == nil; {
		l = strings.Trim(l, " \n")
		if !strings.HasPrefix(l, *iface) {
			l, err = buff.ReadString('\n')
			continue
		}
		re := regexp.MustCompile(" +")
		s := strings.Split(re.ReplaceAllString(l, " "), " ")
		rx, err := strconv.ParseUint(s[1], 10, 64)
		if err != nil {
			return 0, 0
		}
		tx, err := strconv.ParseUint(s[9], 10, 64)
		if err != nil {
			return 0, 0
		}
		return rx, tx
	}
	return 0, 0
}

func getRoutedInterface() *net.Interface {
	val, _ := nettest.RoutedInterface("ip", net.FlagUp|net.FlagBroadcast)
	return val
}

// format converts a number of bytes in KiB or MiB.
func format(counter, prevCounter uint64, window float64) string {
	if prevCounter == 0 {
		return "B"
	}
	r := float64(counter-prevCounter) / window
	if r < 1024 {
		return fmt.Sprintf("%.0f B", r)
	}
	if r < 1024*1024 {
		return fmt.Sprintf("%.0f KiB", r/1024)
	}
	return fmt.Sprintf("%.1f MiB", r/1024/1024)
}

func main() {
	flag.Parse()
	if *iface == "" {
		autoInterface := getRoutedInterface().Name
		fmt.Printf("Auto detected %s\n", autoInterface)
		*iface = autoInterface
	}
	prevRx, prevTx := uint64(0), uint64(0)
	prev := time.Now()
	for {

		time.Sleep(1 * time.Second)
		now := time.Now()
		window := now.Sub(prev).Seconds()
		prev = now
		rx, tx := stats()
		rxRate := format(rx, prevRx, window)
		txRate := format(tx, prevTx, window)
		arrowDownFmt := "<span size=\"10.5pt\"> </span>" // ↑
		arrowUPFmt := "<span size=\"10.5pt\"> </span>"   // ↓
		if *showDownload {
			output := Output{
				Text:    fmt.Sprintf("%8s/s %s", rxRate, arrowDownFmt),
				Tooltip: fmt.Sprintf("%8s/s %s", txRate, arrowUPFmt),
			}
			jsonOutput, err := json.Marshal(output)
			if err != nil {
				log.Fatalf("JSON marshaling failed: %s", err)
			}
			fmt.Println(string(jsonOutput))
		}
		if *showUpload {
			output := Output{
				Text:    fmt.Sprintf("%8s/s %s", txRate, arrowDownFmt),
				Tooltip: fmt.Sprintf("%8s/s %s", rxRate, arrowUPFmt),
			}
			jsonOutput, err := json.Marshal(output)
			if err != nil {
				log.Fatalf("JSON marshaling failed: %s", err)
			}
			fmt.Println(string(jsonOutput))
		}
		if !*showDownload && !*showUpload {
			output := Output{
				Text:    fmt.Sprintf("%8s/s %s %8s/s %s", rxRate, arrowDownFmt, txRate, arrowUPFmt),
				Tooltip: fmt.Sprintf("%8s/s %s %8s/s %s", rxRate, arrowDownFmt, txRate, arrowUPFmt),
			}
			jsonOutput, err := json.Marshal(output)
			if err != nil {
				log.Fatalf("JSON marshaling failed: %s", err)
			}
			fmt.Println(string(jsonOutput))
		}
		prevRx, prevTx = rx, tx
	}

}
