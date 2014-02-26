package dockerfile

import (
	"bytes"
	"fmt"
)

type DockerfileWriter interface {
	WriteAdd(from, to string)
	WriteFrom(from string)
	WriteRun(cmd string)
	WriteUser(user string)
	WriteEnv(key, val string)
	WriteWorkdir(workdir string)
	WriteEntrypoint(entrypoint string)

	Bytes() []byte
}

type Dockerfile struct {
	bytes.Buffer
}

func New() *Dockerfile {
	return &Dockerfile{}
}

func (d *Dockerfile) WriteAdd(from, to string) {
	d.WriteString(fmt.Sprintf("ADD %s %s\n", from, to))
}

func (d *Dockerfile) WriteFrom(from string) {
	d.WriteString(fmt.Sprintf("FROM %s\n", from))
}

func (d *Dockerfile) WriteRun(cmd string) {
	d.WriteString(fmt.Sprintf("RUN %s\n", cmd))
}

func (d *Dockerfile) WriteUser(user string) {
	d.WriteString(fmt.Sprintf("USER %s\n", user))
}

func (d *Dockerfile) WriteEnv(key, val string) {
	d.WriteString(fmt.Sprintf("ENV %s %s\n", key, val))
}

func (d *Dockerfile) WriteWorkdir(workdir string) {
	d.WriteString(fmt.Sprintf("WORKDIR %s\n", workdir))
}

func (d *Dockerfile) WriteEntrypoint(entrypoint string) {
	d.WriteString(fmt.Sprintf("ENTRYPOINT %s\n", entrypoint))
}
