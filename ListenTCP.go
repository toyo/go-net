package net

import (
	traditionalnat "net"
	"strconv"
	"strings"

	"github.com/huin/goupnp/dcps/internetgateway1"
	"github.com/pkg/errors"
)

// ListenTCP is listener which support UPnP
func ListenTCP(network string, laddr *traditionalnat.TCPAddr) (l *traditionalnat.TCPListener, err error) {
	if network != `tcp` {
		return nil, errors.New(`NAT: Not TCP`)
	}

	l, err = traditionalnat.ListenTCP("tcp", laddr)
	if err != nil {
		return
	}

	addrs, err := traditionalnat.InterfaceAddrs()
	for i := 0; i < len(addrs); i++ {
		if !strings.HasPrefix(addrs[i].String(), `192.168`) {
			addrs = append(addrs[:i], addrs[i+1:]...)
			i--
		}
	}
	if len(addrs) != 0 {
		if addrmask := strings.Split(addrs[0].String(), `/`); len(addrmask) == 2 {
			myaddr := addrmask[0]
			clients, errorss, err := internetgateway1.NewWANPPPConnection1Clients()
			if err == nil {
				if len(errorss) == 0 {
					if len(clients) == 1 {
						scpd, err := clients[0].ServiceClient.Service.RequestSCPD()
						if err == nil {
							if scpd == nil || scpd.GetAction("AddPortMapping") != nil {
								err = clients[0].AddPortMapping("", uint16(laddr.Port) /* internal */, "TCP", uint16(laddr.Port) /* external */, myaddr, true, "Test port mapping", 86400 /* sec */)
								if err == nil {
									logln("[DEBUG] AddPortMapping: ", laddr.Port)
									laddr0, _ := traditionalnat.ResolveTCPAddr(`tcp`, `8.8.8.8:6911`)
									laddr0.Port = laddr.Port
									*laddr = *laddr0
									logln(`[DEBUG] Use NAT: Listen on `, laddr, ` (8.8.8.8 is fake)`)
									return l, nil
								} else {
									logln(`NAT: No port on global address available`, err)
								}
							} else {
								logln("NAT: AddPortMapping not exist.")
							}
						} else {
							logln("NAT: Error requesting service SCPD", err)
						}
					} else {
						logln("NAT: Too many/No UPnP Router: " + strconv.Itoa(len(clients)))
					}
				} else {
					logln("NAT: UPnP Error", errorss)
				}
			} else {
				logln("NAT: Error on UPnP device discovering: ", err)
			}
		} else {
			logln(`NAT: IP address of this machine is unexpected IP/Mask ` + addrs[0].String())
		}
	}
	laddr0, _ := traditionalnat.ResolveTCPAddr(`tcp`, l.Addr().String())
	*laddr = *laddr0
	logln(`[DEBUG] Without NAT: Listen on `, laddr)
	return l, nil
}
