package main

import (
	"github.com/proxy-wasm/proxy-wasm-go-sdk/proxywasm"
	"github.com/proxy-wasm/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {}

func init() {
	proxywasm.SetPluginContext(func(contextID uint32) types.PluginContext {
		return &pluginContext{}
	})
}

type pluginContext struct {
	types.DefaultPluginContext
}

func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpHandler{}
}

type httpHandler struct {
	types.DefaultHttpContext
}

func (ctx *httpHandler) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	// Coin flip: 50/50 chance of 200 or 500
	// Use hash of ALL headers for randomness
	hashInput := ""
	
	// Get all header pairs and concatenate them
	headerPairs, _ := proxywasm.GetHttpRequestHeaders()
	for _, pair := range headerPairs {
		hashInput += pair[0] + ":" + pair[1] + ";"
	}
	
	// Simple hash function
	hash := uint32(0)
	for _, b := range []byte(hashInput) {
		hash = hash*31 + uint32(b)
	}
	
	coinFlip := int(hash % 2)
	
	if coinFlip == 0 {
		// Heads: return 200
		proxywasm.LogInfo("WASM filter: coin flip = heads, returning 200")
		if err := proxywasm.SendHttpResponse(200, nil, []byte("WASM filter: heads = 200\n"), -1); err != nil {
			proxywasm.LogError("Failed to send response: " + err.Error())
			return types.ActionContinue
		}
	} else {
		// Tails: return 500
		proxywasm.LogInfo("WASM filter: coin flip = tails, returning 500")
		if err := proxywasm.SendHttpResponse(500, nil, []byte("WASM filter: tails = 500\n"), -1); err != nil {
			proxywasm.LogError("Failed to send response: " + err.Error())
			return types.ActionContinue
		}
	}
	
	return types.ActionPause
}