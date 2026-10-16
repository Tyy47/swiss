package shortcuts

import (
	"os/exec"
	"swiss/utils"

	"github.com/Tyy47/clibox/argbin"
)

// gitAdd creates a exec.Command and run's "git add .". Returns an exec error if unable to add files.
func gitAdd() error {

	// Create the add command
	addCmd := exec.Command("git", "add", ".")

	// Run the add command
	if err := addCmd.Run(); err != nil {
		return err
	}

	return nil
}

// gitCommit creates a commit based on added files with a provided message
func gitCommit(commitMessage string) error {

	// Create commit command
	commitCmd := exec.Command("git", "commit", "-m", commitMessage)

	// Run git commit command and return the potential error
	if err := commitCmd.Run(); err != nil {
		return err
	}

	return nil
}

// gitPush pushes any commits created by get to the remote origin url
func gitPush() error {

	// Create the push command
	pushCmd := exec.Command("git", "push")

	utils.ToggleOutputForCMD(pushCmd)

	// Execute the push command
	if err := pushCmd.Run(); err != nil {
		return err
	}

	return nil
}

// gitPull pulls down any changes from the remote repository
func gitPull() error {

	// Create the pull command
	pullCmd := exec.Command("git", "pull")

	// Execute the pull command
	if err := pullCmd.Run(); err != nil {
		return err
	}

	return nil
}

// gitSync fetches data from the remote repository to update the local one, then displays the current status of the branch.
func gitSync() error {

	// Create the fetch command
	fetchCmd := exec.Command("git", "fetch")

	// Run the fetch command
	if err := fetchCmd.Run(); err != nil {
		return err
	}

	// Create the status command
	statusCmd := exec.Command("git", "status")

	utils.ToggleOutputForCMD(statusCmd)

	// Run the status command
	if err := statusCmd.Run(); err != nil {
		return err
	}

	return nil
}

// gitCommitCommand creates the "commit" command. commit is a compressed and shorthand form of commiting to a git repository.
func gitCommitCommand() *argbin.Command {
	return &argbin.Command{
		Name:       "commit",
		HelpMenu:   "commit shortcut",
		TakesValue: true,
		Execute: func(cmdCtx *argbin.CommandContext) error {

			// Runs git add command
			if err := gitAdd(); err != nil {
				return err
			}

			// Run commit command
			if err := gitCommit(cmdCtx.ParsedValue); err != nil {
				return err
			}

			// Success message
			utils.Output.Success("git commit created.")

			return nil
		},
	}
}

// gitPushCommand creates the "push" command. push is a compressed shorthand form of adding files, commiting, and pushing to a git repository.
func gitPushCommand() *argbin.Command {
	return &argbin.Command{
		Name:       "push",
		HelpMenu:   "push shortcut",
		TakesValue: true,
		Execute: func(cmdCtx *argbin.CommandContext) error {

			// Runs the git add command
			if err := gitAdd(); err != nil {
				return err
			}

			utils.Output.Success("Added files to commit.")

			// Runs git commit -m ctx.ParsedValue
			if err := gitCommit(cmdCtx.ParsedValue); err != nil {
				return err
			}

			utils.Output.Success("Commit created with message.")

			// Runs the git push command
			if err := gitPush(); err != nil {
				return err
			}

			utils.Output.Success("Pushed commit to repository.")

			return nil
		},
	}
}

// gitSyncCommand creates the "sync" command. sync updates the local repository with changes from the remote repo.
func gitSyncCommand() *argbin.Command {
	return &argbin.Command{
		Name:     "sync",
		HelpMenu: "sync shortcut",
		Execute: func(cmdCtx *argbin.CommandContext) error {

			// Runs git fetch & git status
			if err := gitSync(); err != nil {
				return err
			}

			// Checks if pull is toggled on,
			// If so, it will pull the changes from the remote repo
			pull, _ := cmdCtx.Values["pull"].(bool)
			if pull {
				if err := gitPull(); err != nil {
					return err
				}
			}

			utils.Output.Success("Local repository updated.")
			return nil
		},
		Flags: argbin.Flags{
			"-p": {
				Execute: func(cmdCtx *argbin.CommandContext) error {
					cmdCtx.Values["pull"] = true
					return nil
				},
			},
			"--pull": {
				Execute: func(cmdCtx *argbin.CommandContext) error {
					cmdCtx.Values["pull"] = true
					return nil
				},
			},
		},
	}
}

func ShortcutCommand() *argbin.Command {
	return &argbin.Command{
		Name:            "sc",
		AdditionalNames: []string{"shortcut"},
		Subcommands: []*argbin.Command{
			gitCommitCommand(),
			gitPushCommand(),
			gitSyncCommand(),
		},
		Flags: argbin.Flags{
			"-h":     utils.HelpFlag(),
			"--help": utils.HelpFlag(),
		},
		HelpMenu: `
╭───────────────────  Swiss  ────────────────────╮
│                                                │
│       The army knife of CLI applications       │
│                                                │
╰────────────────────────────────────────────────╯
Shortcut module - Commands that are multiple commands into one.

Commands:
	commit <message : string>: Adds all changed files to commit with a message.
	push [message : string]: Adds all files, commits changes with a message, then pushes to your repository.
	sync: Fetch's all changes to the repository and prints a status message with changes to the repository.

Flags:
	-h --help: Opens the help menu.`,
	}
}
