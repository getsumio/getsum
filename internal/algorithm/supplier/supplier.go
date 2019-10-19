package supplier

type Supplier interface {
	Run() (string, error)
}

type BaseSupplier struct {
	Algorithm int
}
