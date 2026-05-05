# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Clone progress indicator now always shows `Cloning... done` on a single line
  regardless of how fast the operation completes

## [2.4.0] - 2026-05-04

### Added

- `--verbose` flag to print additional diagnostic output to stderr throughout
  the clone flow (repository info, fork decisions, remote setup)
- Print the clone destination path to stdout at the end of every clone operation for easy piping

### Changed

- Git transfer progress is now hidden by default; use `--verbose` to stream it

## [2.3.0] - 2026-05-04

### Changed

- When forking a repository, gh-get now clones the original into the canonical
  location, renames `origin` to `upstream`, and adds the fork as `origin`

## [2.2.0] - 2026-03-30

### Added

- Clone progress is now streamed to the terminal in real time using go-git,
  resolving the hang on large repositories (#11)
- The repository argument now accepts any GitHub URL, including URLs pointing
  to branches, tags, commits, files, pull requests, and issues
  (e.g. `https://github.com/owner/repo/tree/my-branch`) (#15)
- `--fork` flag to fork a repository before cloning. When omitted and the
  user does not have write access, gh-get prompts whether to fork. If the
  repository does not allow forking, the original is cloned with a warning (#6)

### Changed

- Upgraded `go-gh` from v1 to v2 and Go toolchain from 1.21 to 1.25

### Fixed

- Corrected shebang line in integration test script

## [2.1.0] - 2024-12-13

### Added

- Repositories can now be specified as full HTTPS URLs
  (e.g. `https://github.com/owner/repo` or `https://github.com/owner/repo.git`)
- Integration tests running against the real GitHub CLI in a Docker container
- `GH_GET_ROOT` environment variable to set an absolute path as the clone root

### Changed

- The default clone root folder was renamed from `src` to `github`

## [2.0.0] - 2024-11-27

### Changed

- Rewrote the extension in Go, replacing the original shell script
- The extension is now distributed as a precompiled binary for all major platforms
- Set up CI and automated release workflow

## [1.0.0] - 2022-07-16

### Added

- Initial implementation as a shell script inspired by [ghq](https://github.com/x-motemen/ghq)
- Clones repositories into `$HOME/github/<owner>/<repo>` by default
- Configurable clone root via `GH_GET_FOLDER` and `GH_GET_ROOT` environment variables

[Unreleased]: https://github.com/britter/gh-get/compare/v2.4.0...HEAD
[2.4.0]: https://github.com/britter/gh-get/compare/v2.3.0...v2.4.0
[2.3.0]: https://github.com/britter/gh-get/compare/v2.2.0...v2.3.0
[2.2.0]: https://github.com/britter/gh-get/compare/v2.1.0...v2.2.0
[2.1.0]: https://github.com/britter/gh-get/compare/v2.0.0...v2.1.0
[2.0.0]: https://github.com/britter/gh-get/compare/v1.0.0...v2.0.0
[1.0.0]: https://github.com/britter/gh-get/releases/tag/v1.0.0
