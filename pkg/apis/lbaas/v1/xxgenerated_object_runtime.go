package v1

import (
	"go.anx.io/go-anxcloud/pkg/api/types"
)

func (obj *Backend) DeepCopy() types.Object {
	// Initialize arrays
	copyOfAutomationRules := make([]RuleInfo, 0, len(obj.AutomationRules))
	for _, v := range obj.AutomationRules {
		copyOfAutomationRules = append(copyOfAutomationRules, *v.DeepCopy())
	}
	

	out := &Backend {
		// Primitives
		CustomerIdentifier: obj.CustomerIdentifier,
		ResellerIdentifier: obj.ResellerIdentifier,
		Identifier: obj.Identifier,
		Name: obj.Name,
		HealthCheck: obj.HealthCheck,
		Mode: obj.Mode,
		ServerTimeout: obj.ServerTimeout,
		LoadBalancer: obj.LoadBalancer,
		

		// DeepCopyable
		
		
		// Arrays
		AutomationRules: copyOfAutomationRules,
		
	}

	

	return out
}

func (obj *Bind) DeepCopy() types.Object {
	// Initialize arrays
	copyOfAutomationRules := make([]RuleInfo, 0, len(obj.AutomationRules))
	for _, v := range obj.AutomationRules {
		copyOfAutomationRules = append(copyOfAutomationRules, *v.DeepCopy())
	}
	

	out := &Bind {
		// Primitives
		CustomerIdentifier: obj.CustomerIdentifier,
		ResellerIdentifier: obj.ResellerIdentifier,
		Identifier: obj.Identifier,
		Name: obj.Name,
		Address: obj.Address,
		Port: obj.Port,
		SSL: obj.SSL,
		SslCertificatePath: obj.SslCertificatePath,
		Frontend: obj.Frontend,
		

		// DeepCopyable
		
		
		// Arrays
		AutomationRules: copyOfAutomationRules,
		
	}

	

	return out
}

func (obj *Frontend) DeepCopy() types.Object {
	// Initialize arrays
	copyOfAutomationRules := make([]RuleInfo, 0, len(obj.AutomationRules))
	for _, v := range obj.AutomationRules {
		copyOfAutomationRules = append(copyOfAutomationRules, *v.DeepCopy())
	}
	

	out := &Frontend {
		// Primitives
		CustomerIdentifier: obj.CustomerIdentifier,
		ResellerIdentifier: obj.ResellerIdentifier,
		Identifier: obj.Identifier,
		Name: obj.Name,
		Mode: obj.Mode,
		ClientTimeout: obj.ClientTimeout,
		

		// DeepCopyable
		LoadBalancer: obj.LoadBalancer.DeepCopy().(*LoadBalancer),
		DefaultBackend: obj.DefaultBackend.DeepCopy().(*Backend),
		
		
		// Arrays
		AutomationRules: copyOfAutomationRules,
		
	}

	

	return out
}

func (obj *LoadBalancer) DeepCopy() types.Object {
	// Initialize arrays
	copyOfAutomationRules := make([]RuleInfo, 0, len(obj.AutomationRules))
	for _, v := range obj.AutomationRules {
		copyOfAutomationRules = append(copyOfAutomationRules, *v.DeepCopy())
	}
	

	out := &LoadBalancer {
		// Primitives
		CustomerIdentifier: obj.CustomerIdentifier,
		ResellerIdentifier: obj.ResellerIdentifier,
		Identifier: obj.Identifier,
		Name: obj.Name,
		IpAddress: obj.IpAddress,
		

		// DeepCopyable
		
		
		// Arrays
		AutomationRules: copyOfAutomationRules,
		
	}

	

	return out
}

func (obj *Server) DeepCopy() types.Object {
	// Initialize arrays
	copyOfAutomationRules := make([]RuleInfo, 0, len(obj.AutomationRules))
	for _, v := range obj.AutomationRules {
		copyOfAutomationRules = append(copyOfAutomationRules, *v.DeepCopy())
	}
	

	out := &Server {
		// Primitives
		CustomerIdentifier: obj.CustomerIdentifier,
		ResellerIdentifier: obj.ResellerIdentifier,
		Identifier: obj.Identifier,
		Name: obj.Name,
		IP: obj.IP,
		Port: obj.Port,
		Check: obj.Check,
		Backend: obj.Backend,
		

		// DeepCopyable
		
		
		// Arrays
		AutomationRules: copyOfAutomationRules,
		
	}

	

	return out
}

func (obj *State) DeepCopy() *State {
	// Initialize arrays
	

	out := &State {
		// Primitives
		ID: obj.ID,
		Text: obj.Text,
		Type: obj.Type,
		

		// DeepCopyable
		
		
		// Arrays
		
	}

	

	return out
}

func (obj *HasState) DeepCopy() *HasState {
	// Initialize arrays
	

	out := &HasState {
		// Primitives
		State: obj.State,
		

		// DeepCopyable
		
		
		// Arrays
		
	}

	

	return out
}

func (obj *RuleInfo) DeepCopy() *RuleInfo {
	// Initialize arrays
	

	out := &RuleInfo {
		// Primitives
		Identifier: obj.Identifier,
		Name: obj.Name,
		

		// DeepCopyable
		
		
		// Arrays
		
	}

	

	return out
}
