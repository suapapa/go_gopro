package ble

/*
func explore(cln goble.Client, p *goble.Profile) string {
	var ret string
	for _, s := range p.Services {
		ret += fmt.Sprintf("    Service: %s %s, Handle (0x%02X)\n", s.UUID, goble.Name(s.UUID), s.Handle)

		for _, c := range s.Characteristics {
			ret += fmt.Sprintf("      Characteristic: %s %s, Property: 0x%02X (%s), Handle(0x%02X), VHandle(0x%02X)\n",
				c.UUID, goble.Name(c.UUID), c.Property, propString(c.Property), c.Handle, c.ValueHandle)
			if (c.Property & goble.CharRead) != 0 {
				b, err := cln.ReadCharacteristic(c)
				if err != nil {
					ret += fmt.Sprintf("Failed to read chr: %s\n", err)
					continue
				}
				ret += fmt.Sprintf("        Value         %x | %q\n", b, b)
			}

			for _, d := range c.Descriptors {
				ret += fmt.Sprintf("        Descriptor: %s %s, Handle(0x%02x)\n", d.UUID, goble.Name(d.UUID), d.Handle)
				b, err := cln.ReadDescriptor(d)
				if err != nil {
					ret += fmt.Sprintf("Failed to read descriptor: %s\n", err)
					continue
				}
				ret += fmt.Sprintf("        Value         %x | %q\n", b, b)
			}

			var sub time.Duration
			if sub != 0 {
				// Don't bother to subscribe the Service Changed chrs.
				if c.UUID.Equal(goble.ServiceChangedUUID) {
					continue
				}

				// Don't touch the Apple-specific Service/Characteristic.
				// Service: D0611E78BBB44591A5F8487910AE4366
				// Characteristic: 8667556C9A374C9184ED54EE27D90049, Property: 0x18 (WN),
				//   Descriptor: 2902, Client Characteristic Configuration
				//   Value         0000 | "\x00\x00"
				if c.UUID.Equal(goble.MustParse("8667556C9A374C9184ED54EE27D90049")) {
					continue
				}

				if (c.Property & goble.CharNotify) != 0 {
					ret += fmt.Sprintf("\n-- Subscribe to notification for %s --\n", sub)
					h := func(req []byte) { ret += fmt.Sprintf("Notified: %q [ % X ]\n", string(req), req) }
					if err := cln.Subscribe(c, false, h); err != nil {
						log.Fatalf("subscribe failed: %s", err)
					}
					time.Sleep(sub)
					if err := cln.Unsubscribe(c, false); err != nil {
						log.Fatalf("unsubscribe failed: %s", err)
					}
					ret += fmt.Sprintf("-- Unsubscribe to notification --\n")
				}
				if (c.Property & goble.CharIndicate) != 0 {
					ret += fmt.Sprintf("\n-- Subscribe to indication of %s --\n", sub)
					h := func(req []byte) { ret += fmt.Sprintf("Indicated: %q [ % X ]\n", string(req), req) }
					if err := cln.Subscribe(c, true, h); err != nil {
						log.Fatalf("subscribe failed: %s", err)
					}
					time.Sleep(sub)
					if err := cln.Unsubscribe(c, true); err != nil {
						log.Fatalf("unsubscribe failed: %s", err)
					}
					ret += fmt.Sprintf("-- Unsubscribe to indication --\n")
				}
			}
		}
		ret += fmt.Sprintf("\n")
	}
	return ret
}

func propString(p goble.Property) string {
	var s string
	for k, v := range map[goble.Property]string{
		goble.CharBroadcast:   "B",
		goble.CharRead:        "R",
		goble.CharWriteNR:     "w",
		goble.CharWrite:       "W",
		goble.CharNotify:      "N",
		goble.CharIndicate:    "I",
		goble.CharSignedWrite: "S",
		goble.CharExtended:    "E",
	} {
		if p&k != 0 {
			s += v
		}
	}
	return s
}
*/
