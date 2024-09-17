package main

import (
	"fmt"

	"github.com/Arcadian-Sky/datakkeeper/internal/app/client"
	"github.com/Arcadian-Sky/datakkeeper/internal/settings"
)

func main() {
	clientConfig := settings.GetClientConfig()
	fmt.Printf("clientConfig: %v\n", clientConfig)
	app := client.NewClientApp(&clientConfig)
	defer app.Conn.Close()
}
