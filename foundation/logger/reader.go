package logger

import (
	"encoding/json"
	"os"
	"strings"
	"sync"

	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/perm"
)

const BUF_LENGTH = 400 //1024 * 500
const LOG_FILE = constants.LOG_FILE

type Log map[string]interface{}

type Reader struct {
	logFile *os.File
	//buf        *bufio.Reader
	lastOffset int64
}

var stringPool sync.Pool = sync.Pool{
	New: func() interface{} {
		return make([]byte, BUF_LENGTH)
	},
}

// creates log reader instance. @logSrc is log file location, empty if use the default log file
func NewReader(logSrc string) (*Reader, error) {
	// init log file
	if logSrc == "" {
		logSrc = LOG_FILE
	}
	file, err := openLogfile(logSrc)
	if err != nil {
		return nil, err
	}
	return &Reader{
		lastOffset: -1,
		logFile:    file,
		//buf:      bufio.NewReader(file),
	}, nil
}

// load and read log file
func openLogfile(src string) (*os.File, error) {
	return os.OpenFile(src, os.O_RDONLY|os.O_CREATE, perm.PUBLIC)
}

// use to refresh file. the file might change last time
func (r *Reader) refreshFile() {
	var lastOffset int64 = -1
	info, err := os.Stat(r.logFile.Name())
	if err == nil {
		lastOffset = info.Size()
	}
	if r.lastOffset > 0 && lastOffset > 0 {
		return
	}
	r.lastOffset = lastOffset

}

// use to load some logs.
func (r *Reader) Load() []map[string]interface{} {
	r.refreshFile()
	for offset := r.lastOffset; offset >= 0; {
		offset -= BUF_LENGTH
		if offset < 0 {
			offset = 0
		}
		chunk := stringPool.Get().([]byte)
		readN, err := r.logFile.ReadAt(chunk, offset)
		if err != nil {
			println(err)
			break
		}
		if readN == 0 {
			println("Finished")
			break
		}
		// skip until newline is found
		if offset > 0 {
			skips := seekToNewline(chunk)
			if skips > 0 {
				chunk = chunk[skips:]
				offset += skips
			}
		}
		r.processChunk(chunk)
		stringPool.Put(chunk)
		if offset <= 0 {
			break
		}
	}
	return nil
}

// seek until newline is found. returns the number of bytes to skip
func seekToNewline(chunk []byte) int64 {
	length := len(chunk)
	for i := 0; i < length; i++ {
		if chunk[i] == '\n' {
			return int64(i)
		}
	}
	return -1
}

// use this to process chunk
func (r *Reader) processChunk(chunk []byte) []Log {
	str := string(chunk)
	splitted := strings.Split(str, "\n")
	logs := make([]Log, len(splitted))
	for i := len(splitted) - 1; i >= 0; i-- {
		var log Log
		if err := json.Unmarshal([]byte(splitted[i]), &log); err != nil {
			continue
		}
		logs[i] = log
		println(splitted[i])
	}
	return logs
}

func (r *Reader) Reset() {

}
