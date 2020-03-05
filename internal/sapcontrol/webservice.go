package sapcontrol

import (
	"context"
	"encoding/xml"

	"github.com/hooklift/gowsdl/soap"
	"github.com/pkg/errors"
)

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

type STATECOLOR string

const (
	STATECOLORGRAY   STATECOLOR = "SAPControl-GRAY"
	STATECOLORGREEN  STATECOLOR = "SAPControl-GREEN"
	STATECOLORYELLOW STATECOLOR = "SAPControl-YELLOW"
	STATECOLORRED    STATECOLOR = "SAPControl-RED"
)

type WebService interface {
	/* Returns a list of all processes directly started by the webservice according to the SAP start profile. */
	GetProcessList() (*GetProcessListResponse, error)
}

type webService struct {
	client *soap.Client
}

// implements WebService.GetProcessList()
func (service *webService) GetProcessList() (*GetProcessListResponse, error) {
	request := &GetProcessList{}
	response := new(GetProcessListResponse)
	err := service.client.CallContext(context.Background(), "''", request, response)
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

// makes the STATECOLOR values more humnan-readable
func StateColorToString(statecolor STATECOLOR) (string, error) {
	switch statecolor {
	case STATECOLORGRAY:
		return "GRAY", nil
	case STATECOLORGREEN:
		return "GREEN", nil
	case STATECOLORYELLOW:
		return "YELLOW", nil
	case STATECOLORRED:
		return "RED", nil
	default:
		return "", errors.New("Invalid STATECOLOR value")
	}
}
