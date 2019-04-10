package main

var EXIT_CODE = map[int]string{
	0:  "OK",
	1:  "Failed to initialize.",
	6:  "Couldn't resolve host. The given remote host was not resolved.",
	7:  "Failed to connect to host.",
	8:  "Weird server reply. The server sent data curl couldn't parse.",
	28: "Operation timeout. The specified time-out period was reached according to the conditions.",
	35: "SSL connect error. The SSL handshaking failed.",
}
