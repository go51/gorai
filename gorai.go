package gorai

type gorai struct {}

var goraiInstance *gorai = nil

func Load() *gorai {
	if goraiInstance != nil {
		return goraiInstance
	}

	goraiInstance = &gorai{}

	return goraiInstance
}