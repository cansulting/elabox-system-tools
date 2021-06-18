package servicecenter

type ServiceData struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	BinaryPath  string `json:"binaryPath"`
}
