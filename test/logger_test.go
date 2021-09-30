package main

import (
	"ela/foundation/logger"
	"testing"
)

func TestLogger(t *testing.T) {
	logger.Init("ela.testing")
	logger.GetInstance().Debug().Msg("Hello")
	t.Log("Sucess")
}
