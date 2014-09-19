package gorocks

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"

import (
	"unsafe"
)

/* DB */
/*
extern rocksdb_t* rocksdb_open_column_families(
    const rocksdb_options_t* options,
    const char* name,
    int num_column_families,
    const char** column_family_names,
    const rocksdb_options_t** column_family_options,
    rocksdb_column_family_handle_t** column_family_handles,
    char** errptr);
*/
func OpenColumnFamilies(dbname string, o *Options, cfnames []string, cfopts []*Options, cfhandles []*ColumnFamilyHandle) (*DB, error) {
	ldbname := C.CString(dbname)
	defer C.free(unsafe.Pointer(ldbname))

	num := len(cfnames)
	ptrcfnames := make([]*C.char, num)
	ptrcfopts := make([]*C.rocksdb_options_t, num)
	ptrcfhandles := make([]*C.rocksdb_column_family_handle_t, num)

	for i := 0; i < num; i++ {
		ptrcfnames[i] = C.CString(cfnames[i])
		ptrcfopts[i] = cfopts[i].Opt
		ptrcfhandles[i] = cfhandles[i].cf
	}

	defer func() {
		for i := 0; i < num; i++ {
			C.free(unsafe.Pointer(ptrcfnames[i]))
		}
	}()

	var errStr *C.char
	rocksdb := C.rocksdb_open_column_families(o.Opt, ldbname, C.int(num), &ptrcfnames[0], &ptrcfopts[0], &ptrcfhandles[0], &errStr)
	if errStr != nil {
		gs := C.GoString(errStr)
		C.free(unsafe.Pointer(errStr))
		return nil, DatabaseError(gs)
	}

	return &DB{rocksdb}, nil
}

func (db *DB) CreateColumnFamily(name string, o *Options) error {
	lname := C.CString(name)
	defer C.free(unsafe.Pointer(lname))
	var errStr *C.char
	C.rocksdb_create_column_family(db.Ldb, o.Opt, lname, &errStr)
	if errStr != nil {
		gs := C.GoString(errStr)
		C.free(unsafe.Pointer(errStr))
		return DatabaseError(gs)
	}
	return nil
}

func (db *DB) DropColumnFamily(cf *ColumnFamilyHandle) error {
	var errStr *C.char
	C.rocksdb_drop_column_family(db.Ldb, cf.cf, &errStr)
	if errStr != nil {
		gs := C.GoString(errStr)
		C.free(unsafe.Pointer(errStr))
		return DatabaseError(gs)
	}
	return nil
}

/* ColumnFamilyHandle */
type ColumnFamilyHandle struct {
	cf *C.rocksdb_column_family_handle_t
}

func (c *ColumnFamilyHandle) Destory() {
	C.rocksdb_column_family_handle_destroy(c.cf)
}

/* Options */

// This specifies the absolute dir path for write-ahead logs (WAL).
// If it is empty, the log files will be in the same dir as data,
//   dbname is used as the data dir by default
// If it is non empty, the log files will be in kept the specified dir.
// When destroying the db,
//   all log files in wal_dir and the dir itself is deleted
func (o *Options) SetWalDir(dir string) {
	ldir := C.CString(dir)
	defer C.free(unsafe.Pointer(ldir))
	C.rocksdb_options_set_wal_dir(o.Opt, ldir)
}

// This specifies the info LOG dir.
// If it is empty, the log files will be in the same dir as data.
// If it is non empty, the log files will be in the specified dir,
// and the db data dir's absolute path will be used as the log file
// name's prefix.
func (o *Options) SetDbLogDir(dir string) {
	ldir := C.CString(dir)
	defer C.free(unsafe.Pointer(ldir))
	C.rocksdb_options_set_db_log_dir(o.Opt, ldir)
}

// Specify the maximal size of the info log file. If the log file
// is larger than `max_log_file_size`, a new info log file will
// be created.
// If max_log_file_size == 0, all logs will be written to one
// log file.
func (o *Options) SetMaxLogFileSize(s int) {
	C.rocksdb_options_set_max_log_file_size(o.Opt, C.size_t(s))
}

// Time for the info log file to roll (in seconds).
// If specified with non-zero value, log file will be rolled
// if it has been active longer than `log_file_time_to_roll`.
// Default: 0 (disabled)
func (o *Options) SetLogFileTimeToRoll(s int) {
	C.rocksdb_options_set_log_file_time_to_roll(o.Opt, C.size_t(s))
}

// Maximal info log files to be kept.
// Default: 1000
func (o *Options) SetKeepLogFileNum(n int) {
	C.rocksdb_options_set_keep_log_file_num(o.Opt, C.size_t(n))
}
