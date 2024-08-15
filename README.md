# branchswapper

Stash git branch names for later use.


## Motivation (Dramatization)

Let's assume you have been working on a feature for the last couple of days and you have a feature branch for that..
```shell
➜  project $ git status
On branch JIRA-TICKET-feature-with-a-very-long-name
```

Suddenly, a wild bug/issue/support-ticket appears and you have to fix it `soon`. Typically, you switch to the branch where the problem popped up, branch out from it and start doing whatever you need to.
```shell
# Stash your current branch and checkout an existing branch
➜  project $ branchswapper release/this-is-yet-not-a-stable-version

➜  project $ git status
On branch release/this-is-yet-not-a-stable-version
```

Some hours later, after closing that horrible bug you come back to the original feature branch you were working on. But, oh Lord, what was the name of the branch? Chill, we got your back.. 
```shell
# List your stashed branches
➜  project $ branchswapper -ls
0: JIRA-TICKET-feature-with-a-very-long-name

# Swap to the desired branch using the index assigned to it
➜  project $ branchswapper -i 0

➜  project $ git status
On branch JIRA-TICKET-feature-with-a-very-long-name
```


## Install

```shell
➜  project $ go install github.com/grevych/branchswapper@latest
```


## Usage

```shell
➜  project $ branchswapper -i 0
NAME:
   branchswap - Stash git branches for later use

USAGE:
   branchswap [global options] command [command options]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --list, --ls             List stashed branches (default: false)
   --index value, -i value  Swap branch by index (default: -1)
   --help, -h               show help
```
