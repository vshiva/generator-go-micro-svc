<%=licenseText%>
package health

// Probe is used to determine the health of a service.
type Probe interface {
	// Healthy checks wether the machine is healthy. It should return nil if it
	// is, otherwise it can return an error with more information.
	Healthy() error
}
