package main

import (
	"log"
	"strconv"

	"github.com/SixofClubsss/dReams/menu"
)

func searchByKey(scid string, key string, s bool) string {
	if menu.Gnomes.Init {
		var sValue []string
		var uValue []uint64

		if s {
			sValue, uValue = menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, key, menu.Gnomes.Indexer.ChainHeight, true)
		} else {
			if i, err := strconv.Atoi(key); err != nil {
				log.Println("[dSlate]", err)
			} else {
				sValue, uValue = menu.Gnomes.Indexer.Backend.GetSCIDValuesByKey(scid, uint64(i), menu.Gnomes.Indexer.ChainHeight, true)
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
	if menu.Gnomes.Init {
		var sValue []string
		var uValue []uint64
		if s {
			sValue, uValue = menu.Gnomes.Indexer.Backend.GetSCIDKeysByValue(scid, value, menu.Gnomes.Indexer.ChainHeight, true)
		} else {
			if i, err := strconv.Atoi(value); err != nil {
				log.Println("[dSlate]", err)
			} else {
				sValue, uValue = menu.Gnomes.Indexer.Backend.GetSCIDKeysByValue(scid, uint64(i), menu.Gnomes.Indexer.ChainHeight, true)
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
