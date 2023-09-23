# Github PageRank

## Description

The Github PageRank project aims to gather data about Github users and their followers to determine the trustworthiness of users based on their connections. The project utilizes the original PageRank algorithm from Google to analyze the network of Github followers and calculate a rank for each user. By identifying patterns and analyzing the network structure, the project aims to detect potential bot users and provide insights into user relationships on Github.

## Features

- Crawling: The project provides a command-line interface (CLI) for crawling Github user data and collecting follower information.
- Database Storage: All collected data is stored in an SQLite database for easy retrieval and analysis.
- PageRank Calculation: The project utilizes the PageRank algorithm to calculate the rank of each user based on their follower connections.
- Trustworthiness Analysis: By analyzing the network structure and user rankings, the project aims to identify potential bot users and determine the trustworthiness of Github users.

## Installation and Usage

To use the Github PageRank project, follow these steps:

1. Clone the repository:

   ```
   git clone https://github.com/jtarchie/github-pagerank.git
   ```

2. Change to the project directory:

   ```
   cd github-pagerank
   ```

3. Install the required dependencies:

   ```
   go mod download
   ```

4. Build the project:

   ```
   go build
   ```

5. Run the CLI commands for crawling and ranking:

   ```
   # Example CLI commands
   ./github-pagerank crawl --help
   ./github-pagerank crawl --Username=<username> --DBFilename=<filename>

   ./github-pagerank rank --help
   ./github-pagerank rank --DBFilename=<filename>
   ```

   Replace `<username>` with the starting username for crawling and `<filename>` with the desired name of the SQLite database file.

## Evaluation

The Github PageRank project offers a powerful tool for analyzing Github user relationships and determining user trustworthiness based on follower connections. By collecting and analyzing data, potential patterns and insights can be discovered. The project showcases the implementation of the original PageRank algorithm and provides a CLI for easy usage.

FIXME: Include any additional details or information on how to run the project, potential limitations, or further improvements that need to be addressed.