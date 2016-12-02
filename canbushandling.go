package main

import (
	"fmt"
	CAN "github.com/brendoncarroll/go-can"
	"log"
	"sync"
)

var cb *CAN.CANBus
var csi []uint32 // subscribed IDs slice
var csi_lock sync.Mutex // CAN subscribed IDs Mutex

func CANStart(iface string) {
	fmt.Printf("canbushandler: Initialisiere CAN-Bus Interface %s\n", iface)
	var err error
	cb, err = CAN.NewCANBus(iface)
	if err != nil {
		fmt.Printf("canbushandler: Error beim aktivieren von CAN-Bus Interface %s\n", iface)
		log.Fatal(err)
	}
	fmt.Printf("canbushandler: Interface %s erfolgreich initialisiert.\n", iface)
	var cf CAN.CANFrame
	fmt.Printf("canbushadler: Betrete jetzt Endlosschleife\n")
	for {
		cb.Read(&cf)
		fmt.Printf("canbushandler: CAN-Frame empfangen (ID:%d). Lock Mutex\n", cf.ID)
		csi_lock.Lock()
		fmt.Printf("canbushandler: Mutex erfolreich gelockt.\n")
		var id_sub = false // zeigt an ob die ID subscribed war oder nicht
		for _, i := range csi {
			if i == cf.ID {
				fmt.Printf("canbushandler: ID %d ist abonniert starte receivehandler\n", cf.ID)
				handleCAN(cf)
				id_sub = true
				break
			}
		}
		if !id_sub {
			fmt.Printf("canbushandler: ID:%d war nicht abonniert. /dev/nulled that frame...\n", cf.ID)
		}
		csi_lock.Unlock()
		fmt.Printf("canbushandler: Mutex unlocked.\n")
	}
}

func CANSubscribe(id uint32) {
	csi_lock.Lock()
	csi = append(csi, id)
	csi_lock.Unlock()
	fmt.Printf("canbushandler: Mutex lock&unlock successful Subscribed to ID:%d\n", id)
}

func CANUnsubscribe(id uint32) {
	var tmp []uint32
	csi_lock.Lock()
	for _, elem := range csi {
		if elem != id {
			tmp = append(tmp, elem)
		}
	}
	csi = tmp
	csi_lock.Unlock()
	fmt.Printf("canbushandler: Mutex lock&unlock successful unsubscribed ID:%d\n", id)
}

func CANPublish(cf CAN.CANFrame) {
	CANUnsubscribe(cf.ID)
	fmt.Println("canbushandler: Sende CAN-Frame: ", cf)
	err := cb.Write(&cf)
	if err != nil {
		fmt.Printf("canbushandler: Error beim Senden des Frames\n")
                log.Fatal(err)
	}
	CANSubscribe(cf.ID)
}
