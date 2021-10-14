// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

package logger

import (
	"encoding/json"
	"os"
	"strings"
	"sync"

	"github.com/cansulting/elabox-system-tools/foundation/constants"
	"github.com/cansulting/elabox-system-tools/foundation/perm"
)

const CHUNK_SIZE = 1024 * 10 // the chunk to read per batch. the reading will be separated into batch to minimize load
const LOG_FILE = constants.LOG_FILE

type Log map[string]interface{}

// this struct provides log reading.
type Reader struct {
	logFile      *os.File
	EndingOffset int64 // the end position in log file
}

// reuseable pool of bytes
var chunkPool sync.Pool = sync.Pool{
	New: func() interface{} {
		return make([]byte, CHUNK_SIZE)
	},
}

// reusable pool of string
var stringPool = sync.Pool{
	New: func() interface{} {
		str := ""
		return str
	},
}

var logPool = sync.Pool{
	New: func() interface{} {
		return Log{}
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
		EndingOffset: -1,
		logFile:      file,
		//buf:      bufio.NewReader(file),
	}, nil
}

// load and read log file
func openLogfile(src string) (*os.File, error) {
	return os.OpenFile(src, os.O_RDONLY|os.O_CREATE, perm.PUBLIC)
}

// use to refresh file. the file might changed last time
func (r *Reader) refreshFile() {
	var newEndingOffset int64 = -1
	info, err := os.Stat(r.logFile.Name())
	if err == nil {
		newEndingOffset = info.Size()
	}
	if r.EndingOffset > 0 && newEndingOffset > 0 {
		return
	}
	r.EndingOffset = newEndingOffset

}

// use to load some logs.
// @start - start position of file reading backwards. start <= 0 means start from latest log
// @length - the number of bytes to read. -1 if read all. 0 if use the single chunk, see the CHUNK_SIZE
// @return - new offset
func (r *Reader) Load(start int64, length int64, filter func(int, Log) bool) int64 {
	r.refreshFile()
	if r.EndingOffset <= 0 {
		return 0
	}
	var from int64 = 0
	var to int64 = r.EndingOffset
	if start > 0 && start < r.EndingOffset {
		to = start
	}
	if length >= 0 && length < CHUNK_SIZE {
		length = CHUNK_SIZE
	} else if length < 0 {
		length = r.EndingOffset
	}
	from = to - length
	if from < 0 {
		from = 0
	}
	//println("Reading " + strconv.Itoa(int(from)) + " - " + strconv.Itoa(int(to)))
	// read file bytes from specified range
	chunkI := 0 // counter for chunk
	var waitG sync.WaitGroup
	offset := to
	for offset >= from {
		// step: initialize chunk
		chunk := chunkPool.Get().([]byte)
		toffset := offset - int64(len(chunk))
		// create a chunk with new size if it doesnt fit
		if toffset < 0 {
			chunkPool.Put(chunk)
			chunk = make([]byte, offset)
			toffset = 0
		}
		offset = toffset

		//  step: read file at offset
		readN, err := r.logFile.ReadAt(chunk, offset)
		if err != nil {
			println(err)
			break
		}
		if readN == 0 {
			println("Finished")
			break
		}
		// step: fix heading
		// the heading might be incomplete
		logs := stringPool.Get().(string)
		if offset > 0 && chunk[0] != '\n' {
			hchunks, newOffset := r.findHeadingFragment(offset)
			offset = newOffset
			chunk = append(hchunks, chunk...)
			logs = string(chunk)
		} else {
			logs = string(chunk)
		}
		waitG.Add(1)
		go func() {
			r.processLogs(logs, chunkI, filter)
			waitG.Done()
		}()
		// if len(chunk) < CHUNK_SIZE {
		// 	chunk = append(chunk, make([]byte, CHUNK_SIZE-len(chunk))...)
		// }
		chunkPool.Put(chunk)
		if offset <= 0 {
			break
		}
		chunkI++
	}
	waitG.Wait()
	return offset
}

// load logs and stop when specific limit of logs returned
// @start - start position of file reading backwards. start <= 0 means start from latest log
// @limit - max log to load. -1 if theres no limit
// @return - total loaded, new offset
func (r *Reader) LoadLimit(start int64, limit int, filter func(Log) bool) (int, int64) {
	r.refreshFile()
	if r.EndingOffset <= 0 {
		return 0, 0
	}
	offset := start
	if start <= 0 {
		offset = r.EndingOffset
	}
	if limit < 0 {
		limit = 1000000
	}
	i := 0
	for i < limit && offset >= 0 {
		offset = r.Load(offset, CHUNK_SIZE, func(chunkI int, l Log) bool {
			if filter(l) {
				i++
				// limit was achieved. quit to inner loop
				if i >= int(limit) {
					return false
				}
			}
			return true
		})
		// limit was achieved
		if i >= int(limit) || offset <= 0 {
			break
		}
		//offset -= CHUNK_SIZE
	}
	return i, offset
}

// seek until newline is found. returns index
func searchNewline(chunk []byte) int {
	length := len(chunk)
	for i := 0; i < length; i++ {
		if chunk[i] == '\n' {
			return i
		}
	}
	return -1
}

var tmpChunk = make([]byte, 20)

// This lookup the missing heading of json value. This specifically searches for newline
// @param offset - the tail position of file.
// @return []byte - the missing heading, int64 - the new tail offset
func (r *Reader) findHeadingFragment(offset int64) ([]byte, int64) {
	var size int64 = 20
	heading := make([]byte, size)
	init := false
	// iterate starting from end of file
	for i := offset - size; ; i -= size {
		// nothing left to process? then read from 0 offset and return
		if i < 0 {
			missingL := size + i
			arry := make([]byte, missingL, size)
			_, err := r.logFile.ReadAt(arry, 0)
			if err != nil {
				panic(err)
			}
			heading = append(arry, heading...)
			return heading, 0
		}
		// read string
		_, err := r.logFile.ReadAt(tmpChunk, i)
		if err != nil {
			panic(err)
		}
		// append the process string
		if init {
			heading = append(tmpChunk, heading...)
		} else {
			copy(heading, tmpChunk)
			init = true
		}
		// search for newline. if found YEHEY!
		foundI := searchNewline(tmpChunk)
		if foundI >= 0 {
			return heading, offset - int64(len(tmpChunk)) + int64(foundI)
		}
	}
	//return heading, offset
}

// use this to process chunk
func (r *Reader) processLogs(chunkStr string, chunkIndex int, filter func(int, Log) bool) {
	splitted := strings.Split(chunkStr, "\n")
	logs := make([]Log, len(splitted))
	hasFilter := false
	if filter != nil {
		hasFilter = true
	}
	for i := len(splitted) - 1; i >= 0; i-- {
		log := logPool.Get().(Log)
		// failed parsing json
		if err := json.Unmarshal([]byte(splitted[i]), &log); err != nil {
			continue
		}
		logs[i] = log
		if hasFilter {
			if !filter(chunkIndex, log) {
				logPool.Put(log)
				return
			}
		}
		logPool.Put(log)
		//println(splitted[i])
	}
	stringPool.Put(chunkStr)
}

func (r *Reader) Reset() {

}
