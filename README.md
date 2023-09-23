# Project Name: GitHub PageRank

## Description:

The GitHub PageRank project is designed to gather data about followers for
GitHub users and apply the PageRank algorithm to determine the trustworthiness
of a user based on their connections. The project aims to identify patterns and
detect bot users through the analysis of user connections.

## Features:

- Crawl and gather data on GitHub users and their followers
- Apply PageRank algorithm to calculate the trustworthiness of a user based on
  their connections
- Identify patterns and detect bot users

## Installation and Usage:

1. Install the project dependencies by running the following command:
   ```bash
   go mod download
   ```
2. Import the necessary packages in your Go code:
   ```go
   import (
       "github.com/jtarchie/github-pagerank/crawl"
       "github.com/jtarchie/github-pagerank/rank"
   )
   ```
3. To use the crawling functionality, create an instance of `crawl.Cmd` and set
   the required fields:
   ```go
   cli := crawl.Cmd{
       DBFilename: "<path_to_database_file>",
       GithubAPIKey: "<github_api_key>",
       WaitInterval: time.Duration(1) * time.Second,
       ResultLimit: 100,
       MaxFollowing: 510,
       Username: "<starting_username>",
   }
   ```
4. To start crawling, call the `Run` method of the `Cmd` instance:
   ```go
   err := cli.Run()
   if err != nil {
       // handle error
   }
   ```

5. To use the ranking functionality, create an instance of `rank.Cmd` and set
   the required fields:
   ```go
   cli := rank.Cmd{
       DBFilename: "<path_to_database_file>",
   }
   ```
6. To run the ranking algorithm and get the results, call the `Run` method of
   the `Cmd` instance:
   ```go
   err := cli.Run()
   if err != nil {
       // handle error
   }
   ```
   The results will be printed on the console, showing the IDs and ranks of the
   users.

## Evaluation:

- This project provides a way to gather and analyze GitHub user data to
  determine trustworthiness based on their connections.
- It uses the PageRank algorithm to calculate the ranks of users.
- The crawling functionality allows you to gather data on users and their
  followers, while the ranking functionality applies the algorithm to calculate
  the ranks.
- The project currently runs on a subset of users and has shown promising
  results in identifying patterns and detecting bot users.
- TODO: Provide instructions or guidelines for using the project with a larger
  dataset or expanding the analysis capabilities.
- TODO: Add more detailed information about the specific patterns and
  characteristics that are used to identify bot users.
- TODO: Include information on the limitations or potential improvements for the
  project, such as handling different types of users and organizations.

## Contributing:

- Contributions are welcome! If you find any issues or have suggestions for
  improvements, please feel free to open an issue or submit a pull request.
- TODO: Add guidelines for contributing to the project.

## License:

- This project is licensed under the MIT License. See the [LICENSE](LICENSE)
  file for more details.
