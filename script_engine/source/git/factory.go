// Package git provides Git repository as a script source.
package git

import (
	"context"
	"fmt"

	"github.com/tx7do/go-scripts/source"
	gitSource "github.com/tx7do/go-scripts/source/git"

	v1 "github.com/tx7do/go-wind-bootstrap/conf/gen/go/bootstrap/v1"
	scriptEngine "github.com/tx7do/go-wind-bootstrap/script_engine"
)

func init() {
	scriptEngine.MustRegisterSourceFactory(v1.Script_Source_GIT, NewSource)
}

// NewSource 根据配置创建 Git 脚本来源。
func NewSource(cfg *v1.Script_Source) (source.Reader, error) {
	if cfg == nil {
		return nil, fmt.Errorf("git source: config is nil")
	}

	opts := cfg.GetGitOptions()
	if opts == nil {
		return nil, fmt.Errorf("git source: git_options is required")
	}

	url := opts.GetUrl()
	if url == "" {
		return nil, fmt.Errorf("git source: url is required")
	}

	var gitOpts []gitSource.Option
	gitOpts = append(gitOpts, gitSource.WithRepoURL(url))

	if branch := opts.GetBranch(); branch != "" {
		gitOpts = append(gitOpts, gitSource.WithBranch(branch))
	}
	if subdir := opts.GetSubdir(); subdir != "" {
		gitOpts = append(gitOpts, gitSource.WithPrefix(subdir))
	}

	return gitSource.New(context.Background(), gitOpts...)
}
