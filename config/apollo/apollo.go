package apollo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/okayes/gopkgs/logger"
)

type Apollo struct {
	url        string
	releaseKey string
	data       *ApolloInfo
	handler    Handler
}

type ApolloInfo struct {
	AppID          string      `json:"appId"`
	Cluster        string      `json:"cluster"`
	NamespaceName  string      `json:"namespaceName"`
	Configurations interface{} `json:"configurations"`
	ReleaseKey     string      `json:"releaseKey"`
}

type Handler interface {
	Handle(data interface{})
}

func NewApollo(url string, data *ApolloInfo, handler Handler) *Apollo {
	return &Apollo{url: url, data: data, handler: handler}
}

func (c *Apollo) LoadConfig(interval time.Duration) {
	if len(c.url) == 0 {
		log.Panicln("apollo config url is empty")
	}

	c.syncConfig()
	go func() {
		for {
			time.Sleep(interval)
			c.syncConfig()
		}
	}()
}

func (c *Apollo) syncConfig() {
	defer func() {
		if err := recover(); err != nil {
			em := fmt.Sprintf("sync apollo config panic: %s, url: %s", err, c.url)
			logger.ErrorMsg(em)
		}
	}()

	url := c.url + "?releaseKey=" + c.releaseKey
	resp, err := http.Get(url)
	if err != nil {
		em := fmt.Sprintf("Get config from apollo error: %s, url: %s", err, url)
		logger.ErrorMsg(em)
		return
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		em := fmt.Sprintf("Read body from apollo response error: %s", err)
		logger.ErrorMsg(em)
		return
	}

	if resp.StatusCode == http.StatusOK {
		err := json.Unmarshal(bytes, c.data)
		if err != nil {
			em := fmt.Sprintf("Unmarshal config from apollo response error: %s, data: %s", err, bytes)
			logger.ErrorMsg(em)
			return
		}

		if c.handler != nil {
			c.handler.Handle(c.data.Configurations)
		}

		c.releaseKey = c.data.ReleaseKey
	} else if resp.StatusCode != http.StatusNotModified {
		em := fmt.Sprintf("Get config from apollo error, response status code: %d", resp.StatusCode)
		logger.ErrorMsg(em)
	}
}
