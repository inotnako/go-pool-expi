package oneof

import "sync"

func MakePool(s int) *Pool {
	p := &Pool{
		size: s,
		run:  make(chan func() int, s),
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
			go func(jid int) {
				p.up.Done()
				//println(`UP JID=`, jid)
				//defer println(`UPJID=`, jid)
				defer p.down.Done()
				for {
					select {
					case j, ok := <-p.run:
						if ok {
							j()
							//println(`JID=`, jid, `VAL=`, )
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
	size int
	run  chan func() int
	up   *sync.WaitGroup
	down *sync.WaitGroup
}

func (p *Pool) Go(f func() int) {
	p.run <- f
}

func (p *Pool) Stop() {
	close(p.run)
	p.down.Wait()
	//println(`CLOSE`)
}

func (p *Pool) Wait() {
	p.down.Wait()
}
