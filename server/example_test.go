package server_test

import (
	"fmt"

	"github.com/LuaSavage/mta-monitoring/server"
)

func ExampleServer_ReadRow() {
	response := []byte{
		'E', 'Y', 'E', '1',
		4, 'm', 't', 'a',
		6, '2', '2', '0', '0', '3',
		8, 'T', 'e', 's', 't', 'S', 'r', 'v',
		7, 'M', 'T', 'A', ':', 'S', 'A',
		5, 'N', 'o', 'n', 'e',
		4, '3', '.', '0',
		2, '0',
		2, '1',
		3, '8', '0',
		1,
	}

	srv := server.NewServer("127.0.0.1", 22003)
	if err := srv.ReadRow(&response); err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(srv.Name)
	fmt.Println(srv.Passworded)
	// Output:
	// TestSrv
	// false
}

func ExampleServer_GetJoinLink() {
	srv := server.NewServer("127.0.0.1", 22003)
	fmt.Println(srv.GetJoinLink())
	// Output:
	// mtasa://127.0.0.1:22003
}
