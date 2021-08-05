package jdb

import (
	"encoding/json"
	"fmt"
	"github.com/oknors/okno/app/cfg"
	"github.com/oknors/okno/pkg/utl"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/jcelliott/lumber"
)

var JDB, _ = NewJDB(cfg.Path, nil)

type (

	// Logger is a generic logger interface
	Logger interface {
		Fatal(string, ...interface{})
		Error(string, ...interface{})
		Warn(string, ...interface{})
		Info(string, ...interface{})
		Debug(string, ...interface{})
		Trace(string, ...interface{})
	}
	// Driver is what is used to interact with the scribble database. It runs
	// transactions, and provides log output
	jdb struct {
		col     string
		mutex   sync.Mutex
		mutexes map[string]*sync.Mutex
		path    string // the directory where scribble will create the database
		log     Logger // the logger scribble will log to
	}
)

// Options uses for specification of working golang-scribble
type Options struct {
	Logger // the logger scribble will use (configurable)
}

// New creates a new scribble database at the desired directory location, and
// returns a *Driver to then use for interacting with the database
func NewJDB(path string, options *Options) (*jdb, error) {
	// a new javazac database, providing the directory where it will be writing to,
	// and a qualified logger if desired

	path = filepath.Clean(path)

	// create default options
	opts := Options{}

	// if options are passed in, use those
	if options != nil {
		opts = *options
	}

	// if no logger is provided, create a default
	if opts.Logger == nil {
		opts.Logger = lumber.NewConsoleLogger(lumber.INFO)
	}

	j := jdb{
		path:    path,
		mutexes: make(map[string]*sync.Mutex),
		log:     opts.Logger,
	}

	// if the database already exists, just use it
	if _, err := os.Stat(path); err == nil {
		opts.Logger.Debug("Using '%s' (database already exists)\n", path)
		return &j, nil
	}

	// if the database doesn't exist create it
	opts.Logger.Debug("Creating scribble database at '%s'...\n", path)
	return &j, os.MkdirAll(path, 0755)
}

// Write locks the database and attempts to write the record to the database under
// the [collection] specified with the [resource] name given
func (j *jdb) Write(collection, resource string, v interface{}) error {

	// ensure there is a place to save record
	if collection == "" {
		return fmt.Errorf("Missing collection - no place to save record!")
	}

	// ensure there is a resource (name) to save record as
	if resource == "" {
		return fmt.Errorf("Missing resource - unable to save record (no name)!")
	}

	mutex := j.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	//
	dir := filepath.Join(j.path, collection)
	fnlPath := filepath.Join(dir, resource)
	tmpPath := fnlPath + ".tmp"

	// create collection pathectory
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	//
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}

	// write marshaled data to the temp file
	if err := ioutil.WriteFile(tmpPath, b, 0644); err != nil {
		return err
	}

	// move final file into place
	return os.Rename(tmpPath, fnlPath)
}

// Read a record from the database
func (j *jdb) Read(collection, resource string, v interface{}) error {

	// ensure there is a place to save record
	if collection == "" {
		return fmt.Errorf("Missing collection - no place to save record!")
	}

	// ensure there is a resource (name) to save record as
	if resource == "" {
		return fmt.Errorf("Missing resource - unable to save record (no name)!")
	}

	//
	record := filepath.Join(j.path, collection, resource)

	// check to see if file exists
	if _, err := stat(record); err != nil {
		return err
	}

	// read record from database
	b, err := ioutil.ReadFile(record)
	if err != nil {
		return err
	}

	// unmarshal data
	return json.Unmarshal(b, &v)
}

// ReadAll records from a collection; this is returned as a slice of strings because
// there is no way of knowing what type the record is.
func (j *jdb) ReadAll(collection string) ([]string, error) {

	// ensure there is a collection to read
	if collection == "" {
		return nil, fmt.Errorf("Missing collection - unable to record location!")
	}

	//
	dir := filepath.Join(j.path, collection)

	// check to see if collection (directory) exists
	if _, err := stat(dir); err != nil {
		return nil, err
	}

	// read all the files in the transaction.Collection; an error here just means
	// the collection is either empty or doesn't exist
	files, _ := ioutil.ReadDir(dir)

	// the files read from the database
	var records []string

	// iterate over each of the files, attempting to read the file. If successful
	// append the files to the collection of read files
	for _, file := range files {
		b, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		// append read file
		records = append(records, string(b))
	}

	// unmarhsal the read files as a comma delimeted byte array
	return records, nil
}

// Delete locks that database and then attempts to remove the collection/resource
// specified by [path]
func (j *jdb) Delete(collection, resource string) error {
	path := filepath.Join(collection, resource)
	//
	mutex := j.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	//
	dir := filepath.Join(j.path, path)

	switch fi, err := stat(dir); {

	// if fi is nil or error is not nil return
	case fi == nil, err != nil:
		return fmt.Errorf("Unable to find file or directory named %v\n", path)

	// remove directory and all contents
	case fi.Mode().IsDir():
		return os.RemoveAll(dir)

	// remove file
	case fi.Mode().IsRegular():
		return os.RemoveAll(dir)
	}

	return nil
}

//
func stat(path string) (fi os.FileInfo, err error) {

	// check for dir, if path isn't a directory check to see if it's a file
	if fi, err = os.Stat(path); os.IsNotExist(err) {
		fi, err = os.Stat(path)
	}

	return
}

// getOrCreateMutex creates a new collection specific mutex any time a collection
// is being modfied to avoid unsafe operations
func (j *jdb) getOrCreateMutex(collection string) *sync.Mutex {

	j.mutex.Lock()
	defer j.mutex.Unlock()

	m, ok := j.mutexes[collection]

	// if the mutex doesn't exist make it
	if !ok {
		m = &sync.Mutex{}
		j.mutexes[collection] = m
	}

	return m
}

// ReadCoins reads in all coin data in and converts to bytes for unmarshalling
func ReadData(path string) [][]byte {
	data, err := JDB.ReadAll(path)
	utl.ErrorLog(err)
	b := make([][]byte, len(data))
	for i := range data {
		b[i] = []byte(data[i])
	}
	return b
}
