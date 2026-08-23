package parser

type Node interface {
	Nodium() int
}

// Types

type Identifier struct {

}

func (X Identifier) Nodium { return 0 }
