package bus

import "github.com/oklog/ulid/v2"

type ServiceRequestData interface {
	ServiceType() string
}

type serviceRequest struct {
	Data ServiceRequestData
	ID   ulid.ULID
}

type ServiceResponseEvent struct {
	RequestID ulid.ULID
	Response  any
}

// EventType implements EventData.
func (s *ServiceResponseEvent) EventType() string {
	return "service_response"
}

var _ EventData = (*ServiceResponseEvent)(nil)
