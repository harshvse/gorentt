package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/jackpal/bencode-go"
)

type BenTorrent struct {
	Announce string      `bencode:"announce"`
	Info     TorrentInfo `bencode:"info"`
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

func (benTorrent BenTorrent) generateInfoHash() [20]byte {

	infoHash := sha1.Sum([]bytes(benTorrent))
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

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}

	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}

	base.RawQuery = params.Encode()
	return base.String(), nil
}

func ReadTorrent(r io.Reader) (*BenTorrent, error) {
	benTorrentFile := BenTorrent{}
	if err := bencode.Unmarshal(r, &benTorrentFile); err != nil {
		return nil, err
	}
	return &benTorrentFile, nil
}

type Peer struct {
	IP   net.IP
	Port uint16
}

func UnmarshalPeers(PeersBin []byte) ([]Peer, error) {
	peerSize := 6 // each seperate peer will have 6 bytes 4 for IP and 2 for Port number
	numPeers := len(PeersBin) / peerSize
	if (len(PeersBin) % peerSize) != 0 {
		err := fmt.Errorf("Recieved Corrupt Peers")
		return nil, err
	}

	peers := make([]Peer, numPeers)

	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		peers[i].IP = net.IP(PeersBin[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(PeersBin[offset+4 : offset+6])
	}
	return peers, nil
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

	clientID := []byte("gorentt-k8hj0wgej6ch")
	url, err := torrentFile.buildTrackerURL([20]byte(clientID), 9888)

	peersBin, err := http.Get(url)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("buildURL: %s", url)
	fmt.Printf("peersBin.Body: %s\n", peersBin.Status)

	// peer := UnmarshalPeers()

	// conn, err := net.DialTimeout("TCP", peer.string(), 3*time.Second)

}
