package main

import (
	"fmt"
	"strings"

	"github.com/jackpal/bencode-go"
)

type Torrent struct {
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

func main() {

	sampleTorrent := `d8:announce41:http://bttracker.debian.org:6969/announce7:comment35:"Debian CD from cdimage.debian.org"13:creation datei1573903810e9:httpseedsl145:https://cdimage.debian.org/cdimage/release/10.2.0//srv/cdbuilder.debian.org/dst/deb-cd/weekly-builds/amd64/iso-cd/debian-10.2.0-amd64-netinst.iso145:https://cdimage.debian.org/cdimage/archive/10.2.0//srv/cdbuilder.debian.org/dst/deb-cd/weekly-builds/amd64/iso-cd/debian-10.2.0-amd64-netinst.isoe4:infod6:lengthi351272960e4:name31:debian-10.2.0-amd64-netinst.iso12:piece lengthi262144e6:pieces26800:�����PS�^�� (binary blob of the hashes of each piece)ee`

	r := strings.NewReader(sampleTorrent)

	torrentFile := Torrent{}
	if err := bencode.Unmarshal(r, &torrentFile); err != nil {
		fmt.Errorf("Failed to read torrent file")
	}

}
