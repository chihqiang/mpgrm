package factory

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"log"
	"time"
	"wangzhiqiang/mpgrm/flags"
	"wangzhiqiang/mpgrm/pkg/credential"
	"wangzhiqiang/mpgrm/pkg/gitx"
)

type Git struct {
	ctx context.Context

	workspace      string
	branches, tags []string

	credential       *credential.Credential
	fullName         string
	targetCredential *credential.Credential
	targetFullName   string
}

// NewDoubleCredentialGit creates a new Git instance using provided source and target credentials.
// It initializes the workspace, branches, tags, and parses the full names of source and target repositories.
func NewDoubleCredentialGit(cmd *cli.Command, credential *credential.Credential, targetCredential *credential.Credential) (*Git, error) {
	// Initialize Git instance with context, workspace, branches, tags, and credentials
	rt := &Git{
		ctx:              context.Background(),    // Background context
		workspace:        flags.GetWorkspace(cmd), // Local workspace directory
		branches:         flags.GetBranches(cmd),  // Branches to operate on
		tags:             flags.GetTags(cmd),      // Tags to operate on
		credential:       credential,              // Source repository credential
		targetCredential: targetCredential,        // Target repository credential
	}

	// Parse full name of the source repository (format: owner/repo or owner/repo/subdir)
	fullName, err := credential.GetFullName()
	if err != nil {
		return rt, err
	}
	rt.fullName = fullName

	// Parse full name of the target repository
	targetFullName, err := targetCredential.GetFullName()
	if err != nil {
		return rt, err
	}
	rt.targetFullName = targetFullName

	// Return the initialized Git instance
	return rt, nil
}
func NewCredentialGit(cmd *cli.Command, credential *credential.Credential) (*Git, error) {
	// Initialize Git instance with context, workspace, branches, tags, and credentials
	rt := &Git{
		ctx:        context.Background(),    // Background context
		workspace:  flags.GetWorkspace(cmd), // Local workspace directory
		branches:   flags.GetBranches(cmd),  // Branches to operate on
		tags:       flags.GetTags(cmd),      // Tags to operate on
		credential: credential,              // Source repository credential
	}
	// Parse full name of the source repository (format: owner/repo or owner/repo/subdir)
	fullName, err := credential.GetFullName()
	if err != nil {
		return rt, err
	}
	rt.fullName = fullName
	// Return the initialized Git instance
	return rt, nil
}

// NewCmdDoubleGit creates a new Git instance based on CLI command flags and credentials.
// It initializes the source and target repository information, including authentication and full repository names.
func NewCmdDoubleGit(ctx context.Context, cmd *cli.Command) (*Git, error) {
	// Create a new Git instance with context, workspace, branches, and tags from command flags
	rt := &Git{
		ctx:       ctx,
		workspace: flags.GetWorkspace(cmd), // Local workspace directory for cloning/pushing
		branches:  flags.GetBranches(cmd),  // Branches to operate on
		tags:      flags.GetTags(cmd),      // Tags to operate on
	}

	// Get source repository URL and credentials from CLI flags
	_, cred, _ := flags.GetFormCredential(cmd, false)
	rt.credential = cred // Assign source repository credential
	// Parse the full name of the source repository (format: owner/repo or owner/repo/subdir)
	fullName, err := cred.GetFullName()
	if err != nil {
		return rt, err
	}
	rt.fullName = fullName // Assign parsed full name

	// Get target repository URL and credentials from CLI flags
	_, targetCred, err := flags.GetTargetCredential(cmd)
	if err != nil {
		return rt, err
	}
	rt.targetCredential = targetCred // Assign target repository credential

	// Parse the full name of the target repository (format: owner/repo or owner/repo/subdir)
	targetFullName, err := targetCred.GetFullName()
	if err != nil {
		return rt, err
	}
	rt.targetFullName = targetFullName // Assign parsed target repository full name

	// Return the initialized Git instance
	return rt, nil
}

func (g *Git) getPath() (string, error) {
	return g.credential.GetCategoryNamWorkspace(credential.WorkspaceCategoryGit, g.workspace)
}

func (g *Git) Clone() error {
	start := time.Now()
	migrate := gitx.NewGitMigrate()
	migrate.WithForm(g.credential)
	// 获取 workspace
	workspace, err := g.getPath()
	if err != nil {
		return err
	}
	log.Printf("Start cloning repository: %s to local %s", g.credential.CloneURL, workspace)
	// Clone 仓库并获取实际分支和标签
	_, _, err = migrate.Clone(workspace, []string{}, []string{})
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	elapsed := time.Since(start)
	log.Printf("Clone %s completed in %s", g.credential.CloneURL, elapsed)
	return nil
}

func (g *Git) Push() error {
	log.Printf("Starting git sync from %s to %s", g.credential.CloneURL, g.targetCredential.CloneURL)
	start := time.Now()
	migrate := gitx.NewGitMigrateDouble(g.credential, g.targetCredential)
	// 获取 workspace
	workspace, err := g.getPath()
	if err != nil {
		log.Printf("Failed to get workspace: %v", err)
		return err
	}
	// 获取分支和标签
	branches := g.branches
	tags := g.tags
	log.Printf("Preparing to clone branches: %v, tags: %v", branches, tags)

	// Clone 仓库并获取实际分支和标签
	actualBranches, actualTags, err := migrate.Clone(workspace, branches, tags)
	if err != nil {
		log.Printf("Failed to clone repository: %v", err)
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	log.Printf("Repository cloned successfully. Cloned branches: %v, tags: %v", actualBranches, actualTags)

	// Push 到目标仓库
	if err := migrate.Push(workspace); err != nil {
		log.Printf("Failed to push repository: %v", err)
		return err
	}

	elapsed := time.Since(start)
	log.Printf("Push to target repository completed successfully, total elapsed time: %s", elapsed)
	return nil
}
