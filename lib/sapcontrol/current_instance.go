package sapcontrol

import (
	"fmt"
	"strconv"
	"math"

	"github.com/pkg/errors"
)

// this structure will be used for static labels, common to all the metrics
type CurrentSapInstance struct {
	SID      string
	Number   int32
	Name     string
	Hostname string
}

func (i *CurrentSapInstance) String() string {
	return fmt.Sprintf("SID: %s, Name: %s, Number: %d, Hostname: %s", i.SID, i.Name, i.Number, i.Hostname)
}

func (s *webService) GetCurrentInstance() (*CurrentSapInstance, error) {
	var err error

	// since the information we want here doesn't change over time, we only want to execute the remote call once
	s.once.Do(func() {
		var response *GetInstancePropertiesResponse
		response, err = s.GetInstanceProperties()

		if err != nil {
			err = errors.Wrap(err, "could not perform GetInstanceProperties query")
			return
		}

		s.currentSapInstance = &CurrentSapInstance{}

		for _, prop := range response.Properties {
			switch prop.Property {
			case "SAPSYSTEM":
				var num int64
				num, err = strconv.ParseInt(prop.Value, 10, 32)
				if err != nil {
					err = errors.Wrap(err, "could not parse instance number to int32")
					return
				}
				if num < math.MinInt32 || num > math.MaxInt32 {
					err = errors.New("parsed instance number out of int32 range")
					return
				}
				s.currentSapInstance.Number = int32(num)
			case "SAPSYSTEMNAME":
				s.currentSapInstance.SID = prop.Value
			case "INSTANCE_NAME":
				s.currentSapInstance.Name = prop.Value
			case "SAPLOCALHOST":
				s.currentSapInstance.Hostname = prop.Value
			}
		}
	})

	return s.currentSapInstance, err
}
