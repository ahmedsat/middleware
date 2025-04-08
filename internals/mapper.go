package internals

import "io"

// Mapper is the interface for the mapping data from one format to another
type Mapper interface {
	// Scan reads the data from the source format as a reader and maps it to the struct
	Scan(r io.Reader) (err error)
	// Validate checks if the data is valid
	Validate() (err error)
	// Value returns the data mapped to needed format as a reader
	Value() (r io.Reader, err error)
}
