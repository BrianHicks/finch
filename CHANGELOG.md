# Changelog

## 0.2.0

**INCOMPATIBLE WITH OLDER VERSIONS**

This should be OK though, nobody uses this but me (Brian Hicks). If you have
been using and enjoying Finch, super! Before upgrading, echo all your tasks out
to a file: `finch select > tasks.txt`. You'll just have to re-import them by
script or by hand into the new version.

 - Rewrite in a more sustainable style
 - Uses JSON on disk instead of msgpack in leveldb (but that may be added again later)

## 0.1.0

 - Initial release
