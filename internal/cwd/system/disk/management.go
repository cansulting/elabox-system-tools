package disk

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/cansulting/elabox-system-tools/foundation/system"
)


func Check() (bool,error){
	cmd := exec.Command("/bin/sh", "-c", "sudo umount -l /dev/sda; sudo fsck -a /home/elabox; sudo mount /dev/sda /home/elabox")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr	
	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return false, err
	}
	fmt.Fprintln(os.Stdout)
	system.SetEnv("ELASHUTDOWNSTATUS", "disk_checked")			
	return true, nil
}