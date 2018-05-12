package gofastcgi

import (
	"io";
	"os";
//	"fmt";
	"bufio";
)

type FlushWriter interface {
	Write(p []byte) (n int, err os.Error);
	Flush() os.Error;
}

type FCGIStreamWriter struct {
	RequestID uint16;
	Type byte;
	io.Writer;
};

func (s *FCGIStreamWriter) Write(p []byte) (n int, err os.Error) {
	var size = 8192;
	var content []byte;

	if len(p) > size {
		// Send no more than 8192 bytes
		content = p[0:size];
	} else {
		content = p;
	}
	
	r := new(Record);
	r.Header.RequestID=s.RequestID;
	r.Header.Type=s.Type;
	r.Header.ContentLength = uint16(len(content));
	r.Header.PaddingLength = 0;
	r.Content = content;
	r.Padding = make([]byte,0);
	err = r.Write(s.Writer);
	return int(r.Header.ContentLength), err;
}

func (s *FCGIStreamWriter) Flush() os.Error {
	return nil;
}

func NewWriter(RequestID uint16, Type byte, w io.Writer, chunked bool) FlushWriter {
	fsw := &FCGIStreamWriter{RequestID, Type, w};
	tmp := FlushWriter(fsw);
	if chunked {
		tmp, _ = bufio.NewWriterSize(tmp, 8192);
	}
	return tmp;
}
