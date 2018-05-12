package gofastcgi


func BE_uint16(v uint16) (data []byte) {
	data = make([]byte, 2);
	data[0] = byte(v>>8&255);
	data[1] = byte(v&255);
	return;
}

func Uint16_BE(data []byte) (v uint16) {
	v = uint16(data[0])<<8;
	v += uint16(data[1]);
	return v
}

func BE_uint32(v uint32) (data []byte) {
	data = make([]byte, 4);
	data[0] = byte(v>>24&255);	
	data[1] = byte(v>>16&255);
	data[2] = byte(v>>8&255);
	data[3] = byte(v&255);
	return;
}


func Uint32_BE(data []byte) (v uint32) {
	v = uint32(data[0])<<24;
	v += uint32(data[1])<<16;
	v += uint32(data[2])<<8;
	v += uint32(data[3]);

	return v
}
