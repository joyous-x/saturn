package pay 


import (
	"fmt"
	"github.com/iGoogle-ink/gopay"
)

func Version() string {
	xlog.Debug(fmt.Println("GoPay(github.com/iGoogle-ink/gopay) Version: %v", gopay.Version)
	return gopay.Version
}

func UnifiedOrder() {
	gopay.
	client.UnifiedOrder()
}