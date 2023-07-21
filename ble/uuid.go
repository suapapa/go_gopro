package ble

import (
	"fmt"
	"strings"

	goble "github.com/go-ble/ble"
)

type uuid string

const (
	WifiAccessPointSSID       uuid = "GP-0002"
	WifiAccessPointPassword   uuid = "GP-0003"
	WifiAccessPointPower      uuid = "GP-0004"
	WifiAccessPointState      uuid = "GP-0005"
	NetworkManagementCommand  uuid = "GP-0091"
	NetworkManagementResponse uuid = "GP-0092"
	Command                   uuid = "GP-0072"
	CommandResponse           uuid = "GP-0073"
	Setting                   uuid = "GP-0074"
	SettingResponse           uuid = "GP-0075"
	Query                     uuid = "GP-0076"
	QueryResponse             uuid = "GP-0077"
)

var (
	// services and characteristics
	// svcGoProWifiAccessPoint    = gpUUID("GP-0001")
	// chrWifiAccessPointSSID     = gpUUID(WifiAccessPointSSID)     // Read / Write
	// chrWifiAccessPointPassword = gpUUID(WifiAccessPointPassword) // Read / Write
	// chrWifiAccessPointPower    = gpUUID(WifiAccessPointPower)    // Write
	// chrWifiAccessPointState    = gpUUID(WifiAccessPointState)    // Read / Indicate

	// svcGoProCamaraManagement     = gpUUID("GP-0090")
	// chrNetworkManagementCommand  = gpUUID(NetworkManagementCommand)  // Write
	// chrNetworkManagementResponse = gpUUID(NetworkManagementResponse) // Notify

	svcUUIDControlAndQuery = goble.UUID16(0xFEA6)
	// chrCommand             = gpUUID(Command)         // Write
	// chrCommandResponse     = gpUUID(CommandResponse) // Notify
	// chrSetting             = gpUUID(Setting)         // Write
	// chrSettingResponse     = gpUUID(SettingResponse) // Notify
	// chrQuery               = gpUUID(Query)           // Write
	// chrQueryResponse       = gpUUID(QueryResponse)   // Notify
)

func gpUUID(id uuid) goble.UUID {
	uuidStr := strings.Replace(string(id), "GP-", "", -1)
	return goble.MustParse(fmt.Sprintf("b5f9%s-aa8d-11e3-9046-0002a5d5c51b", uuidStr))
}
