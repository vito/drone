package fakedockerfile

import (
	"reflect"
)

type FakeDockerfile struct {
	WrittenEntries []DockerfileEntry
}

type DockerfileEntry struct {
	Action string
	Args   []string
}

func New() *FakeDockerfile {
	return &FakeDockerfile{}
}

func (d *FakeDockerfile) WriteAdd(from, to string) {
	d.appendEntry("ADD", from, to)
}

func (d *FakeDockerfile) WriteFrom(from string) {
	d.appendEntry("FROM", from)
}

func (d *FakeDockerfile) WriteRun(cmd string) {
	d.appendEntry("RUN", cmd)
}

func (d *FakeDockerfile) WriteUser(user string) {
	d.appendEntry("USER", user)
}

func (d *FakeDockerfile) WriteEnv(key, val string) {
	d.appendEntry("ENV", key, val)
}

func (d *FakeDockerfile) WriteWorkdir(workdir string) {
	d.appendEntry("WORKDIR", workdir)
}

func (d *FakeDockerfile) WriteEntrypoint(entrypoint string) {
	d.appendEntry("ENTRYPOINT", entrypoint)
}

func (d *FakeDockerfile) Bytes() []byte {
	return []byte("")
}

func (d *FakeDockerfile) IsWritten(action string, args ...string) bool {
	for _, e := range d.WrittenEntries {
		if e.Action == action && reflect.DeepEqual(e.Args, args) {
			return true
		}
	}

	return false
}

func (d *FakeDockerfile) appendEntry(action string, args ...string) {
	d.WrittenEntries = append(d.WrittenEntries, DockerfileEntry{
		Action: action,
		Args:   args,
	})
}
