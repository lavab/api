# Changelog

All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## [Unreleased][unreleased]
### Changed
  - Removed .drone and .travis

## [2.0.2] - 2015-05-19
### Added
 - Added a check whether an address mapping is used in the username
   reservation.
 - SockJS API client for headless serverside client development.
 - Onboarding emails that introduce users to the service.
 - Multiple identity support (also known as email aliases).

### Changed
 - Moved from traditional `go get`-based flow to dependency vendoring
   using [godep](https://github.com/tools/godep).
 - Disabled most of the log output to make it easier to analyze.
 - Matching for 10k most used passwords replaced with a bloom filter
   containing 17.5m leaked passwords from various hacks.

### Fixed
 - Cursor leakage all over the `db` package.
 - thread.update changing date_modified field of the model, which
   resulted in invalid ordering of the emails in the web client.
   Emails in "spam" being shown as unread on the sidebar (new label
   fetching query).
 - Incorrect difference checker in thread.update.

## [2.0.1] - 2015-04-15
### Added
 - Address mapping table for account's name-to-id lookups.
 - Username length check during registration.

### Changed
 - New index creation code (multiple compound and multi indexes).

### Fixed
 - Lack of Message-ID header causing Lavaboom emails to be flagged as
   spam.

## 2.0.0 - 2015-04-02
### Added
 - Initial release of Lavaboom API 2.0

[unreleased]: https://github.com/lavab/api/compare/2.0.2...HEAD
[2.0.2]: https://github.com/lavab/api/compare/2.0.2...2.0.1
[2.0.1]: https://github.com/lavab/api/compare/2.0.1...0.2.0