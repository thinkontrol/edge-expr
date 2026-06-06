package rvnode

import "net"

type NetInterface struct {
	Name  string   `json:"name" mapstructure:"name"`
	Mac   string   `json:"mac" mapstructure:"mac"`
	Addrs []string `json:"addrs" mapstructure:"addrs"`
	IP    net.IP   `json:"ip" mapstructure:"ip"`
}

func (n *NetInterface) Equal(other *NetInterface) bool {
	if n == nil && other == nil {
		return true
	}
	if n == nil || other == nil {
		return false
	}
	return n.Name == other.Name &&
		n.Mac == other.Mac &&
		sliceStringEqual(n.Addrs, other.Addrs) &&
		n.IP.Equal(other.IP)
}

func sliceStringEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type GSMInfo struct {
	IpV4Up       bool   `json:"ipv4_up" mapstructure:"ipv4_up"`
	IpV6Up       bool   `json:"ipv6_up" mapstructure:"ipv6_up"`
	SerialNumber string `json:"serial_number" mapstructure:"serial_number"`
	IMEI         string `json:"imei" mapstructure:"imei"`
	ICCID        string `json:"iccid" mapstructure:"iccid"`
	RSSI         int    `json:"rssi" mapstructure:"rssi"`
	State        string `json:"state" mapstructure:"state"`
	CNum         string `json:"cnum" mapstructure:"cnum"`
	StateNum     int    `json:"state_num" mapstructure:"state_num"`
	Operator     string `json:"operator" mapstructure:"operator"`
	AccessTech   string `json:"access_tech" mapstructure:"access_tech"`
	Ipv4         string `json:"ipv4" mapstructure:"ipv4"`
	Ipv6         string `json:"ipv6" mapstructure:"ipv6"`
}

func (g *GSMInfo) Equal(other *GSMInfo) bool {
	if g == nil && other == nil {
		return true
	}
	if g == nil || other == nil {
		return false
	}
	return g.IpV4Up == other.IpV4Up &&
		g.IpV6Up == other.IpV6Up &&
		g.SerialNumber == other.SerialNumber &&
		g.IMEI == other.IMEI &&
		g.ICCID == other.ICCID &&
		g.RSSI == other.RSSI &&
		g.State == other.State &&
		g.CNum == other.CNum &&
		g.StateNum == other.StateNum &&
		g.Operator == other.Operator &&
		g.AccessTech == other.AccessTech &&
		g.Ipv4 == other.Ipv4 &&
		g.Ipv6 == other.Ipv6
}

type WifiInfo struct {
	SSID        string `json:"ssid" mapstructure:"ssid"`
	Quality     int    `json:"quality" mapstructure:"quality"`
	SignalLevel int    `json:"signal_level" mapstructure:"signal_level"`
	NoiseLevel  int    `json:"noise_level" mapstructure:"noise_level"`
}

func (w *WifiInfo) Equal(other *WifiInfo) bool {
	if w == nil && other == nil {
		return true
	}
	if w == nil || other == nil {
		return false
	}
	return w.SSID == other.SSID &&
		w.Quality == other.Quality &&
		w.SignalLevel == other.SignalLevel &&
		w.NoiseLevel == other.NoiseLevel
}

type LocalBridge struct {
	DeviceName string `json:"device_name" mapstructure:"DeviceName_str"`
	HubNameLB  string `json:"hub_name_lb" mapstructure:"HubNameLB_str"`
	Online     bool   `json:"online" mapstructure:"Online_bool"`
	Active     bool   `json:"active" mapstructure:"Active_bool"`
	TapMode    bool   `json:"tap_mode" mapstructure:"TapMode_bool"`
}

type LocalBridgeList struct {
	LocalBridgeList []*LocalBridge `json:"local_bridge_list" mapstructure:"LocalBridgeList"`
}

type Hub struct {
	HubName string `json:"hub_name" mapstructure:"HubName_str"`
	Online  bool   `json:"online" mapstructure:"Online_bool"`
}

type HubList struct {
	NumHub  int    `json:"num_hub" mapstructure:"NumHub_u32"`
	HubList []*Hub `json:"hub_list" mapstructure:"HubList"`
}

type LocalBridgeSupport struct {
	IsBridgeSupportedOs bool `json:"is_bridge_supported_os" mapstructure:"IsBridgeSupportedOs_bool"`
	IsWinPcapNeeded     bool `json:"is_win_pcap_needed" mapstructure:"IsWinPcapNeeded_bool"`
}

type Ethernet struct {
	DeviceName            string `json:"device_name" mapstructure:"DeviceName_str"`
	NetworkConnectionName string `json:"network_connection_name" mapstructure:"NetworkConnectionName_utf"`
}

type EthernetList struct {
	EthList []*Ethernet
}

type Config struct {
	Host     string `json:"host" mapstructure:"host"`
	HubName  string `json:"hub_name" mapstructure:"hub_name"`
	Port     int    `json:"port" mapstructure:"port"`
	Password string `json:"password" mapstructure:"password"`
}

type Host struct {
	Ipv4   string   `json:"ipv4" mapstructure:"ipv4"`
	Mac    string   `json:"mac" mapstructure:"mac"`
	Vendor string   `json:"vendor" mapstructure:"vendor"`
	Names  []string `json:"names" mapstructure:"names"`
}

type VPNConnectData struct {
	HubHost  string `json:"hub_host" mapstructure:"hub_host"`
	HubPort  int    `json:"hub_port" mapstructure:"hub_port"`
	HubName  string `json:"hub_name" mapstructure:"hub_name"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}
