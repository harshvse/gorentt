package main

import (
	"io"
)

type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49) // 1(len(h.pstr)) + pstr(BitTorrent Protocol) + 8 (zeroes) + 20(info_hash) + 20(peer_id)

	buf[0] = byte(len(h.Pstr))

	curr := 1
	curr += copy(buf[curr:], []byte(h.Pstr))
	curr += copy(buf[curr:], make([]byte, 8))
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.PeerID[:])
	return buf
}

func Deserialize(r io.Reader) (*Handshake, error) {
	handShake := Handshake{}

	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	pstrLength := int(buf[0])

	handShake.Pstr = string(buf[1 : pstrLength+1])

	handShake.InfoHash = [20]byte(buf[pstrLength+1+8 : pstrLength+21+8])
	handShake.PeerID = [20]byte(buf[pstrLength+21+8:])

	return &handShake, nil
}
