package hid

import (
	"errors"
	"machine"
	"machine/usb"
)

// from usb-hid.go
var (
	ErrHIDInvalidPort    = errors.New("invalid USB port")
	ErrHIDInvalidCore    = errors.New("invalid USB core")
	ErrHIDReportTransfer = errors.New("failed to transfer HID report")
)

const (
	hidEndpoint = usb.HID_ENDPOINT_IN

	REPORT_TYPE_INPUT   = 1
	REPORT_TYPE_OUTPUT  = 2
	REPORT_TYPE_FEATURE = 3
)

type hidDevicer interface {
	Handler() bool
}

var devices [5]hidDevicer
var size int

// SetHandler sets the handler. Only the first time it is called, it
// calls machine.EnableHID for USB configuration
func SetHandler(d hidDevicer) {
	if size == 0 {
		machine.EnableHID(handler, nil, setupHandler)
	}

	devices[size] = d
	size++
}

func handler() {
	for _, d := range devices {
		if d == nil {
			continue
		}
		if done := d.Handler(); done {
			return
		}
	}
}

var DefaultSetupHandler = setupHandler

func setupHandler(setup usb.Setup) bool {
	ok := false
	if setup.BmRequestType == usb.SET_REPORT_TYPE && setup.BRequest == usb.SET_IDLE {
		machine.SendZlp()
		ok = true
	}
	return ok
}

// SendUSBPacket sends a HIDPacket.
func SendUSBPacket(b []byte) {
	machine.SendUSBInPacket(hidEndpoint, b)
}
