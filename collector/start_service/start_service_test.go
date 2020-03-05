package start_service

import (
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/SUSE/sap_host_exporter/internal/sapcontrol"
	"github.com/SUSE/sap_host_exporter/test/mock_sapcontrol"
)

func TestNewCollector(t *testing.T) {
	_, err := NewCollector("foobar")

	assert.Nil(t, err)
}

func TestProcessesMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)
	mockWebService.EXPECT().GetProcessList().Return(&sapcontrol.GetProcessListResponse{
		Process: &sapcontrol.ArrayOfOSProcess{
			Item: []*sapcontrol.OSProcess{
				{
					Name:        "enserver",
					Description: "foobar",
					Dispstatus:  sapcontrol.STATECOLORGREEN,
					Textstatus:  "Running",
					Starttime:   "",
					Elapsedtime: "",
					Pid:         30787,
				},
				{
					Name:        "msg_server",
					Description: "foobar2",
					Dispstatus:  sapcontrol.STATECOLORGREEN,
					Textstatus:  "Running",
					Starttime:   "",
					Elapsedtime: "",
					Pid:         30786,
				},
			},
		},
	}, nil)

	expectedMetrics := `
	# HELP sap_start_service_processes The processes started by the SAP Start Service
	# TYPE sap_start_service_processes gauge
	sap_start_service_processes{dispstatus="GREEN",name="enserver",pid="30787",textstatus="Running"} 1
	sap_start_service_processes{dispstatus="GREEN",name="msg_server",pid="30786",textstatus="Running"} 1
	`

	var err error
	collector, err := NewCollector("foobar")
	assert.NoError(t, err)
	collector.webService = mockWebService

	err = testutil.CollectAndCompare(collector, strings.NewReader(expectedMetrics), "sap_start_service_processes")
	assert.NoError(t, err)
}
