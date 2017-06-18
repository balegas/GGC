# Gorgeous Go Crawler
A simple crawler application made to learn Go.

GGC searches for all links and static files in a user provided domain.
It supports two modes for searching links. In **Well behaved** mode, the crawler follows the rules specified under the Robots.txt file. In **Aggreessive** mode, the crawler searches for all urls within the domain.

Parameters:

-d domain\_name: the domain name (Required).
-m mode: A for aggressive, W for well-behaved (default: W).
-t ms: time delay between requests, use 0 for no delay, -1 for getting the value from Robots.txt, or no delay if the site has no Robots.txt file (default: -1).
-a user\_agent: user agent to be used in each request (default: GGC).
-w number\_workers: the maximum number of workers to be used to analyze pages in parallel (default: 1).
-r max\_depth: maximum depth for searching an url. -1 for no constraint (default: -1).



