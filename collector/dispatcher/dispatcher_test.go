package dispatcher

import (
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"

	"github.com/SUSE/sap_host_exporter/lib/sapcontrol"
	"github.com/SUSE/sap_host_exporter/test/mock_sapcontrol"
)

func TestNewCollector(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)

	_, err := NewCollector(mockWebService)

	assert.Nil(t, err)
}

func TestWorkProcessQueueStatsMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)
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
	mockWebService.EXPECT().GetCurrentInstance().Return(&sapcontrol.CurrentSapInstance{
		SID:      "HA1",
		Number:   0,
		Name:     "ASCS",
		Hostname: "sapha1as",
	}, nil)

	expectedMetrics := `
	# HELP sap_dispatcher_queue_high Work process peak queue length
	# TYPE sap_dispatcher_queue_high counter
	sap_dispatcher_queue_high{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/BTC"} 2
	sap_dispatcher_queue_high{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/DIA"} 5
	sap_dispatcher_queue_high{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/ENQ"} 0
	sap_dispatcher_queue_high{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/NOWP"} 3
	sap_dispatcher_queue_high{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/SPO"} 1
	sap_dispatcher_queue_high{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/UP2"} 1
	sap_dispatcher_queue_high{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/UPD"} 2
	sap_dispatcher_queue_high{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ICM/Intern"} 1
	# HELP sap_dispatcher_queue_max Work process maximum queue length
	# TYPE sap_dispatcher_queue_max gauge
	sap_dispatcher_queue_max{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/BTC"} 14000
	sap_dispatcher_queue_max{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/DIA"} 14000
	sap_dispatcher_queue_max{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/ENQ"} 14000
	sap_dispatcher_queue_max{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/NOWP"} 14000
	sap_dispatcher_queue_max{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/SPO"} 14000
	sap_dispatcher_queue_max{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/UP2"} 14000
	sap_dispatcher_queue_max{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/UPD"} 14000
	sap_dispatcher_queue_max{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ICM/Intern"} 6000
	# HELP sap_dispatcher_queue_now Work process current queue length
	# TYPE sap_dispatcher_queue_now gauge
	sap_dispatcher_queue_now{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/BTC"} 0
	sap_dispatcher_queue_now{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/DIA"} 0
	sap_dispatcher_queue_now{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/ENQ"} 0
	sap_dispatcher_queue_now{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/NOWP"} 0
	sap_dispatcher_queue_now{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/SPO"} 0
	sap_dispatcher_queue_now{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/UP2"} 0
	sap_dispatcher_queue_now{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/UPD"} 0
	sap_dispatcher_queue_now{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ICM/Intern"} 0
	# HELP sap_dispatcher_queue_reads Work process queue reads
	# TYPE sap_dispatcher_queue_reads counter
	sap_dispatcher_queue_reads{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/BTC"} 10464
	sap_dispatcher_queue_reads{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/DIA"} 447173
	sap_dispatcher_queue_reads{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/ENQ"} 0
	sap_dispatcher_queue_reads{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/NOWP"} 249133
	sap_dispatcher_queue_reads{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/SPO"} 38366
	sap_dispatcher_queue_reads{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/UP2"} 3488
	sap_dispatcher_queue_reads{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/UPD"} 3491
	sap_dispatcher_queue_reads{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ICM/Intern"} 34877
	# HELP sap_dispatcher_queue_writes Work process queue writes
	# TYPE sap_dispatcher_queue_writes counter
	sap_dispatcher_queue_writes{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/BTC"} 10464
	sap_dispatcher_queue_writes{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/DIA"} 447173
	sap_dispatcher_queue_writes{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/ENQ"} 0
	sap_dispatcher_queue_writes{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/NOWP"} 249133
	sap_dispatcher_queue_writes{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/SPO"} 38366
	sap_dispatcher_queue_writes{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/UP2"} 3488
	sap_dispatcher_queue_writes{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ABAP/UPD"} 3491
	sap_dispatcher_queue_writes{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0",type="ICM/Intern"} 34877
`

	var err error
	collector, err := NewCollector(mockWebService)
	assert.NoError(t, err)

	err = testutil.CollectAndCompare(collector, strings.NewReader(expectedMetrics))
	assert.NoError(t, err)
}

func TestWorkProcessQueueStatsMetricWithEmptyData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)
	mockWebService.EXPECT().GetQueueStatistic().Return(&sapcontrol.GetQueueStatisticResponse{}, nil)
	mockWebService.EXPECT().GetCurrentInstance().Return(&sapcontrol.CurrentSapInstance{}, nil)

	var err error
	collector, err := NewCollector(mockWebService)
	assert.NoError(t, err)

	err = testutil.CollectAndCompare(collector, strings.NewReader(""))
	assert.NoError(t, err)
}
