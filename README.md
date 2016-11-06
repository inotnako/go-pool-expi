# go-pool-expi
exp with make go pool of workers


## Pool with personal channel listener

* Up pool size=1000 - `5.558913391s`
* Processed 1000 tasks - `5.460518467s`
* Down that pool - `299.447µs`

Code
```
func MakePool(s int) *Pool {
	p := &Pool{
		size: s,
		proc: make([]chan func() int, s),
		up:   &sync.WaitGroup{},
		down: &sync.WaitGroup{},
	}
	p.up.Add(s + 1)
	p.down.Add(s)
	go func() {
		defer p.up.Done()
		for i := 0; i < s; i++ {
			p.proc[i] = make(chan func() int, 1)
			go func(jid int) {
				p.up.Done()
				defer p.down.Done()
				in := p.proc[jid]
				for {
					select {
					case j, ok := <-in:
						if ok {
							j()
							//println(`JID=`, jid, `VAL=`, j())
						} else {
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

...
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

```

## Pool one channel/ditributer of tasks

* Up pool size=1000 - `6.274303933s`
* Processed 1000 tasks - `3.841723ms`
* Down that pool - `282.432µs`

Code
```

func MakePool(s int) *Pool {
	p := &Pool{
		size: s,
		run:  make(chan func() int, s),
		up:   &sync.WaitGroup{},
		down: &sync.WaitGroup{},
	}
	p.up.Add(s + 1)
	p.down.Add(s)
	go func() {
		defer p.up.Done()
		for i := 0; i < s; i++ {
			go func(jid int) {
				p.up.Done()
				defer p.down.Done()
				for {
					select {
					case j, ok := <-p.run:
						if ok {
							j()
						} else {
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

func (p *Pool) Go(f func() int) {
	p.run <- f
}
```
