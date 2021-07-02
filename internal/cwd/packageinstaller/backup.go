package main

import (
	"archive/zip"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
)

const CONFIG_FILENAME = "config"

type Backup struct {
	PackageId  string       // package id of backup
	Version    string       // version of this backup
	Files      []BackupFile // files for backup
	sourceFile string       // location of this backup
	FileCount  uint16       // number of files in this backup
	archive    *zip.Writer
}

type BackupFile struct {
	Id  string `json:"id"`
	Src string `json:"src"`
}

// generate the backup to target location. call close after
func (this *Backup) Create(target string) error {
	log.Println("Creating backup @ " + target)
	this.sourceFile = target
	// step: create zip file
	backupFile, err := os.Create(target)
	if err != nil {
		return err
	}
	this.archive = zip.NewWriter(backupFile)
	this.Files = make([]BackupFile, 0, 5)
	return nil
}

// load the backup file
func (this *Backup) LoadAndApply(src string) error {
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
		// step: if file is config then convert and initialize this
		if file.Name == CONFIG_FILENAME {
			configBytes, _readErr := io.ReadAll(zipFile)
			if _readErr != nil {
				return &BackupError{errorStr: CONFIG_FILENAME + " coulnt load file." + _readErr.Error()}
			}
			err = json.Unmarshal(configBytes, this)
			if err != nil {
				return &BackupError{errorStr: CONFIG_FILENAME + " coudnt be load from backup. " + err.Error()}
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
	for _, file := range this.Files {
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

func (this *Backup) AddFile(src string) error {
	// add file to list
	id := strconv.Itoa(int(this.FileCount))
	this.Files = append(this.Files, BackupFile{
		Id:  id,
		Src: src,
	})
	this.FileCount++
	// add file to archive
	compressedFile, err := this.archive.Create(id)
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

// save this backup
func (this *Backup) Close() error {
	// ste: files is empty dont create backup
	if this.FileCount == 0 {
		this.archive.Close()
		os.Remove(this.sourceFile)
		return nil
	}
	defer this.archive.Close()
	// step: write config
	configBytes, err := json.Marshal(this)
	if err != nil {
		return err
	}
	compressedConfig, err := this.archive.Create(CONFIG_FILENAME)
	_, errConfig := compressedConfig.Write(configBytes)
	if errConfig != nil {
		return err
	}
	return nil
}
