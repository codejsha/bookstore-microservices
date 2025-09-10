package aggregate

type ShipmentAggregate struct {
	Uid                    string
	OrderUid               string
	CarrierUid             string
	Status                 ShipmentStatus
	TrackingNumber         string
	DestinationCity        string
	DestinationState       string
	DestinationCountryCode string
	DestinationPostalCode  string
	PlannedDeliveryAt      string
	ActualDeliveryAt       string
	CreatedAt              string
	UpdatedAt              string
}

type ShipmentStatus int32

const (
	SHIPMENTSTATUS_UNKNOWN          ShipmentStatus = 0
	SHIPMENTSTATUS_PLANNED          ShipmentStatus = 1
	SHIPMENTSTATUS_DISPATCHED       ShipmentStatus = 2
	SHIPMENTSTATUS_PICKED_UP        ShipmentStatus = 3
	SHIPMENTSTATUS_IN_TRANSIT       ShipmentStatus = 4
	SHIPMENTSTATUS_OUT_FOR_DELIVERY ShipmentStatus = 5
	SHIPMENTSTATUS_DELIVERED        ShipmentStatus = 6
	SHIPMENTSTATUS_FAILED           ShipmentStatus = 7
	SHIPMENTSTATUS_CANCELLED        ShipmentStatus = 8
)
