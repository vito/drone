package handler

import (
	"fmt"
	"net/http"

	"github.com/drone/drone/pkg/channel"
	"github.com/drone/drone/pkg/database"
	. "github.com/drone/drone/pkg/model"
	"github.com/drone/drone/pkg/queue"
	"github.com/drone/go-github/github"
	"github.com/drone/drone/pkg/build/script"
)

// Display a specific Commit.
func CommitShow(w http.ResponseWriter, r *http.Request, u *User, repo *Repo) error {
	hash := r.FormValue(":commit")
	labl := r.FormValue(":label")

	// get the commit from the database
	commit, err := database.GetCommitHash(hash, repo.ID)
	if err != nil {
		return err
	}

	// get the builds from the database. a commit can have
	// multiple sub-builds (or matrix builds)
	builds, err := database.ListBuilds(commit.ID)
	if err != nil {
		return err
	}

	admin, err := database.IsRepoAdmin(u, repo)
	if err != nil {
		return err
	}

	data := struct {
		User    *User
		Repo    *Repo
		Commit  *Commit
		Build   *Build
		Builds  []*Build
		Token   string
		IsAdmin bool
	}{u, repo, commit, builds[0], builds, "", admin}

	// get the specific build requested by the user. instead
	// of a database round trip, we can just loop through the
	// list and extract the requested build.
	for _, b := range builds {
		if b.Slug == labl {
			data.Build = b
			break
		}
	}

	// generate a token to connect with the websocket
	// handler and stream output, if the build is running.
	data.Token = channel.Token(fmt.Sprintf(
		"%s/%s/%s/commit/%s/builds/%s", repo.Host, repo.Owner, repo.Name, commit.Hash, builds[0].Slug))

	// render the repository template.
	return RenderTemplate(w, "repo_commit.html", &data)
}

type CommitRebuildHandler struct {
	queue *queue.Queue
}

func NewCommitRebuildHandler(queue *queue.Queue) *CommitRebuildHandler {
	return &CommitRebuildHandler{
		queue: queue,
	}
}

// Rebuild a commit
func (h *CommitRebuildHandler) CommitRebuild(w http.ResponseWriter, r *http.Request, u *User, repo *Repo) error {
	hash := r.FormValue(":commit")
	labl := r.FormValue(":label")
	host := r.FormValue(":host")

	// get the commit from the database
	commit, err := database.GetCommitHash(hash, repo.ID)
	if err != nil {
		return err
	}

	// get the builds from the database. a commit can have
	// multiple sub-builds (or matrix builds)
	builds, err := database.ListBuilds(commit.ID)
	if err != nil {
		return err
	}

	build := builds[0]

	// get the specific build requested by the user. instead
	// of a database round trip, we can just loop through the
	// list and extract the requested build.
	for _, b := range builds {
		if b.Slug == labl {
			build = b
			break
		}
	}

	// get the github settings from the database
	settings := database.SettingsMust()

	// get the drone.yml file from GitHub
	client := github.New(u.GithubToken)
	client.ApiUrl = settings.GitHubApiUrl

	content, err := client.Contents.FindRef(repo.Owner, repo.Name, ".drone.yml", commit.Hash) // TODO should this really be the hash??
	if err != nil {
		println(err.Error())
		RenderText(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}

	raw, err := content.DecodeContent()
	if err != nil {
		msg := "Could not decode the yaml from GitHub.	Check that your .drone.yml is a valid yaml file.\n"
		if err := saveFailedBuild(commit, msg); err != nil {
			return RenderText(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return RenderText(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	// parse the build script
	buildscript, err := script.ParseBuild(raw, repo.Params)
	if err != nil {
		// TODO if the YAML is invalid we should create a commit record
		// with an ERROR status so that the user knows why a build wasn't
		// triggered in the system
		RenderText(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return nil
	}

	h.queue.Add(&queue.BuildTask{Repo: repo, Commit: commit, Build: build, Script: buildscript})

	if labl != "" {
		http.Redirect(w, r, fmt.Sprintf("/%s/%s/%s/commit/%s/build/%s", host, repo.Owner, repo.Name, hash, labl), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/%s/%s/%s/commit/%s", host, repo.Owner, repo.Name, hash), http.StatusSeeOther)
	}
	return nil
}
