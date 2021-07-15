package utils

import (
	"archive/zip"
	"ela/foundation/errors"
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
	Id  string `json:"id"`
	Src string `json:"src"`
}

// generate the backup to target location. call close after
func (instance *Backup) Create(target string) error {
	log.Println("Creating backup @ " + target)
	instance.sourceFile = target
	// step: create zip file
	backupFile, err := os.Create(target)
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
	// iterate through compressed files
	for _, file := range zipFile.File {
		zipFile, err := file.Open()
		defer zipFile.Close()
		if err != nil {
			return err
		}
		// step: if file is config then convert and initialize instance
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
		defer newFile.Close()
		if err != nil {
			return err
		}
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
		err := os.Rename(tempDir+"/"+file.Id, file.Src)
		if err != nil {
			log.Println("Backup:LoadAndApply skippable ", err.Error())
		}
	}
	return nil
}

func (instance *Backup) AddFile(src string) error {
	// add file to list
	id := strconv.Itoa(int(instance.FileCount))
	instance.Files = append(instance.Files, BackupFile{
		Id:  id,
		Src: src,
	})
	instance.FileCount++
	// add file to archive
	compressedFile, err := instance.archive.Create(id)
	if err != nil {
		return err
	}
	// read bytes from file and save to archive
	uncompressedFile, err := os.OpenFile(src, os.O_RDONLY, 0764)
	if err != nil {
		return err
	}
	_, errCopy := io.Copy(compressedFile, uncompressedFile)
	if errCopy != nil {
		return errCopy
	}
	uncompressedFile.Close()
	return nil
}

func (instance *Backup) AddFiles(srcDir string) error {
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err := instance.AddFile(srcDir + "/" + file.Name()); err != nil {
			return err
		}
	}
	return nil
}

// save instance backup
func (instance *Backup) Close() error {
	// ste: files is empty dont create backup
	if instance.FileCount == 0 {
		instance.archive.Close()
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
