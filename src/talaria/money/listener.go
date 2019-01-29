package money

import (
	"github.com/Comcast/webpa-common/device"
)

// DecorateListener decorates a Listener with money
//
// Everytime a dispatcher is called money is checked in two ways
// 1. If a RecieveTracker channel recieves a tracker.
// 2. If a dispatchTo Event contains a money map that is then forwarded to the HTTPSpanner decorator to be written to headers.
type Listener func(*device.Event)

// DecorateListener decorates all listeners in device's listener list.  If an event's message contains money, money
// is sent to the http decorator through a money channel
func DecorateListeners(l []device.Listener, tc MoneyContainer) []device.Listener {
	for i, v := range l {
		v = func(e *device.Event) {
			// listen for channel
			select {
			case <-tc.start:
				// return the span map to the money decorator to be injected into an http response.
				tc.contents <- e.Contents

				v(e)

			default:

				v(e)
			}
		}

	}

	return l
}
