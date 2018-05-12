package gofastcgi


import (
	"io";
	"os";
//	"log";
)

/////
// Header struct
/////
type Header struct {
	Version byte;
	Type byte;
	RequestID uint16;
	ContentLength uint16;
	PaddingLength byte;
	Reserved byte;
}

func (header *Header) Read(r io.Reader) os.Error {
	tmp := make([]byte, 2);
	data := make([]byte, 8);

	_, err := io.ReadFull(r, data);
	if err != nil { return err };

	header.Version = data[0];
	header.Type = data[1];

	tmp[0] = data[2];
	tmp[1] = data[3];	
	
	header.RequestID = Uint16_BE(tmp);

	tmp[0] = data[4];
	tmp[1] = data[5];	
	
	header.ContentLength = Uint16_BE(tmp);
	header.PaddingLength = data[6];
	header.Reserved = data[7];

	return nil;
}

func (header *Header) Bytes() []byte {
	var tmp []byte;
	data := make([]byte, 8);
	
	data[0] = header.Version;
	data[1] = header.Type;

	tmp = BE_uint16(header.RequestID);
	data[2] = tmp[0];
	data[3] = tmp[1];
	
	tmp = BE_uint16(header.ContentLength);
	data[4] = tmp[0];
	data[5] = tmp[1];

	data[6] = header.PaddingLength;
	data[7] = header.Reserved;
	
	return data;
}
func (header *Header) Write(w io.Writer) os.Error {
	data := header.Bytes();
	_, err := w.Write(data);
	return err;
}
////
// Content
////

type Content []byte;

func getParamLen(data []byte, pos uint32) (uint32, uint32) {
	tmp := data[pos];
	if tmp >> 7 == 1 {
		return Uint32_BE(data[pos:pos+4])& 0x7FFFFFFF, pos+4;
	} 
	return uint32(tmp), pos+1;
}

func (content Content) GetParams() map[string]string {
	params := make(map[string]string);
	if len(content) == 0 { return params }

	pos := uint32(0);
	for pos < uint32(len(content)) {

		var nameLen, valueLen uint32;
		var name, val []byte;

		nameLen, pos = getParamLen(content, pos);
		valueLen, pos = getParamLen(content, pos);

		name = content[pos:pos+nameLen];
		pos += nameLen;

		val = content[pos:pos+valueLen];
		pos += valueLen;
		params[string(name)] = string(val);

	}


	return params;
}

/////
// Record struct
/////
type Record struct {
	Header;
	Content Content;
	Padding []byte;
}


func (record *Record) Read(r io.Reader) os.Error {

	// Read in the header
//	log.Stderrf("Reading Header");

	if err := record.Header.Read(r); err != nil{
		return err;
	}
//	log.Stderrf("Got Header: %v", record.Header);

	/////
	// Read in the Content and Padding
	/////
	// Read the Content
//	log.Stderrf("Content: Reading %s bytes", record.Header.ContentLength);
	record.Content = make([]byte, record.Header.ContentLength);
	n, err := io.ReadFull(r, record.Content); 

	if err != nil || uint16(n) != record.ContentLength {
		return err;
	}
//	log.Stderrf("Content: Read %s bytes", n);

	// Read the Padding
	record.Padding = make([]byte, record.Header.PaddingLength);
	if n, err := io.ReadFull(r, record.Padding); err != nil || uint8(n) != record.PaddingLength {
		return err;
	}

	return nil;
}

func (record *Record) Write(w io.Writer) (os.Error) {
	var err os.Error;
	// Write out the Header
	if err = record.Header.Write(w); err != nil { return err; }

	// Write out the Content
	if _, err = w.Write(record.Content); err != nil { return err; }
	
	// Write out Padding
	if _, err = w.Write(record.Padding); err != nil { return err; }

	return nil;
}
////
// RecordTypes
////
type UnknownTypeBody struct {
	Type byte;
	Reserved [7]byte;
}

func (body *UnknownTypeBody) Read(r io.Reader) os.Error {
	data := make([]byte, 8);
	_, err := io.ReadFull(r,data);
	if err != nil { return err };

	body.Type = data[0];
	body.Reserved[0] = data[1];
	body.Reserved[1] = data[2];
	body.Reserved[2] = data[3];
	body.Reserved[3] = data[4];
	body.Reserved[4] = data[5];
	body.Reserved[5] = data[6];
	body.Reserved[6] = data[7];	
	return nil;
}

func (body *UnknownTypeBody) Bytes() []byte {
	data := make([]byte, 8);
	data[0] = body.Type;
	return data;
}
func (body *UnknownTypeBody) Write(w io.Writer) os.Error {
	_, err := w.Write(body.Bytes());
	return err;
}

type BeginRequestBody struct {
	Role uint16;
	Flags byte;
	Reserved [5]byte;
}

func (body *BeginRequestBody) Read(r io.Reader) os.Error {
	data := make([]byte, 8);
	tmp := make([]byte, 2);
	_, err := io.ReadFull(r, data);
	if err != nil { return err };

	tmp[0] = data[0];
	tmp[1] = data[1];	
	
	body.Role = Uint16_BE(tmp);
	body.Flags = data[2];

	body.Reserved[0] = data[3];
	body.Reserved[1] = data[4];
	body.Reserved[2] = data[5];
	body.Reserved[3] = data[6];
	body.Reserved[4] = data[7];

	return nil;
}

func (body *BeginRequestBody) Bytes() []byte {
	data := make([]byte, 8);
	
	tmp := BE_uint16(body.Role);
	
	data[0] = tmp[0];
	data[1] = tmp[1];
	data[2] = body.Flags;
	
	return data;
}

func (body *BeginRequestBody) Write(w io.Writer) os.Error {
	_, err := w.Write(body.Bytes());
	return err;

}

type EndRequestBody struct {
	Status uint32;
	ProtocolStatus byte;
	Reserved [3]byte;
}

func (body *EndRequestBody) Read(r io.Reader) os.Error {
	data := make([]byte, 8);
	tmp := make([]byte, 4);
	
	_, err := io.ReadFull(r, data);
	if err != nil { return err; }

	tmp[0] = data[0];
	tmp[1] = data[1];
	tmp[2] = data[2];
	tmp[3] = data[3];	
	
	body.Status = Uint32_BE(tmp);
	body.ProtocolStatus = data[4];

	body.Reserved[0] = data[5];
	body.Reserved[1] = data[6];
	body.Reserved[2] = data[7];

	return nil;
}
func (body *EndRequestBody) Bytes() []byte {
	d := make([]byte,8);
	tmp := BE_uint32(body.Status);

	d[0] = tmp[0];
	d[1] = tmp[1];
	d[2] = tmp[2];
	d[3] = tmp[3];
	d[4] = body.ProtocolStatus;
	return d
}

func (body *EndRequestBody) Write(w io.Writer) os.Error {
	_, err := w.Write(body.Bytes());
	return err;
}
