package session

import (
	"errors"
	"log"
	mathrand "math/rand"
)

type Session struct {
	SessionID  uint32
	Reactions  map[uint32]float64
	CurrentPID uint32
	leap       uint32
	start      uint32
	i          uint32
}

var Sessions = make(map[uint32]*Session)

func generateOddNumber() uint32 {
	num := mathrand.Uint32()
	if num%2 == 0 {
		num += 1
	}
	return num
}

func New() (*Session, error) {
	var value uint32

	for {
		value = mathrand.Uint32()

		if _, ok := Sessions[value]; !ok {
			break
		}
	}

	retSession := &Session{}
	retSession.SessionID = value
	retSession.Reactions = make(map[uint32]float64)
	retSession.leap = generateOddNumber()
	retSession.start = generateOddNumber()
	retSession.i = 1

	Sessions[retSession.SessionID] = retSession

	return retSession, nil
}

func (s *Session) UpdateReactionPID(pid uint32, reaction float64) error {
	s.Reactions[pid] = reaction
	return nil
}

func (s *Session) UpdateReaction(reaction float64) error {
	return s.UpdateReactionPID(s.CurrentPID, reaction)
}

func (s *Session) CalculateRating(newPID uint32) (float64, error) {
	var projectedRating float64
	var count float64

	// sums up similarity * Reactions for every rated image and the new image
	for ratedPID, reaction := range s.Reactions {
		similarity, err := lookupSimilarity(ratedPID, newPID)
		if err != nil {
			return 0, nil
		}

		projectedRating += similarity * reaction
		count += 1
	}

	//normalizes it
	projectedRating /= count

	return projectedRating, nil
}

func lookupSimilarity(a uint32, b uint32) (float64, error) {

	// var small, big uint32
	// if a < b {
	// 	small = a
	// 	big = b
	// } else {
	// 	small = b
	// 	big = a
	// }

	// var similarity float64

	// err := db.Db.QueryRow("SELECT similarity FROM similarities WHERE small=? AND big=?", small, big).Scan(&similarity)
	// if err != nil {
	// 	return 0, err
	// }
	// return similarity, nil

	return 0.5, nil
}

func (s *Session) NextPictureUrlLimit(limit uint32) (string, error) {
	first := true
	for {
		if !first {
			if s.i == 0 {
				return "", errors.New("Exhaused database")
			}
			s.i += 1
		} else {
			first = false
		}
		s.CurrentPID = s.leap*s.i + s.start

		if s.CurrentPID > limit {
			continue
		}

		url, err := getPictureUrl(s.CurrentPID)
		if err != nil {
			log.Println(err)
			continue
		}

		s.i += 1
		return url, nil

	}

}

func (s *Session) NextPictureUrl() (string, error) {
	var uintMax uint32 = 0
	uintMax -= 1
	return s.NextPictureUrlLimit(uintMax)
}

func getPictureUrl(pid uint32) (string, error) {
	// var url string
	// err := db.Db.QueryRow("SELECT url FROM urls WHERE pid=?", pid).Scan(&url)
	// if err != nil {
	// 	return 0, err
	// }
	// return url, nil

	return "https://www.marktai.com/img/corgi_bg.jpg", nil
}
