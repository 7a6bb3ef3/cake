package pool

func NewTickets(maxConnnum int) *Tickets {
	p := &Tickets{}
	p.localTks = make(chan struct{}, maxConnnum)
	return p
}

type Tickets struct {
	localTks chan struct{}
}

func (p *Tickets) GetTicket() {
	select {
	case p.localTks <- struct{}{}:
		// TODO timeout
	}
}

func (p *Tickets) FreeTicket() {
	_ = <-p.localTks
}
