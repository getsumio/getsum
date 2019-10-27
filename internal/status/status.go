package status

type StatusType uint8

type Status struct {
	Type     StatusType
	Value    string
	Checksum string
}

const (
	PREPARED StatusType = iota
	ALLOCATED
	DOWNLOAD
	FETCHED
	STARTED
	RUNNING
	COMPLETED
	ERROR
	TIMEDOUT
	MISMATCH
	TERMINATED
)

var statusStr = []string{
	"PREPARED",
	"ALLOCATED",
	"DOWNLOAD",
	"FETCHED",
	"STARTED",
	"RUNNING",
	"COMPLETED",
	"ERROR",
	"TIMEDOUT",
	"MISMATCH",
	"TERMINATED",
}

func (s StatusType) Name() string {
	return statusStr[s]
}

func (s StatusType) Ordinal() uint8 {
	return uint8(s)
}
