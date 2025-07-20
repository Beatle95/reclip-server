package internal

import (
	"strconv"
	"strings"
)

type Version struct {
	major uint32
	minor uint32
	patch uint32
}

var applicationVersionString string

func GetApplicationVersion() Version {
	splitted := strings.Split(applicationVersionString[1:], ".")
	if len(splitted) != 3 {
		panic("Version string is not set for application")
	}

	major, err := strconv.Atoi(splitted[0])
	if err != nil {
		panic("Incorrect major part of version string")
	}
	minor, err := strconv.Atoi(splitted[1])
	if err != nil {
		panic("Incorrect minor part of version string")
	}
	patch, err := strconv.Atoi(splitted[2])
	if err != nil {
		panic("Incorrect patch part of version string")
	}
	return Version{
		major: uint32(major),
		minor: uint32(minor),
		patch: uint32(patch),
	}
}

func GetApplicationVersionString() string {
	return applicationVersionString[1:]
}
