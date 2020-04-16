package sapcontrol

import (
	"strconv"

	"github.com/pkg/errors"
)

// this structure will be used for static labels, common to all the metrics
type CurrentSapInstance struct {
	SID      string
	Number   int32
	Name     string
	Hostname string
}

func NewCurrentInstance(service WebService) (*CurrentSapInstance, error) {
	sapInstance := &CurrentSapInstance{}
	response, err := service.GetInstanceProperties()
	if err != nil {
		return nil, errors.Wrap(err, "SAPControl webservice error")
	}
	for _, prop := range response.Properties {
		switch prop.Property {
		case "SAPSYSTEM":
			num, err := strconv.Atoi(prop.Value)
			if err != nil {
				return nil, errors.Wrap(err, "could not parse instance number to int")
			}
			sapInstance.Number = int32(num)
		case "SAPSYSTEMNAME":
			sapInstance.SID = prop.Value
		case "INSTANCE_NAME":
			sapInstance.Name = prop.Value
		case "SAPLOCALHOST":
			sapInstance.Hostname = prop.Value
		}
	}

	return sapInstance, nil
}
