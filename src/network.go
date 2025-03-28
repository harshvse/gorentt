package main

import (
	"net/url"
	"strconv"
)

func (t *TorrentFile) buildTrackerURL(peerID [20]byte, port uint16) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("info_hash", string(t.InfoHash[:]))
	params.Add("peer_id", string(peerID[:]))
	params.Add("port", strconv.Itoa(int(port)))
	params.Add("uploaded", "0")
	params.Add("downloaded", "0")
	params.Add("compact", "1")
	params.Add("left", strconv.Itoa(t.Length))

	base.RawQuery = params.Encode()
	return base.String(), nil
}

func (p *Peer) ToString() string {
	url := p.IP.String() + ":" + strconv.Itoa(int(p.Port))
	return url
}
