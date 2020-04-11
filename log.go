package bluzelle

import (
	"fmt"
)

// Debugf level formatted messagctx.Logger.
func (ctx *Client) Debugf(msg string, v ...interface{}) {
	if ctx.Options.Debug {
		ctx.Logger.Debug(fmt.Sprintf(msg, v...))
	}
}

// Infof level formatted messagctx.Logger.
func (ctx *Client) Infof(msg string, v ...interface{}) {
	if ctx.Options.Debug {
		ctx.Logger.Info(fmt.Sprintf(msg, v...))
	}
}

// Warnf level formatted messagctx.Logger.
func (ctx *Client) Warnf(msg string, v ...interface{}) {
	if ctx.Options.Debug {
		ctx.Logger.Warn(fmt.Sprintf(msg, v...))
	}
}

// Errorf level formatted messagctx.Logger.
func (ctx *Client) Errorf(msg string, v ...interface{}) {
	if ctx.Options.Debug {
		ctx.Logger.Error(fmt.Sprintf(msg, v...))
	}
}

// Fatalf level formatted messagctx.Logger.
func (ctx *Client) Fatalf(msg string, v ...interface{}) {
	if ctx.Options.Debug {
		ctx.Logger.Fatal(fmt.Sprintf(msg, v...))
	}
}
