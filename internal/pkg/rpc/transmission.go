package rpc

import (
	"fmt"
	"github.com/hekmon/transmissionrpc"
	log "github.com/sirupsen/logrus"
)

type TransmissionConnection struct {
	*transmissionrpc.Client
}

func NewTransmissionConnection(hostname string, port int, username, password string, useHTTPS bool) (*TransmissionConnection, error) {
	c, err := transmissionrpc.New(hostname, username, password,
		&transmissionrpc.AdvancedConfig{
			HTTPS: useHTTPS,
			Port: uint16(port),
		})
	tc := TransmissionConnection{c}
	return &tc, err
}

func (tc *TransmissionConnection) GetTorrentList(onlyActive bool) ([]*transmissionrpc.Torrent, error) {
	return tc.TorrentGetAll()
}

func (tc *TransmissionConnection) AddTorrent(url string) (*transmissionrpc.Torrent, error) {
	torrent, err := tc.TorrentAdd(&transmissionrpc.TorrentAddPayload{
		Filename: &url,
	})
	return torrent, err
}

func (tc *TransmissionConnection) RemoveTorrent(id int, deleteData bool) error {
	ids := []int64{int64(id)}
	rp := transmissionrpc.TorrentRemovePayload{
		IDs:             ids,
		DeleteLocalData: deleteData,
	}

	err := tc.TorrentRemove(&rp)
		if err == nil {
			log.Info(fmt.Sprintf("Removed torrent ID %d. Delete data: %v\n", id, deleteData))
		}
	return err
}

func (tc *TransmissionConnection) PauseTorrent(id int) error {
	ids := []int64{int64(id)}
	err := tc.TorrentStopIDs(ids)
	if err == nil {
		log.Info(fmt.Sprintf("Stopped torrent ID %d.\n", id))
	}
	return err
}

func (tc *TransmissionConnection) StartTorrent(id int) error {
	ids := []int64{int64(id)}
	err := tc.TorrentStartIDs(ids)
	if err == nil {
		log.Info(fmt.Sprintf("Started torrent ID %d.\n", id))
	}
	return err
}

func (tc *TransmissionConnection) IsConnected() (bool, int64, int64, error) {
	ok, serverVersion, serverMinimumVersion, err := tc.RPCVersion()
	return ok, serverVersion, serverMinimumVersion, err
}