package logging

//LogDetails structure for all logrus error and warnings
type LogDetails struct {
	Code    string
	Message string
	Details interface{}
	Error   error
}
