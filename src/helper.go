package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"

	"github.com/jackpal/bencode-go"
)

func (i *TorrentInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("what in the world are these pieces harsh! %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (benTorrent BenTorrent) generateInfoHash() ([20]byte, error) {
	// Create a buffer to hold bencoded info dictionary
	var buf bytes.Buffer

	// Encode the 'info' dictionary into the buffer
	err := bencode.Marshal(&buf, benTorrent.Info)
	if err != nil {
		return [20]byte{}, fmt.Errorf("failed to bencode info dict: %w", err)
	}

	// Compute SHA-1 hash of the bencoded info dictionary
	infoHash := sha1.Sum(buf.Bytes())

	return infoHash, nil
}
