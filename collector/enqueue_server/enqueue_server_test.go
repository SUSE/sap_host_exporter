package enqueue_server

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

func TestEnqueueServerMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)
	mockWebService.EXPECT().EnqGetStatistic().Return(&sapcontrol.EnqGetStatisticResponse{
		OwnerNow:           1,
		OwnerHigh:          2,
		OwnerMax:           3,
		OwnerState:         sapcontrol.STATECOLOR_GREEN,
		ArgumentsNow:       4,
		ArgumentsHigh:      5,
		ArgumentsMax:       6,
		ArgumentsState:     sapcontrol.STATECOLOR_GRAY,
		LocksNow:           7,
		LocksHigh:          8,
		LocksMax:           9,
		LocksState:         sapcontrol.STATECOLOR_YELLOW,
		EnqueueRequests:    10,
		EnqueueRejects:     11,
		EnqueueErrors:      12,
		DequeueRequests:    13,
		DequeueErrors:      14,
		DequeueAllRequests: 15,
		CleanupRequests:    16,
		BackupRequests:     17,
		ReportingRequests:  18,
		CompressRequests:   19,
		VerifyRequests:     20,
		LockTime:           21,
		LockWaitTime:       22,
		ServerTime:         23,
		ReplicationState:   sapcontrol.STATECOLOR_RED,
	}, nil)
	mockWebService.EXPECT().GetCurrentInstance().Return(&sapcontrol.CurrentSapInstance{
		SID:      "HA1",
		Number:   0,
		Name:     "ASCS",
		Hostname: "sapha1as",
	}, nil)

	expectedMetrics := `
	# HELP sap_enqueue_server_arguments_high Peak number of lock arguments that have been stored simultaneously in the lock table
	# TYPE sap_enqueue_server_arguments_high counter
	sap_enqueue_server_arguments_high{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 5
	# HELP sap_enqueue_server_arguments_max Maximum number of lock arguments that can be stored in the lock table
	# TYPE sap_enqueue_server_arguments_max gauge
	sap_enqueue_server_arguments_max{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 6
	# HELP sap_enqueue_server_arguments_now Current number of lock arguments in the lock table
	# TYPE sap_enqueue_server_arguments_now gauge
	sap_enqueue_server_arguments_now{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 4
	# HELP sap_enqueue_server_arguments_state General state of lock arguments
	# TYPE sap_enqueue_server_arguments_state gauge
	sap_enqueue_server_arguments_state{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 1
	# HELP sap_enqueue_server_backup_requests Number of requests forwarded to the update process
	# TYPE sap_enqueue_server_backup_requests counter
	sap_enqueue_server_backup_requests{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 17
	# HELP sap_enqueue_server_cleanup_requests Requests to release of all the locks of an application server
	# TYPE sap_enqueue_server_cleanup_requests counter
	sap_enqueue_server_cleanup_requests{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 16
	# HELP sap_enqueue_server_compress_requests Internal use
	# TYPE sap_enqueue_server_compress_requests counter
	sap_enqueue_server_compress_requests{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 19
	# HELP sap_enqueue_server_dequeue_all_requests Requests to release of all the locks of an LUW
	# TYPE sap_enqueue_server_dequeue_all_requests counter
	sap_enqueue_server_dequeue_all_requests{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 15
	# HELP sap_enqueue_server_dequeue_errors Lock release errors
	# TYPE sap_enqueue_server_dequeue_errors counter
	sap_enqueue_server_dequeue_errors{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 14
	# HELP sap_enqueue_server_dequeue_requests Lock release requests
	# TYPE sap_enqueue_server_dequeue_requests counter
	sap_enqueue_server_dequeue_requests{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 13
	# HELP sap_enqueue_server_enqueue_errors Lock acquisition errors
	# TYPE sap_enqueue_server_enqueue_errors counter
	sap_enqueue_server_enqueue_errors{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 12
	# HELP sap_enqueue_server_enqueue_rejects Rejected lock requests
	# TYPE sap_enqueue_server_enqueue_rejects counter
	sap_enqueue_server_enqueue_rejects{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 11
	# HELP sap_enqueue_server_enqueue_requests Lock acquisition requests
	# TYPE sap_enqueue_server_enqueue_requests counter
	sap_enqueue_server_enqueue_requests{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 10
	# HELP sap_enqueue_server_lock_time Total time spent in lock operations
	# TYPE sap_enqueue_server_lock_time counter
	sap_enqueue_server_lock_time{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 21
	# HELP sap_enqueue_server_lock_wait_time Total waiting time of all work processes for accessing lock table
	# TYPE sap_enqueue_server_lock_wait_time counter
	sap_enqueue_server_lock_wait_time{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 22
	# HELP sap_enqueue_server_locks_high Peak number of elementary locks that have been stored simultaneously in the lock table
	# TYPE sap_enqueue_server_locks_high counter
	sap_enqueue_server_locks_high{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 8
	# HELP sap_enqueue_server_locks_max Maximum number of elementary locks that can be stored in the lock table
	# TYPE sap_enqueue_server_locks_max gauge
	sap_enqueue_server_locks_max{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 9
	# HELP sap_enqueue_server_locks_now Current number of elementary locks in the lock table
	# TYPE sap_enqueue_server_locks_now gauge
	sap_enqueue_server_locks_now{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 7
	# HELP sap_enqueue_server_locks_state General state of elementary locks
	# TYPE sap_enqueue_server_locks_state gauge
	sap_enqueue_server_locks_state{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 3
	# HELP sap_enqueue_server_owner_high Peak number of lock owners that have been stored simultaneously in the lock table
	# TYPE sap_enqueue_server_owner_high counter
	sap_enqueue_server_owner_high{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 2
	# HELP sap_enqueue_server_owner_max Maximum number of lock owner IDs that can be stored in the lock table
	# TYPE sap_enqueue_server_owner_max gauge
	sap_enqueue_server_owner_max{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 3
	# HELP sap_enqueue_server_owner_now Current number of lock owners in the lock table
	# TYPE sap_enqueue_server_owner_now gauge
	sap_enqueue_server_owner_now{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 1
	# HELP sap_enqueue_server_owner_state General state of lock owners
	# TYPE sap_enqueue_server_owner_state gauge
	sap_enqueue_server_owner_state{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 2
	# HELP sap_enqueue_server_replication_state General state of lock server replication
	# TYPE sap_enqueue_server_replication_state gauge
	sap_enqueue_server_replication_state{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 4
	# HELP sap_enqueue_server_reporting_requests Number of reading operations on the lock table
	# TYPE sap_enqueue_server_reporting_requests counter
	sap_enqueue_server_reporting_requests{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 18
	# HELP sap_enqueue_server_server_time Total time spent in lock operations by all processes in the enqueue server
	# TYPE sap_enqueue_server_server_time counter
	sap_enqueue_server_server_time{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 23
	# HELP sap_enqueue_server_verify_requests Internal use
	# TYPE sap_enqueue_server_verify_requests counter
	sap_enqueue_server_verify_requests{SID="HA1",instance_hostname="sapha1as",instance_name="ASCS",instance_number="0"} 20
	`

	var err error
	collector, err := NewCollector(mockWebService)
	assert.NoError(t, err)

	err = testutil.CollectAndCompare(collector, strings.NewReader(expectedMetrics))
	assert.NoError(t, err)
}
