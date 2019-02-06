[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Ftb0hdan%2Fmicroservices-keeper.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Ftb0hdan%2Fmicroservices-keeper?ref=badge_shield)

*microservices-keeper* is a tool to keep track of Code Change Decisions [CCD] as part of microservices maintenance.
I've coined up this term after reading Michael Nygard's article - see [About](#about).
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
microservices-keeper --message "we have decided to write some tests" --gituser yourGithubUser
```

or

```sh
cat somebigmessage.txt|microservices-keeper --gituser yourGithubUser
```


Slack integration
------------

Receiving slack messages (@channel message, where bot is added) using WebSockets:

```sh
microservices-keeper --slack-token xoxb-xxxx --slack-modes 10 --gituser yourGithubUser
```


using Events API (typicall app_mention capability):

```sh
/microservices-keeper --slack-token xoxb-xxxx --slack-modes 01 --slack-verification-token RRxxxx --gituser yourGithubUser
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
Loosely based on following:
- http://thinkrelevance.com/blog/2011/11/15/documenting-architecture-decisions
- https://github.com/npryce/adr-tools



## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Ftb0hdan%2Fmicroservices-keeper.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Ftb0hdan%2Fmicroservices-keeper?ref=badge_large)