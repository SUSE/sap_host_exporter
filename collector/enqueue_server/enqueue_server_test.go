package enqueue_server

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)

	_, err := NewCollector(mockWebService)

	assert.Nil(t, err)
}

func TestProcessesMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)
	mockWebService.EXPECT().EnqGetStatistic().Return(&sapcontrol.EnqStatisticResponse{
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

	expectedMetrics := `
	# HELP sap_enqueue_server_arguments_high TBD
	# TYPE sap_enqueue_server_arguments_high gauge
	sap_enqueue_server_arguments_high 5
	# HELP sap_enqueue_server_arguments_max TBD
	# TYPE sap_enqueue_server_arguments_max gauge
	sap_enqueue_server_arguments_max 6
	# HELP sap_enqueue_server_arguments_now TBD
	# TYPE sap_enqueue_server_arguments_now gauge
	sap_enqueue_server_arguments_now 4
	# HELP sap_enqueue_server_arguments_state TBD
	# TYPE sap_enqueue_server_arguments_state gauge
	sap_enqueue_server_arguments_state 1
	# HELP sap_enqueue_server_backup_requests TBD
	# TYPE sap_enqueue_server_backup_requests gauge
	sap_enqueue_server_backup_requests 17
	# HELP sap_enqueue_server_cleanup_requests TBD
	# TYPE sap_enqueue_server_cleanup_requests gauge
	sap_enqueue_server_cleanup_requests 16
	# HELP sap_enqueue_server_compress_requests TBD
	# TYPE sap_enqueue_server_compress_requests gauge
	sap_enqueue_server_compress_requests 19
	# HELP sap_enqueue_server_dequeue_all_requests TBD
	# TYPE sap_enqueue_server_dequeue_all_requests gauge
	sap_enqueue_server_dequeue_all_requests 15
	# HELP sap_enqueue_server_dequeue_errors TBD
	# TYPE sap_enqueue_server_dequeue_errors gauge
	sap_enqueue_server_dequeue_errors 14
	# HELP sap_enqueue_server_dequeue_requests TBD
	# TYPE sap_enqueue_server_dequeue_requests gauge
	sap_enqueue_server_dequeue_requests 13
	# HELP sap_enqueue_server_enqueue_errors TBD
	# TYPE sap_enqueue_server_enqueue_errors gauge
	sap_enqueue_server_enqueue_errors 12
	# HELP sap_enqueue_server_enqueue_rejects TBD
	# TYPE sap_enqueue_server_enqueue_rejects gauge
	sap_enqueue_server_enqueue_rejects 11
	# HELP sap_enqueue_server_enqueue_requests TBD
	# TYPE sap_enqueue_server_enqueue_requests gauge
	sap_enqueue_server_enqueue_requests 10
	# HELP sap_enqueue_server_lock_time TBD
	# TYPE sap_enqueue_server_lock_time gauge
	sap_enqueue_server_lock_time 21
	# HELP sap_enqueue_server_lock_wait_time TBD
	# TYPE sap_enqueue_server_lock_wait_time gauge
	sap_enqueue_server_lock_wait_time 22
	# HELP sap_enqueue_server_locks_high TBD
	# TYPE sap_enqueue_server_locks_high gauge
	sap_enqueue_server_locks_high 8
	# HELP sap_enqueue_server_locks_max TBD
	# TYPE sap_enqueue_server_locks_max gauge
	sap_enqueue_server_locks_max 9
	# HELP sap_enqueue_server_locks_now TBD
	# TYPE sap_enqueue_server_locks_now gauge
	sap_enqueue_server_locks_now 7
	# HELP sap_enqueue_server_locks_state TBD
	# TYPE sap_enqueue_server_locks_state gauge
	sap_enqueue_server_locks_state 3
	# HELP sap_enqueue_server_owner_high TBD
	# TYPE sap_enqueue_server_owner_high gauge
	sap_enqueue_server_owner_high 2
	# HELP sap_enqueue_server_owner_max TBD
	# TYPE sap_enqueue_server_owner_max gauge
	sap_enqueue_server_owner_max 3
	# HELP sap_enqueue_server_owner_now TBD
	# TYPE sap_enqueue_server_owner_now gauge
	sap_enqueue_server_owner_now 1
	# HELP sap_enqueue_server_owner_state TBD
	# TYPE sap_enqueue_server_owner_state gauge
	sap_enqueue_server_owner_state 2
	# HELP sap_enqueue_server_replication_state TBD
	# TYPE sap_enqueue_server_replication_state gauge
	sap_enqueue_server_replication_state 4
	# HELP sap_enqueue_server_reporting_requests TBD
	# TYPE sap_enqueue_server_reporting_requests gauge
	sap_enqueue_server_reporting_requests 18
	# HELP sap_enqueue_server_server_time TBD
	# TYPE sap_enqueue_server_server_time gauge
	sap_enqueue_server_server_time 23
	# HELP sap_enqueue_server_verify_requests TBD
	# TYPE sap_enqueue_server_verify_requests gauge
	sap_enqueue_server_verify_requests 20
	`

	var err error
	collector, err := NewCollector(mockWebService)
	assert.NoError(t, err)

	err = testutil.CollectAndCompare(collector, strings.NewReader(expectedMetrics))
	assert.NoError(t, err)
}
