package vpn

// Connection is an interface that represent an object that can connect to a Network
type Connection interface {
	//EnsureConnected checks if you are logged to the connection.
	//
	//If the connection can not be established, should return an error.
	//
	//The connection process must be executed synchronously, if this function returns no error, that's mean that the
	//connection is established.
	EnsureConnected() error

	//Name of the connection for log purpose
	Name() string
}
