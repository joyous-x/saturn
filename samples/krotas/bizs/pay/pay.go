package pay

import (
	"github.com/iGoogle-ink/gopay"
	"github.com/joyous-x/saturn/common/xlog"
)

// Version ...
func Version() string {
	xlog.Debug("GoPay(github.com/iGoogle-ink/gopay) Version: %v", gopay.Version)
	return gopay.Version
}

// UnifiedOrder ...
func UnifiedOrder() {
	// gopay.
	// client.UnifiedOrder()
}
