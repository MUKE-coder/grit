package scaffold

// apiGateGo emits internal/authz/gate.go: the policy layer.
//
// This project could already answer "may this user touch widgets at all" from
// the permission catalogue, and "is this their row" from the owner column.
// Neither can say anything conditional on the record: a post editable while it
// is a draft and not after, an order approvable by anyone except its raiser.
// Those were hand-written ifs at the top of a handler, which is where they stop
// being findable and start being forgotten in the second handler that needs
// them.
func apiGateGo() string { return tmpl("api/authz/gate.go") }

// apiGateTestGo is the gate's own test suite, which ships into the project
// because the property it protects is one an upgrade could break for everybody:
// an ability nobody wrote a rule for must allow.
func apiGateTestGo() string { return tmpl("api/authz/gate_test.go") }
