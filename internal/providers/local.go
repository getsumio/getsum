package providers

type LocalProvider struct {
	BaseProvider
}

func (l *LocalProvider) Run() (string, error) {
	l.Status = "CALCULATING"
	defer complete(l)
	return l.Supplier.Run()
}

func (l *LocalProvider) Data() *BaseProvider {
	return &l.BaseProvider
}

func (l *LocalProvider) Close() {

}

func complete(l *LocalProvider) {
	l.Status = "COMPLETED"
}
