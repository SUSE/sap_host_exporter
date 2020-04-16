package sapcontrol

import (
	"encoding/xml"

	"github.com/hooklift/gowsdl/soap"
	"github.com/pkg/errors"
)

//go:generate mockgen -destination ../../test/mock_sapcontrol/webservice.go github.com/SUSE/sap_host_exporter/internal/sapcontrol WebService

// the main interface exposed by this package
type WebService interface {
	/* Returns a list of all processes directly started by the webservice according to the SAP start profile. */
	GetProcessList() (*GetProcessListResponse, error)

	/* Returns enque statistic. */
	EnqGetStatistic() (*EnqGetStatisticResponse, error)

	/* Returns a list of queue information of work processes and icm (similar to dpmon). */
	GetQueueStatistic() (*GetQueueStatisticResponse, error)

	/* Checks high availability configuration and status of the system. */
	HACheckConfig() (*HACheckConfigResponse, error)

	/* Checks HA failover third party configuration and status of an instnace. */
	HACheckFailoverConfig() (*HACheckFailoverConfigResponse, error)

	/* Returns HA failover third party information. */
	HAGetFailoverConfig() (*HAGetFailoverConfigResponse, error)

	/* Returns a list of SAP instances of the SAP system. */
	GetSystemInstanceList() (*GetSystemInstanceListResponse, error)

	/* Returns a list of available instance features and information how to get it. */
	GetInstanceProperties() (*GetInstancePropertiesResponse, error)
}

type HACheckCategory string
type HAVerificationState string
type HAVerificationStateCode int
type STATECOLOR string
type STATECOLOR_CODE int

const (
	HA_CHECK_CATEGORY_SAP_CONFIGURATION HACheckCategory = "SAPControl-SAP-CONFIGURATION"
	HA_CHECK_CATEGORY_SAP_STATE         HACheckCategory = "SAPControl-SAP-STATE"
	HA_CHECK_CATEGORY_HA_CONFIGURATION  HACheckCategory = "SAPControl-HA-CONFIGURATION"
	HA_CHECK_CATEGORY_HA_STATE          HACheckCategory = "SAPControl-HA-STATE"
)

const (
	HA_VERIFICATION_STATE_SUCCESS      HAVerificationState     = "SAPControl-HA-SUCCESS"
	HA_VERIFICATION_STATE_WARNING      HAVerificationState     = "SAPControl-HA-WARNING"
	HA_VERIFICATION_STATE_ERROR        HAVerificationState     = "SAPControl-HA-ERROR"
	HA_VERIFICATION_STATE_CODE_SUCCESS HAVerificationStateCode = 0
	HA_VERIFICATION_STATE_CODE_WARNING HAVerificationStateCode = 1
	HA_VERIFICATION_STATE_CODE_ERROR   HAVerificationStateCode = 2
)

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

type EnqGetStatistic struct {
	XMLName xml.Name `xml:"urn:SAPControl EnqGetStatistic"`
}

type EnqGetStatisticResponse struct {
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

type GetInstanceProperties struct {
	XMLName xml.Name `xml:"urn:SAPControl GetInstanceProperties"`
}

type GetInstancePropertiesResponse struct {
	XMLName    xml.Name            `xml:"urn:SAPControl GetInstancePropertiesResponse"`
	Properties []*InstanceProperty `xml:"properties>item,omitempty" json:"properties>item,omitempty"`
}

type GetProcessList struct {
	XMLName xml.Name `xml:"urn:SAPControl GetProcessList"`
}

type GetProcessListResponse struct {
	XMLName   xml.Name     `xml:"urn:SAPControl GetProcessListResponse"`
	Processes []*OSProcess `xml:"process>item,omitempty" json:"process>item,omitempty"`
}

type GetQueueStatistic struct {
	XMLName xml.Name `xml:"urn:SAPControl GetQueueStatistic"`
}

type GetQueueStatisticResponse struct {
	XMLName xml.Name            `xml:"urn:SAPControl GetQueueStatisticResponse"`
	Queues  []*TaskHandlerQueue `xml:"queue>item,omitempty" json:"queue>item,omitempty"`
}

type GetSystemInstanceList struct {
	XMLName xml.Name `xml:"urn:SAPControl GetSystemInstanceList"`
	Timeout int32    `xml:"timeout,omitempty" json:"timeout,omitempty"`
}

type GetSystemInstanceListResponse struct {
	XMLName   xml.Name       `xml:"urn:SAPControl GetSystemInstanceListResponse"`
	Instances []*SAPInstance `xml:"instance>item,omitempty" json:"instance>item,omitempty"`
}

type HACheck struct {
	State       HAVerificationState `xml:"state,omitempty" json:"state,omitempty"`
	Category    HACheckCategory     `xml:"category,omitempty" json:"category,omitempty"`
	Description string              `xml:"description,omitempty" json:"description,omitempty"`
	Comment     string              `xml:"comment,omitempty" json:"comment,omitempty"`
}

type HACheckConfig struct {
	XMLName xml.Name `xml:"urn:SAPControl HACheckConfig"`
}

type HACheckConfigResponse struct {
	XMLName xml.Name   `xml:"urn:SAPControl HACheckConfigResponse"`
	Checks  []*HACheck `xml:"check>item,omitempty" json:"check>item,omitempty"`
}

type HACheckFailoverConfig struct {
	XMLName xml.Name `xml:"urn:SAPControl HACheckFailoverConfig"`
}

type HACheckFailoverConfigResponse struct {
	XMLName xml.Name   `xml:"urn:SAPControl HACheckFailoverConfigResponse"`
	Checks  []*HACheck `xml:"check>item,omitempty" json:"check>item,omitempty"`
}

type HAGetFailoverConfigResponse struct {
	XMLName               xml.Name `xml:"urn:SAPControl HAGetFailoverConfigResponse"`
	HAActive              bool     `xml:"HAActive,omitempty" json:"HAActive,omitempty"`
	HAProductVersion      string   `xml:"HAProductVersion,omitempty" json:"HAProductVersion,omitempty"`
	HASAPInterfaceVersion string   `xml:"HASAPInterfaceVersion,omitempty" json:"HASAPInterfaceVersion,omitempty"`
	HADocumentation       string   `xml:"HADocumentation,omitempty" json:"HADocumentation,omitempty"`
	HAActiveNode          string   `xml:"HAActiveNode,omitempty" json:"HAActiveNode,omitempty"`
	HANodes               []string `xml:"HANodes>item,omitempty" json:"HANodes,omitempty"`
}

type HAGetFailoverConfig struct {
	XMLName xml.Name `xml:"urn:SAPControl HAGetFailoverConfig"`
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

type InstanceProperty struct {
	Property     string `xml:"property,omitempty" json:"property,omitempty"`
	Propertytype string `xml:"propertytype,omitempty" json:"propertytype,omitempty"`
	Value        string `xml:"value,omitempty" json:"value,omitempty"`
}

type SAPInstance struct {
	Hostname      string     `xml:"hostname,omitempty" json:"hostname,omitempty"`
	InstanceNr    int32      `xml:"instanceNr,omitempty" json:"instanceNr,omitempty"`
	HttpPort      int32      `xml:"httpPort,omitempty" json:"httpPort,omitempty"`
	HttpsPort     int32      `xml:"httpsPort,omitempty" json:"httpsPort,omitempty"`
	StartPriority string     `xml:"startPriority,omitempty" json:"startPriority,omitempty"`
	Features      string     `xml:"features,omitempty" json:"features,omitempty"`
	Dispstatus    STATECOLOR `xml:"dispstatus,omitempty" json:"dispstatus,omitempty"`
}

type TaskHandlerQueue struct {
	Type   string `xml:"Typ,omitempty" json:"Typ,omitempty"`
	Now    int32  `xml:"Now,omitempty" json:"Now,omitempty"`
	High   int32  `xml:"High,omitempty" json:"High,omitempty"`
	Max    int32  `xml:"Max,omitempty" json:"Max,omitempty"`
	Writes int32  `xml:"Writes,omitempty" json:"Writes,omitempty"`
	Reads  int32  `xml:"Reads,omitempty" json:"Reads,omitempty"`
}

type webService struct {
	client *soap.Client
}

// constructor of a WebService interface
func NewWebService(client *soap.Client) WebService {
	return &webService{
		client: client,
	}
}

// implements WebService.GetInstanceProperties()
func (service *webService) GetInstanceProperties() (*GetInstancePropertiesResponse, error) {
	request := &GetInstanceProperties{}
	response := &GetInstancePropertiesResponse{}
	err := service.client.Call("''", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
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

// implements WebService.GetSystemInstanceList()
func (service *webService) GetSystemInstanceList() (*GetSystemInstanceListResponse, error) {
	request := &GetSystemInstanceList{}
	response := &GetSystemInstanceListResponse{}
	err := service.client.Call("''", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// implements WebService.EnqGetStatistic()
func (service *webService) EnqGetStatistic() (*EnqGetStatisticResponse, error) {
	request := &EnqGetStatistic{}
	response := &EnqGetStatisticResponse{}
	err := service.client.Call("''", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// implements WebService.GetQueueStatistic()
func (service *webService) GetQueueStatistic() (*GetQueueStatisticResponse, error) {
	request := &GetQueueStatistic{}
	response := &GetQueueStatisticResponse{}
	err := service.client.Call("''", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// implements WebService.HACheckConfig()
func (service *webService) HACheckConfig() (*HACheckConfigResponse, error) {
	request := &HACheckConfig{}
	response := &HACheckConfigResponse{}
	err := service.client.Call("''", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// implements WebService.HACheckFailoverConfig()
func (service *webService) HACheckFailoverConfig() (*HACheckFailoverConfigResponse, error) {
	request := &HACheckFailoverConfig{}
	response := &HACheckFailoverConfigResponse{}
	err := service.client.Call("''", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// implements WebService.HAGetFailoverConfig()
func (service *webService) HAGetFailoverConfig() (*HAGetFailoverConfigResponse, error) {
	request := &HAGetFailoverConfig{}
	response := &HAGetFailoverConfigResponse{}
	err := service.client.Call("''", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
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
		return -1, errors.New("Invalid STATECOLOR value")
	}
}

// makes HACheckCategory values more human-readable
func HaCheckCategoryToString(category HACheckCategory) (string, error) {
	switch category {
	case HA_CHECK_CATEGORY_HA_CONFIGURATION:
		return "HA-CONFIGURATION", nil
	case HA_CHECK_CATEGORY_HA_STATE:
		return "HA-STATE", nil
	case HA_CHECK_CATEGORY_SAP_CONFIGURATION:
		return "SAP-CONFIGURATION", nil
	case HA_CHECK_CATEGORY_SAP_STATE:
		return "SAP-STATE", nil
	default:
		return "", errors.New("Invalid HACheckCategory value")
	}
}

// makes HAVerificationState values more metric friendly
func HaVerificationStateToFloat(state HAVerificationState) (float64, error) {
	switch state {
	case HA_VERIFICATION_STATE_SUCCESS:
		return float64(HA_VERIFICATION_STATE_CODE_SUCCESS), nil
	case HA_VERIFICATION_STATE_WARNING:
		return float64(HA_VERIFICATION_STATE_CODE_WARNING), nil
	case HA_VERIFICATION_STATE_ERROR:
		return float64(HA_VERIFICATION_STATE_CODE_ERROR), nil
	default:
		return -1, errors.New("Invalid HAVerificationState value")
	}
}
