package collector_register

import (
	"github.com/SUSE/sap_host_exporter/collector/dispatcher"
	"github.com/SUSE/sap_host_exporter/internal/sapcontrol"
	"github.com/SUSE/sap_host_exporter/test/mock_sapcontrol"
	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestActivationDispatcher(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)

	//RegisterOptionalCollectors(mockWebService)

	mockWebService.EXPECT().GetProcessList().Return(&sapcontrol.GetProcessListResponse{
		Processes: []*sapcontrol.OSProcess{
			{
				Name:        "enserver",
				Description: "foobar",
				Dispstatus:  sapcontrol.STATECOLOR_GREEN,
				Textstatus:  "Running",
				Starttime:   "",
				Elapsedtime: "",
				Pid:         30787,
			},
			{
				// activate the collector
				Name:        "disp+work",
				Description: "foobar2",
				Dispstatus:  sapcontrol.STATECOLOR_YELLOW,
				Textstatus:  "Stopping",
				Starttime:   "",
				Elapsedtime: "",
				Pid:         30786,
			},
		},
	}, nil)

	mockWebService.EXPECT().GetQueueStatistic().Return(&sapcontrol.GetQueueStatisticResponse{
		Queues: []*sapcontrol.TaskHandlerQueue{
			{Type: "ABAP/NOWP", High: 3, Max: 14000, Writes: 249133, Reads: 249133},
			{Type: "ABAP/DIA", High: 5, Max: 14000, Writes: 447173, Reads: 447173},
			{Type: "ABAP/UPD", High: 2, Max: 14000, Writes: 3491, Reads: 3491},
			{Type: "ABAP/ENQ", Max: 14000},
			{Type: "ABAP/BTC", High: 2, Max: 14000, Writes: 10464, Reads: 10464},
			{Type: "ABAP/SPO", High: 1, Max: 14000, Writes: 38366, Reads: 38366},
			{Type: "ABAP/UP2", High: 1, Max: 14000, Writes: 3488, Reads: 3488},
			{Type: "ICM/Intern", High: 1, Max: 6000, Writes: 34877, Reads: 34877},
		},
	}, nil)

	expectedMetrics := `
	# HELP sap_dispatcher_queue_high Work process peak queue length
	# TYPE sap_dispatcher_queue_high counter
	sap_dispatcher_queue_high{type="ABAP/BTC"} 2
	sap_dispatcher_queue_high{type="ABAP/DIA"} 5
	sap_dispatcher_queue_high{type="ABAP/ENQ"} 0
	sap_dispatcher_queue_high{type="ABAP/NOWP"} 3
	sap_dispatcher_queue_high{type="ABAP/SPO"} 1
	sap_dispatcher_queue_high{type="ABAP/UP2"} 1
	sap_dispatcher_queue_high{type="ABAP/UPD"} 2
	sap_dispatcher_queue_high{type="ICM/Intern"} 1
	# HELP sap_dispatcher_queue_max Work process maximum queue length
	# TYPE sap_dispatcher_queue_max gauge
	sap_dispatcher_queue_max{type="ABAP/BTC"} 14000
	sap_dispatcher_queue_max{type="ABAP/DIA"} 14000
	sap_dispatcher_queue_max{type="ABAP/ENQ"} 14000
	sap_dispatcher_queue_max{type="ABAP/NOWP"} 14000
	sap_dispatcher_queue_max{type="ABAP/SPO"} 14000
	sap_dispatcher_queue_max{type="ABAP/UP2"} 14000
	sap_dispatcher_queue_max{type="ABAP/UPD"} 14000
	sap_dispatcher_queue_max{type="ICM/Intern"} 6000
	# HELP sap_dispatcher_queue_now Work process current queue length
	# TYPE sap_dispatcher_queue_now gauge
	sap_dispatcher_queue_now{type="ABAP/BTC"} 0
	sap_dispatcher_queue_now{type="ABAP/DIA"} 0
	sap_dispatcher_queue_now{type="ABAP/ENQ"} 0
	sap_dispatcher_queue_now{type="ABAP/NOWP"} 0
	sap_dispatcher_queue_now{type="ABAP/SPO"} 0
	sap_dispatcher_queue_now{type="ABAP/UP2"} 0
	sap_dispatcher_queue_now{type="ABAP/UPD"} 0
	sap_dispatcher_queue_now{type="ICM/Intern"} 0
	# HELP sap_dispatcher_queue_reads Work process queue reads
	# TYPE sap_dispatcher_queue_reads counter
	sap_dispatcher_queue_reads{type="ABAP/BTC"} 10464
	sap_dispatcher_queue_reads{type="ABAP/DIA"} 447173
	sap_dispatcher_queue_reads{type="ABAP/ENQ"} 0
	sap_dispatcher_queue_reads{type="ABAP/NOWP"} 249133
	sap_dispatcher_queue_reads{type="ABAP/SPO"} 38366
	sap_dispatcher_queue_reads{type="ABAP/UP2"} 3488
	sap_dispatcher_queue_reads{type="ABAP/UPD"} 3491
	sap_dispatcher_queue_reads{type="ICM/Intern"} 34877
	# HELP sap_dispatcher_queue_writes Work process queue writes
	# TYPE sap_dispatcher_queue_writes counter
	sap_dispatcher_queue_writes{type="ABAP/BTC"} 10464
	sap_dispatcher_queue_writes{type="ABAP/DIA"} 447173
	sap_dispatcher_queue_writes{type="ABAP/ENQ"} 0
	sap_dispatcher_queue_writes{type="ABAP/NOWP"} 249133
	sap_dispatcher_queue_writes{type="ABAP/SPO"} 38366
	sap_dispatcher_queue_writes{type="ABAP/UP2"} 3488
	sap_dispatcher_queue_writes{type="ABAP/UPD"} 3491
	sap_dispatcher_queue_writes{type="ICM/Intern"} 34877
`

	var err error
	collector, err := dispatcher.NewCollector(mockWebService)
	assert.NoError(t, err)

	err = testutil.CollectAndCompare(collector, strings.NewReader(expectedMetrics))
	assert.NoError(t, err)

}