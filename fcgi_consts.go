package gofastcgi

const (
	// Set various FastCGI constants
	// Maximum number of requests that can be handled
	FCGI_HDR_LEN=8;
	FCGI_MAX_REQS=1;
	FCGI_MAX_CONNS = 1;

	// Supported version of the FastCGI protocol
	FCGI_VERSION_1 = 1;

	// Boolean: can this application multiplex connections?
	FCGI_MPXS_CONNS=0;

	// Record types
	FCGI_BEGIN_REQUEST = 1 ; FCGI_ABORT_REQUEST = 2 ; FCGI_END_REQUEST   = 3;
	FCGI_PARAMS        = 4 ; FCGI_STDIN         = 5 ; FCGI_STDOUT        = 6;
	FCGI_STDERR        = 7 ; FCGI_DATA          = 8 ; FCGI_GET_VALUES    = 9;
	FCGI_GET_VALUES_RESULT = 10;
	FCGI_UNKNOWN_TYPE = 11;
	FCGI_MAXTYPE = FCGI_UNKNOWN_TYPE;

	// Types of management records
	//ManagementTypes = [FCGI_GET_VALUES];

	FCGI_NULL_REQUEST_ID=0;

	// Masks for flags component of FCGI_BEGIN_REQUEST
	FCGI_KEEP_CONN = 1;

	// Values for role component of FCGI_BEGIN_REQUEST
	FCGI_RESPONDER = 1 ; FCGI_AUTHORIZER = 2 ; FCGI_FILTER = 3;

	// Values for protocolStatus component of FCGI_END_REQUEST
	FCGI_REQUEST_COMPLETE = 0;  // Request completed nicely
	FCGI_CANT_MPX_CONN    = 1;  // This app can't multiplex
    FCGI_OVERLOADED       = 2;  // New request rejected; too busy
	FCGI_UNKNOWN_ROLE     = 3;  // Role value not known
)