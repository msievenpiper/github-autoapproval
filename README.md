# Github Auto Approval 

This repo is just a basic go package used to automatically approve different packages for work. It requires you to be authenticated using githubs `gh` tool which you can install using `brew install gh` if you have homebrew installed.

To trigger auth if you haven't used done so already use the `gh auth login` command and it will prompt you to auth via the command line through githubs website.

## Usage 
You can use the project to specify repositories and what branch to target, this will then perform either a mock approval or direct auto approval. This is designed for automated workflows where the output is generated from a script or github action. Therefore, this shouldn't be used generally.

```bash
./github-autoapproval --branch=example msievenpiper/github-autoapproval
```

or testing using 

```bash
./github-autoapproval --branch=example --probe msievenpiper/github-autoapproval
```

or multiple repos
```bash
./github-autoapproval --branch=example msievenpiper/github-autoapproval msievenpiper/example
```

### Auto merging

```bash
./github-autoapproval --branch=example --merge msievenpiper/github-autoapproval
```

## Config files

Instead of passing flags every time, you can store settings in a JSON config file and load it with `--config`.

### Creating a config

Run the interactive wizard:

```bash
./github-autoapproval --init
```

You'll be prompted for a name (or path), branch, repos, probe mode, and auto-merge. Config files are saved to the `configs/` folder by default (gitignored, so credentials/repo lists stay local).

### Using a config

```bash
./github-autoapproval --config=configs/my-project.json
```

CLI flags always override config values when both are provided:

```bash
# Use config but override the branch
./github-autoapproval --config=configs/my-project.json --branch=other-branch
```

### Config file format

```json
{
  "branch": "my-feature-branch",
  "repos": ["owner/repo-a", "owner/repo-b"],
  "probe": false,
  "merge": true
}
```

## Dev
```bash
go run index.go
```

## Build

```bash 
go build
```
