package internal

type Version struct {
	major        uint32
	minor        uint32
	build_number uint32
}

var AppVersion = Version{major: 0, minor: 0, build_number: 1}
