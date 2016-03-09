# CHANGELOG

## v0.1.2 (March 9, 2016)

- Introduce integration tests for verifying a release.
- Updated to work with recent breaking changes in AWS SDK.
- Re-initialize logger after encountering InvalidSequenceTokenException.

## v0.1.1 (October 13, 2015)

- Drop messages that consist of only whitespace to prevent entire batches from being rejected by Cloudwatch - [8c992d3](https://github.com/bradgignac/logspout-cloudwatch/commit/8c992d358ed89fc078ce2fe38fe6a1ee0caf6ce5)

## v0.1.0 (September 27, 2015)

- Initial Release
