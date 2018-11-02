package nstar

// #include "setns.c"
import "C"
import "os"

func main() {
	err := C.nsenter(C.int(os.Getpid()))
	if err != nil {
		panic(err.Error())
	}
}
