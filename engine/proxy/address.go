//
//Created by xuzhuoxi
//on 2019-03-13.
//@author xuzhuoxi
//
package proxy

import "sync"

//难住地址与id的双向映射
type IAddressProxy interface {
	//能过地址找id
	GetId(address string) (id string, ok bool)
	//能过id找地址
	GetAddress(id string) (address string, ok bool)
	//把id和地址加入映射表
	MapIdAddress(id string, address string)
	//移除id相关映射
	RemoveById(id string)
	//移除地址相关映射
	RemoveByAddress(address string)
	//重置
	Reset()
}

func NewIAddressProxy() IAddressProxy {
	return NewAddressProxy()
}

func NewAddressProxy() *AddressProxy {
	return &AddressProxy{idAddr: make(map[string]string), addrId: make(map[string]string)}
}

type AddressProxy struct {
	idAddr map[string]string
	addrId map[string]string
	mu     sync.RWMutex
}

func (p *AddressProxy) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.idAddr = make(map[string]string)
	p.addrId = make(map[string]string)
}

func (p *AddressProxy) GetId(address string) (id string, ok bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	id, ok = p.addrId[address]
	return
}

func (p *AddressProxy) GetAddress(id string) (address string, ok bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	id, ok = p.idAddr[id]
	return
}

func (p *AddressProxy) MapIdAddress(id string, address string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.removeId(id)
	p.removeAddress(address)
	p.idAddr[id] = address
	p.addrId[address] = id
}

func (p *AddressProxy) RemoveById(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.removeId(id)
}

func (p *AddressProxy) RemoveByAddress(address string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.removeAddress(address)
}

func (p *AddressProxy) removeId(id string) {
	if address, ok := p.idAddr[id]; ok {
		delete(p.addrId, address)
		delete(p.idAddr, id)
	}
}

func (p *AddressProxy) removeAddress(address string) {
	if id, ok := p.addrId[address]; ok {
		delete(p.idAddr, id)
		delete(p.addrId, address)
	}
}
