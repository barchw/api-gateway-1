package processing

import (
	v1beta1 "istio.io/api/networking/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Action int

const (
	Create Action = iota
	Update
	Delete
)

func (s Action) String() string {
	switch s {
	case Create:
		return "create"
	case Update:
		return "update"
	case Delete:
		return "delete"
	}
	return "unknown"
}

type ObjectChange struct {
	Action Action
	Obj    client.Object
}

func NewObjectCreateAction(obj client.Object) *ObjectChange {
	return &ObjectChange{
		Action: Create,
		Obj:    obj,
	}
}

func NewObjectUpdateAction(obj client.Object) *ObjectChange {
	return &ObjectChange{
		Action: Update,
		Obj:    obj,
	}
}

func NewObjectDeleteAction(obj client.Object) *ObjectChange {
	return &ObjectChange{
		Action: Delete,
		Obj:    obj,
	}
}

// CorsConfig is an internal representation of v1alpha3.CorsPolicy object
type CorsConfig struct {
	AllowOrigins []*v1beta1.StringMatch
	AllowMethods []string
	AllowHeaders []string
}

type ReconciliationConfig struct {
	OathkeeperSvc     string
	OathkeeperSvcPort uint32
	CorsConfig        *CorsConfig
	AdditionalLabels  map[string]string
	DefaultDomainName string
	ServiceBlockList  map[string][]string
	DomainAllowList   []string
	HostBlockList     []string
}
