package main

type BackupError struct {
	errorStr string
}

func (b *BackupError) Error() string {
	return b.errorStr
}
