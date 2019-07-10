package async

import (
	"bufio"
	"io"
)

func CallbackReader(reader io.Reader, consumer func(*string, error) bool) {

	r := bufio.NewReader(reader)

	go func() {
		for {
			result, _, err := r.ReadLine()
			var rc bool
			if err != nil {
				rc = consumer(nil, err)
			} else {
				s := string(result)
				rc = consumer(&s, nil)
			}
			if !rc {
				return
			}
		}
	}()
}

func ChannelReader(reader io.Reader, data chan string) {

	CallbackReader(reader, func(s *string, e error) bool {
		if e != nil {
			close(data)
			return false
		} else {
			data <- *s
			return true
		}
	})

}
