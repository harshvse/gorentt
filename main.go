package main

import (
	"bytes"
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

func (benTorrent BenTorrent) ToTorrentFile() (*TorrentFile, error) {
	torrentFile := TorrentFile{}

	torrentFile.Announce = benTorrent.Announce

	torrentFile.name = benTorrent.Info.Name

	infoHash, err := benTorrent.generateInfoHash()

	if err != nil {
		return nil, err
	}

	torrentFile.InfoHash = infoHash

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

	// ✅ Correctly encode info_hash (using URL-encoded raw bytes)
	params := url.Values{}
	params.Add("info_hash", string(t.InfoHash[:])) // No extra encoding needed
	params.Add("peer_id", string(peerID[:]))
	params.Add("port", strconv.Itoa(int(port)))
	params.Add("uploaded", "0")
	params.Add("downloaded", "0")
	params.Add("compact", "1")
	params.Add("left", strconv.Itoa(t.Length))

	base.RawQuery = params.Encode() // Let Go handle correct encoding
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

	fmt.Println(benTorrentFile.Announce)
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
	fmt.Printf("peersBin.Body: %s\n", peersBin.Body)

	// peer := UnmarshalPeers()

	// conn, err := net.DialTimeout("TCP", peer.string(), 3*time.Second)

}
