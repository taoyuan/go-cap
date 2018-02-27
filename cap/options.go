package cap

import "go-awt/tools/iu"

type Option struct {
	Name      string
	Shorthand string
	CapOpt    string
	Type      string
	Default   interface{}
	Usage     string
	Resolve   func(interface{}) (interface{}, error)
}

var ResolveIface = func(v interface{}) (interface{}, error) {
	return iu.ResolveIface(v.(string))
}

var Options = []Option{
	{
		Name:      "iface",
		Shorthand: "i",
		CapOpt:    "",
		Type:      "string",
		Default:   "default",
		Usage:     "WiFi interface",
		Resolve:   ResolveIface,
	},
	{
		Name:      "iface-share",
		Shorthand: "r",
		CapOpt:    "",
		Type:      "string",
		Default:   "",
		Usage:     "Interface with internet",
		Resolve:   ResolveIface,
	},
	{
		Name:      "ssid",
		Shorthand: "s",
		CapOpt:    "",
		Type:      "string",
		Default:   "CAP",
		Usage:     "Access point name",
	},
	{
		Name:      "passphrase",
		Shorthand: "p",
		CapOpt:    "",
		Type:      "string",
		Default:   "12345678",
		Usage:     "Passphrase for access point",
	},
	{
		Name:      "channel",
		Shorthand: "c",
		CapOpt:    "-c",
		Type:      "string",
		Usage:     "Channel number (default: 1)",
	},
	{
		Name:      "wpa",
		Shorthand: "w",
		CapOpt:    "-w",
		Type:      "string",
		Usage:     "Use 1 for WPA, use 2 for WPA2, use 1+2 for both (default: 1+2)",
	},
	{
		Name:      "no-share",
		Shorthand: "n",
		CapOpt:    "-n",
		Type:      "bool",
		Usage:     "Disable Internet sharing (if you use this, don't pass the <interface-with-internet> argument)",
	},
	{
		Name:      "method",
		Shorthand: "m",
		CapOpt:    "-m",
		Type:      "bool",
		Usage:     "Method for Internet sharing. Use: 'nat' for NAT (default); 'bridge' for bridging; 'none' for no Internet sharing (equivalent to -n)",
	},
	{
		Name:      "psk",
		Shorthand: "",
		CapOpt:    "--psk",
		Type:      "string",
		Usage:     "Use 64 hex digits pre-shared-key instead of passphrase",
	},
	{
		Name:      "hidden",
		Shorthand: "",
		CapOpt:    "--hidden",
		Type:      "bool",
		Usage:     "Make the Access Point hidden (do not broadcast the SSID)",
	},
	{
		Name:      "mac-filter",
		Shorthand: "",
		CapOpt:    "--mac-filter",
		Type:      "bool",
		Usage:     "Enable MAC address filtering",
	},
	{
		Name:      "mac-filter-accept",
		Shorthand: "",
		CapOpt:    "--mac-filter-accept",
		Type:      "string",
		Usage:     "Location of MAC address filter list (defaults to /etc/hostapd/hostapd.accept)",
	},
	{
		Name:      "redirect-to-localhost",
		Shorthand: "",
		CapOpt:    "--redirect-to-localhost",
		Type:      "bool",
		Usage:     "If -n is set, redirect every web request to localhost (useful for public information networks)",
	},
	{
		Name:      "hostapd-debug",
		Shorthand: "",
		CapOpt:    "--hostapd-debug",
		Type:      "string",
		Usage:     "With level between 1 and 2, passes arguments -d or -dd to hostapd for debugging.",
	},
	{
		Name:      "isolate-clients",
		Shorthand: "",
		CapOpt:    "--isolate-clients",
		Type:      "bool",
		Usage:     "Disable communication between clients",
	},
	{
		Name:      "ieee80211n",
		Shorthand: "",
		CapOpt:    "--ieee80211n",
		Type:      "bool",
		Usage:     "Enable IEEE 802.11n (HT)",
	},
	{
		Name:      "ieee80211ac",
		Shorthand: "",
		CapOpt:    "--ieee80211ac",
		Type:      "bool",
		Usage:     "Enable IEEE 802.11ac (VHT)",
	},
	{
		Name:      "htcapab",
		Shorthand: "",
		CapOpt:    "--ht_capab",
		Type:      "string",
		Usage:     "HT capabilities (default: [HT40+])",
	},
	{
		Name:      "vhtcapab",
		Shorthand: "",
		CapOpt:    "--vht_capab",
		Type:      "string",
		Usage:     "VHT capabilities",
	},
	{
		Name:      "country",
		Shorthand: "",
		CapOpt:    "--country",
		Type:      "string",
		Usage:     "Set two-letter country code for regularity (example: US)",
	},
	{
		Name:      "freq-band",
		Shorthand: "",
		CapOpt:    "--freq-band",
		Type:      "string",
		Usage:     "Set frequency band. Valid inputs: 2.4, 5 (default: 2.4)",
	},
	{
		Name:      "diver",
		Shorthand: "",
		CapOpt:    "--diver",
		Type:      "string",
		Usage:     "Choose your WiFi adapter driver (default: nl80211)",
	},
	{
		Name:      "no-virt",
		Shorthand: "",
		CapOpt:    "--no-virt",
		Type:      "bool",
		Usage:     "Do not create virtual interface",
	},
	{
		Name:      "no-haveged",
		Shorthand: "",
		CapOpt:    "--no-haveged",
		Type:      "bool",
		Usage:     "Do not run 'haveged' automatically when needed",
	},
	{
		Name:      "fix-unmanaged",
		Shorthand: "",
		CapOpt:    "--fix-unmanaged",
		Type:      "bool",
		Usage:     "If NetworkManager shows your interface as unmanaged after you close create_ap, then use this option to switch your interface back to managed",
	},
	{
		Name:      "mac",
		Shorthand: "",
		CapOpt:    "--mac",
		Type:      "string",
		Usage:     "Set MAC address",
	},
	{
		Name:      "dhcp-dns",
		Shorthand: "",
		CapOpt:    "--dhcp-dns",
		Type:      "string",
		Usage:     "Set DNS returned by DHCP: --dhcp-dns <IP1[,IP2]>",
	},
	{
		Name:      "no-dns",
		Shorthand: "",
		CapOpt:    "--no-dns",
		Type:      "bool",
		Usage:     "Disable dnsmasq DNS server",
	},
	{
		Name:      "no-dnsmasq",
		Shorthand: "",
		CapOpt:    "--no-dnsmasq",
		Type:      "bool",
		Usage:     "Disable dnsmasq server completely",
	},
	{
		Name:      "gateway",
		Shorthand: "",
		CapOpt:    "-g",
		Type:      "string",
		Usage:     "IPv4 Gateway for the Access Point (default: 192.168.12.1)",
	},
	{
		Name:      "etc-hosts",
		Shorthand: "",
		CapOpt:    "-d",
		Type:      "bool",
		Usage:     "Disable dnsmasq server completely",
	},
	{
		Name:      "addn-hosts",
		Shorthand: "",
		CapOpt:    "-e",
		Type:      "bool",
		Usage:     "DNS server will take into account additional hosts file",
	},
}
