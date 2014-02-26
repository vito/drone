package repo

import (
	"testing"
)

func TestIsRemote(t *testing.T) {
	repos := []struct {
		path   string
		remote bool
	}{
		{"git://github.com/foo/far", true},
		{"git://github.com/foo/far.git", true},
		{"git@github.com:foo/far", true},
		{"git@github.com:foo/far.git", true},
		{"http://github.com/foo/far.git", true},
		{"https://github.com/foo/far.git", true},
		{"ssh://baz.com/foo/far.git", true},
		{"/var/lib/src", false},
		{"/home/ubuntu/src", false},
		{"src", false},
	}

	for _, r := range repos {
		repo := Repo{Path: r.path}
		if remote := repo.IsRemote(); remote != r.remote {
			t.Errorf("IsRemote %s was %v, expected %v", r.path, remote, r.remote)
		}
	}
}

func TestIsGit(t *testing.T) {
	repos := []struct {
		path   string
		remote bool
	}{
		{"git://github.com/foo/far", true},
		{"git://github.com/foo/far.git", true},
		{"git@github.com:foo/far", true},
		{"git@github.com:foo/far.git", true},
		{"http://github.com/foo/far.git", true},
		{"https://github.com/foo/far.git", true},
		{"ssh://baz.com/foo/far.git", true},
		{"svn://gcc.gnu.org/svn/gcc/branches/gccgo", false},
		{"https://code.google.com/p/go", false},
	}

	for _, r := range repos {
		repo := Repo{Path: r.path}
		if remote := repo.IsGit(); remote != r.remote {
			t.Errorf("IsGit %s was %v, expected %v", r.path, remote, r.remote)
		}
	}
}

func TestShouldRunPrivileged(t *testing.T) {
	if !(Repo{Privileged: true, PR: ""}.ShouldRunPrivileged()) {
		t.Errorf("ShouldRunPrivileged should be true for repos with no pull requests")
	}

	if (Repo{Privileged: true, PR: "foo"}.ShouldRunPrivileged()) {
		t.Errorf("ShouldRunPrivileged should be false for repos with a pull request")
	}

	if (Repo{Privileged: false, PR: ""}.ShouldRunPrivileged()) {
		t.Errorf("ShouldRunPrivileged should be false for repos with it disabled")
	}

	if (Repo{Privileged: false, PR: "foo"}.ShouldRunPrivileged()) {
		t.Errorf("ShouldRunPrivileged should be false for repos with it disabled")
	}
}
