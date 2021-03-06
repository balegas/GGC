# Gorgeous Go Crawler
A simple crawler application made to learn Go.

GGC searches for all links and static files in a user-provided domain.
It supports two modes for searching links. In **well behaved** mode, the crawler follows the rules specified under the Robots.txt file. In **aggreessive** mode, the crawler searches for all urls within the domain.

Current Version only support **aggressive** mode. To run type:

```bash
 go build && ./GGC -d=[DURATION] -f=[PATH_OUTPUT_FILE | stdout] -w=[NUM_WORKERS] -b=[BUFFER_SIZE] -t=[THINK_TIME] [domain_name]+
```

**NOTE**: The crawler will visit any url that has domain\_name as suffix. For example: sites.google.com or fakegoogle.com are allowed. A strict single-domain policy is also available.
