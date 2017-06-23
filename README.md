# Gorgeous Go Crawler
A simple crawler application made to learn Go.

GGC searches for all links and static files in a user provided domain.
It supports two modes for searching links. In **well behaved** mode, the crawler follows the rules specified under the Robots.txt file. In **aggreessive** mode, the crawler searches for all urls within the domain.

Current Version only support **aggressive** mode. To run type:

```bash
 go build && ./GGC [LIST OF DOMAIN NAMES] [DURATION]
```
