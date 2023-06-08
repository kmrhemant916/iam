package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Connection(username string, password string, host string, port string) (*amqp.Connection, error) {
    connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/",
        username,
        password,
        host,
        port,
    )
    conn, err := amqp.Dial(connectionString)
    if err != nil {
        return nil, err
    }
    return conn, nil
}
