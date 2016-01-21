# Contributing

We love pull requests from everyone. By participating in this project, you agree to abide by the kvexpress [code of conduct](http://todogroup.org/opencodeofconduct/#kvexpress/darron@froese.org).

Fork, then clone the repo:

    git clone git@github.com:your-username/kvexpress.git

Set up your machine:

    brew cask install consul
    brew install go
    make deps
    make

Make sure the tests pass:

    make test

Make your change. Add tests for your change. Make the tests pass:

    make test

Push to your fork and [submit a pull request][pr].

[pr]: https://github.com/DataDog/kvexpress/compare/

At this point you're waiting on us. We like to at least comment on pull requests within a few business days (and, typically, one or two business days). We may suggest some changes or improvements or alternatives.

Some things that will increase the chance that your pull request is accepted:

* Write tests.
* Make sure your code uses standard Golang conventions. [go-plus](https://atom.io/packages/go-plus) works great with Atom.
* [Clear is better than clever.](http://go-proverbs.github.io/)
* Write a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).
* Communicate with us clearly so that we understand what you're trying to accomplish.
