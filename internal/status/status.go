package status

type StatusType uint8

type Status struct {
	Type     StatusType
	Value    string
	Checksum string
}

const (
	PREPARED   StatusType = iota
	ALLOCATED  StatusType = iota
	MISMATCH   StatusType = iota
	ERROR      StatusType = iota
	FETCHED    StatusType = iota
	DOWNLOAD   StatusType = iota
	STARTED    StatusType = iota
	TIMEDOUT   StatusType = iota
	COMPLETED  StatusType = iota
	RUNNING    StatusType = iota
	TERMINATED StatusType = iota
)

var statusStr = []string{
	"PREPARED",
	"ALLOCATED",
	"MISMATCH",
	"ERROR",
	"FETCHED",
	"DOWNLOAD",
	"STARTED",
	"TIMEDOUT",
	"COMPLETED",
	"RUNNING",
	"TERMINATED",
}

func (s StatusType) Name() string {
	return statusStr[s]
}

func (s StatusType) Ordinal() uint8 {
	return uint8(s)
}
