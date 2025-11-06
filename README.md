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

## Dev
```bash
go run index.go
```

## Build

```bash 
go build
```
