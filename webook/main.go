package wire

func main() {
	server := InitWebServer()

	server.Run("127.0.0.1:8080")

}
