// Package types represents the interface "wasi:filesystem.types".
//
// WASI filesystem is a filesystem API primarily intended to let users run WASI
// programs that access their files on their existing filesystems, without
// significant overhead.
//
// It is intended to be roughly portable between Unix-family platforms and
// Windows, though it does not hide many of the major differences.
//
// Paths are passed as interface-type `string`s, meaning they must consist of
// a sequence of Unicode Scalar Values (USVs). Some filesystems may contain
// paths which are not accessible by this API.
//
// The directory separator in WASI is always the forward-slash (`/`).
//
// All paths in WASI are relative paths, and are interpreted relative to a
// `descriptor` referring to a base directory. If a `path` argument to any WASI
// function starts with `/`, or if any step of resolving a `path`, including
// `..` and symbolic link steps, reaches a directory outside of the base
// directory, or reaches a symlink to an absolute or rooted path in the
// underlying filesystem, the function fails with `error-code::not-permitted`.
//
// For more information about WASI path resolution and sandboxing, see
// [WASI filesystem path resolution].
//
// [WASI filesystem path resolution]: https://github.com/WebAssembly/wasi-filesystem/blob/main/path-resolution.md
package types

import (
	"github.com/ydnar/wasm-tools-go/cm"
	"github.com/ydnar/wasm-tools-go/wasi/clocks/wallclock"
	"github.com/ydnar/wasm-tools-go/wasi/io/streams"
)

type (
	InputStream  = streams.InputStream
	OutputStream = streams.OutputStream
	Error        = streams.Error
	DateTime     = wallclock.DateTime
)

// FileSize represents the type "wasi:filesystem/types.filesize".
//
// File size or length of a region within a file.
type FileSize = uint64

// DescriptorType represents the enum "wasi:filesystem.types.descriptor-type".
//
// The type of a filesystem object referenced by a descriptor.
//
// Note: This was called `filetype` in earlier versions of WASI.
type DescriptorType uint8

const (
	DescriptorTypeUnknown         DescriptorType = iota // The type of the descriptor or file is unknown or is different from any of the other types specified.
	DescriptorTypeBlockDevice                           // The descriptor refers to a block device inode.
	DescriptorTypeCharacterDevice                       // The descriptor refers to a character device inode.
	DescriptorTypeDirectory                             // The descriptor refers to a directory inode.
	DescriptorTypeFIFO                                  // The descriptor refers to a named pipe.
	DescriptorTypeSymbolicLink                          // The file refers to a symbolic link inode.
	DescriptorTypeRegularFile                           // The descriptor refers to a regular file inode.
	DescriptorTypeSocket                                // The descriptor refers to a socket.
)

// DescriptorFlags represents the flags "wasi:filesystem/types.descriptor-flags".
//
// Descriptor flags.
//
// Note: This was called `fdflags` in earlier versions of WASI.
type DescriptorFlags uint8

const (
	// Read mode: Data can be read.
	DescriptorFlagsRead DescriptorFlags = 1 << iota

	// Write mode: Data can be written to.
	DescriptorFlagsWrite

	// Request that writes be performed according to synchronized I/O file
	// integrity completion. The data stored in the file and the file's
	// metadata are synchronized. This is similar to `O_SYNC` in POSIX.
	//
	// The precise semantics of this operation have not yet been defined for
	// WASI. At this time, it should be interpreted as a request, and not a
	// requirement.
	DescriptorFlagsFileIntegritySync

	// Request that writes be performed according to synchronized I/O data
	// integrity completion. Only the data stored in the file is
	// synchronized. This is similar to `O_DSYNC` in POSIX.
	//
	// The precise semantics of this operation have not yet been defined for
	// WASI. At this time, it should be interpreted as a request, and not a
	// requirement.
	DescriptorFlagsDataIntegritySync

	// Requests that reads be performed at the same level of integrety
	// requested for writes. This is similar to `O_RSYNC` in POSIX.
	//
	// The precise semantics of this operation have not yet been defined for
	// WASI. At this time, it should be interpreted as a request, and not a
	// requirement.
	DescriptorFlagsRequestedWriteSync

	// Mutating directories mode: Directory contents may be mutated.
	//
	// When this flag is unset on a descriptor, operations using the
	// descriptor which would create, rename, delete, modify the data or
	// metadata of filesystem objects, or obtain another handle which
	// would permit any of those, shall fail with `error-code::read-only` if
	// they would otherwise succeed.
	//
	// This may only be set on directories.
	DescriptorFlagsMutateDirectory
)

// DescriptorStat represents the record "wasi:filesystem/types.descriptor-stat".
//
// File attributes.
//
// Note: This was called `filestat` in earlier versions of WASI.
type DescriptorStat struct {
	// File type.
	Type DescriptorType

	// Number of hard links to the file.
	LinkCount LinkCount

	// For regular files, the file size in bytes. For symbolic links, the
	// length in bytes of the pathname contained in the symbolic link.
	Size FileSize

	// Last data access timestamp.
	//
	// If the `option` is none, the platform doesn't maintain an access
	// timestamp for this file.
	DataAccessTimestamp cm.Option[DateTime]

	// Last data modification timestamp.
	//
	// If the `option` is none, the platform doesn't maintain a
	// modification timestamp for this file.
	DataModificationTimestamp cm.Option[DateTime]

	// Last file status-change timestamp.
	//
	// If the `option` is none, the platform doesn't maintain a
	// status-change timestamp for this file.
	StatusChangeTimestamp cm.Option[DateTime]
}

// PathFlags represents the flags "wasi:filesystem/types.path-flags".
//
// Flags determining the method of how paths are resolved.
type PathFlags uint8

const (
	// As long as the resolved path corresponds to a symbolic link, it is
	// expanded.
	PathFlagsSymlinkFollow PathFlags = 1 << iota
)

// OpenFlags represents the flags "wasi:filesystem/types.open-flags".
//
// Open flags used by `open-at`.
type OpenFlags uint8

const (
	// Create file if it does not exist, similar to `O_CREAT` in POSIX.
	OpenFlagsCreate OpenFlags = 1 << iota

	// Fail if not a directory, similar to `O_DIRECTORY` in POSIX.
	OpenFlagsDirectory

	// Fail if file already exists, similar to `O_EXCL` in POSIX.
	OpenFlagsExclusive

	// Truncate file to size 0, similar to `O_TRUNC` in POSIX.
	OpenFlagsTruncate
)

// LinkCount represents the type "wasi:filesystem/types.link-count".
//
// Number of hard links to an inode.
type LinkCount = uint64

// NewTimestamp represents the variant "wasi:filesystem/types.new-timestamp".
//
// When setting a timestamp, this gives the value to set it to.
type NewTimestamp struct {
	v cm.Variant[uint8, DateTime, struct{}]
}

// NewTimestampNoChange returns a NewTimestamp with variant case "no-change".
func NewTimestampNoChange() NewTimestamp {
	var result NewTimestamp
	cm.Set(&result.v, 0, struct{}{})
	return result
}

// NoChange represents variant case "no-change".
//
// Leave the timestamp set to its previous value.
func (self *NewTimestamp) NoChange() bool {
	return self.v.Is(0)
}

// NewTimestampNow returns a NewTimestamp with variant case "now".
func NewTimestampNow() NewTimestamp {
	var result NewTimestamp
	cm.Set(&result.v, 1, struct{}{})
	return result
}

// Now represents variant case "now".
//
// Leave the timestamp set to its previous value.
func (self *NewTimestamp) Now() bool {
	return self.v.Is(1)
}

// Timestamp represents variant case "timestamp(datetime)".
//
// Set the timestamp to the given value.
func (self *NewTimestamp) Timestamp() (DateTime, bool) {
	return cm.Get[DateTime](&self.v, 2)
}

// NewTimestampTimestamp returns a NewTimestamp with variant case "timestamp(datetime)".
func NewTimestampTimestamp(v DateTime) NewTimestamp {
	var result NewTimestamp
	cm.Set(&result.v, 2, v)
	return result
}

// DirectoryEntry represents the record "wasi:filesystem/types.directory-entry".
//
// A directory entry.
type DirectoryEntry struct {
	// The type of the file referred to by this directory entry.
	Type DescriptorType

	// The name of the object.
	Name string
}

// ErrorCode represents the enum "wasi:filesystem/types.error-code".
//
// Error codes returned by functions, similar to `errno` in POSIX.
// Not all of these error codes are returned by the functions provided by this
// API; some are used in higher-level library layers, and others are provided
// merely for alignment with POSIX.
type ErrorCode uint8

const (
	// Permission denied, similar to `EACCES` in POSIX.
	ErrorCodeAccess ErrorCode = iota

	// Resource unavailable, or operation would block, similar to `EAGAIN` and `EWOULDBLOCK` in POSIX.
	ErrorCodeWouldBlock

	// Connection already in progress, similar to `EALREADY` in POSIX.
	ErrorCodeAlready

	// Bad descriptor, similar to `EBADF` in POSIX.
	ErrorCodeBadDescriptor

	// Device or resource busy, similar to `EBUSY` in POSIX.
	ErrorCodeBusy

	// Resource deadlock would occur, similar to `EDEADLK` in POSIX.
	ErrorCodeDeadlock

	// Storage quota exceeded, similar to `EDQUOT` in POSIX.
	ErrorCodeQuota

	// File exists, similar to `EEXIST` in POSIX.
	ErrorCodeExist

	// File too large, similar to `EFBIG` in POSIX.
	ErrorCodeFileTooLarge

	// Illegal byte sequence, similar to `EILSEQ` in POSIX.
	ErrorCodeIllegalByteSequence

	// Operation in progress, similar to `EINPROGRESS` in POSIX.
	ErrorCodeInProgress

	// Interrupted function, similar to `EINTR` in POSIX.
	ErrorCodeInterrupted

	// Invalid argument, similar to `EINVAL` in POSIX.
	ErrorCodeInvalid

	// I/O error, similar to `EIO` in POSIX.
	ErrorCodeIo

	// Is a directory, similar to `EISDIR` in POSIX.
	ErrorCodeIsDirectory

	// Too many levels of symbolic links, similar to `ELOOP` in POSIX.
	ErrorCodeLoop

	// Too many links, similar to `EMLINK` in POSIX.
	ErrorCodeTooManyLinks

	// Message too large, similar to `EMSGSIZE` in POSIX.
	ErrorCodeMessageSize

	// Filename too long, similar to `ENAMETOOLONG` in POSIX.
	ErrorCodeNameTooLong

	// No such device, similar to `ENODEV` in POSIX.
	ErrorCodeNoDevice

	// No such file or directory, similar to `ENOENT` in POSIX.
	ErrorCodeNoEntry

	// No locks available, similar to `ENOLCK` in POSIX.
	ErrorCodeNoLock

	// Not enough space, similar to `ENOMEM` in POSIX.
	ErrorCodeInsufficientMemory

	// No space left on device, similar to `ENOSPC` in POSIX.
	ErrorCodeInsufficientSpace

	// Not a directory or a symbolic link to a directory, similar to `ENOTDIR` in POSIX.
	ErrorCodeNotDirectory

	// Directory not empty, similar to `ENOTEMPTY` in POSIX.
	ErrorCodeNotEmpty

	// State not recoverable, similar to `ENOTRECOVERABLE` in POSIX.
	ErrorCodeNotRecoverable

	// Not supported, similar to `ENOTSUP` and `ENOSYS` in POSIX.
	ErrorCodeUnsupported

	// Inappropriate I/O control operation, similar to `ENOTTY` in POSIX.
	ErrorCodeNoTTY

	// No such device or address, similar to `ENXIO` in POSIX.
	ErrorCodeNoSuchDevice

	// Value too large to be stored in data type, similar to `EOVERFLOW` in POSIX.
	ErrorCodeOverflow

	// Operation not permitted, similar to `EPERM` in POSIX.
	ErrorCodeNotPermitted

	// Broken pipe, similar to `EPIPE` in POSIX.
	ErrorCodePipe

	// ReadOnly file system, similar to `EROFS` in POSIX.
	ErrorCodeReadOnly

	// Invalid seek, similar to `ESPIPE` in POSIX.
	ErrorCodeInvalidSeek

	// Text file busy, similar to `ETXTBSY` in POSIX.
	ErrorCodeTextFileBusy

	// Cross-device link, similar to `EXDEV` in POSIX.
	ErrorCodeCrossDevice
)

// Advice represents the enum "wasi:filesystem/types.advice".
//
// File or memory access pattern advisory information.
type Advice uint8

const (
	// The application has no advice to give on its behavior with respect
	// to the specified data.
	Normal Advice = iota

	// The application expects to access the specified data sequentially
	// from lower offsets to higher offsets.
	Sequential

	// The application expects to access the specified data in a random
	// order.
	Random

	// The application expects to access the specified data in the near
	// future.
	WillNeed

	// The application expects that it will not access the specified data
	// in the near future.
	DontNeed

	// The application expects to access the specified data once and then
	// not reuse it thereafter.
	NoReuse
)

// MetadataHashValue represents the record "wasi:filesystem/types.metadata-hash-value".
//
// A 128-bit hash value, split into parts because wasm doesn't have a
// 128-bit integer type.
type MetadataHashValue struct {
	// 64 bits of a 128-bit hash value.
	Lower uint64

	// Another 64 bits of a 128-bit hash value.
	Upper uint64
}

/*
package wasi:filesystem@0.2.0-rc-2023-11-10;
interface types {
    // A descriptor is a reference to a filesystem object, which may be a file,
    // directory, named pipe, special file, or other object on which filesystem
    // calls may be made.
    resource descriptor {
        // Return a stream for reading from a file, if available.
        //
        // May fail with an error-code describing why the file cannot be read.
        //
        // Multiple read, write, and append streams may be active on the same open
        // file and they do not interfere with each other.
        //
        // Note: This allows using `read-stream`, which is similar to `read` in POSIX.
        read-via-stream: func(
            // The offset within the file at which to start reading.
            offset: filesize,
        ) -> result<input-stream, error-code>;

        // Return a stream for writing to a file, if available.
        //
        // May fail with an error-code describing why the file cannot be written.
        //
        // Note: This allows using `write-stream`, which is similar to `write` in
        // POSIX.
        write-via-stream: func(
            // The offset within the file at which to start writing.
            offset: filesize,
        ) -> result<output-stream, error-code>;

        // Return a stream for appending to a file, if available.
        //
        // May fail with an error-code describing why the file cannot be appended.
        //
        // Note: This allows using `write-stream`, which is similar to `write` with
        // `O_APPEND` in in POSIX.
        append-via-stream: func() -> result<output-stream, error-code>;

        // Provide file advisory information on a descriptor.
        //
        // This is similar to `posix_fadvise` in POSIX.
        advise: func(
            // The offset within the file to which the advisory applies.
            offset: filesize,
            // The length of the region to which the advisory applies.
            length: filesize,
            // The advice.
            advice: advice
        ) -> result<_, error-code>;

        // Synchronize the data of a file to disk.
        //
        // This function succeeds with no effect if the file descriptor is not
        // opened for writing.
        //
        // Note: This is similar to `fdatasync` in POSIX.
        sync-data: func() -> result<_, error-code>;

        // Get flags associated with a descriptor.
        //
        // Note: This returns similar flags to `fcntl(fd, F_GETFL)` in POSIX.
        //
        // Note: This returns the value that was the `fs_flags` value returned
        // from `fdstat_get` in earlier versions of WASI.
        get-flags: func() -> result<descriptor-flags, error-code>;

        // Get the dynamic type of a descriptor.
        //
        // Note: This returns the same value as the `type` field of the `fd-stat`
        // returned by `stat`, `stat-at` and similar.
        //
        // Note: This returns similar flags to the `st_mode & S_IFMT` value provided
        // by `fstat` in POSIX.
        //
        // Note: This returns the value that was the `fs_filetype` value returned
        // from `fdstat_get` in earlier versions of WASI.
        get-type: func() -> result<descriptor-type, error-code>;

        // Adjust the size of an open file. If this increases the file's size, the
        // extra bytes are filled with zeros.
        //
        // Note: This was called `fd_filestat_set_size` in earlier versions of WASI.
        set-size: func(size: filesize) -> result<_, error-code>;

        // Adjust the timestamps of an open file or directory.
        //
        // Note: This is similar to `futimens` in POSIX.
        //
        // Note: This was called `fd_filestat_set_times` in earlier versions of WASI.
        set-times: func(
            // The desired values of the data access timestamp.
            data-access-timestamp: new-timestamp,
            // The desired values of the data modification timestamp.
            data-modification-timestamp: new-timestamp,
        ) -> result<_, error-code>;

        // Read from a descriptor, without using and updating the descriptor's offset.
        //
        // This function returns a list of bytes containing the data that was
        // read, along with a bool which, when true, indicates that the end of the
        // file was reached. The returned list will contain up to `length` bytes; it
        // may return fewer than requested, if the end of the file is reached or
        // if the I/O operation is interrupted.
        //
        // In the future, this may change to return a `stream<u8, error-code>`.
        //
        // Note: This is similar to `pread` in POSIX.
        read: func(
            // The maximum number of bytes to read.
            length: filesize,
            // The offset within the file at which to read.
            offset: filesize,
        ) -> result<tuple<list<u8>, bool>, error-code>;

        // Write to a descriptor, without using and updating the descriptor's offset.
        //
        // It is valid to write past the end of a file; the file is extended to the
        // extent of the write, with bytes between the previous end and the start of
        // the write set to zero.
        //
        // In the future, this may change to take a `stream<u8, error-code>`.
        //
        // Note: This is similar to `pwrite` in POSIX.
        write: func(
            // Data to write
            buffer: list<u8>,
            // The offset within the file at which to write.
            offset: filesize,
        ) -> result<filesize, error-code>;

        // Read directory entries from a directory.
        //
        // On filesystems where directories contain entries referring to themselves
        // and their parents, often named `.` and `..` respectively, these entries
        // are omitted.
        //
        // This always returns a new stream which starts at the beginning of the
        // directory. Multiple streams may be active on the same directory, and they
        // do not interfere with each other.
        read-directory: func() -> result<directory-entry-stream, error-code>;

        // Synchronize the data and metadata of a file to disk.
        //
        // This function succeeds with no effect if the file descriptor is not
        // opened for writing.
        //
        // Note: This is similar to `fsync` in POSIX.
        sync: func() -> result<_, error-code>;

        // Create a directory.
        //
        // Note: This is similar to `mkdirat` in POSIX.
        create-directory-at: func(
            // The relative path at which to create the directory.
            path: string,
        ) -> result<_, error-code>;

        // Return the attributes of an open file or directory.
        //
        // Note: This is similar to `fstat` in POSIX, except that it does not return
        // device and inode information. For testing whether two descriptors refer to
        // the same underlying filesystem object, use `is-same-object`. To obtain
        // additional data that can be used do determine whether a file has been
        // modified, use `metadata-hash`.
        //
        // Note: This was called `fd_filestat_get` in earlier versions of WASI.
        stat: func() -> result<descriptor-stat, error-code>;

        // Return the attributes of a file or directory.
        //
        // Note: This is similar to `fstatat` in POSIX, except that it does not
        // return device and inode information. See the `stat` description for a
        // discussion of alternatives.
        //
        // Note: This was called `path_filestat_get` in earlier versions of WASI.
        stat-at: func(
            // Flags determining the method of how the path is resolved.
            path-flags: path-flags,
            // The relative path of the file or directory to inspect.
            path: string,
        ) -> result<descriptor-stat, error-code>;

        // Adjust the timestamps of a file or directory.
        //
        // Note: This is similar to `utimensat` in POSIX.
        //
        // Note: This was called `path_filestat_set_times` in earlier versions of
        // WASI.
        set-times-at: func(
            // Flags determining the method of how the path is resolved.
            path-flags: path-flags,
            // The relative path of the file or directory to operate on.
            path: string,
            // The desired values of the data access timestamp.
            data-access-timestamp: new-timestamp,
            // The desired values of the data modification timestamp.
            data-modification-timestamp: new-timestamp,
        ) -> result<_, error-code>;

        // Create a hard link.
        //
        // Note: This is similar to `linkat` in POSIX.
        link-at: func(
            // Flags determining the method of how the path is resolved.
            old-path-flags: path-flags,
            // The relative source path from which to link.
            old-path: string,
            // The base directory for `new-path`.
            new-descriptor: borrow<descriptor>,
            // The relative destination path at which to create the hard link.
            new-path: string,
        ) -> result<_, error-code>;

        // Open a file or directory.
        //
        // The returned descriptor is not guaranteed to be the lowest-numbered
        // descriptor not currently open/ it is randomized to prevent applications
        // from depending on making assumptions about indexes, since this is
        // error-prone in multi-threaded contexts. The returned descriptor is
        // guaranteed to be less than 2**31.
        //
        // If `flags` contains `descriptor-flags::mutate-directory`, and the base
        // descriptor doesn't have `descriptor-flags::mutate-directory` set,
        // `open-at` fails with `error-code::read-only`.
        //
        // If `flags` contains `write` or `mutate-directory`, or `open-flags`
        // contains `truncate` or `create`, and the base descriptor doesn't have
        // `descriptor-flags::mutate-directory` set, `open-at` fails with
        // `error-code::read-only`.
        //
        // Note: This is similar to `openat` in POSIX.
        open-at: func(
            // Flags determining the method of how the path is resolved.
            path-flags: path-flags,
            // The relative path of the object to open.
            path: string,
            // The method by which to open the file.
            open-flags: open-flags,
            // Flags to use for the resulting descriptor.
            %flags: descriptor-flags,
        ) -> result<descriptor, error-code>;

        // Read the contents of a symbolic link.
        //
        // If the contents contain an absolute or rooted path in the underlying
        // filesystem, this function fails with `error-code::not-permitted`.
        //
        // Note: This is similar to `readlinkat` in POSIX.
        readlink-at: func(
            // The relative path of the symbolic link from which to read.
            path: string,
        ) -> result<string, error-code>;

        // Remove a directory.
        //
        // Return `error-code::not-empty` if the directory is not empty.
        //
        // Note: This is similar to `unlinkat(fd, path, AT_REMOVEDIR)` in POSIX.
        remove-directory-at: func(
            // The relative path to a directory to remove.
            path: string,
        ) -> result<_, error-code>;

        // Rename a filesystem object.
        //
        // Note: This is similar to `renameat` in POSIX.
        rename-at: func(
            // The relative source path of the file or directory to rename.
            old-path: string,
            // The base directory for `new-path`.
            new-descriptor: borrow<descriptor>,
            // The relative destination path to which to rename the file or directory.
            new-path: string,
        ) -> result<_, error-code>;

        // Create a symbolic link (also known as a "symlink").
        //
        // If `old-path` starts with `/`, the function fails with
        // `error-code::not-permitted`.
        //
        // Note: This is similar to `symlinkat` in POSIX.
        symlink-at: func(
            // The contents of the symbolic link.
            old-path: string,
            // The relative destination path at which to create the symbolic link.
            new-path: string,
        ) -> result<_, error-code>;

        // Unlink a filesystem object that is not a directory.
        //
        // Return `error-code::is-directory` if the path refers to a directory.
        // Note: This is similar to `unlinkat(fd, path, 0)` in POSIX.
        unlink-file-at: func(
            // The relative path to a file to unlink.
            path: string,
        ) -> result<_, error-code>;

        // Test whether two descriptors refer to the same filesystem object.
        //
        // In POSIX, this corresponds to testing whether the two descriptors have the
        // same device (`st_dev`) and inode (`st_ino` or `d_ino`) numbers.
        // wasi-filesystem does not expose device and inode numbers, so this function
        // may be used instead.
        is-same-object: func(other: borrow<descriptor>) -> bool;

        // Return a hash of the metadata associated with a filesystem object referred
        // to by a descriptor.
        //
        // This returns a hash of the last-modification timestamp and file size, and
        // may also include the inode number, device number, birth timestamp, and
        // other metadata fields that may change when the file is modified or
        // replaced. It may also include a secret value chosen by the
        // implementation and not otherwise exposed.
        //
        // Implementations are encourated to provide the following properties:
        //
        //  - If the file is not modified or replaced, the computed hash value should
        //    usually not change.
        //  - If the object is modified or replaced, the computed hash value should
        //    usually change.
        //  - The inputs to the hash should not be easily computable from the
        //    computed hash.
        //
        // However, none of these is required.
        metadata-hash: func() -> result<metadata-hash-value, error-code>;

        // Return a hash of the metadata associated with a filesystem object referred
        // to by a directory descriptor and a relative path.
        //
        // This performs the same hash computation as `metadata-hash`.
        metadata-hash-at: func(
            // Flags determining the method of how the path is resolved.
            path-flags: path-flags,
            // The relative path of the file or directory to inspect.
            path: string,
        ) -> result<metadata-hash-value, error-code>;
    }

    // A stream of directory entries.
    resource directory-entry-stream {
        // Read a single directory entry from a `directory-entry-stream`.
        read-directory-entry: func() -> result<option<directory-entry>, error-code>;
    }

    // Attempts to extract a filesystem-related `error-code` from the stream
    // `error` provided.
    //
    // Stream operations which return `stream-error::last-operation-failed`
    // have a payload with more information about the operation that failed.
    // This payload can be passed through to this function to see if there's
    // filesystem-related information about the error to return.
    //
    // Note that this function is fallible because not all stream-related
    // errors are filesystem-related errors.
    filesystem-error-code: func(err: borrow<error>) -> option<error-code>;
}

*/
