package sapcontrol

import (
	"encoding/xml"

	"github.com/hooklift/gowsdl/soap"
	"github.com/pkg/errors"
)

type STATECOLOR string
type STATECOLOR_CODE int

type GetProcessList struct {
	XMLName xml.Name `xml:"urn:SAPControl GetProcessList"`
}

type GetProcessListResponse struct {
	XMLName xml.Name          `xml:"urn:SAPControl GetProcessListResponse"`
	Process *ArrayOfOSProcess `xml:"process,omitempty" json:"process,omitempty"`
}

type OSProcess struct {
	Name        string     `xml:"name,omitempty" json:"name,omitempty"`
	Description string     `xml:"description,omitempty" json:"description,omitempty"`
	Dispstatus  STATECOLOR `xml:"dispstatus,omitempty" json:"dispstatus,omitempty"`
	Textstatus  string     `xml:"textstatus,omitempty" json:"textstatus,omitempty"`
	Starttime   string     `xml:"starttime,omitempty" json:"starttime,omitempty"`
	Elapsedtime string     `xml:"elapsedtime,omitempty" json:"elapsedtime,omitempty"`
	Pid         int32      `xml:"pid,omitempty" json:"pid,omitempty"`
}

type ArrayOfOSProcess struct {
	Item []*OSProcess `xml:"item,omitempty" json:"item,omitempty"`
}

type EnqGetStatistic struct {
	XMLName xml.Name `xml:"urn:SAPControl EnqGetStatistic"`
}

type EnqStatisticResponse struct {
	XMLName            xml.Name   `xml:"urn:SAPControl EnqStatistic"`
	OwnerNow           int32      `xml:"owner-now,omitempty" json:"owner-now,omitempty"`
	OwnerHigh          int32      `xml:"owner-high,omitempty" json:"owner-high,omitempty"`
	OwnerMax           int32      `xml:"owner-max,omitempty" json:"owner-max,omitempty"`
	OwnerState         STATECOLOR `xml:"owner-state,omitempty" json:"owner-state,omitempty"`
	ArgumentsNow       int32      `xml:"arguments-now,omitempty" json:"arguments-now,omitempty"`
	ArgumentsHigh      int32      `xml:"arguments-high,omitempty" json:"arguments-high,omitempty"`
	ArgumentsMax       int32      `xml:"arguments-max,omitempty" json:"arguments-max,omitempty"`
	ArgumentsState     STATECOLOR `xml:"arguments-state,omitempty" json:"arguments-state,omitempty"`
	LocksNow           int32      `xml:"locks-now,omitempty" json:"locks-now,omitempty"`
	LocksHigh          int32      `xml:"locks-high,omitempty" json:"locks-high,omitempty"`
	LocksMax           int32      `xml:"locks-max,omitempty" json:"locks-max,omitempty"`
	LocksState         STATECOLOR `xml:"locks-state,omitempty" json:"locks-state,omitempty"`
	EnqueueRequests    int64      `xml:"enqueue-requests,omitempty" json:"enqueue-requests,omitempty"`
	EnqueueRejects     int64      `xml:"enqueue-rejects,omitempty" json:"enqueue-rejects,omitempty"`
	EnqueueErrors      int64      `xml:"enqueue-errors,omitempty" json:"enqueue-errors,omitempty"`
	DequeueRequests    int64      `xml:"dequeue-requests,omitempty" json:"dequeue-requests,omitempty"`
	DequeueErrors      int64      `xml:"dequeue-errors,omitempty" json:"dequeue-errors,omitempty"`
	DequeueAllRequests int64      `xml:"dequeue-all-requests,omitempty" json:"dequeue-all-requests,omitempty"`
	CleanupRequests    int64      `xml:"cleanup-requests,omitempty" json:"cleanup-requests,omitempty"`
	BackupRequests     int64      `xml:"backup-requests,omitempty" json:"backup-requests,omitempty"`
	ReportingRequests  int64      `xml:"reporting-requests,omitempty" json:"reporting-requests,omitempty"`
	CompressRequests   int64      `xml:"compress-requests,omitempty" json:"compress-requests,omitempty"`
	VerifyRequests     int64      `xml:"verify-requests,omitempty" json:"verify-requests,omitempty"`
	LockTime           float64    `xml:"lock-time,omitempty" json:"lock-time,omitempty"`
	LockWaitTime       float64    `xml:"lock-wait-time,omitempty" json:"lock-wait-time,omitempty"`
	ServerTime         float64    `xml:"server-time,omitempty" json:"server-time,omitempty"`
	ReplicationState   STATECOLOR `xml:"replication-state,omitempty" json:"replication-state,omitempty"`
}

const (
	STATECOLOR_GRAY        STATECOLOR      = "SAPControl-GRAY"
	STATECOLOR_GREEN       STATECOLOR      = "SAPControl-GREEN"
	STATECOLOR_YELLOW      STATECOLOR      = "SAPControl-YELLOW"
	STATECOLOR_RED         STATECOLOR      = "SAPControl-RED"
	STATECOLOR_CODE_GRAY   STATECOLOR_CODE = 1
	STATECOLOR_CODE_GREEN  STATECOLOR_CODE = 2
	STATECOLOR_CODE_YELLOW STATECOLOR_CODE = 3
	STATECOLOR_CODE_RED    STATECOLOR_CODE = 4
)

type WebService interface {
	/* Returns a list of all processes directly started by the webservice according to the SAP start profile. */
	GetProcessList() (*GetProcessListResponse, error)

	/* Returns enque statistic. */
	EnqGetStatistic() (*EnqStatisticResponse, error)
}

type webService struct {
	client *soap.Client
}

// implements WebService.GetProcessList()
func (service *webService) GetProcessList() (*GetProcessListResponse, error) {
	request := &GetProcessList{}
	response := &GetProcessListResponse{}
	err := service.client.Call("''", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// implements WebService.EnqGetStatistic()
func (service *webService) EnqGetStatistic() (*EnqStatisticResponse, error) {
	request := &EnqGetStatistic{}
	response := &EnqStatisticResponse{}
	err := service.client.Call("''", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// constructor of a WebService interface
func NewWebService(client *soap.Client) WebService {
	return &webService{
		client: client,
	}
}

// makes the STATECOLOR values more human-readable
func StateColorToString(statecolor STATECOLOR) (string, error) {
	switch statecolor {
	case STATECOLOR_GRAY:
		return "GRAY", nil
	case STATECOLOR_GREEN:
		return "GREEN", nil
	case STATECOLOR_YELLOW:
		return "YELLOW", nil
	case STATECOLOR_RED:
		return "RED", nil
	default:
		return "", errors.New("Invalid STATECOLOR value")
	}
}

// makes the STATECOLOR values more metric friendly
func StateColorToFloat(statecolor STATECOLOR) (float64, error) {
	switch statecolor {
	case STATECOLOR_GRAY:
		return float64(STATECOLOR_CODE_GRAY), nil
	case STATECOLOR_GREEN:
		return float64(STATECOLOR_CODE_GREEN), nil
	case STATECOLOR_YELLOW:
		return float64(STATECOLOR_CODE_YELLOW), nil
	case STATECOLOR_RED:
		return float64(STATECOLOR_CODE_RED), nil
	default:
		return 0, errors.New("Invalid STATECOLOR value")
	}
}
