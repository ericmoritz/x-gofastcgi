package main

import (
	"fmt";
	"bytes";
	"os";
	"io";
	"gofastcgi";
	"strings";
)

type RWByteBuffer struct {
	In *bytes.Buffer;
	Out *bytes.Buffer;
}

func (rw *RWByteBuffer)Read(p []byte) (n int, err os.Error) {
	return rw.Out.Read(p)
}

func (rw *RWByteBuffer)Write(p []byte) (n int, err os.Error) {
	return rw.In.Write(p)
}

func main() {
	filedata, err := io.ReadFile("./request.bin");
	if err  != nil { fmt.Print(err); return; }

	rw := &RWByteBuffer{
		bytes.NewBuffer(make([]byte, 8000)),
		bytes.NewBuffer(filedata)
		};
		
	req := new(gofastcgi.Request);
	
	req.Read(rw);
	req.Stdout.Write(strings.Bytes("Content-Type: text/plain\n\n"));
	req.Stdout.Write(strings.Bytes("Hello,World\n"));
	req.Write(rw, 0);

	//fmt.Print(string(rw.In.Bytes()));
}