package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

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

type PeerBin struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

type Peer struct {
	IP   net.IP
	Port uint16
}

func ReadBenTorrent(r io.Reader) (*BenTorrent, error) {
	benTorrentFile := BenTorrent{}
	if err := bencode.Unmarshal(r, &benTorrentFile); err != nil {
		return nil, err
	}
	return &benTorrentFile, nil
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



func UnmarshalPeers(PeersBin []byte) ([]Peer, error) {
	peerSize := 6 // each seperate peer will have 6 bytes 4 for IP and 2 for Port number
	fmt.Println(len(PeersBin))
	numPeers := len(PeersBin) / peerSize
	if (len(PeersBin) % peerSize) != 0 {
		err := fmt.Errorf("recieved corrupt peers")
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
