package vpn

import (
	"fmt"
	"time"

	"github.com/go-ping/ping"
)

func GetPing(addr string) (time.Duration, error) {
	pinger, err := ping.NewPinger(addr)
	if err != nil {
		return 0, err
	}
	pinger.Count = 3
	pinger.Interval = time.Millisecond * 200
	pinger.Timeout = time.Second * 2
	pinger.SetPrivileged(true)

	err = pinger.Run()
	if err != nil {
		return 0, err
	}

	stats := pinger.Statistics()
	if stats.PacketsRecv == 0 {
		return 0, fmt.Errorf("no packets received")
	}

	return stats.AvgRtt, nil
}