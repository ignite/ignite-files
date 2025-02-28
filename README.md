## Static files for ignite

This repository is used to store large binary files separately from the main project repository to keep it clean and
manageable.

### Git LFS

To handle large files efficiently, we use **Git Large File Storage (Git LFS)**.

##### Setting Up Git LFS

To work with this repository, you need to install Git LFS on your local machine.

Follow the installation instructions from GitHub:

- [GitHub Docs: Installing Git LFS](https://docs.github.com/en/repositories/working-with-files/managing-large-files/installing-git-large-file-storage)

Alternatively, you can install it using a package manager:

- **Mac (Homebrew)**: `brew install git-lfs`
- **Linux**: Use your package manager (e.g., `sudo apt install git-lfs`)
- **Windows**: Use [Git for Windows](https://gitforwindows.org/), which includes Git LFS

More details
here: [git-lfs/git-lfs#3041 (comment)](https://github.com/git-lfs/git-lfs/issues/3041#issuecomment-533879953)

##### Pulling with LFS

After cloning, ensure all LFS files are fetched:

```sh
git lfs pull
```

##### Handling Remote URL Issues

Some users have reported issues when using SSH remotes. If you encounter errors like:

```
batch request: missing protocol: "<unknown>"
```

Try using an HTTPS remote instead of SSH:

```sh
git remote set-url origin https://github.com/ignite/ignite-files.git
```

## Further Reading

- [GitHub Docs: Git LFS](https://docs.github.com/en/repositories/working-with-files/managing-large-files/collaboration-with-git-large-file-storage)
