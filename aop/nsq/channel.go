package nsq

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
)

const (
	Chat  uint32 = 0
	Logic uint32 = 1
)

var nsqInsMgr *NsqInstanceMgr

type NsqInstanceMgr struct {
	nsqs *sync.Map
}

func Init(lookupAddr string, poolSize int) {
	if poolSize == 0 {
		poolSize = 3
	}
	nsqInsMgr = &NsqInstanceMgr{nsqs: &sync.Map{}}
	insAddrs := strings.Split(lookupAddr, "|")
	for i := 0; i < len(insAddrs); i++ {
		nsqInsMgr.addNsqInstance(uint32(i), insAddrs[i], poolSize)
	}
}

func (mgr *NsqInstanceMgr) addNsqInstance(insType uint32, addr string, poolSize int) {
	_, ok := mgr.nsqs.Load(insType)
	if ok {
		return
	}

	ins := newNsqInstance(addr, poolSize)

	mgr.nsqs.Store(insType, ins)
}

func (mgr *NsqInstanceMgr) delNsqInstance(insType uint32) {
	_, ok := mgr.nsqs.Load(insType)
	if !ok {
		return
	}

	mgr.nsqs.Delete(insType)
}

func (mgr *NsqInstanceMgr) getNsqInstance(insType uint32) (*NsqInstance, error) {
	ins, ok := mgr.nsqs.Load(insType)
	if !ok {
		return nil, fmt.Errorf("NSQ实例类型错误")
	}

	nsq, ok := ins.(*NsqInstance)

	if !ok {
		return nil, fmt.Errorf("NSQ实例断言错误")
	}

	return nsq, nil
}

type NsqInstance struct {
	lookups   []string
	producers []*nsq.Producer
}

func newNsqInstance(lookupAddr string, poolSize int) *NsqInstance {
	addrs := strings.Split(lookupAddr, ",")
	ins := &NsqInstance{}
	ins.producers = make([]*nsq.Producer, 0, 3*poolSize)
	ins.lookups = append(ins.lookups, addrs...)
	ins.newProducerPool(poolSize)
	return ins
}

func (ins *NsqInstance) getProducer() *nsq.Producer {
	var nsqIns *nsq.Producer
	pLen := len(ins.producers)
	if pLen > 0 {
		randIdx := rand.Intn(pLen)
		nsqIns = ins.producers[randIdx]
	}
	return nsqIns
}

func newProducer(addr string) (*nsq.Producer, error) {
	config := nsq.NewConfig()
	pro, err := nsq.NewProducer(addr, config)
	if err != nil {
		return nil, err
	}
	pro.SetLogger(nil, nsq.LogLevelDebug)
	return pro, nil
}

func (ins *NsqInstance) newProducerPool(poolSize int) {
	nsqds := ins.getAvailableTCPAddr()
	nsqdLen := len(nsqds)
	if nsqdLen <= 0 {
		return
	}

	for i := 0; i < poolSize; i++ {
		for _, addr := range nsqds {
			producer, err := newProducer(addr)
			if err != nil {
				continue
			}

			ins.producers = append(ins.producers, producer)
		}
	}
}

type NsqNodeData struct {
	RemoteAddr    string      `json:"remote_address"`
	HostName      string      `json:"hostname"`
	BroadcastAddr string      `json:"broadcast_address"`
	TCPPort       int         `json:"tcp_port"`
	HTTPPort      int         `json:"http_port"`
	Version       string      `json:"version"`
	Tombstones    interface{} `json:"tombstones"`
	Topics        interface{} `json:"topics"`
}

type NsqdNodesData struct {
	Producers []*NsqNodeData `json:"producers"`
}

func (ins *NsqInstance) getAvailableTCPAddr() []string {
	var nsqdAddrs []string
	for _, lookupAddr := range ins.lookups {
		querURL := fmt.Sprintf("http://%v/nodes", lookupAddr)
		resp, err := http.Get(querURL)
		if err != nil {
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			continue
		}

		var nsqds NsqdNodesData
		if err := json.Unmarshal(body, &nsqds); err != nil {
			resp.Body.Close()
			continue
		}

		for _, nsqd := range nsqds.Producers {
			addr := fmt.Sprintf("%v:%v", nsqd.BroadcastAddr, nsqd.TCPPort)
			nsqdAddrs = append(nsqdAddrs, addr)
		}
		resp.Body.Close()
	}
	return nsqdAddrs
}

func (ins *NsqInstance) newConsumer(topic, channl string, handler nsq.Handler) *nsq.Consumer {
	config := nsq.NewConfig()
	config.MaxInFlight = 5
	config.LookupdPollInterval = 30 * time.Second
	consumer, err := nsq.NewConsumer(topic, channl, config)
	if err != nil {
		return nil
	}
	consumer.SetLogger(nil, nsq.LogLevelError)
	consumer.AddHandler(handler)
	if err := consumer.ConnectToNSQLookupds(ins.lookups); err != nil {
		return nil
	}
	return consumer
}

func (ins *NsqInstance) createTopic(topic string) error {
	nsqds := ins.getAllNsqdHTTPAddr()
	if len(nsqds) <= 0 {
		return errors.New("nsqd节点数量为0")
	}
	for _, nsqd := range nsqds {
		ins.createTopicInEveryNsqdNode(topic, nsqd)
	}
	return nil
}

func (ins *NsqInstance) createTopicInEveryNsqdNode(topic string, nodeAddr string) {
	topurl := "http://" + nodeAddr + "/topic/create?topic=" + topic
	request, err := http.NewRequest("POST", topurl, nil)
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
	}
	resp.Body.Close()
}

func (ins *NsqInstance) deleteTopic(topic string) error {
	nsqds := ins.getAllNsqdHTTPAddr()
	if len(nsqds) <= 0 {
		return errors.New("nsqd节点数量为0")
	}
	for _, nsqd := range nsqds {
		ins.deleteTopicInEveryNsqdNode(topic, nsqd)
	}
	return nil
}

func (ins *NsqInstance) deleteTopicInEveryNsqdNode(topic string, nodeAddr string) {
	url := "http://" + nodeAddr + "/topic/delete?topic=" + topic
	request, err := http.NewRequest("POST", url, nil)
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
	}
	resp.Body.Close()
}

func (ins *NsqInstance) createChannel(topic, channel string) error {
	nsqds := ins.getAllNsqdHTTPAddr()
	if len(nsqds) <= 0 {
		return errors.New("nsqd节点数量为0")
	}
	for _, nsqd := range nsqds {
		ins.createChannelInEveryNsqdNode(topic, channel, nsqd)
	}
	return nil
}

func (ins *NsqInstance) createChannelInEveryNsqdNode(topic, channel, nodeAddr string) {
	chaurl := "http://" + nodeAddr + "/channel/create?topic=" + topic + "&channel=" + channel
	request, err := http.NewRequest("POST", chaurl, nil)
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
	}
	resp.Body.Close()
}

// deleteChannel 删除订阅通道
func (ins *NsqInstance) deleteChannel(topic, channel string) error {
	nsqds := ins.getAllNsqdHTTPAddr()
	if len(nsqds) <= 0 {
		return errors.New("nsqd节点数量为0")
	}
	for _, nsqd := range nsqds {
		ins.deleteChannlInEveryNsqdNode(topic, channel, nsqd)
	}
	return nil
}

func (ins *NsqInstance) deleteChannlInEveryNsqdNode(topic, channl, nodeAddr string) {
	url := "http://" + nodeAddr + "/channel/delete?topic=" + topic + "&channel=" + channl
	request, err := http.NewRequest("POST", url, nil)
	var resp *http.Response
	resp, err = http.DefaultClient.Do(request)
	if err != nil {
	}
	resp.Body.Close()
}

func (ins *NsqInstance) getAllNsqdHTTPAddr() []string {
	var nsqdAddrs []string
	for _, lookupAddr := range ins.lookups {
		querURL := fmt.Sprintf("http://%v/nodes", lookupAddr)
		resp, err := http.Get(querURL)
		if err != nil {
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			continue
		}
		var nsqds NsqdNodesData
		if err := json.Unmarshal(body, &nsqds); err != nil {
			resp.Body.Close()
			continue
		}

		for _, nsqd := range nsqds.Producers {
			addr := fmt.Sprintf("%v:%v", nsqd.BroadcastAddr, nsqd.HTTPPort)
			nsqdAddrs = append(nsqdAddrs, addr)
		}
		resp.Body.Close()

	}
	return nsqdAddrs
}

// NewConsumer 根据主题创建消费者
func NewConsumer(insType uint32, topic, channl string, handler nsq.Handler) *nsq.Consumer {
	nsqIns, err := nsqInsMgr.getNsqInstance(insType)
	if err != nil {
		return nil
	}

	return nsqIns.newConsumer(topic, channl, handler)
}

// PublishAsync 发布消息
func PublishAsync(insType uint32, topic string, data []byte, doneChan chan *nsq.ProducerTransaction) error {
	nsqIns, err := nsqInsMgr.getNsqInstance(insType)
	if err != nil {
		return err
	}
	producer := nsqIns.getProducer()
	if producer == nil {
		return errors.New("获取producer失败")
	}
	return producer.PublishAsync(topic, data, doneChan)
}

// DeferredPublishAsync 发布消息
func DeferredPublishAsync(insType uint32, topic string, data []byte,
	doneChan chan *nsq.ProducerTransaction, delay time.Duration) error {
	nsqIns, err := nsqInsMgr.getNsqInstance(insType)
	if err != nil {
		return err
	}
	producer := nsqIns.getProducer()
	if producer == nil {
		return errors.New("获取producer失败")
	}
	return producer.DeferredPublishAsync(topic, delay, data, doneChan)
}

// CreateTopic 创建主题
func CreateTopic(insType uint32, topic string) error {
	nsqIns, err := nsqInsMgr.getNsqInstance(insType)
	if err != nil {
		return err
	}

	return nsqIns.createTopic(topic)
}

// DeleteTopic 删除主题
func DeleteTopic(insType uint32, topic string) error {
	nsqIns, err := nsqInsMgr.getNsqInstance(insType)
	if err != nil {
		return err
	}
	return nsqIns.deleteTopic(topic)
}

// CreateChannel ...
func CreateChannel(insType uint32, topic, channel string) error {
	nsqIns, err := nsqInsMgr.getNsqInstance(insType)
	if err != nil {
		return err
	}

	return nsqIns.createChannel(topic, channel)
}

// DeleteChannel ...
func DeleteChannel(insType uint32, topic, channel string) error {
	nsqIns, err := nsqInsMgr.getNsqInstance(insType)
	if err != nil {
		return err
	}
	return nsqIns.deleteChannel(topic, channel)
}
