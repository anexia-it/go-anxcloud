package v1

import (
	"go.anx.io/go-anxcloud/pkg/api/types"
)

func (obj *Record) DeepCopy() types.Object {
	// Initialize arrays
	

	out := &Record {
		// Primitives
		Identifier: obj.Identifier,
		ZoneName: obj.ZoneName,
		Immutable: obj.Immutable,
		Name: obj.Name,
		RData: obj.RData,
		Region: obj.Region,
		TTL: obj.TTL,
		Type: obj.Type,
		

		// DeepCopyable
		
		
		// Arrays
		
	}

	

	return out
}

func (obj *Revision) DeepCopy() *Revision {
	// Initialize arrays
	copyOfRecords := make([]Record, 0, len(obj.Records))
	for _, v := range obj.Records {
		copyOfRecords = append(copyOfRecords, *v.DeepCopy().(*Record))
	}
	

	out := &Revision {
		// Primitives
		CreatedAt: obj.CreatedAt,
		Identifier: obj.Identifier,
		ModifiedAt: obj.ModifiedAt,
		Serial: obj.Serial,
		State: obj.State,
		

		// DeepCopyable
		
		
		// Arrays
		Records: copyOfRecords,
		
	}

	

	return out
}

func (obj *DNSServer) DeepCopy() *DNSServer {
	// Initialize arrays
	

	out := &DNSServer {
		// Primitives
		Server: obj.Server,
		Alias: obj.Alias,
		

		// DeepCopyable
		
		
		// Arrays
		
	}

	

	return out
}

func (obj *Zone) DeepCopy() types.Object {
	// Initialize arrays
	copyOfNotifyAllowedIPs := make([]string, 0, len(obj.NotifyAllowedIPs))
	for _, v := range obj.NotifyAllowedIPs {
		copyOfNotifyAllowedIPs = append(copyOfNotifyAllowedIPs, v)
	}
	copyOfDNSServers := make([]DNSServer, 0, len(obj.DNSServers))
	for _, v := range obj.DNSServers {
		copyOfDNSServers = append(copyOfDNSServers, *v.DeepCopy())
	}
	copyOfRevisions := make([]Revision, 0, len(obj.Revisions))
	for _, v := range obj.Revisions {
		copyOfRevisions = append(copyOfRevisions, *v.DeepCopy())
	}
	

	out := &Zone {
		// Primitives
		Name: obj.Name,
		IsMaster: obj.IsMaster,
		DNSSecMode: obj.DNSSecMode,
		AdminEmail: obj.AdminEmail,
		Refresh: obj.Refresh,
		Retry: obj.Retry,
		Expire: obj.Expire,
		TTL: obj.TTL,
		MasterNS: obj.MasterNS,
		Customer: obj.Customer,
		CreatedAt: obj.CreatedAt,
		UpdatedAt: obj.UpdatedAt,
		PublishedAt: obj.PublishedAt,
		IsEditable: obj.IsEditable,
		ValidationLevel: obj.ValidationLevel,
		DeploymentLevel: obj.DeploymentLevel,
		CurrentRevision: obj.CurrentRevision,
		

		// DeepCopyable
		
		
		// Arrays
		NotifyAllowedIPs: copyOfNotifyAllowedIPs,
		DNSServers: copyOfDNSServers,
		Revisions: copyOfRevisions,
		
	}

	

	return out
}
