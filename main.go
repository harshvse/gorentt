package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"

	"github.com/jackpal/bencode-go"
)

type BenTorrent struct {
	Announce     string      `bencode:"announce"`
	Comment      string      `bencode:"comment"`
	CreationDate int         `bencode:"creation date"`
	Info         TorrentInfo `bencode:"info"`
}

type TorrentInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	Pieces      [][20]byte
	PieceLength int
	Length      int
	name        string
}

func (i *TorrentInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		err := fmt.Errorf("What in the world are these pieces harsh! %d", len(buf))
		return nil, err
	}
	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)

	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}

func (benTorrent BenTorrent) generateInfoHash() [20]byte {

	infoHash := sha1.New()
	infoHash.Write([]byte(benTorrent.Info.Name))
	infoHash.Write([]byte(benTorrent.Info.Pieces))
	infoHash.Write([]byte(benTorrent.Info.Pieces))
	infoHash.Write([]byte(string(benTorrent.Info.PieceLength)))

	return [20]byte(infoHash.Sum(nil))
}

func (benTorrent BenTorrent) ToTorrentFile() (*TorrentFile, error) {
	torrentFile := TorrentFile{}

	torrentFile.Announce = benTorrent.Announce

	torrentFile.name = benTorrent.Info.Name

	torrentFile.InfoHash = benTorrent.generateInfoHash()

	torrentFile.PieceLength = benTorrent.Info.PieceLength

	pieces, err := benTorrent.Info.splitPieceHashes()
	torrentFile.Pieces = pieces
	if err != nil {
		return nil, err
	}
	return &torrentFile, nil
}
func ReadTorrent(r io.Reader) (*BenTorrent, error) {
	benTorrentFile := BenTorrent{}
	if err := bencode.Unmarshal(r, &benTorrentFile); err != nil {
		return nil, err
	}
	return &benTorrentFile, nil
}

func main() {

	file, err := os.Open("debian.torrent")
	if err != nil {
		fmt.Printf("Failed to open torrent file. Error: %v\n", err)
		return
	}
	defer file.Close()
	benTorrentFile, err := ReadTorrent(file)

	if err != nil {
		fmt.Printf("Failed to read torrent file. Error: %s", err.Error())
		return
	}

	torrentFile, err := benTorrentFile.ToTorrentFile()
	if err != nil {
		fmt.Printf("Failed to convert into torrent file. Error: %s", err)
		return
	}

	fmt.Println(torrentFile.Announce)
	return
}
