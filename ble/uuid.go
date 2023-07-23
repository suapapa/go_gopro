package ble

import (
	"fmt"
	"strings"
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
	svcUUIDGoProWifiAccessPoint    = gpUUID("GP-0001")
	chrUUIDWifiAccessPointSSID     = gpUUID(WifiAccessPointSSID)     // Read / Write
	chrUUIDWifiAccessPointPassword = gpUUID(WifiAccessPointPassword) // Read / Write
	chrUUIDWifiAccessPointPower    = gpUUID(WifiAccessPointPower)    // Write
	chrUUIDWifiAccessPointState    = gpUUID(WifiAccessPointState)    // Read / Indicate

	svcGoProCamaraManagement         = gpUUID("GP-0090")
	chrUUIDNetworkManagementCommand  = gpUUID(NetworkManagementCommand)  // Write
	chrUUIDNetworkManagementResponse = gpUUID(NetworkManagementResponse) // Notify

	svcUUIDControlAndQuery = "FEA6"
	chrUUIDCommand         = gpUUID(Command)         // Write
	chrUUIDCommandResponse = gpUUID(CommandResponse) // Notify
	chrUUIDSetting         = gpUUID(Setting)         // Write
	chrUUIDSettingResponse = gpUUID(SettingResponse) // Notify
	chrUUIDQuery           = gpUUID(Query)           // Write
	chrUUIDQueryResponse   = gpUUID(QueryResponse)   // Notify
)

func gpUUID(id uuid) string {
	uuidStr := strings.Replace(string(id), "GP-", "", -1)
	// return goble.MustParse(fmt.Sprintf("b5f9%s-aa8d-11e3-9046-0002a5d5c51b", uuidStr))
	return fmt.Sprintf("b5f9%s-aa8d-11e3-9046-0002a5d5c51b", uuidStr)
}
