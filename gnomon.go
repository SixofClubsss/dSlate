package main

import (
	"strconv"

	"github.com/civilware/Gnomon/structures"
	"github.com/dReam-dApps/dReams/gnomes"
	"github.com/sirupsen/logrus"
)

// Log output to stdout
var logger = structures.Logger.WithFields(logrus.Fields{})

// Gnomon instance for dSlate
var gnomon = gnomes.NewGnomes()

func searchByKey(scid string, key string, s bool) string {
	if gnomon.IsInitialized() {
		var sValue []string
		var uValue []uint64

		if s {
			sValue, uValue = gnomon.GetSCIDValuesByKey(scid, key)
		} else {
			if i, err := strconv.Atoi(key); err != nil {
				logger.Errorln("[dSlate]", err)
			} else {
				sValue, uValue = gnomon.GetSCIDValuesByKey(scid, uint64(i))
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
	if gnomon.IsInitialized() {
		var sValue []string
		var uValue []uint64
		if s {
			sValue, uValue = gnomon.GetSCIDKeysByValue(scid, value)
		} else {
			if i, err := strconv.Atoi(value); err != nil {
				logger.Errorln("[dSlate]", err)
			} else {
				sValue, uValue = gnomon.GetSCIDKeysByValue(scid, uint64(i))
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
