package appman

func Initialize(commandline bool) error {
	if !commandline {
		InitializeStartups()
	}
	return nil
}
