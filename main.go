package main

import (
	"github.com/antonikonovalov/nn/drive/pool/each"
	"github.com/antonikonovalov/nn/drive/pool/oneof"
	"sync"
	"time"
)

func main() {
	t := time.Now()
	for i := 0; i < 1000; i++ {
		func() int {
			return i
		}()
	}
	println(`SPEED ITER 1000 i:`, time.Since(t).String())
	t = time.Now()
	epool := each.MakePool(1000)
	println(`UP EACH 1000 procs:`, time.Since(t).String())

	wg := &sync.WaitGroup{}
	wg.Add(1000)
	tp := time.Now()
	for i := 0; i < 1000; i++ {
		epool.Go(func() int {
			defer wg.Done()
			return i
		})
	}
	wg.Wait()
	println(`PROC EACH 1000 tasks:`, time.Since(tp).String())

	t = time.Now()
	epool.Stop()
	println(`DOWN EACH 1000 procs:`, time.Since(t).String())

	//-------------------------------------------------------

	t = time.Now()
	opool := oneof.MakePool(1000)
	println(`UP ONEOF 1000 procs:`, time.Since(t).String())

	wgo := &sync.WaitGroup{}
	wgo.Add(1000)
	tpo := time.Now()
	for i := 0; i < 1000; i++ {
		opool.Go(func() int {
			defer wgo.Done()
			return i
		})
	}
	wgo.Wait()
	println(`PROC ONEOF 1000 tasks:`, time.Since(tpo).String())

	t = time.Now()
	opool.Stop()
	println(`DOWN ONEOF 1000 procs:`, time.Since(t).String())

	wgb := &sync.WaitGroup{}
	wgb.Add(1000)
	tpb := time.Now()
	for i := 0; i < 1000; i++ {
		go func() int {
			defer wgb.Done()
			return i
		}()
	}
	wgb.Wait()
	println(`NATIVE MIX  - UP/DOWN/PROC 1000 tasks:`, time.Since(tpb).String())
}
