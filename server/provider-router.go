package server

type Route interface {
	Routes() []RouteConfig
}

type ProviderRouter struct {
	Routes []RouteConfig
}

func (b *ProviderRouter) RegisterHandlers(handlers ...interface{}) {
	for _, h := range handlers {
		if rp, ok := h.(Route); ok {
			b.Routes = append(b.Routes, rp.Routes()...)
		}
	}
}
