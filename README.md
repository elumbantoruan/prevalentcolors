# prevalent colors

The system calculates 3 most prevalent colors in the RGB scheme in hexadecimal format and write the result into a CSV file  and  
must be able to scale with limited resources given of an input file.

## Prevalent colors computation

Utilize open source from [EdlinOrg/prominentcolor](https://github.com/EdlinOrg/prominentcolor) to find the K most dominant colors in an image

## Data reader

To handle huge load with a limited resources given of the input file, the system utilizes [Memory-mapped file(Mmap)](https://en.wikipedia.org/wiki/Memory-mapped_file) so applications can treat the mapped portion as it if were primary memory.

## Project structure

The project is structured as a tree below

``` go
├── README.md
├── cmd
│   └── main.go
├── data
│   ├── input
│   │   └── input.txt
│   └── output
├── go.mod
├── go.sum
└── pkg
    ├── datareader
    │   └── data_reader.go
    ├── datawriter
    │   └── data_writer.go
    └── imageprocessor
        └── image_processor.go
```

### cmd

cmd contains package main, which is an entry point of the application.  
It takes few arguments such as path for input file, output file, and error log.  
It uses other packages as reference such as `pkg/datareader`, `pkg/datawriter`, and `pkg/imageprocessor`  

### pkg/datareader  

pkg/datareader contains `MMapReader` type which utilizes `syscall.Mmap` to allow memory mapping file.  It is a concrete type of a `DataReader` interface with the following contracts  

```go
type DataReader interface {
    Read() bool
    Close()
    Data() []byte
    Err() error
}
```

Invoke, `NewMMapReader` function to create an instance of `MMapReader`  

``` go
func NewMMapReader(filename string, chunkSize ...int) (DataReader, error)
```

In order to maximize the throughput, given of a resource, one must specify the chunkSize (option), which later on will be used as an `offset` and `size` parameters of an `syscall.Mmap`.  
If chunkSize is not being provided, then it take a default pagesize

```go
func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error)
```

The `offset` will start from zero and it will be incremented each time the `Read` function is called, while the `size` will be constant until the very end.  
The `chunkSize` has to be the multiple of the page size.  In most system, the page size will be 4KB.  
During the chunk process, the number of bytes that are returned may not have a fully portion of segment of data.  Since each line represents a URL follow by the line feed, it's possible that the whole URL is not completed.  To verify whether the line is complete, it will check if it contains the linefeed `(\n)`
So, the `backBuffer` is used to capture the incomplete URL, and the data that will be returned will be truncated at the end.  The `backBuffer` will be prepended at the next `Read` operation

### pkg/datawriter

pkg/datawriter appends bytes for a given output file.  

### pkg/imageprocessor

pkg/imageprocessor utilizes open source library from [EdlinOrg/prominentcolor](https://github.com/EdlinOrg/prominentcolor) to find the K most dominant colors in an image
