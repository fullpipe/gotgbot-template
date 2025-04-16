package tools

import (
	"errors"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

func BestPhoto(ps []gotgbot.PhotoSize) (gotgbot.PhotoSize, error) {
	if len(ps) == 0 {
		return gotgbot.PhotoSize{}, errors.New("no photos")
	}

	size := ps[0].FileSize
	best := ps[0]

	for _, p := range ps {
		if p.FileSize > size {
			size = p.FileSize
			best = p
		}
	}

	return best, nil
}
