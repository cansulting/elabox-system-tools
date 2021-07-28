package utils

import (
	"archive/zip"
	"ela/foundation/errors"
	"ela/foundation/perm"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const CONFIG_FILENAME = "config"

type Backup struct {
	PackageId  string       // package id of backup
	Version    string       // version of instance backup
	Files      []BackupFile // files for backup
	sourceFile string       // location of instance backup
	FileCount  uint16       // number of files in instance backup
	archive    *zip.Writer
}

type BackupFile struct {
	Id   string      `json:"id"`
	Src  string      `json:"src"`
	Perm os.FileMode `json:"perm"`
}

func (bkfile BackupFile) open() *os.File {
	file, _ := os.OpenFile(bkfile.Src, os.O_CREATE|os.O_WRONLY, bkfile.Perm)
	return file
}

func (instance *Backup) GetSource() string {
	return instance.sourceFile
}

// generate the backup to target location. call close after
func (instance *Backup) Create(target string) error {
	log.Println("Creating backup @ "+target, "packageId", instance.PackageId)
	instance.sourceFile = target
	// step: create zip file
	backupFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm.PUBLIC)
	if err != nil {
		return err
	}
	instance.archive = zip.NewWriter(backupFile)
	instance.Files = make([]BackupFile, 0, 5)
	return nil
}

// load the backup file
func (instance *Backup) LoadAndApply(src string) error {
	// init cache dir
	tempDir, _ := os.Getwd()
	tempDir += "/temp"
	os.MkdirAll(tempDir, 0765)
	defer os.RemoveAll(tempDir)
	// step: open zip
	zipFile, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	// iterate through compressed files and find the config file
	for _, file := range zipFile.File {
		zipFile, err := file.Open()
		if err != nil {
			return err
		}
		defer zipFile.Close()
		// step: if file is config then load and initialize instance
		if file.Name == CONFIG_FILENAME {
			configBytes, _readErr := io.ReadAll(zipFile)
			if _readErr != nil {
				return errors.SystemNew(CONFIG_FILENAME+" coulnt load file.", _readErr)
			}
			err = json.Unmarshal(configBytes, instance)
			if err != nil {
				return errors.SystemNew(CONFIG_FILENAME+" coudnt be load from backup.", _readErr)
			}
			continue
		}
		// save file to tempDir
		newFile, err := os.Create(tempDir + "/" + file.Name)
		if err != nil {
			return err
		}
		defer newFile.Close()
		_, errWrite := io.Copy(newFile, zipFile)
		if errWrite != nil {
			return err
		}
	}
	// step: move files to source
	for _, file := range instance.Files {
		os.MkdirAll(file.Src, 0765)
		// remove old files
		if _, err := os.Stat(file.Src); err == nil {
			os.Remove(file.Src)
		}
		// copy from temp to target
		fromName := tempDir + "/" + file.Id
		from, err := os.Open(fromName)
		_, err2 := io.Copy(file.open(), from)
		if err2 != nil {
			log.Println("Backup:LoadAndApply skippable ", err2.Error())
		}
		// remove temp
		err3 := os.Remove(fromName)
		if err != nil {
			log.Println("Backup:LoadAndApply skippable ", err3.Error())
		}
	}
	return nil
}

func (instance *Backup) AddFile(src string) error {
	// add file to list
	id := strconv.Itoa(int(instance.FileCount))
	// add file to archive
	compressedFile, err := instance.archive.Create(id)
	if err != nil {
		return err
	}
	// read bytes from file and save to archive
	uncompressedFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer uncompressedFile.Close()
	_, errCopy := io.Copy(compressedFile, uncompressedFile)
	if errCopy != nil {
		return errCopy
	}
	info, _ := uncompressedFile.Stat()
	// update backup info
	instance.Files = append(instance.Files, BackupFile{
		Id:   id,
		Src:  src,
		Perm: info.Mode(),
	})
	instance.FileCount++
	return nil
}

func (instance *Backup) AddFiles(srcDir string) error {
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() {
			if err := instance.AddFile(srcDir + "/" + file.Name()); err != nil {
				return err
			}
			continue
		}
		instance.AddFiles(srcDir + "/" + file.Name())
	}
	return nil
}

// save instance backup
func (instance *Backup) Close() error {
	// ste: files is empty dont create backup
	if instance.FileCount == 0 {
		if instance.archive != nil {
			instance.archive.Close()
		}
		os.Remove(instance.sourceFile)
		return nil
	}
	defer instance.archive.Close()
	// step: write config
	configBytes, err := json.Marshal(instance)
	if err != nil {
		return err
	}
	compressedConfig, err := instance.archive.Create(CONFIG_FILENAME)
	_, errConfig := compressedConfig.Write(configBytes)
	if errConfig != nil {
		return err
	}
	return nil
}
