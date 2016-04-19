package app

import (
	"time"
	"encoding/binary"
	"github.com/revel/revel"
	"golang.org/x/net/icmp"
)

type ServerPinger struct {
	
}

type Endpoint struct {
	Id			int
	HostId 		int
	Address		[]byte
	DateAdded	int
	Status		string
}

func NewServerPinger() *ServerPinger {
	return &ServerPinger{}
}

func (this *ServerPinger) StartPolling() {
	go func() {
		for {
			sem := make(chan bool)
			t := int(time.Now().Unix())
			go this.pollDatabase(sem)

			select {
			case <-sem:
				d := int(time.Now().Unix()) - t
				time.Sleep(time.Duration(60 - d) * time.Second)
				break
			case <-time.After(60 * time.Second):
				revel.ERROR.Printf("Is the server overloaded? Couldn't finish a single ping iteration in 1 minute!\n")
				// Drain sem to prevent memory leaks
				go func() {
					<-sem
				}()
				break
			}
		}
	}()
}

func (this *ServerPinger) pollDatabase(sem chan bool) {
	hosts, err := DB.Query("SELECT id FROM hosts")
	if err != nil {
		revel.ERROR.Printf("Failed to query hosts: %s\n", err.Error())
		return
	}

	for hosts.Next() {
		var hostId int
		hosts.Scan(&hostId)

		eps, err := DB.Query("SELECT id, address, date_added, status FROM endpoints WHERE host_id = ?", hostId)
		if err != nil {
			revel.ERROR.Printf("Failed to query endpoints: %s\n", err.Error())
			continue
		}
		endpoints := make([]*Endpoint, 0)

		for eps.Next() {
			e := &Endpoint{}
			e.HostId = hostId
			ipAddr := uint32(0)
			eps.Scan(&e.Id, &ipAddr, &e.DateAdded, &e.Status)
			binary.BigEndian.PutUint32(e.Address, ipAddr)
		}

		if len(endpoints) == 0 {
			// Priority scan queue
			return
		}

		go func() {
			for _,ep := range endpoints {
				checkEndpoint(ep)
			}
		}()
	}

	sem <- true
}

func (this *ServerPinger) checkEndpoint(endpoint *Endpoint) {

}