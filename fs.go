/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package fs

import (
	"encoding/json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io"
	"os"
)

/**
FileService interface is used to inject this module to applications
 */
type FileService interface {
	JsonFileService
	ProtoFileService
	CsvFileService

	/*
	Gets current buffer size, default value is 64k
	 */
	BufferSize() int

	/*
	Sets current buffer size, that would be used on each file opening or creation. Particularly useful for gzip files.
	 */
	SetBufferSize(rwBufSize int)

	/*
	Gets JSON marshal options
	 */
	MarshalOptions() protojson.MarshalOptions

	/*
	Sets JSON marshal options
	*/
	SetMarshalOptions(protojson.MarshalOptions)

	/*
	Gets JSON unmarshal options
	*/
	UnmarshalOptions() protojson.UnmarshalOptions

	/*
	Sets JSON unmarshal options
	*/
	SetUnmarshalOptions(protojson.UnmarshalOptions)
}

/**
Base interface for JSON files r/w operations
*/
type JsonFileService interface {

	/*
	Creates new JSON stream.
	 */
	NewJsonStream(fd io.Writer, withGzip bool) JsonWriter

	/*
	Creates new JSON file in local file system. If file path ends with `.gz` extension it would be compressed.
	*/
	NewJsonFile(filePath string) (JsonWriter, error)

	/*
	Opens JSON stream from reader.
	 */
	JsonStream(fr io.Reader, withGzip bool) (JsonReader, error)

	/*
	Opens JSON file from local file system. If file path ends with `.gz` extension it would be decompressed.
	*/
	OpenJsonFile(filePath string) (JsonReader, error)

	/*
	Opens JSON file from descriptor.
	 */
	JsonFile(fd *os.File) (JsonReader, error)

	/*
	Splits one single JSON file in to parts. Partition function would be called to format file name for each part.
	 */
	SplitJsonFile(inputFilePath string, limit int, partitionFn func (int) string) ([]string, error)

	/*
	Joins JSON files in to one.
	 */
	JoinJsonFiles(outputFilePath string, parts []string) error
}

/**
Base interface to write content in to JSON file.
 */
type JsonWriter interface {

	/*
	Writes already formatted JSON message
	 */
	WriteRaw(message json.RawMessage) error

	/*
	Writes golang object that supports serialization to JSON format.
	 */
    Write(object interface{}) error

    /*
    Closes stream and flashes underline buffers
     */
	Close() error

}

/**
Base interface to read content from JSON file.
 */
type JsonReader interface {

	/*
	Reads single row from JSON file, assuming that lines are separated by `\n` character.
	 */
	ReadRaw() (json.RawMessage, error)

	/*
	Reads single raw from JSON file in to golang object. Golang object must support JSON serialization.
	 */
	Read(holder interface{}) error

	/*
	Closes stream and underline buffers.
	 */
	Close() error

}

/**
Base interface for protobuf files r/w operations.
This file serialization implementation uses BigEndian 32 unsigned integer as a header for each serialized protobuf object equal to the size of it.
 */
type ProtoFileService interface {

	/*
	Opens protofile stream.
	 */
	ProtoStream(r io.Reader, withGzip bool) (ProtoReader, error)

	/*
	Opens protofile stream from load file system. If file path ends with `.gz` extension it would be decompressed.
	*/
	OpenProtoFile(filePath string) (ProtoReader, error)

	/*
	Opens protofile stream from file object
	 */
	ProtoFile(fd *os.File) (ProtoReader, error)

	/*
	Creates new protofile stream.
	 */
	NewProtoStream(fd io.Writer, withGzip bool) ProtoWriter

	/*
	Creates new protofile stream. If file path ends with `.gz` extension it would be compressed.
	*/
	NewProtoBuf(gzipEnabled bool) (ProtoWriter, error)

	/*
	Creates new protofile stream in local file system. If file path ends with `.gz` extension it would be compressed.
	 */
	NewProtoFile(filePath string) (ProtoWriter, error)

	/*
	Splits one single protofile in to parts. Partition function would be called to format file name for each part.
	*/
	SplitProtoFile(inputFilePath string, holder proto.Message, limit int, partFn func (int) string) ([]string, error)

	/*
	Joins protofiles in to one.
	*/
	JoinProtoFiles(outputFilePath string, row proto.Message, parts []string) error

}

/**
Base interface to write content in to proto file.
*/
type ProtoWriter interface {

	/**
	Writes message to the stream
	 */
	Write(message proto.Message) ([]byte, error)

	/*
	Closes stream and flashes underline buffers
	*/
	Close() error

}

/**
Base interface to read content from JSON file.
*/
type ProtoReader interface {

	/*
	Reads size header and single protobuf object.
	*/
	ReadTo(message proto.Message) error

	/*
	Closes stream and underline buffers.
	*/
	Close() error

}

/**
Base interface for csv files r/w operations
*/
type CsvFileService interface {

	/*
	Creates new CSV file stream.
	*/
	NewCsvStream(fw io.Writer, withGzip bool, valueProcessors ...CsvValueProcessor) CsvWriter

	/*
	Creates new CSV file stream in local file system. If file path ends with `.gz` extension it would be compressed.
	*/
	NewCsvFile(filePath string, valueProcessors ...CsvValueProcessor) (CsvWriter, error)

	/*
	Opens CSV file stream.
	*/
	OpenCsvStream(fr io.Reader, withGzip bool, valueProcessors ...CsvValueProcessor) (CsvStream, error)

	/*
	Opens CSV file stream from load file system. If file path ends with `.gz` extension it would be decompressed.
	*/
	OpenCsvFile(filePath string, valueProcessors ...CsvValueProcessor) (CsvReader, error)

	/*
	Opens CSV file stream from file object.
	*/
	CsvFileReader(fd *os.File, valueProcessors ...CsvValueProcessor) (CsvReader, error)

	/*
	Creates CSV file scheme from the header.
	*/
	NewCsvSchema(header []string) CsvSchema

	/*
	Splits one single CSV in to parts. Partition function would be called to format file name for each part.
	*/
	SplitCsvFile(inputFilePath string, limit int, partFn func (int) string) ([]string, error)

	/*
	Joins CSV files in to one.
	*/
	JoinCsvFiles(outputFilePath string, parts []string) error

}

/**
Base interface processor that pre-process value on reading or writing, keeping certain compatibility of CSV file with other systems.
 */
type CsvValueProcessor  func(string) string

/**
Base interface to write content in to CSV file.
*/
type CsvWriter interface {

	/**
	Writes values to the stream
	*/
	Write(values ...string) error

	/*
	Closes stream and flashes underline buffers
	*/
	Close() error

}

/**
Base interface to read content from CSV stream.
*/
type CsvStream interface {

	/*
	Reads single row from CSV file, assuming that lines are separated by `\n` character.
	*/
	Read() ([]string, error)

	/*
	Closes stream and underline buffers.
	*/
	Close() error
}

/**
Base interface to read content from CSV file.
*/
type CsvReader interface {

	/*
	Reads first row from CSV file, assuming that lines are separated by `\n` character. Uses first row as a header.
	*/
	ReadHeader() (CsvFile, error)

	/*
	Reads single row from CSV file, assuming that lines are separated by `\n` character.
	*/
	Read() ([]string, error)

	/*
	Closes stream and underline buffers.
	*/
	Close() error
}

/**
Base interface of CSV file scheme.
 */
type CsvSchema interface {

	/*
	Wraps record in to CsvRecord interface.
	 */
	Record(record []string) CsvRecord

}

/**
Base interface of CSV file record that is using CSV file scheme.
*/
type CsvRecord interface {

	/**
	Gets the list of values of the record.
	 */
	Record() []string

	/*
	Gets particular column from the record. If column not found then returns def value.
	 */
	Field(name string, def string) string

	/**
	Gets all columns with values.
	 */
	Fields() map[string]string

}

/**
Base interface for CSV file.
 */
type CsvFile interface {

	/*
	Get CSV file header.
	 */
	Header() []string

	/*
	Gets index of columns.
	 */
	Index() map[string]int

	/*
	Reads next record. Return EOF error if no more records in file.
	 */
	Next() (CsvRecord, error)

}