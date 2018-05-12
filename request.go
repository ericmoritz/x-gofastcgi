package gofastcgi


import (
	"io";
	"os";
	"bytes";
//	"log";
)

type Request struct {
	RequestID uint16;

	// Bodies from the records
	BeginRequestBody;
	UnknownTypeBody;
	EndRequestBody;

	// The good stuff
	Env map[string]string;
	Stdin  *bytes.Buffer;
	Stdout FlushWriter;
	Stderr FlushWriter;
};

func (req *Request) Read(s io.ReadWriter) (err os.Error) {
	remaining := 1;
	inbytes := make([]byte, 8192); // 8192 is what libfcgi uses
	//outbytes := make([]byte, 8192);
    //errbytes := make([]byte, 512); // 512 is what libfcgi uses
	req.RequestID = 0;

	req.Stdin = bytes.NewBuffer(inbytes);

	timer := new(TimeIt);

	for remaining != 0 {
		var r Record;
		err = r.Read(s);

		if err != nil {
			return err;
		}
		
		
		switch {
		// Start the request
		case r.Header.Type == FCGI_BEGIN_REQUEST && req.RequestID == 0:
			req.RequestID = r.Header.RequestID;
			req.BeginRequestBody.Read(bytes.NewBuffer(r.Content));

			switch req.BeginRequestBody.Role {
			case FCGI_AUTHORIZER:
				remaining=1;
			case FCGI_RESPONDER:
				remaining=2;
			case FCGI_FILTER:
				remaining=3;
			}
			
		case r.Header.Type == FCGI_PARAMS:
			req.Env = r.Content.GetParams();
			if len(r.Content) == 0 {
				remaining--;
			}

		case r.Header.Type == FCGI_STDIN:
			if len(r.Content) > 0 {
				req.Stdin.Write(r.Content);
			} else {
				remaining--;
			}
		case r.Header.Type == FCGI_DATA:
			timer.Start();			
			if len(r.Content) == 0 {
				remaining--;
			}
		}
	}

	req.Stdout = NewWriter(req.RequestID, FCGI_STDOUT, s, true); //bytes.NewBuffer(outbytes);
	req.Stderr = NewWriter(req.RequestID, FCGI_STDERR, s, true); //bytes.NewBuffer(outbytes);		
	return nil
}


func (req *Request) Write(w io.ReadWriter, status uint32) (os.Error) {
	var r *Record;

	// Flush out any left over data
	req.Stdout.Flush();
	req.Stderr.Flush();

	// Set the end of stream record
	r = new(Record);
	r.Header.RequestID=req.RequestID;
	r.Header.Type=FCGI_STDOUT;
	r.Content = make([]byte,0);
	r.Header.ContentLength = 0;
	r.Write(w);

	// Set the end of stream record
	r = new(Record);
	r.Header.RequestID=req.RequestID;
	r.Header.Type=FCGI_STDERR;
	r.Content = make([]byte,0);
	r.Header.ContentLength = 0;
	r.Write(w);

	r = new(Record);
	r.Header.Type = FCGI_END_REQUEST;
	r.Header.RequestID = req.RequestID;
	
	body := new(EndRequestBody);
	body.Status = status;
	body.ProtocolStatus = FCGI_REQUEST_COMPLETE;

	r.Content = body.Bytes();
	r.Header.ContentLength = uint16(len(r.Content));

	if err := r.Write(w); err != nil { return err; }

	return nil;
}