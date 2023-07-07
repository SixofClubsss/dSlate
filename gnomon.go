package main

import (
	"strconv"

	"github.com/civilware/Gnomon/structures"
	"github.com/dReam-dApps/dReams/menu"
	"github.com/sirupsen/logrus"
)

var logger = structures.Logger.WithFields(logrus.Fields{})

func searchByKey(scid string, key string, s bool) string {
	if menu.Gnomes.Init {
		var sValue []string
		var uValue []uint64

		if s {
			sValue, uValue = menu.Gnomes.GetSCIDValuesByKey(scid, key)
		} else {
			if i, err := strconv.Atoi(key); err != nil {
				logger.Errorln("[dSlate]", err)
			} else {
				sValue, uValue = menu.Gnomes.GetSCIDValuesByKey(scid, uint64(i))
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
			sValue, uValue = menu.Gnomes.GetSCIDKeysByValue(scid, value)
		} else {
			if i, err := strconv.Atoi(value); err != nil {
				logger.Errorln("[dSlate]", err)
			} else {
				sValue, uValue = menu.Gnomes.GetSCIDKeysByValue(scid, uint64(i))
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
