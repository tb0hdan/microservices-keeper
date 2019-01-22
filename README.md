*microservices-keeper* is a tool to keep track of Code Change Decisions [CCD] as part of microservices maintenance.
I've coined up this term after reading Michael Nygard's article (see (About)[#about] ).
Basically, it's works with git repository as a form of centralized storage,
where all necessary changes get committed so the developer can walk through
the list and confirm that code edits have been made. Work in progress.


Installation
------------

The recommended way to install *microservices-keeper* is:

```sh
go get -u github.com/tb0hdan/microservices-keeper
```

Usage
------------

```sh
microservices-keeper -help
```

Example CLI usage:

```sh
microservices-keeper bla-bla
```

Development
------------

It is recommended to use Git via SSH while developing this project:

```sh
git config --global url.ssh://git@github.com.insteadof https://github.com
```

See: https://github.com/golang/go/issues/26894


About
------------
Loosely based on http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions
and https://github.com/npryce/adr-tools
