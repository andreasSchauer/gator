# Gator

Gator is a CLI app that lets you follow and aggregate RSS feeds and display their posts in the terminal.

## Installation

You need Go and Postgres installed. Afterwards, run `go install https://github.com/andreasSchauer/gator` to create the gator binary.

You also need to create a .gatorconfig.json file in your home directory with the following contents:

`{"db_url":"postgres://USERNAME:@localhost:5432/gator?sslmode=disable"}`

Replace "USERNAME" with your username obviously.

## Commands

### gator register

Usage: `gator login <username>`

Registers a user by the name of username and add them to the database.


### gator login

Usage: `gator login <username>`

Log into a registered user's account. There is no authorization in this project.


### gator users

Usage: `gator users`

List all the registered users in the database.


### gator reset

Usage: `gator reset`

Delete every account and feed from the database.


### gator addfeed

Usage: `gator addfeed <feedName> <feedURL>`

Adds an RSS feed to the database and automatically follows it for the logged in user.


### gator feeds

Usage: `gator feeds`

List all the feeds in the database.


### gator follow

Usage: `gator follow <feedURL>`

Follow a registered feed.


### gator unfollow

Usage: `gator unfollow <feedURL>`

Unfollow a registered feed.


### gator following

Usage: `gator following`

List all the feeds the logged in user follows.


### gator agg

Usage: `gator agg <timeBetweenRequests>`

Aggregate the posts from the feeds the user is following to display them with the `browse` command.

timeBetweenRequests expects a time string like "3s" or "3m". This sets the interval at which the aggregator is looking for new posts.

Example: ´gator agg 3m´


### gator browse

Usage: `gator browse [postsLimit]`

Display posts previously aggregated with the `agg` command.

postsLimit expects an integer that limits the number of posts being shown. The default value is 2 and the maximum is 50.