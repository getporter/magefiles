# Contributing Guide

This is part of the [Porter][porter] project. If you are a new contributor,
check out our [New Contributor Guide][new-contrib]. The Porter [Contributing
Guide][contrib] also has lots of information about how to interact with the
project.

[porter]: https://github.com/getporter/porter
[new-contrib]: https://porter.sh/contribute
[contrib]: https://porter.sh/src/CONTRIBUTING.md

---

* [Initial setup](#initial-setup)
* [Magefile explained](#magefile-explained)

---

# Initial setup

You need to have [porter installed](https://porter.sh/install) first. Then run
`mage build install`. This will build and install the mixin into your porter
home directory.

## Magefile explained

Here are common [Magefile](https://magefile.org) tasks:

* `build` compiles all the go code.
* `test` runs all the tests.
