package gofastcgi

import (
	"bytes";
	"testing";
	"reflect";
//	"strings";
)

var hdr_data = []byte{
	1, // Version
	2, // Type
	0x06,
	0x01, // Request ID: 1537
	0x01,
	0x2C, // Content Length: 300 bytes
	10, // Padding Length: 10 bytes
	1, // Reserved
};

var hdr_expect = Header{
	1,
	2,
	1537,
	300,
	10,
	1
};

var rec_data = []byte {
	1, // Version
	2, // Type
	0x06,
	0x01, // Request ID: 1537
	0,
	4, // Content Length: 4 bytes
	1, // Padding Length: 1 bytes
	1, // Reserved
	'E','r','i','c', // Content
	' ' // Padding
};

var rec_expect = Record{
	Header{1,2,1537,4,1,1},
	[]byte{'E','r','i','c'},
	[]byte{' '}
};

func TestHeaderRead(t *testing.T) {
	var header Header;

	err := header.Read(bytes.NewBuffer(hdr_data));
	if err != nil {
		t.Errorf("Read header fail: %v", err);
	}
	if !reflect.DeepEqual(hdr_expect, header) {
		t.Errorf("Result != Expected: %v != %v", header, hdr_expect);
	}
}
func TestCorruptHeader(t *testing.T) {
	var header Header;
	bad_header := make([]byte, 5);
	
	err := header.Read(bytes.NewBuffer(bad_header));
	if err == nil {
		t.Errorf("Corrupt header didn't error.");
	}
}

func TestRecordRead(t *testing.T) {
	var record Record;

	err := record.Read(bytes.NewBuffer(rec_data));
	if err != nil {
		t.Errorf("Read record fail: %v", err);
	}
	if !reflect.DeepEqual(rec_expect, record) {
		t.Errorf("Result != Expected: %v != %v", record, rec_expect);
	}
}

func TestCorruptRecord(t *testing.T) {
	var record Record;
	bad_header := make([]byte, 5);
	
	err := record.Read(bytes.NewBuffer(bad_header));
	if err == nil {
		t.Errorf("Corrupt record didn't error.");
	}
}


func TestGetParamLen(t *testing.T) {
	var lenbytes = []byte{4}; // 4
	var expect = uint32(4);
	
	result, _ := getParamLen(lenbytes,0);
	expect = 4;
	
	if result != expect {
		t.Errorf("getParam reported the wrong param length: %v, expected %v", result, expect);		
	}
}

func TestGetParamLength32(t *testing.T) {
	var lenbytes = []byte{0x80, 0, 0, 4}; // FCGI's encoding for a 32bit length of the value 4
	var expect = uint32(4);
	
	result, _ := getParamLen(lenbytes,0);
	expect = 4;
	if result != expect {
		t.Errorf("getParam reported the wrong param length: %v, expected %v", result, expect);
	}
}

func TestGetParams8(t *testing.T) {
	expected := make(map[string]string);
	expected["Name"] = "Eric";

	var content = Content([]byte{4,4,'N','a','m','e','E','r','i','c'});
	result := content.GetParams();
	
	if result["Name"] != expected["Name"] {
		t.Errorf("TestGetParam8: %v != %v", result, expected);
	}
}


func TestGetParams32(t *testing.T) {
	expected := make(map[string]string);
	expected["Name"] = "Eric";

	var content = Content([]byte{0x80,0,0,4,0x80,0,0,4,'N','a','m','e','E','r','i','c'});
	result := content.GetParams();
	
	if result["Name"] != expected["Name"] {
		t.Errorf("TestGetParam32: %v != %v", result, expected);
	}
}
