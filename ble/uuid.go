package ble

import (
	"fmt"
	"strings"

	goble "github.com/go-ble/ble"
)

type Characteristic string

const (
	WifiAccessPointSSID       Characteristic = "GP-0002"
	WifiAccessPointPassword   Characteristic = "GP-0003"
	WifiAccessPointPower      Characteristic = "GP-0004"
	WifiAccessPointState      Characteristic = "GP-0005"
	NetworkManagementCommand  Characteristic = "GP-0091"
	NetworkManagementResponse Characteristic = "GP-0092"
	Command                   Characteristic = "GP-0072"
	CommandResponse           Characteristic = "GP-0073"
	Setting                   Characteristic = "GP-0074"
	SettingResponse           Characteristic = "GP-0075"
	Query                     Characteristic = "GP-0076"
	QueryResponse             Characteristic = "GP-0077"
)

var (
	// services and characteristics
	svcGoProWifiAccessPoint    = gpUUID("GP-0001")
	chrWifiAccessPointSSID     = gpUUID(WifiAccessPointSSID)     // Read / Write
	chrWifiAccessPointPassword = gpUUID(WifiAccessPointPassword) // Read / Write
	chrWifiAccessPointPower    = gpUUID(WifiAccessPointPower)    // Write
	chrWifiAccessPointState    = gpUUID(WifiAccessPointState)    // Read / Indicate

	svcGoProCamaraManagement     = gpUUID("GP-0090")
	chrNetworkManagementCommand  = gpUUID(NetworkManagementCommand)  // Write
	chrNetworkManagementResponse = gpUUID(NetworkManagementResponse) // Notify

	svcUUIDControlAndQuery = goble.UUID16(0xFEA6)
	chrCommand             = gpUUID(Command)         // Write
	chrCommandResponse     = gpUUID(CommandResponse) // Notify
	chrSetting             = gpUUID(Setting)         // Write
	chrSettingResponse     = gpUUID(SettingResponse) // Notify
	chrQuery               = gpUUID(Query)           // Write
	chrQueryResponse       = gpUUID(QueryResponse)   // Notify
)

func gpUUID(uuid Characteristic) goble.UUID {
	uuidStr := strings.Replace(string(uuid), "GP-", "", -1)
	return goble.MustParse(fmt.Sprintf("b5f9%s-aa8d-11e3-9046-0002a5d5c51b", uuidStr))
}

func makeCharacteristicMap(p *goble.Profile) map[Characteristic]*goble.Characteristic {
	chrs := make(map[Characteristic]*goble.Characteristic)
	for _, s := range p.Services {
		for _, c := range s.Characteristics {
			switch {
			case c.UUID.Equal(chrWifiAccessPointSSID):
				chrs[WifiAccessPointSSID] = c
			case c.UUID.Equal(chrWifiAccessPointPassword):
				chrs[WifiAccessPointPassword] = c
			case c.UUID.Equal(chrWifiAccessPointPower):
				chrs[WifiAccessPointPower] = c
			case c.UUID.Equal(chrWifiAccessPointState):
				chrs[WifiAccessPointState] = c

			case c.UUID.Equal(chrNetworkManagementCommand):
				chrs[NetworkManagementCommand] = c
			case c.UUID.Equal(chrNetworkManagementResponse):
				chrs[NetworkManagementResponse] = c

			case c.UUID.Equal(chrCommand):
				chrs[Command] = c
			case c.UUID.Equal(chrCommandResponse):
				chrs[CommandResponse] = c
			case c.UUID.Equal(chrSetting):
				chrs[Setting] = c
			case c.UUID.Equal(chrSettingResponse):
				chrs[SettingResponse] = c
			case c.UUID.Equal(chrQuery):
				chrs[Query] = c
			case c.UUID.Equal(chrQueryResponse):
				chrs[QueryResponse] = c
			}
		}
	}
	return chrs
}
