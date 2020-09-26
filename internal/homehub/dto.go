package homehub

import "encoding/xml"

type (
	// The GetStatusResponse type represents the wan_conn.xml response given
	// by the BT home hub.
	GetStatusResponse struct {
		XMLName           xml.Name `xml:"status"`
		Text              string   `xml:",chardata"`
		WanConnStatusList struct {
			Text  string `xml:",chardata"`
			Type  string `xml:"type,attr"`
			Value string `xml:"value,attr"`
		} `xml:"wan_conn_status_list"`
		WanConnVolumeList struct {
			Text  string `xml:",chardata"`
			Type  string `xml:"type,attr"`
			Value string `xml:"value,attr"`
		} `xml:"wan_conn_volume_list"`
		WanLinestatusRateList struct {
			Text  string `xml:",chardata"`
			Type  string `xml:"type,attr"`
			Value string `xml:"value,attr"`
		} `xml:"wan_linestatus_rate_list"`
		WlanChannelList struct {
			Text  string `xml:",chardata"`
			Type  string `xml:"type,attr"`
			Value string `xml:"value,attr"`
		} `xml:"wlan_channel_list"`
		Curlinkstatus struct {
			Text  string `xml:",chardata"`
			Type  string `xml:"type,attr"`
			Value string `xml:"value,attr"`
		} `xml:"curlinkstatus"`
		Sysuptime struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"sysuptime"`
		StatusRate struct {
			Text  string `xml:",chardata"`
			Type  string `xml:"type,attr"`
			Value string `xml:"value,attr"`
		} `xml:"status_rate"`
		WanActiveIdx struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"wan_active_idx"`
		LinkStatus struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"link_status"`
		IP4InfoList struct {
			Text  string `xml:",chardata"`
			Type  string `xml:"type,attr"`
			Value string `xml:"value,attr"`
		} `xml:"ip4_info_list"`
		IP6LLAList struct {
			Text  string `xml:",chardata"`
			Type  string `xml:"type,attr"`
			Value string `xml:"value,attr"`
		} `xml:"ip6_lla_list"`
		IP6GUAList struct {
			Text  string `xml:",chardata"`
			Type  string `xml:"type,attr"`
			Value string `xml:"value,attr"`
		} `xml:"ip6_gua_list"`
		IP6RDNSList struct {
			Text  string `xml:",chardata"`
			Type  string `xml:"type,attr"`
			Value string `xml:"value,attr"`
		} `xml:"ip6_rdns_list"`
		Locktime struct {
			Text  string `xml:",chardata"`
			Value string `xml:"value,attr"`
		} `xml:"locktime"`
	}
)
