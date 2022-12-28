package broadcast

import "errors"



func Init() error {
	if err := registerRecievers(); err != nil {
		return errors.New("failed to register broadcast recievers. inner: " + err.Error())
	}
	return nil
}
