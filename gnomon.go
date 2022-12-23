package main

import (
	"log"
	"strconv"
	"time"

	"github.com/SixofClubsss/dReams/menu"
	"github.com/SixofClubsss/dReams/table"
	"github.com/civilware/Gnomon/indexer"
)

type gnomon struct {
	Start   bool
	Init    bool
	SCIDS   uint64
	Indexer *indexer.Indexer
}

var Gnomes gnomon

func startGnomon(ep string) {
	Gnomes.Start = true
	log.Println("Starting Gnomon.")
	backend := menu.GnomonDB()

	last_height := backend.GetLastIndexHeight()
	daemon_endpoint := ep
	runmode := "daemon"
	mbl := false
	closeondisconnect := false

	table.Assets.Asset_map = make(map[string]string)
	Gnomes.Indexer = indexer.NewIndexer(backend, []string{}, last_height, daemon_endpoint, runmode, mbl, closeondisconnect, true)
	Gnomes.Indexer.StartDaemonMode(1)
	time.Sleep(3 * time.Second)
	Gnomes.Init = true

	Gnomes.Start = false
}

func StopGnomon(gi bool) {
	if gi && !menu.GnomonClosing() {
		log.Println("Putting Gnomon to Sleep.")
		Gnomes.Indexer.Close()
		Gnomes.Init = false
		gnomonEnabled.SetSelected("Off")
		time.Sleep(1 * time.Second)
		log.Println("Gnomon is Sleeping.")
	}
}

func searchByKey(scid string, key string, s bool) string {
	if Gnomes.Init {
		var sValue []string
		var uValue []uint64

		if s {
			sValue, uValue = Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, key, Gnomes.Indexer.ChainHeight, true)
		} else {
			if i, err := strconv.Atoi(key); err != nil {
				log.Println(err)
			} else {
				sValue, uValue = Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, uint64(i), Gnomes.Indexer.ChainHeight, true)
			}
		}

		if sValue != nil {
			return sValue[0]
		}

		if uValue != nil {
			return strconv.Itoa(int(uValue[0]))
		}
	}

	return "nil"
}

func searchByValue(scid string, value string, s bool) string {
	if Gnomes.Init {
		var sValue []string
		var uValue []uint64
		if s {
			sValue, uValue = Gnomes.Indexer.Backend.GetSCIDKeysByValue(scid, value, Gnomes.Indexer.ChainHeight, true)
		} else {
			if i, err := strconv.Atoi(value); err != nil {
				log.Println(err)
			} else {
				sValue, uValue = Gnomes.Indexer.Backend.GetSCIDKeysByValue(scid, uint64(i), Gnomes.Indexer.ChainHeight, true)
			}
		}

		if sValue != nil {
			return sValue[0]
		}

		if sValue != nil {
			return strconv.Itoa(int(uValue[0]))
		}
	}

	return "nil"
}
