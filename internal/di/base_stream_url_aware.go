package di

import "net/url"

type BaseStreamURLAware struct {
	baseStreamURL *url.URL
}

func (b *BaseStreamURLAware) SetBaseStreamURL(url *url.URL) {
	b.baseStreamURL = url
}

func (b *BaseStreamURLAware) GetBaseStreamURL() *url.URL {
	return b.baseStreamURL
}
