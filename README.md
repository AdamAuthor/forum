# Forum Application

This is a full-fledged forum web application built with Go for the backend, SQLite for storage, and vanilla JavaScript for the frontend, formatted as a Single Page Application (SPA). It also comes with Docker support for easy deployment.

## This project consists in creating a web forum that allows

- communication between users.
- associating categories to posts.
- liking and disliking posts and comments.
- filtering posts.

## Communication

In order for users to communicate between each other, they will have to be able to create posts and comments.

- Only registered users will be able to create posts and comments.
- When registered users are creating a post they can associate one or more categories to it.
  - The implementation and choice of the categories is up to you.
- The posts and comments should be visible to all users (registered or not).
- Non-registered users will only be able to see posts and comments.

## Features

- Categories and Posts: Users can create categories and make posts within them.
- Comments: Users can comment on posts.
- Like/Dislike System: Posts can be liked or disliked by users.
- User Activity Tracking: Users can see posts they've liked or created.
- Authentication and Authorization: Basic authentication and registration system using session management.
- Data Storage: SQLite is used as the database for storing information.
- Deletion: Users can delete their own posts.

## Tech Stack

- Backend: Go (Golang)
- Frontend: Vanilla JavaScript
- Database: SQLite
- Containerization: Docker
- Build Automation: Makefile
- Getting Started
- Clone the Repository:

## Getting Started

1. Clone the Repository:

``` bash
git clone <https://github.com/AdamAuthor/forum.git>
cd forum
```

2. Set Up Environment:

- Ensure you have Go installed.
- Ensure Docker is installed.

3. Build and Run:

```bash
make build
make run
```

4. Access the Application:

- Open your browser and navigate to <http://localhost:8080>.

## Usage

- Registration/Login: Users can register and log in to the forum.
- Creating Categories and Posts: Once logged in, users can create categories and make posts within them.
- Commenting: Users can comment on posts.
- Liking/Disliking: Posts can be liked or disliked by users.
- Viewing Activity: Users can view posts they've liked or created.
- Deleting Posts: Users can delete their own posts.

## This project will help me learn about

- The basics of web :
  - HTML
  - HTTP
  - Sessions and cookies
- Using and [setting up Docker](https://docs.docker.com/get-started/)
  - Containerizing an application
  - Compatibility/Dependency
  - Creating images
- SQL language
  - Manipulation of databases
- The basics of encryptions
