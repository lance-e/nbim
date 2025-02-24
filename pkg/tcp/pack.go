package tcp

import netreactors "github.com/lance-e/net-reactors"

// tcp
func Unpacking(buf *netreactors.Buffer) []byte {

}

// tcp
func Packing(data []byte) *netreactors.Buffer {
	buf := netreactors.NewBuffer()
	len := data[:4]
}
