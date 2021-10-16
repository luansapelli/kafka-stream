package notification

import "io"

type ServiceNotification interface {
	io.Writer
}
