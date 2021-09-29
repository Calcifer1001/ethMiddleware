package main

import (
	"flag"
	"fmt"

	diatoken "github.com/ethMiddlewareTest/config/dia-token"
	unitoken "github.com/ethMiddlewareTest/config/uni-token"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/event"
	log "github.com/sirupsen/logrus"
)

const (
	// wsDial     = "wss://mainnet.infura.io/ws/v3/995a689a50d84bfb9c7ab144cd720bfd"
	// wsDial     = "wss://mainnet.infura.io/ws/v3/251a25bd10b8460fa040bb7202e22571"
	// DIAAddress = "0x84cA8bc7997272c7CfB4D0Cd3D55cd942B3c9419"
	// DIAAddress = "0x1f9840a85d5af5bf1d1762f925bdaddc4201f984" UNISWAP
	DIAAddress = "0xdac17f958d2ee523a2206206994597c13d831ec7" // USDT

	UNIAddress = "0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984"
)

var (
	// dial = flag.String("wsDial", "https://mainnet.infura.io/v3/9bdd9b1d1270497795af3f522ad85091", "Ethereum client")
	dial = flag.String("wsDial", "ws://localhost:8080", "Ethereum client")
	// dial    = flag.String("wsDial", "wss://mainnet.infura.io/ws/v3/9bdd9b1d1270497795af3f522ad85091", "Ethereum client")
	address = flag.String("Address", DIAAddress, "Subscribe to transfers of which Ethereum asset.")
)

func init() {
	flag.Parse()
	flag.Usage()
	fmt.Println("Available addresses: ", []string{DIAAddress, UNIAddress})
	// var dialer = *websocket.DefaultDialer
	// dialer.TLSClientConfig = &tls.TLSConfig{InsecureSkipVerify: true}
	// http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

}

func main() {
	log.Infof("Start watching transfers of asset %s \n", *address)

	wsClient, err := ethclient.Dial(*dial)

	if err != nil {
		log.Fatal(err)
	}

	sink, sub, err := getDIASink(wsClient)
	if err != nil {
		log.Error("get sink: ", err)
	}
	// sink, sub, err := getUNISink(wsClient)
	// if err != nil {
	// 	log.Error("get sink: ", err)
	// }

	for {
		select {
		case subErr := <-sub.Err():
			log.Error("subscription: ", subErr)
		case rawSwap, ok := <-sink:
			if ok {
				log.Infof("got swap from %s to %s with value %v \n", rawSwap.From, rawSwap.To, rawSwap.Value)
				// log.Infof("got swap from %s to %s with value %v \n", rawSwap.From, rawSwap.To, rawSwap.Amount)
			} else {
				log.Info("channel not ok")
			}
		}
	}

}

func getDIASink(wsClient *ethclient.Client) (chan *diatoken.DIATokenTransfer, event.Subscription, error) {
	sink := make(chan *diatoken.DIATokenTransfer)
	cc, err := diatoken.NewDIATokenFilterer(common.HexToAddress(DIAAddress), wsClient)
	if err != nil {
		log.Error("getting token filterer: ", err)
		return sink, nil, err
	}
	sub, err := cc.WatchTransfer(
		&bind.WatchOpts{},
		sink,
		[]common.Address{},
		[]common.Address{},
	)
	if err != nil {
		log.Error("watch transfer: ", err)
		return sink, sub, err
	}
	return sink, sub, err
}

func getUNISink(wsClient *ethclient.Client) (chan *unitoken.UniTransfer, event.Subscription, error) {
	sink := make(chan *unitoken.UniTransfer)
	cc, err := unitoken.NewUniFilterer(common.HexToAddress(UNIAddress), wsClient)
	if err != nil {
		log.Error("getting token filterer: ", err)
		return sink, nil, err
	}
	sub, err := cc.WatchTransfer(
		&bind.WatchOpts{},
		sink,
		[]common.Address{},
		[]common.Address{},
	)
	if err != nil {
		log.Error("watch transfer: ", err)
		return sink, sub, err
	}
	return sink, sub, err
}
