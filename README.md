# Forum Application

This is a full-fledged forum web application built with Go for the backend, SQLite for storage, and vanilla JavaScript for the frontend, formatted as a Single Page Application (SPA). It also comes with Docker support for easy deployment.

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
