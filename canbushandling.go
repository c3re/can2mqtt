// Package can2mqtt contains some tools for briding a CAN-Interface
// and a mqtt-network
package can2mqtt

import (
	"fmt"
	CAN "github.com/brendoncarroll/go-can"
	"log"
	"sync"
)

var cb *CAN.CANBus
var csi []uint32        // subscribed IDs slice
var csi_lock sync.Mutex // CAN subscribed IDs Mutex

func canStart(iface string) {
	if dbg {
		fmt.Printf("canbushandler: initializing CAN-Bus interface %s\n", iface)
	}
	var err error
	cb, err = CAN.NewCANBus(iface)
	if err != nil {
		if dbg {
			fmt.Printf("canbushandler: error while activating CAN-Bus interface: %s\n", iface)
		}
		log.Fatal(err)
	}
	if dbg {
		fmt.Printf("canbushandler: successfully initialized CAN-Bus interface %s.\n", iface)
	}
	var cf CAN.CANFrame
	if dbg {
		fmt.Printf("canbushadler: entering infinite loop\n")
	}
	for {
		cb.Read(&cf)
		if dbg {
			fmt.Printf("canbushandler: received CAN-Frame: (ID:%d). Locking mutex\n", cf.ID)
		}
		csi_lock.Lock()
		if dbg {
			fmt.Printf("canbushandler: mutex was locked successfully.\n")
		}
		var id_sub = false // indicates, wether the id was subscribed or not
		for _, i := range csi {
			if i == cf.ID {
				if dbg {
					fmt.Printf("canbushandler: ID %d is in subscribed list, calling receivehadler.\n", cf.ID)
				}
				go handleCAN(cf)
				id_sub = true
				break
			}
		}
		if !id_sub {
			if dbg {
				fmt.Printf("canbushandler: ID:%d was not subscribed. /dev/nulled that frame...\n", cf.ID)
			}
		}
		csi_lock.Unlock()
		if dbg {
			fmt.Printf("canbushandler: unlocked mutex.\n")
		}
	}
}

func canSubscribe(id uint32) {
	csi_lock.Lock()
	csi = append(csi, id)
	csi_lock.Unlock()
	if dbg {
		fmt.Printf("canbushandler: mutex lock+unlock successful. subscribed to ID:%d\n", id)
	}
}

func canUnsubscribe(id uint32) {
	var tmp []uint32
	csi_lock.Lock()
	for _, elem := range csi {
		if elem != id {
			tmp = append(tmp, elem)
		}
	}
	csi = tmp
	csi_lock.Unlock()
	if dbg {
		fmt.Printf("canbushandler: mutex lock+unlock successful. unsubscribed ID:%d\n", id)
	}
}

func canPublish(cf CAN.CANFrame) {
	canUnsubscribe(cf.ID)
	if dbg {
		fmt.Println("canbushandler: sending CAN-Frame: ", cf)
	}
	err := cb.Write(&cf)
	if err != nil {
		if dbg {
			fmt.Printf("canbushandler: error while transmitting the CAN-Frame.\n")
		}
		log.Fatal(err)
	}
	canSubscribe(cf.ID)
}
