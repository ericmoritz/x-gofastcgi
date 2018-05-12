package gofastcgi

import (
	"testing";
//	"bytes";
	"strings";
	"reflect";
//	"io";
	"os";
	"fmt";
)

type StrictWriter struct {
	data []byte;
	pos int;
}

func (s *StrictWriter) Write(p []byte) (n int, err os.Error) {
	count := 0;

	if len(p) == 0 { return 0, nil }
	if s.pos >= len(s.data) { return 0, os.NewError("Strict Writer IndexError") }

	for _, v := range p {
		s.data[s.pos] = v;
		s.pos++;
		count++;
	}
	return count, nil;
}

func TestFCGIStreamWriter(t *testing.T) {
	expected := []byte{0,FCGI_STDOUT, 0, 1, 0, 14, 0, 0,
		'T','h','i','s',' ','i','s',' ','a',' ','t','e','s','t'};
	content := "This is a test";

	
	buffer := &StrictWriter{make([]byte, 22), 0};
	stdout := NewWriter(1, FCGI_STDOUT, buffer, false);

	stdout.Write(strings.Bytes(content));
	result := buffer.data;
	
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("%v != %v", expected, result);
	}
}

type MockWriter struct {
	i int
};

func (s *MockWriter) Write(p []byte) (n int, err os.Error) {
	fmt.Printf("%v %v: %v\n", len(p), s.i, string(p));
	s.i++;
	return len(p), nil;
}
