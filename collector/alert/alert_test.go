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

func TestHaCheckMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockWebService := mock_sapcontrol.NewMockWebService(ctrl)

	mockHACheckResponse := &sapcontrol.HACheckConfigResponse{
		Check: &sapcontrol.ArrayOfHACheck{
			Item: []*sapcontrol.HACheck{
				{State: sapcontrol.HA_VERIFICATION_STATE_ERROR, Category: sapcontrol.HA_CHECK_CATEGORY_HA_STATE, Description: "foo", Comment: "bar"},
				{State: sapcontrol.HA_VERIFICATION_STATE_WARNING, Category: sapcontrol.HA_CHECK_CATEGORY_SAP_STATE, Description: "foo2", Comment: "bar2"},
				{State: sapcontrol.HA_VERIFICATION_STATE_SUCCESS, Category: sapcontrol.HA_CHECK_CATEGORY_SAP_CONFIGURATION, Description: "foo3", Comment: "bar3"},
			},
		},
	}
	mockWebService.EXPECT().HACheckConfig().Return(mockHACheckResponse, nil)

	var err error
	collector, err := NewCollector(mockWebService)
	assert.NoError(t, err)

	expectedMetrics := `
	# HELP sap_alert_ha_check High Availability system configuration and status checks
	# TYPE sap_alert_ha_check gauge
	sap_alert_ha_check{category="HA-STATE",comment="bar",description="foo",state="ERROR"} 2
	sap_alert_ha_check{category="SAP-STATE",comment="bar2",description="foo2",state="WARNING"} 1
	sap_alert_ha_check{category="SAP-CONFIGURATION",comment="bar3",description="foo3",state="SUCCESS"} 0
`

	err = testutil.CollectAndCompare(collector, strings.NewReader(expectedMetrics))
	assert.NoError(t, err)
}
