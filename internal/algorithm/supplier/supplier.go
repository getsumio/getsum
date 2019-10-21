package supplier

type Supplier interface {
	Run() error
	Status() *Status
	Terminate() error
}

type BaseSupplier struct {
	Algorithm string
	status    *Status
}

type Status struct {
	Status string
	Value  string
}
