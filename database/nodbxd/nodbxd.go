package nodbxd

import (
	"fmt"
	"time"
	".."
)

type NoDBxD struct {}

func (n NoDBxD) Store(entries []database.Entry) {
	for _, entry := range entries {
		fmt.Printf("\n[%s] New data:\n", time.Now().Format(time.Stamp))
		fmt.Printf(" - %s %s %s %d %d\n",
			entry.Mac().String(),
			entry.Ipv4().String(),
			entry.Ipv6().String(),
			entry.GetDownloadSpeed(),
			entry.GetUploadSpeed())
	}
}

func (n NoDBxD) Close() {}
