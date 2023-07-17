package ble

import (
	"fmt"
	"strings"

	goble "github.com/go-ble/ble"
)

var (
	// services and characteristics
	serviceGoProWifiAccessPoint           = gpUUID("GP-0001")
	characteristicWifiAccessPointSSID     = gpUUID("GP-0002") // Read / Write
	characteristicWifiAccessPointPassword = gpUUID("GP-0003") // Read / Write
	characteristicWifiAccessPointPower    = gpUUID("GP-0004") // Write
	characteristicWifiAccessPointState    = gpUUID("GP-0004") // Read / Indicate

	serviceGoProCamaraManagement            = gpUUID("GP-0090")
	characteristicNetworkManagementCommand  = gpUUID("GP-0091") // Write
	characteristicNetworkManagementResponse = gpUUID("GP-0092") // Notify

	serviceUUIDControlAndQuery    = goble.UUID16(0xFEA6)
	characteristicCommand         = gpUUID("GP-0072") // Write
	characteristicCommandResponse = gpUUID("GP-0073") // Notify
	characteristicSetting         = gpUUID("GP-0074") // Write
	characteristicSettingResponse = gpUUID("GP-0075") // Notify
	characteristicQuery           = gpUUID("GP-0076") // Write
	characteristicQueryResponse   = gpUUID("GP-0077") // Notify
)

func gpUUID(uuid string) goble.UUID {
	uuid = strings.Replace(uuid, "GP-", "", -1)
	return goble.MustParse(fmt.Sprintf("b5f9%s-aa8d-11e3-9046-0002a5d5c51b", uuid))
}

type Characteristic byte

const (
	WifiAccessPointSSID Characteristic = iota
	WifiAccessPointPassword
	WifiAccessPointPower
	WifiAccessPointState
	NetworkManagementCommand
	NetworkManagementResponse
	Command
	CommandResponse
	Setting
	SettingResponse
	Query
	QueryResponse
)

func makeCharacteristicMap(p *goble.Profile) map[Characteristic]*goble.Characteristic {
	chrs := make(map[Characteristic]*goble.Characteristic)
	for _, s := range p.Services {
		for _, c := range s.Characteristics {
			switch {
			case c.UUID.Equal(characteristicWifiAccessPointSSID):
				chrs[WifiAccessPointSSID] = c
			case c.UUID.Equal(characteristicWifiAccessPointPassword):
				chrs[WifiAccessPointPassword] = c
			case c.UUID.Equal(characteristicWifiAccessPointPower):
				chrs[WifiAccessPointPower] = c
			case c.UUID.Equal(characteristicWifiAccessPointState):
				chrs[WifiAccessPointState] = c

			case c.UUID.Equal(characteristicNetworkManagementCommand):
				chrs[NetworkManagementCommand] = c
			case c.UUID.Equal(characteristicNetworkManagementResponse):
				chrs[NetworkManagementResponse] = c

			case c.UUID.Equal(characteristicCommand):
				chrs[Command] = c
			case c.UUID.Equal(characteristicCommandResponse):
				chrs[CommandResponse] = c
			case c.UUID.Equal(characteristicSetting):
				chrs[Setting] = c
			case c.UUID.Equal(characteristicSettingResponse):
				chrs[SettingResponse] = c
			case c.UUID.Equal(characteristicQuery):
				chrs[Query] = c
			case c.UUID.Equal(characteristicQueryResponse):
				chrs[QueryResponse] = c
			}
		}
	}
	return chrs
}
