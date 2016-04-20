# Change Log

## [v1.12](https://github.com/DataDog/kvexpress/tree/v1.12) (2016-04-20)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v1.11...v1.12)

**Fixed bugs:**

- If you're using compression... [\#92](https://github.com/DataDog/kvexpress/issues/92)
- When using the compression option the lines metric is fixed at 1. [\#91](https://github.com/DataDog/kvexpress/issues/91)

## [v1.11](https://github.com/DataDog/kvexpress/tree/v1.11) (2016-03-14)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v1.10...v1.11)

**Fixed bugs:**

- Writing to a file should be atomic [\#88](https://github.com/DataDog/kvexpress/issues/88)
- Atomic rename instead of truncate and write in place. [\#89](https://github.com/DataDog/kvexpress/pull/89) ([darron](https://github.com/darron))

## [v1.10](https://github.com/DataDog/kvexpress/tree/v1.10) (2016-03-09)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v1.9...v1.10)

**Implemented enhancements:**

- The SHA256 of a blank key is always the same. [\#87](https://github.com/DataDog/kvexpress/issues/87)
- kvexpress can't currently make directories [\#86](https://github.com/DataDog/kvexpress/issues/86)
- Makefile should insert platform rather than assuming OSX. [\#77](https://github.com/DataDog/kvexpress/issues/77)
- Fix that ugliness with `fmt` width. [\#75](https://github.com/DataDog/kvexpress/issues/75)
- Fix the wercker autotest. [\#72](https://github.com/DataDog/kvexpress/issues/72)

**Fixed bugs:**

- The SHA256 of a blank key is always the same. [\#87](https://github.com/DataDog/kvexpress/issues/87)
- kvexpress can't currently make directories [\#86](https://github.com/DataDog/kvexpress/issues/86)

**Closed issues:**

- Update the docs. [\#71](https://github.com/DataDog/kvexpress/issues/71)

## [v1.9](https://github.com/DataDog/kvexpress/tree/v1.9) (2016-01-14)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v1.8...v1.9)

**Fixed bugs:**

- Doesn't force the file mode if it differs [\#78](https://github.com/DataDog/kvexpress/issues/78)
- 1.3 Bug [\#51](https://github.com/DataDog/kvexpress/issues/51)

## [v1.8](https://github.com/DataDog/kvexpress/tree/v1.8) (2016-01-06)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v1.7...v1.8)

**Implemented enhancements:**

- Add ability to stop `out` process on a single node. [\#54](https://github.com/DataDog/kvexpress/issues/54)
- Be able to deploy code to a subset of nodes. [\#53](https://github.com/DataDog/kvexpress/issues/53)

**Fixed bugs:**

- Silly bug caused by me never just running 'kvexpress' by itself. [\#73](https://github.com/DataDog/kvexpress/issues/73)

**Merged pull requests:**

- Refactor Consul connection. All unit and integration tests pass. [\#74](https://github.com/DataDog/kvexpress/pull/74) ([darron](https://github.com/darron))

## [v1.7](https://github.com/DataDog/kvexpress/tree/v1.7) (2015-12-31)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v1.6...v1.7)

**Implemented enhancements:**

- Does AutoEnable\(\) work? [\#67](https://github.com/DataDog/kvexpress/issues/67)
- When there's a Panic worthy event. [\#66](https://github.com/DataDog/kvexpress/issues/66)
- Add Reason to written `.locked` file. [\#65](https://github.com/DataDog/kvexpress/issues/65)
- Load Dogstatsd if dd-agent is installed. [\#64](https://github.com/DataDog/kvexpress/issues/64)
- Check for leading `/` on lock and unlock operations. [\#63](https://github.com/DataDog/kvexpress/issues/63)
- Add event for not long enough on input. [\#62](https://github.com/DataDog/kvexpress/issues/62)

**Closed issues:**

- Is there a better way to deal with -v? [\#45](https://github.com/DataDog/kvexpress/issues/45)

## [v1.6](https://github.com/DataDog/kvexpress/tree/v1.6) (2015-12-11)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v1.5...v1.6)

**Implemented enhancements:**

- Redirect logger to STDOUT. [\#61](https://github.com/DataDog/kvexpress/issues/61)

**Fixed bugs:**

- If you add a leading slash when you specify the prefix it can't write and panics. [\#60](https://github.com/DataDog/kvexpress/issues/60)

**Closed issues:**

- Remove passing Prefix around. [\#59](https://github.com/DataDog/kvexpress/issues/59)
- Remove passing of DogStatsd global. [\#58](https://github.com/DataDog/kvexpress/issues/58)

## [v1.5](https://github.com/DataDog/kvexpress/tree/v1.5) (2015-12-04)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v1.4...v1.5)

**Implemented enhancements:**

- Can we have additional per node or per role configuration options? [\#55](https://github.com/DataDog/kvexpress/issues/55)
- Add metric to throw at Datadog when a `panic` worthy event happens. [\#52](https://github.com/DataDog/kvexpress/issues/52)
- Add compression? [\#21](https://github.com/DataDog/kvexpress/issues/21)

**Closed issues:**

- Remove the additional `kvexpress` Golang package. [\#56](https://github.com/DataDog/kvexpress/issues/56)
- Don't like how the permissions are output through syslog on OS X [\#1](https://github.com/DataDog/kvexpress/issues/1)

**Merged pull requests:**

- Refactor how we generate Direction in the logs. So much nicer. [\#57](https://github.com/DataDog/kvexpress/pull/57) ([darron](https://github.com/darron))

## [v1.4](https://github.com/DataDog/kvexpress/tree/v1.4) (2015-11-23)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v1.3...v1.4)

## [v1.3](https://github.com/DataDog/kvexpress/tree/v1.3) (2015-11-16)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v1.2...v1.3)

**Implemented enhancements:**

- Add `chown` command [\#49](https://github.com/DataDog/kvexpress/issues/49)
- Strip the top 3 lines from the diff. [\#48](https://github.com/DataDog/kvexpress/issues/48)
- Add `url` command. [\#47](https://github.com/DataDog/kvexpress/issues/47)
- Do we add minimum length as the third key? [\#39](https://github.com/DataDog/kvexpress/issues/39)
- Make diff output for `in` command better. [\#36](https://github.com/DataDog/kvexpress/issues/36)

**Fixed bugs:**

- Add `chown` command [\#49](https://github.com/DataDog/kvexpress/issues/49)

## [v1.2](https://github.com/DataDog/kvexpress/tree/v1.2) (2015-11-10)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v1.1...v1.2)

**Fixed bugs:**

- Something strange happening with CONSUL\_TOKEN ENV variable. [\#43](https://github.com/DataDog/kvexpress/issues/43)

## [v1.1](https://github.com/DataDog/kvexpress/tree/v1.1) (2015-11-08)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v1.0...v1.1)

**Implemented enhancements:**

- Would be nice to have a config file. [\#24](https://github.com/DataDog/kvexpress/issues/24)

**Closed issues:**

- Update Docs - they aren't accurate right now. [\#46](https://github.com/DataDog/kvexpress/issues/46)

## [v1.0](https://github.com/DataDog/kvexpress/tree/v1.0) (2015-11-08)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v0.9...v1.0)

**Implemented enhancements:**

- Time guage over length of kvexpress run. [\#40](https://github.com/DataDog/kvexpress/issues/40)

**Closed issues:**

- Add logs with timer that shows elapsed time. [\#42](https://github.com/DataDog/kvexpress/issues/42)
- Logging is too verbose. [\#35](https://github.com/DataDog/kvexpress/issues/35)
- Profile and see where it's spending time. [\#16](https://github.com/DataDog/kvexpress/issues/16)

## [v0.9](https://github.com/DataDog/kvexpress/tree/v0.9) (2015-11-05)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v0.8...v0.9)

**Implemented enhancements:**

- Refactor Consul connection logic. [\#22](https://github.com/DataDog/kvexpress/issues/22)
- Don't write the file if the checksum matches.  [\#9](https://github.com/DataDog/kvexpress/issues/9)
- Add Datadog event logging? [\#5](https://github.com/DataDog/kvexpress/issues/5)

**Closed issues:**

- Don't connect to Datadog if API and APP keys aren't set. [\#38](https://github.com/DataDog/kvexpress/issues/38)
- Add Datadog event logging to `stop` command. [\#37](https://github.com/DataDog/kvexpress/issues/37)
- Let's not log the token. That's a bit of a fail. [\#34](https://github.com/DataDog/kvexpress/issues/34)
- Doesn't copy the file to .last if the checksum in the KV store matches. [\#33](https://github.com/DataDog/kvexpress/issues/33)
- Add some tests. [\#3](https://github.com/DataDog/kvexpress/issues/3)

## [v0.8](https://github.com/DataDog/kvexpress/tree/v0.8) (2015-11-03)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v0.7...v0.8)

**Implemented enhancements:**

- Add `kvexpress out --ignore-stop-command` [\#32](https://github.com/DataDog/kvexpress/issues/32)
- Add ability to clear out local cache files via command. [\#31](https://github.com/DataDog/kvexpress/issues/31)
- Add ability to set stop comment on the cli. [\#30](https://github.com/DataDog/kvexpress/issues/30)

**Closed issues:**

- Remove file output from syslog [\#29](https://github.com/DataDog/kvexpress/issues/29)

## [v0.7](https://github.com/DataDog/kvexpress/tree/v0.7) (2015-10-30)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v0.6...v0.7)

**Closed issues:**

- Blank lines are annoying. [\#28](https://github.com/DataDog/kvexpress/issues/28)
- When there's a folder the path doesn't work right. [\#27](https://github.com/DataDog/kvexpress/issues/27)
- Show `stop` message in syslog message. [\#17](https://github.com/DataDog/kvexpress/issues/17)

## [v0.6](https://github.com/DataDog/kvexpress/tree/v0.6) (2015-10-29)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v0.5...v0.6)

**Implemented enhancements:**

- Add dogstatsd logging. [\#4](https://github.com/DataDog/kvexpress/issues/4)

**Closed issues:**

- Add Makefile command to startup local Consul. [\#26](https://github.com/DataDog/kvexpress/issues/26)
- Track the byte sizes of the config files we're putting in. [\#20](https://github.com/DataDog/kvexpress/issues/20)
- Add `in` command. [\#12](https://github.com/DataDog/kvexpress/issues/12)
- Detail how the keys are organized in Consul. [\#2](https://github.com/DataDog/kvexpress/issues/2)

## [v0.5](https://github.com/DataDog/kvexpress/tree/v0.5) (2015-10-27)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v0.4...v0.5)

**Closed issues:**

- diff for `in` command doesn't seem to work. [\#15](https://github.com/DataDog/kvexpress/issues/15)
- Make sure `out` command follows `stop` key addition. [\#14](https://github.com/DataDog/kvexpress/issues/14)
- consul api panics if it doesn't find the stop key [\#13](https://github.com/DataDog/kvexpress/issues/13)
- Refactor functions to `kvexpress` package. [\#11](https://github.com/DataDog/kvexpress/issues/11)

## [v0.4](https://github.com/DataDog/kvexpress/tree/v0.4) (2015-10-26)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v0.3...v0.4)

## [v0.3](https://github.com/DataDog/kvexpress/tree/v0.3) (2015-10-25)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v0.2...v0.3)

**Implemented enhancements:**

- Command to run after file is written.  [\#10](https://github.com/DataDog/kvexpress/issues/10)
- Tokens [\#6](https://github.com/DataDog/kvexpress/issues/6)

**Closed issues:**

- make version 0.1 [\#8](https://github.com/DataDog/kvexpress/issues/8)

## [v0.2](https://github.com/DataDog/kvexpress/tree/v0.2) (2015-10-24)
[Full Changelog](https://github.com/DataDog/kvexpress/compare/v0.1...v0.2)

**Implemented enhancements:**

- --version [\#7](https://github.com/DataDog/kvexpress/issues/7)

## [v0.1](https://github.com/DataDog/kvexpress/tree/v0.1) (2015-10-23)


\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*