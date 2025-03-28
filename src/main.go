package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/jackpal/bencode-go"
)

func main() {

	//TODO replace this with being able to pass the location of the torrent file
	file, err := os.Open("debian.torrent")
	if err != nil {
		fmt.Printf("Failed to open torrent file. Error: %v\n", err)
		return
	}
	defer file.Close()

	// torrent file is bencode encoded convert it to BencodeTorrent
	benTorrentFile, err := ReadBenTorrent(file)
	if err != nil {
		fmt.Printf("Failed to read torrent file. Error: %s", err.Error())
		return
	}

	// convert the bentorrent file to torrent file which includes info_hash
	torrentFile, err := benTorrentFile.ToTorrentFile()
	if err != nil {
		fmt.Printf("Failed to convert into torrent file. Error: %s", err)
		return
	}

	// TODO create a random torrent client ID for each user
	clientID := []byte("gorentt-k8hj0wgej6ch")
	url, err := torrentFile.buildTrackerURL([20]byte(clientID), 9888)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// send GET request to the announce URL which returns Interval and Binary of peer ip and port
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	// convert the contents of body to []byte
	bencodeResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// unmarshal the Interval and peers from bencode to a struct
	peersBin := PeerBin{}
	bencodeResponseReader := bytes.NewReader(bencodeResponse)
	if err := bencode.Unmarshal(bencodeResponseReader, &peersBin); err != nil {
		fmt.Println("failed to read peers bin from response")
	}

	// read the IP and Port for each peer and convert it to a []Peers
	peer, err := UnmarshalPeers([]byte(peersBin.Peers))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(peer)

	// conn, err := net.DialTimeout("TCP", peer.string(), 3*time.Second)

}
