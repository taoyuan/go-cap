package cap

import (
	"testing"
	"github.com/smartystreets/goconvey/convey"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/prashantv/gostub"
	"strings"
)

const OUTPUT_STARTED = `WARN: Your adapter does not fully support AP virtual interface, enabling --no-virt
WARN: If AP doesn't work, please read: howto/realtek.md
Config dir: /tmp/create_ap.wlan0.conf.Tg1VJRNW
PID: 3233
Network Manager found, set wlan0 as unmanaged device... DONE
Sharing Internet using method: nat
hostapd command-line interface: hostapd_cli -p /tmp/create_ap.wlan0.conf.Tg1VJRNW/hostapd_ctrl
Configuration file: /tmp/create_ap.wlan0.conf.Tg1VJRNW/hostapd.conf
Using interface wlan0 with hwaddr e8:4e:06:34:ff:db and ssid "MyAccessPoint"
wlan0: interface state UNINITIALIZED->ENABLED
wlan0: AP-ENABLED`

const OUTPUT_DONE = `wlan0: interface state ENABLED->DISABLED
wlan0: AP-DISABLED
nl80211: deinit ifname=wlan0 disabled_11b_rates=0

Doing cleanup.. done`

const OUTPUT_ERROR = `WARN: Your adapter does not fully support AP virtual interface, enabling --no-virt
ERROR: 'eth1' is not an interface`

func TestCreateAP(t *testing.T) {
	convey.Convey("Run with map success", t, func() {
		stub := gostub.Stub(&CommandCreateAP, "test/create_ap_ok")
		defer stub.Reset()

		ap, err := CreateAP(map[string]interface{}{
			"iface": "wlan0",
		})
		assert.Nil(t, err)
		assert.NotNil(t, ap)

		ap.Start()
		time.Sleep(time.Duration(1) * time.Second)

		output := ap.Output()
		expected := strings.Split(OUTPUT_STARTED, "\n")
		assert.Equal(t, expected, output)

		ap.Stop()
		time.Sleep(time.Duration(1) * time.Second)

		output = ap.Output()
		expected = strings.Split(OUTPUT_DONE, "\n")
		assert.Equal(t, expected, output)

		err = ap.Wait()
		assert.Nil(t, err)
	})

	convey.Convey("Run success", t, func() {
		stub := gostub.Stub(&CommandCreateAP, "test/create_ap_ok")
		defer stub.Reset()

		ap, err := CreateAP(map[string]interface{}{
			"iface": "wlan0",
		})
		assert.Nil(t, err)
		assert.NotNil(t, ap)

		ap.Start()
		time.Sleep(time.Duration(1) * time.Second)

		output := ap.Output()
		expected := strings.Split(OUTPUT_STARTED, "\n")
		assert.Equal(t, expected, output)

		ap.Stop()
		time.Sleep(time.Duration(1) * time.Second)

		output = ap.Output()
		expected = strings.Split(OUTPUT_DONE, "\n")
		assert.Equal(t, expected, output)

		err = ap.Wait()
		assert.Nil(t, err)
	})

	convey.Convey("Run error", t, func() {
		stub := gostub.Stub(&CommandCreateAP, "test/create_ap_error")
		defer stub.Reset()

		ap, err := CreateAP(map[string]interface{}{
			"iface": "wlan0",
		})
		assert.Nil(t, err)
		assert.NotNil(t, ap)

		ap.Start()
		time.Sleep(time.Duration(1) * time.Second)

		output := ap.Output()
		expected := strings.Split(OUTPUT_ERROR, "\n")
		assert.Equal(t, expected, output)

		err = ap.Wait()
		assert.NotNil(t, err)
	})
}
