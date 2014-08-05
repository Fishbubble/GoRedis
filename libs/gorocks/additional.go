package gorocks

// #include <stdlib.h>
// #include "rocksdb/c.h"
import "C"

import (
	"unsafe"
)

/* Options */

func (o *Options) SetWalDir(dir string) {
	ldir := C.CString(dir)
	defer C.free(unsafe.Pointer(ldir))
	C.rocksdb_options_set_wal_dir(o.Opt, ldir)
}

func (o *Options) SetDbLogDir(dir string) {
	ldir := C.CString(dir)
	defer C.free(unsafe.Pointer(ldir))
	C.rocksdb_options_set_db_log_dir(o.Opt, ldir)
}
