package client

// package client renders a interface to process file.

// File operations are divided into two families.
// The first family is stream operation, include GetReader and GetWriter.
// The other one is normal operation, include Delete, Copy, Exists, etc.

// Each operation needs a timeout parameter, which type is time.Duration.
// If timeout is zero or negative, the operation do not start timeout mechanism.
// If timeout is positive, when timeout occurs, the operation returns
// immediately, regardless of whether it is completed.

// Example for write and read a file.
//
// flag.Parse()
// client.Initialize()
// writer, err := client.GetWriter(domain, size, fn, biz, user, timeout)
// if err != nil {
//   return nil
// }
// defer writer.Close()

//// do write like io.WriteCloser.

//// get file id

// if dwriter, ok := *DFSWriter.(writer); ok {
//   fid, err := dwriter.GetFileInfoAndClose()
// }

//// read a file
// reader, err := client.GetReader(fid, domain, timeout)
// if err != nil {
//   ...
// }
// defer reader.Close()
//
//// do reader like io.ReadCloser

//// deleta a file
// err := client.Delete(fid, domain, timeout)
