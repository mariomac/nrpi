package plugin

type Plugin interface {
	DoSomething(<-chan interface{})
}
