package retrier

import (
	"log"
	"time"
)

const (
	maxRetires = 3
)

type GopherSocialRetrier struct {
}

func NewGopherSocialRetreir() *GopherSocialRetrier {
	return &GopherSocialRetrier{}
}

func (r *GopherSocialRetrier) Retry(fn func() error, email string, isSandBox bool) (int, error) {
	for i := range maxRetires {
		if !isSandBox {

			err := fn()
			if err != nil {
				log.Printf("failed to send email to %+v, %d attempt of %d\n", email, i+1, maxRetires)
				log.Printf("Error: %+v\n", err.Error())
				//exponential back off
				time.Sleep(time.Second * time.Duration(i+1))
				continue
			}
			log.Printf("Email %+v sent with status code %d\n", email, 200)
			return 200, nil
		}

	}
	return -1, nil
}
