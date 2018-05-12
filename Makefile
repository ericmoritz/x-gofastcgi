# Copyright 2009 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.$(GOARCH)

TARG=gofastcgi
GOFILES=\
	fcgi_consts.go\
	byte_serialize.go\
	iorecord.go\
	record.go\
	timer.go\
	request.go


include $(GOROOT)/src/Make.pkg
