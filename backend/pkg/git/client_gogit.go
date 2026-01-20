package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
)

// GoGitClient 基于 go-git 的 Git 客户端
type GoGitClient struct {
	repoURL  string
	auth     *git.CloneOptions
	memStore *memory.Storage
}

// NewGoGitClient 创建 go-git 客户端
func NewGoGitClient(repoURL, token string) *GoGitClient {
	opts := &git.CloneOptions{
		URL: repoURL,
	}

	// 配置 Basic Auth
	if token != "" {
		// 对于 GitHub/GitLab Token，通常作为 Password 使用
		// Username 可以是 oauth2 (GitLab) 或其他，这里使用 "oauth2" 作为通用值
		opts.Auth = &http.BasicAuth{
			Username: "oauth2",
			Password: token,
		}
	}

	return &GoGitClient{
		repoURL:  repoURL,
		auth:     opts,
		memStore: memory.NewStorage(),
	}
}

// GetDiff 获取两次提交之间的差异
// 注意：go-git 获取两次提交的差异需要完整的历史或至少包含这两个提交
// 这里为了简单，我们可能需要 Clone 整个仓库或使用 Shallow Clone
func (c *GoGitClient) GetDiff(base, head string) ([]DiffFile, error) {
	// Clone 仓库到内存
	// 使用 Shallow Clone 可能会有问题，如果 base/head 不在 depth 范围内
	// 这里尝试完整 Clone (对于大仓库会慢)
	r, err := git.Clone(c.memStore, nil, c.auth)
	if err != nil {
		return nil, fmt.Errorf("failed to clone repo: %w", err)
	}

	// 获取 Commit 对象
	headHash := plumbing.NewHash(head)
	baseHash := plumbing.NewHash(base)

	headCommit, err := r.CommitObject(headHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get head commit: %w", err)
	}

	baseCommit, err := r.CommitObject(baseHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get base commit: %w", err)
	}

	// 获取 Tree
	headTree, err := headCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get head tree: %w", err)
	}

	baseTree, err := baseCommit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get base tree: %w", err)
	}

	// 比较
	patch, err := baseTree.Patch(headTree)
	if err != nil {
		return nil, fmt.Errorf("failed to get patch: %w", err)
	}

	return c.convertPatchToDiffFiles(patch), nil
}

// GetSingleCommitDiff 获取单次提交的差异
func (c *GoGitClient) GetSingleCommitDiff(commitSHA string) ([]DiffFile, error) {
	// Clone 仓库
	r, err := git.Clone(c.memStore, nil, c.auth)
	if err != nil {
		return nil, fmt.Errorf("failed to clone repo: %w", err)
	}

	commitHash := plumbing.NewHash(commitSHA)
	commit, err := r.CommitObject(commitHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %w", err)
	}

	// 获取父提交
	var parentTree *object.Tree
	if commit.NumParents() > 0 {
		parent, err := commit.Parent(0)
		if err != nil {
			return nil, fmt.Errorf("failed to get parent commit: %w", err)
		}
		parentTree, err = parent.Tree()
		if err != nil {
			return nil, fmt.Errorf("failed to get parent tree: %w", err)
		}
	} else {
		// 初始提交，比较空 Tree
		parentTree = &object.Tree{}
	}

	currentTree, err := commit.Tree()
	if err != nil {
		return nil, fmt.Errorf("failed to get current tree: %w", err)
	}

	patch, err := parentTree.Patch(currentTree)
	if err != nil {
		return nil, fmt.Errorf("failed to get patch: %w", err)
	}

	return c.convertPatchToDiffFiles(patch), nil
}

// convertPatchToDiffFiles 转换 patch 为通用 DiffFile 格式
func (c *GoGitClient) convertPatchToDiffFiles(patch *object.Patch) []DiffFile {
	var files []DiffFile
	for _, fp := range patch.FilePatches() {
		from, to := fp.Files()
		
		var filename string
		var status string
		
		if from == nil && to != nil {
			filename = to.Path()
			status = "added"
		} else if from != nil && to == nil {
			filename = from.Path()
			status = "deleted"
		} else if from != nil && to != nil {
			filename = to.Path() // Usually same as from
			if from.Path() != to.Path() {
				status = "renamed"
			} else {
				status = "modified"
			}
		}

		chunks := fp.Chunks()
		var diffBuilder strings.Builder
		additions := 0
		deletions := 0

		for _, chunk := range chunks {
			content := chunk.Content()
			op := chunk.Type()

			// Count lines
			lines := strings.Count(content, "\n")
			if len(content) > 0 && content[len(content)-1] != '\n' {
				lines++
			}

			// Build diff string
			linesArr := strings.Split(strings.TrimSuffix(content, "\n"), "\n")

			for _, line := range linesArr {
				if op == diff.Add {
					additions++ // Count again per line to be safe, or just use the block count
					diffBuilder.WriteString("+" + line + "\n")
				} else if op == diff.Delete {
					deletions++
					diffBuilder.WriteString("-" + line + "\n")
				} else {
					diffBuilder.WriteString(" " + line + "\n")
				}
			}
		}

		// Re-calculate stats from chunks logic if needed, but the loop above does it too.
		// However, iterating lines is safer for stats if chunk content has mixed newlines.
		// Let's rely on the loop counters.
		// Reset counters to be accurate based on lines we processed
		// Actually, let's just use the chunk stats logic I had before but per line.

		files = append(files, DiffFile{
			Filename:  filename,
			Status:    status,
			Patch:     diffBuilder.String(),
			Additions: additions,
			Deletions: deletions,
			Changes:   additions + deletions,
		})
	}
	return files
}
