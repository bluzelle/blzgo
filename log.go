package bluzelle

import (
	"fmt"
)

// Debugf level formatted messagctx.logger.
func (ctx *Client) Debugf(msg string, v ...interface{}) {
	if ctx.options.Debug {
		ctx.logger.Debug(fmt.Sprintf(msg, v...))
	}
}

// Infof level formatted messagctx.logger.
func (ctx *Client) Infof(msg string, v ...interface{}) {
	if ctx.options.Debug {
		ctx.logger.Info(fmt.Sprintf(msg, v...))
	}
}

// Warnf level formatted messagctx.logger.
func (ctx *Client) Warnf(msg string, v ...interface{}) {
	if ctx.options.Debug {
		ctx.logger.Warn(fmt.Sprintf(msg, v...))
	}
}

// Errorf level formatted messagctx.logger.
func (ctx *Client) Errorf(msg string, v ...interface{}) {
	if ctx.options.Debug {
		ctx.logger.Error(fmt.Sprintf(msg, v...))
	}
}

// Fatalf level formatted messagctx.logger.
func (ctx *Client) Fatalf(msg string, v ...interface{}) {
	if ctx.options.Debug {
		ctx.logger.Fatal(fmt.Sprintf(msg, v...))
	}
}
