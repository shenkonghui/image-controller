package controller

import (
	"github.com/shenkonghui/image-controller/pkg/controller/imageconfig"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, imageconfig.Add)
}
