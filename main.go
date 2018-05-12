
package main

import (
	"fmt";
	"gofastcgi";
	"net";
	"os";
	"strings";
//	"log";
	"io";
	"flag";
	"runtime";
)

type EchoReader struct {
	reader io.ReadWriter;
}

func (r *EchoReader) Read(p []byte) (n int, err os.Error) {
	n, err = r.reader.Read(p);

	// Spit out data into stdout
	//os.Stdout.Write(p);
	return n, err;
}
func (r *EchoReader) Write(p []byte) (n int, err os.Error) {
	return r.reader.Write(p);
}

type Job func(request *gofastcgi.Request);

func Serve(c net.Conn, job Job) {
	var req gofastcgi.Request;

	r := &EchoReader{c};
	req.Read(r);
	job(&req);
	req.Write(r, 0);
	c.Close();
}

func ListenAndServe(addr string, job Job) os.Error {

	l, err := net.Listen("tcp", ":3031");
	if err != nil {
		return err
	}
	for {
		
        rw, e := l.Accept();

		go Serve(rw, job);

        if e != nil {
            return e
			}
        if err != nil {
            continue;
		}

    }
	panic("WTF?");
}


func Index(req *gofastcgi.Request) {
	req.Stdout.Write(strings.Bytes("Content-Type: text/plain\n\n"));
	req.Stdout.Write(strings.Bytes("\n\n"));
	req.Stdout.Write(strings.Bytes("Hello, World!\n"));
}


var procs = flag.Int("procs", 1, "number of CPUs") // Q=17, R=18
func main() {
	flag.Parse();

	runtime.GOMAXPROCS(*procs);

	err := ListenAndServe(":3031", Index);
	if err != nil {
		fmt.Printf("Fuck: %v", err);
	}
}