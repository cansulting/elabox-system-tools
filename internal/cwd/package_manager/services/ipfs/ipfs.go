package ipfs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"sync"

	gopath "path"

	files "github.com/ipfs/go-ipfs-files"
	ipfshttp "github.com/ipfs/go-ipfs-http-client"
	iface "github.com/ipfs/interface-go-ipfs-core"
	ipath "github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

var locker sync.Mutex = sync.Mutex{}

func DownloadJson(cid string, val interface{}) error {
	var buf bytes.Buffer
	if err := DownloadFile(context.TODO(), cid, &buf, nil, nil, nil); err != nil {
		return err
	}
	if err := json.Unmarshal(buf.Bytes(), &val); err != nil {
		return errors.New("failed to parse json, " + err.Error())
	}
	return nil
}

func DownloadAndSaveFile(ctx context.Context, ipfsPath string, outputPath string, progressWriter io.Writer, fileNode func(files.Node)) error {
	locker.Lock()
	defer locker.Unlock()
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	//writer := io.MultiWriter(file, progressWriter)
	return DownloadFile(ctx, ipfsPath, file, progressWriter, nil, fileNode)
}

func DownloadFile(ctx context.Context, ipfsPath string, writer io.Writer, progressWriter io.Writer, peers []string, file func(files.Node)) error {
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// initialize http api
	ipfs, err := http(cctx)
	if err != nil {
		return err
	}
	// connect to peers
	go connect(cctx, ipfs, peers)

	path, err := parsePath(ipfsPath)
	if err != nil {
		return errors.New("invalid ipfs path")
	}
	// download the file
	out, err := ipfs.Unixfs().Get(cctx, path)
	if err != nil {
		if err == context.Canceled {
			return nil
		}
		return err
	}
	if file != nil {
		file(out)
	}
	// read buffer
	switch nodeType := out.(type) {
	case files.File:
		if progressWriter == nil {
			if _, err := io.Copy(writer, nodeType); err != nil {
				return err
			}
		} else {
			if _, err := io.Copy(writer, io.TeeReader(nodeType, progressWriter)); err != nil {
				return err
			}
		}
	default:
		return errors.New("unknown file type")
	}
	return nil
}

// parse/resolve to a valid ipfs path
func parsePath(path string) (ipath.Path, error) {
	ipfsPath := ipath.New(path)
	if ipfsPath.IsValid() == nil {
		return ipfsPath, nil
	}

	u, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("%q could not be parsed: %s", path, err)
	}

	switch proto := u.Scheme; proto {
	case "ipfs", "ipld", "ipns":
		ipfsPath = ipath.New(gopath.Join("/", proto, u.Host, u.Path))
	case "http", "https":
		ipfsPath = ipath.New(u.Path)
	default:
		return nil, fmt.Errorf("%q is not recognized as an IPFS path", path)
	}
	return ipfsPath, ipfsPath.IsValid()
}

func http(ctx context.Context) (iface.CoreAPI, error) {
	httpApi, err := ipfshttp.NewLocalApi()
	if err != nil {
		return nil, err
	}
	err = httpApi.Request("version").Exec(ctx, nil)
	if err != nil {
		return nil, err
	}
	return httpApi, nil
}

// connect to peers
func connect(ctx context.Context, ipfs iface.CoreAPI, peers []string) error {
	var wg sync.WaitGroup
	pinfos := make(map[peer.ID]*peer.AddrInfo, len(peers))
	for _, addrStr := range peers {
		addr, err := ma.NewMultiaddr(addrStr)
		if err != nil {
			return err
		}
		pii, err := peer.AddrInfoFromP2pAddr(addr)
		if err != nil {
			return err
		}
		pi, ok := pinfos[pii.ID]
		if !ok {
			pi = &peer.AddrInfo{ID: pii.ID}
			pinfos[pi.ID] = pi
		}
		pi.Addrs = append(pi.Addrs, pii.Addrs...)
	}

	wg.Add(len(pinfos))
	for _, pi := range pinfos {
		go func(pi *peer.AddrInfo) {
			defer wg.Done()
			log.Printf("attempting to connect to peer: %q\n", pi)
			err := ipfs.Swarm().Connect(ctx, *pi)
			if err != nil {
				log.Printf("failed to connect to %s: %s", pi.ID, err)
			}
			log.Printf("successfully connected to %s\n", pi.ID)
		}(pi)
	}
	wg.Wait()
	return nil
}
