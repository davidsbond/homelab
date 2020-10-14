package synology

import (
	"fmt"
)

// GetSystemInfoResponse is the response DTO when calling the system info endpoint on the
// NAS.
type GetSystemInfoResponse struct {
	ApacheIsDefault bool `json:"apache_is_default"`
	BsThrExceed     bool `json:"bs_thr_exceed"`
	Data            struct {
		UpsEnable bool `json:"ups_enable"`
	} `json:"data"`
	Disks []struct {
		AdvProgress        string `json:"adv_progress"`
		AdvStatus          string `json:"adv_status"`
		BelowRemainLifeThr bool   `json:"below_remain_life_thr"`
		Compatibility      string `json:"compatibility"`
		Container          struct {
			Order                int    `json:"order"`
			Str                  string `json:"str"`
			SupportPwrBtnDisable bool   `json:"supportPwrBtnDisable"`
			Type                 string `json:"type"`
		} `json:"container"`
		Device             string `json:"device"`
		DisableSecera      bool   `json:"disable_secera"`
		DiskType           string `json:"diskType"`
		DiskCode           string `json:"disk_code"`
		EraseTime          int    `json:"erase_time"`
		ExceedBadSectorThr bool   `json:"exceed_bad_sector_thr"`
		Firm               string `json:"firm"`
		HasSystem          bool   `json:"has_system"`
		ID                 string `json:"id"`
		IhmTesting         bool   `json:"ihm_testing"`
		Is4Kn              bool   `json:"is4Kn"`
		IsSsd              bool   `json:"isSsd"`
		IsSynoPartition    bool   `json:"isSynoPartition"`
		IsErasing          bool   `json:"is_erasing"`
		LongName           string `json:"longName"`
		Model              string `json:"model"`
		Name               string `json:"name"`
		NumID              int    `json:"num_id"`
		Order              int    `json:"order"`
		OverviewStatus     string `json:"overview_status"`
		PciSlot            int    `json:"pciSlot"`
		PerfTesting        bool   `json:"perf_testing"`
		PortType           string `json:"portType"`
		RemainLife         int    `json:"remain_life"`
		RemoteInfo         struct {
			Compatibility string `json:"compatibility"`
			Unc           int    `json:"unc"`
		} `json:"remote_info"`
		Serial          string  `json:"serial"`
		SizeTotal       string  `json:"size_total"`
		SmartProgress   string  `json:"smart_progress"`
		SmartStatus     string  `json:"smart_status"`
		SmartTestLimit  int     `json:"smart_test_limit"`
		SmartTesting    bool    `json:"smart_testing"`
		Status          string  `json:"status"`
		Support         bool    `json:"support"`
		Temp            float64 `json:"temp"`
		TestingProgress string  `json:"testing_progress"`
		TestingType     string  `json:"testing_type"`
		TrayStatus      string  `json:"tray_status"`
		Unc             int     `json:"unc"`
		UsedBy          string  `json:"used_by"`
		Vendor          string  `json:"vendor"`
	} `json:"disks"`
	DNS        string      `json:"dns"`
	EsataVols  interface{} `json:"esata_vols"`
	Gateway    string      `json:"gateway"`
	Interfaces []struct {
		ID     string `json:"id"`
		Ipaddr string `json:"ipaddr"`
		Ipv6   []struct {
			IPv6Addr   string `json:"IPv6Addr"`
			PrefixLeng int    `json:"PrefixLeng"`
			Scope      string `json:"Scope"`
		} `json:"ipv6"`
		Mac  string `json:"mac"`
		Mask string `json:"mask"`
		Type string `json:"type"`
	} `json:"interfaces"`
	IsSystemCrashed bool `json:"is_system_crashed"`
	MultiPower      struct {
		SupportRp bool `json:"support_rp"`
	} `json:"multiPower"`
	Optime             string      `json:"optime"`
	SdcardVols         interface{} `json:"sdcard_vols"`
	SecurityScan       string      `json:"securityScan"`
	Systempwarn        bool        `json:"systempwarn"`
	TemperatureWarning bool        `json:"temperature_warning"`
	UsbVols            interface{} `json:"usb_vols"`
	VolInfo            []struct {
		Desc       string `json:"desc"`
		InodeFree  string `json:"inode_free"`
		InodeTotal string `json:"inode_total"`
		Name       string `json:"name"`
		Status     string `json:"status"`
		TotalSize  string `json:"total_size"`
		UsedSize   string `json:"used_size"`
		Volume     string `json:"volume"`
	} `json:"vol_info"`
	VolWarnings struct {
		DiskInodeWarningPercent       string `json:"disk_inode_warning_percent"`
		DiskWarningPercent            string `json:"disk_warning_percent"`
		EsataPartitionWarningPercent  string `json:"esata_partition_warning_percent"`
		SdcardPartitionWarningPercent string `json:"sdcard_partition_warning_percent"`
		UsbPartitionWarningPercent    string `json:"usb_partition_warning_percent"`
	} `json:"vol_warnings"`
}

// The APIError type is a combination of different possible API error formats the API can return. It isn't
// standardised, this is a best effort to catch errors.
type APIError struct {
	Err struct {
		Code int `json:"code"`
	} `json:"error"`
	Errno struct {
		Key     string `json:"key"`
		Section string `json:"section"`
	} `json:"errno"`
	Success bool `json:"success"`
}

func (a APIError) Error() string {
	if a.Err.Code > 0 {
		return fmt.Sprintf("error code: %v", a.Err.Code)
	}

	return a.Errno.Key
}

// LoginResponse is the response DTO after successfully authenticating.
type LoginResponse struct {
	Data struct {
		Sid string `json:"sid"`
	} `json:"data"`
	Success bool `json:"success"`
}
