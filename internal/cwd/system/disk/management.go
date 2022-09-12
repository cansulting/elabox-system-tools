package disk

import (
	"os"
	"os/exec"

	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/foundation/system"
)


func Check() (bool,error){
	logger.GetInstance().Info().Msg("Checking disk")	
	cmd := exec.Command("/bin/sh", "-c", "sudo umount -l /dev/sda; sudo fsck -a /home/elabox; sudo mount /dev/sda /home/elabox")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr	
	err := cmd.Run()
	if err != nil {
		logger.GetInstance().Error().Err(err).Msg("Failed checking disk " + err.Error())
		return false, err
	}
	system.SetEnv("ELASHUTDOWNSTATUS", "disk_checked")		
	logger.GetInstance().Info().Msg("Disk checked")	
	return true, nil
}