package reflectx

import "reflect"

// RPCMethodsOf require receiver value used for net/rpc.Register and
// returns all it RPC methods (detected in same way as net/rpc does).
func RPCMethodsOf(v interface{}) []string {
	return suitableMethods(reflect.TypeOf(v))
}
