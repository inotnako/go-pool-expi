package each

import "sync"

func MakePool(s int) *Pool {
	p := &Pool{
		size: s,
		proc: make([]chan func() int, s),
		up:   &sync.WaitGroup{},
		down: &sync.WaitGroup{},
	}
	p.up.Add(s + 1)
	p.down.Add(s)
	//println(`UP`)
	go func() {
		//defer println(`exit from UPer`)
		defer p.up.Done()
		for i := 0; i < s; i++ {
			p.proc[i] = make(chan func() int, 1)
			go func(jid int) {
				p.up.Done()
				//println(`UP JID=`, jid)
				//defer println(`UPJID=`, jid)
				defer p.down.Done()
				in := p.proc[jid]
				for {
					select {
					case j, ok := <-in:
						if ok {
							j()
							//println(`JID=`, jid, `VAL=`, j())
						} else {
							//println(`exit JID=`, jid)
							return
						}
					default:
					}
				}
			}(i)
		}

	}()
	p.up.Wait()
	return p
}

type Pool struct {
	size    int
	proc    []chan func() int
	current int
	up      *sync.WaitGroup
	down    *sync.WaitGroup
}

func (p *Pool) next() int {
	if p.current+1 == p.size {
		p.current = 0
	} else {
		p.current += 1
	}

	//println("\tNEXT=", p.current)
	return p.current
}

func (p *Pool) Go(f func() int) {
	p.proc[p.next()] <- f
}

func (p *Pool) Stop() {
	for pid := range p.proc {
		close(p.proc[pid])
	}

	p.down.Wait()
	//println(`CLOSE`)
}

func (p *Pool) Wait() {
	p.down.Wait()
}
