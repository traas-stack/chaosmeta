package model

import "fmt"

type DaemonSetObject struct {
	DaemonSetName string
	Namespace     string
}

func (d *DaemonSetObject) GetObjectName() string {
	return fmt.Sprintf("%s%s%s%s%s", "deployment", ObjectNameSplit, d.Namespace, ObjectNameSplit, d.DaemonSetName)
}
