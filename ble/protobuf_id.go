package ble

const (
	// Newtwork Management actions will be sent to GP-0091 and
	// responses will be received from GP-0092
	featureNetworkManagement = 0x02
	actionStartScan          = 0x02
	responseStartScan        = 0x82
	responseAsyncStartScan   = 0x0B
	actionGetApEntries       = 0x03
	responseGetApEntries     = 0x83
	actionConnect            = 0x04
	responseConnect          = 0x84
	responseAsyncConnect     = 0x0C
	actionConnectNew         = 0x05
	responseConnectNew       = 0x85
	responseAsyncConnectNew  = 0x0C

	// Command actions will be sent to GP-0072 and
	// responses will be received from GP-0073
	featureCommand                = 0xF1
	actionSetCameraContolStatus   = 0x69
	responseSetCameraContolStatus = 0xE9
	actionSetTurboActive          = 0x6B
	responseSetTurboActive        = 0xEB
	actionReleaseNetwork          = 0x78
	responseReleaseNetwork        = 0xF8
	actionSetLiveStream           = 0x79
	responseSetLiveStream         = 0xF9

	// Query actions will be sent to GP-0076 and
	// responses will be received from GP-0077
	featureQuery                     = 0xF5
	actionGetPresetStatus            = 0x72
	responseGetPresetStatus          = 0xF2
	responseAsyncGetPresetStatus     = 0xF3
	actionGetLiveStreamStatus        = 0x74
	responseGetLiveStreamStatus      = 0xF4
	responseAsyncGetLiveStreamStatus = 0xF5
)
