package status

type StatusType uint8

type Status struct {
	Type     StatusType `json:"type"`
	Value    string     `json:"value"`
	Checksum string     `json:"checksum"`
}

const (
	PREPARED StatusType = iota
	ALLOCATED
	DOWNLOAD
	FETCHED
	STARTED
	RESUMING
	RUNNING
	COMPLETED
	SUSPENDED
	ERROR
	TIMEDOUT
	MISMATCH
	VALIDATED
	TERMINATED
)

var statusStr = []string{
	"PREPARED",
	"ALLOCATED",
	"DOWNLOAD",
	"FETCHED",
	"STARTED",
	"RESUMING",
	"RUNNING",
	"COMPLETED",
	"SUSPENDED",
	"ERROR",
	"TIMEDOUT",
	"MISMATCH",
	"VALIDATED",
	"TERMINATED",
}

func (s StatusType) Name() string {
	return statusStr[s]
}

func (s StatusType) Ordinal() uint8 {
	return uint8(s)
}
