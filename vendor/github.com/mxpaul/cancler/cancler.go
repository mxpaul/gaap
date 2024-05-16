package cancler

import "context"

type Cancler struct {
	ctx        context.Context
	cancelFunc func()
}

func NewCancler(ctx context.Context) *Cancler {
	ctx, cancel := context.WithCancel(ctx)
	return &Cancler{ctx: ctx, cancelFunc: cancel}
}

func (c *Cancler) Ctx() context.Context  { return c.ctx }
func (c *Cancler) Done() <-chan struct{} { return c.ctx.Done() }

func (c *Cancler) Cancel() {
	c.cancelFunc()
}

func (c *Cancler) Caceled() bool {
	select {
	case <-c.ctx.Done():
		return true
	default:
		return false
	}
}
