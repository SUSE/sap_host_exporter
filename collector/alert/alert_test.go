package alert

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

func TestHACheckMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)

	mockHACheckConfigResponse := &sapcontrol.HACheckConfigResponse{
		Check: &sapcontrol.ArrayOfHACheck{
			Item: []*sapcontrol.HACheck{
				{State: sapcontrol.HA_VERIFICATION_STATE_ERROR, Category: sapcontrol.HA_CHECK_CATEGORY_HA_STATE, Description: "foo", Comment: "bar"},
				{State: sapcontrol.HA_VERIFICATION_STATE_WARNING, Category: sapcontrol.HA_CHECK_CATEGORY_SAP_STATE, Description: "foo2", Comment: "bar2"},
				{State: sapcontrol.HA_VERIFICATION_STATE_SUCCESS, Category: sapcontrol.HA_CHECK_CATEGORY_SAP_CONFIGURATION, Description: "foo3", Comment: "bar3"},
			},
		},
	}
	mockHACheckFailoverConfigResponse := &sapcontrol.HACheckFailoverConfigResponse{
		Check: &sapcontrol.ArrayOfHACheck{
			Item: []*sapcontrol.HACheck{
				{State: sapcontrol.HA_VERIFICATION_STATE_SUCCESS, Category: sapcontrol.HA_CHECK_CATEGORY_SAP_CONFIGURATION, Description: "foo4", Comment: "bar4"},
			},
		},
	}
	mockWebService.EXPECT().HACheckConfig().Return(mockHACheckConfigResponse, nil)
	mockWebService.EXPECT().HACheckFailoverConfig().Return(mockHACheckFailoverConfigResponse, nil)
	mockWebService.EXPECT().HAGetFailoverConfig().Return(&sapcontrol.HAGetFailoverConfigResponse{}, nil)

	var err error
	collector, err := NewCollector(mockWebService)
	assert.NoError(t, err)

	expectedMetrics := `
	# HELP sap_alert_ha_check High Availability system configuration and status checks
	# TYPE sap_alert_ha_check gauge
	sap_alert_ha_check{category="HA-STATE",comment="bar",description="foo",state="ERROR"} 2
	sap_alert_ha_check{category="SAP-STATE",comment="bar2",description="foo2",state="WARNING"} 1
	sap_alert_ha_check{category="SAP-CONFIGURATION",comment="bar3",description="foo3",state="SUCCESS"} 0
	sap_alert_ha_check{category="SAP-CONFIGURATION",comment="bar4",description="foo4",state="SUCCESS"} 0
`

	err = testutil.CollectAndCompare(collector, strings.NewReader(expectedMetrics), "sap_alert_ha_check")
	assert.NoError(t, err)
}

func TestHACheckMetricsWithEmptyData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)

	mockWebService.EXPECT().HACheckConfig().Return(&sapcontrol.HACheckConfigResponse{}, nil)
	mockWebService.EXPECT().HACheckFailoverConfig().Return(&sapcontrol.HACheckFailoverConfigResponse{}, nil)
	mockWebService.EXPECT().HAGetFailoverConfig().Return(&sapcontrol.HAGetFailoverConfigResponse{}, nil)

	var err error
	collector, err := NewCollector(mockWebService)
	assert.NoError(t, err)

	err = testutil.CollectAndCompare(collector, strings.NewReader(""), "sap_alert_ha_check")
	assert.NoError(t, err)
}

func TestHAFailoverActiveMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)

	mockWebService.EXPECT().HACheckConfig().Return(&sapcontrol.HACheckConfigResponse{}, nil)
	mockWebService.EXPECT().HACheckFailoverConfig().Return(&sapcontrol.HACheckFailoverConfigResponse{}, nil)
	mockWebService.EXPECT().HAGetFailoverConfig().Return(&sapcontrol.HAGetFailoverConfigResponse{
		HAActive: true,
	}, nil)

	var err error
	collector, err := NewCollector(mockWebService)
	assert.NoError(t, err)

	expectedMetrics := `
	# HELP sap_alert_ha_failover_active Whether or not High Availability Failover is active
	# TYPE sap_alert_ha_failover_active gauge
	sap_alert_ha_failover_active 1
`
	err = testutil.CollectAndCompare(collector, strings.NewReader(expectedMetrics), "sap_alert_ha_failover_active")
	assert.NoError(t, err)
}

func TestHAFailoverActiveMetricWithFalseValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)

	mockWebService.EXPECT().HACheckConfig().Return(&sapcontrol.HACheckConfigResponse{}, nil)
	mockWebService.EXPECT().HACheckFailoverConfig().Return(&sapcontrol.HACheckFailoverConfigResponse{}, nil)
	mockWebService.EXPECT().HAGetFailoverConfig().Return(&sapcontrol.HAGetFailoverConfigResponse{
		HAActive: false,
	}, nil)

	var err error
	collector, err := NewCollector(mockWebService)
	assert.NoError(t, err)

	expectedMetrics := `
	# HELP sap_alert_ha_failover_active Whether or not High Availability Failover is active
	# TYPE sap_alert_ha_failover_active gauge
	sap_alert_ha_failover_active 0
`
	err = testutil.CollectAndCompare(collector, strings.NewReader(expectedMetrics), "sap_alert_ha_failover_active")
	assert.NoError(t, err)
}
