import { FetchCategories, FetchMyPosts, FetchLikedPosts } from "../scripts/fetch.js";
import { ShowAlert } from "../scripts/showAlert.js";
import { CheckAuthentication } from "../scripts/checkAuth.js";

export async function PostInfo(postsContainer, name, id) {
    try {
        const postDetails = await fetchPostDetails(id);
        const postDetailsElement = createPostDetailsElement(postDetails);
        postsContainer.appendChild(postDetailsElement);

        const iconClose = postDetailsElement.querySelector('.post-close');
        iconClose.addEventListener('click', async () => {
            CheckAuthentication();
            await handlePostClose(name);
        });

        const postLikes = await fetchPostLikes(id);
        const postLikesElement = createPostLikesElement(postLikes);
        postsContainer.appendChild(postLikesElement);

        handleLikeButtons(postLikesElement, id, postsContainer, name);

        const commentForm = createCommentForm();
        commentForm.addEventListener('submit', async (event) => {
            CheckAuthentication();
            event.preventDefault();
            const commentContent = document.getElementById('comment-content').value;
            await createComment(id, commentContent);
            postsContainer.innerHTML = '';
            PostInfo(postsContainer, name, id);
        });
        postsContainer.appendChild(commentForm);

        const comments = await fetchPostComments(id);
        postsContainer.appendChild(await createCommentsContainer(comments, postsContainer, name, id));

    } catch (error) {
        console.error('Error fetching post information:', error);
    }
}

async function fetchPostDetails(id) {
    const response = await fetch(`http://localhost:8080/post?id=${id}`, { method: 'GET', mode: 'cors' });
    if (!response.ok) throw new Error('Failed to fetch post details');
    return await response.json();
}

function createPostDetailsElement(postDetails) {
    const postDetailsElement = document.createElement('div');
    postDetailsElement.classList.add('post-details');
    const createdAt = new Date(postDetails.created_at);
    const formattedDate = `${createdAt.toLocaleDateString()} ${createdAt.getHours()}:${createdAt.getMinutes()}`;
    postDetailsElement.innerHTML = `
        <span class="post-close">
            <ion-icon name="close-outline"></ion-icon>
        </span>
        <h2>${postDetails.title}</h2>
        <p>${postDetails.content}</p>
        <p>Posted by: ${postDetails.author_name}</p>
        <p>Created at: ${formattedDate}</p>
        <p>Categories: ${postDetails.categories.join(', ')}</p>
    `;
    return postDetailsElement;
}

async function fetchPostLikes(id) {
    const response = await fetch(`http://localhost:8080/likes?post_id=${id}`, { method: 'GET', mode: 'cors' });
    if (!response.ok) throw new Error('Failed to fetch post likes');
    return await response.json();
}

function createPostLikesElement(postLikes) {
    const postLikesElement = document.createElement('div');
    postLikesElement.classList.add('post-likes');
    postLikesElement.innerHTML = `
        <span class="like">
            <ion-icon name="heart-outline"></ion-icon>
        </span>
        <label>${postLikes.like}</label>
        <span class="dislike">
            <ion-icon name="heart-dislike-outline"></ion-icon>
        </span>
        <label>${postLikes.dislike}</label>
    `;
    return postLikesElement;
}

async function handleLikeButtons(postLikesElement, id, postsContainer, name) {
    const likeButton = postLikesElement.querySelector('.like');
    const dislikeButton = postLikesElement.querySelector('.dislike');

    async function checkLikedStatus() {
        const isLiked = await checkIfLiked(id);
        if (isLiked === null) {
            return null;
        } else if (isLiked.like_type === 1) {
            return 'liked';
        } else if (isLiked.like_type === 2) {
            return 'disliked';
        }
    }

    const likedStatus = await checkLikedStatus();
    if (likedStatus === 'liked') {
        likeButton.classList.add('liked');
    } else if (likedStatus === 'disliked') {
        dislikeButton.classList.add('disliked');
    }

    likeButton.addEventListener('click', async () => {
        CheckAuthentication();
        if (!likeButton.classList.contains('liked')) {
            await handleUnlike(id);
            await handleLike(id, 1);
            likeButton.classList.add('liked');
            if (dislikeButton.classList.contains('disliked')) {
                dislikeButton.classList.remove('disliked');
            }
        } else {
            await handleUnlike(id);
            likeButton.classList.remove('liked');
        }
        postsContainer.innerHTML = '';
        PostInfo(postsContainer, name, id);
    });

    dislikeButton.addEventListener('click', async () => {
        CheckAuthentication();
        if (!dislikeButton.classList.contains('disliked')) {
            await handleUnlike(id);
            await handleLike(id, 2);
            dislikeButton.classList.add('disliked');
            if (likeButton.classList.contains('liked')) {
                likeButton.classList.remove('liked');
            }
        } else {
            await handleUnlike(id);
            dislikeButton.classList.remove('disliked');
        }
        postsContainer.innerHTML = '';
        PostInfo(postsContainer, name, id);
    });
}

async function checkIfLiked(postId) {
    try {
        const response = await fetch(`http://localhost:8080/is-liked?post_id=${postId}`, { method: 'GET', mode: 'cors' });
        return await response.json();
    } catch (error) {
        return null;
    }
}

export async function handleUnlike(postId) {
    const url = `http://localhost:8080/unlike-post?post_id=${postId}`;
    const options = {
        method: 'POST',
        mode: 'cors'
    };

    try {
        const response = await fetch(url, options);
        if (!response.ok) {
            if (response.status === 401) {
                ShowAlert("Unauthorized: Please login to unlike this post.");
            }
        }
    } catch (error) {
        console.error(error);
    }
}

export async function handleLike(postId, likeType) {
    const url = 'http://localhost:8080/like-post';
    const options = {
        method: 'POST',
        mode: 'cors',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ post_id: postId, like_type: likeType })
    };

    try {
        const response = await fetch(url, options);
        if (!response.ok) {
            if (response.status === 401) {
                ShowAlert("Unauthorized: Please login to like this post.");
            }
        }
    } catch (error) {
        console.error(error);
    }
}

async function fetchPostComments(id) {
    const response = await fetch(`http://localhost:8080/comments?post_id=${id}`, { method: 'GET', mode: 'cors' });
    if (!response.ok) throw new Error('Failed to fetch post comments');
    return await response.json();
}

async function createCommentsContainer(comments, postsContainer, name, postId) {
    const commentsContainer = document.createElement('div');
    commentsContainer.classList.add('comments-container');
    if (!comments) {
        const commentElement = document.createElement('div');
        commentElement.classList.add('comment');
        commentElement.innerHTML = `<h2>No comments</h2>`;
        commentsContainer.appendChild(commentElement);
    } else {
        for (const comment of comments) {
            const commentElement = document.createElement('div');
            const createdAt = new Date(comment.created_at);
            const formattedDate = `${createdAt.toLocaleDateString()} ${createdAt.getHours()}:${createdAt.getMinutes()}`;
            commentElement.classList.add('comment');
            commentElement.innerHTML = `
                <p>${comment.content}</p>
                <p>Commented by: ${comment.author_name}</p>
                <p>Created at: ${formattedDate}</p>
            `;
            const commentLikes = await fetchCommentLikes(comment.id);
            const commentLikesElement = createCommentLikesElement(commentLikes);
            await handleCommentLikeButtons(commentLikesElement, comment.id, postsContainer, name, postId);
            commentElement.appendChild(commentLikesElement);
            commentsContainer.appendChild(commentElement);
        }
    }
    return commentsContainer;
}


function createCommentForm() {
    const commentForm = document.createElement('form');
    commentForm.classList.add('comment-form');
    commentForm.innerHTML = `
        <textarea id="comment-content" placeholder="Your comment..." required></textarea>
        <button type="submit">Post Comment</button>
    `;

    commentForm.addEventListener('submit', (event) => {
        CheckAuthentication();
        const commentContent = document.getElementById('comment-content').value.trim();
        if (!commentContent) {
            event.preventDefault();
            ShowAlert("Please enter a comment before posting.");
        }
    });

    return commentForm;
}


export async function createComment(postId, content) {
    const url = 'http://localhost:8080/create-comment';
    const options = {
        method: 'POST',
        mode: 'cors',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ post_id: postId, content: content })
    };

    try {
        const response = await fetch(url, options);
        if (!response.ok) {
            if (response.status === 401) {
                ShowAlert("Unauthorized: Please login to post a comment.");
            }
        }
    } catch (error) {
        console.error(error);
    }
}

async function handlePostClose(name) {
    switch (name) {
        case 'main':
            await FetchCategories();
            break;
        case 'my':
            await FetchMyPosts();
            break;
        case 'like':
            await FetchLikedPosts();
            break;
        default:
            console.error('Invalid post type');
    }
}

async function fetchCommentLikes(id) {
    const response = await fetch(`http://localhost:8080/comment-likes?comment_id=${id}`, { method: 'GET', mode: 'cors' });
    if (!response.ok) throw new Error('Failed to fetch post likes');
    let respBody
    try {
        respBody = await response.json()
    } catch (e) {
        console.log(e)
        return
    }
    return respBody;
}

function createCommentLikesElement(commentLikes) {
    const postLikesElement = document.createElement('div');
    postLikesElement.classList.add('comment-likes');
    postLikesElement.innerHTML = `
        <span class="like">
            <ion-icon name="heart-outline"></ion-icon>
        </span>
        <label>${commentLikes.like}</label>
        <span class="dislike">
            <ion-icon name="heart-dislike-outline"></ion-icon>
        </span>
        <label>${commentLikes.dislike}</label>
    `;
    return postLikesElement;
}

async function handleCommentLikeButtons(commentLikesElement, id, postsContainer, name, postId) {
    const likeButton = commentLikesElement.querySelector('.like');
    const dislikeButton = commentLikesElement.querySelector('.dislike');

    async function checkLikedStatus() {
        const isLiked = await checkIfCommentLiked(id);
        if (isLiked === null) {
            return null;
        } else if (isLiked.like_type === 1) {
            return 'liked';
        } else if (isLiked.like_type === 2) {
            return 'disliked';
        }
    }

    const likedStatus = await checkLikedStatus();
    if (likedStatus === 'liked') {
        likeButton.classList.add('liked');
    } else if (likedStatus === 'disliked') {
        dislikeButton.classList.add('disliked');
    }

    likeButton.addEventListener('click', async () => {
        CheckAuthentication();
        if (!likeButton.classList.contains('liked')) {
            await handleCommentUnlike(id);
            await handleCommentLike(id, 1);
            likeButton.classList.add('liked');
            if (dislikeButton.classList.contains('disliked')) {
                dislikeButton.classList.remove('disliked');
            }
        } else {
            await handleCommentUnlike(id);
            likeButton.classList.remove('liked');
        }
        postsContainer.innerHTML = '';
        PostInfo(postsContainer, name, postId);
    });

    dislikeButton.addEventListener('click', async () => {
        CheckAuthentication();
        if (!dislikeButton.classList.contains('disliked')) {
            await handleCommentUnlike(id);
            await handleCommentLike(id, 2);
            dislikeButton.classList.add('disliked');
            if (likeButton.classList.contains('liked')) {
                likeButton.classList.remove('liked');
            }
        } else {
            await handleCommentUnlike(id);
            dislikeButton.classList.remove('disliked');
        }
        postsContainer.innerHTML = '';
        PostInfo(postsContainer, name, postId);
    });
}

async function handleCommentLike(commentId, likeType) {
    const url = 'http://localhost:8080/like-comment';
    const options = {
        method: 'POST',
        mode: 'cors',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ comment_id: commentId, like_type: likeType })
    };

    try {
        const response = await fetch(url, options);
        if (!response.ok) {
            if (response.status === 401) {
                ShowAlert("Unauthorized: Please login to like this comment.");
            } else {
                const errorMessage = await response.text();
                ShowAlert("Error: " + errorMessage);
            }
        }
    } catch (error) {
        console.error(error);
    }
}

async function handleCommentUnlike(commentId) {
    const url = `http://localhost:8080/unlike-comment?comment_id=${commentId}`;
    const options = {
        method: 'POST',
        mode: 'cors'
    };

    try {
        const response = await fetch(url, options);
        if (!response.ok) {
            if (response.status === 401) {
                ShowAlert("Unauthorized: Please login to unlike this comment.");
            } else {
                const errorMessage = await response.text();
                ShowAlert("Error: " + errorMessage);
            }
        }
    } catch (error) {
        console.error(error);
    }
}

async function checkIfCommentLiked(commentId) {
    try {
        const response = await fetch(`http://localhost:8080/is-comment-liked?comment_id=${commentId}`, { method: 'GET', mode: 'cors' });
        if (!response.ok) {
            throw new Error('Failed to fetch comment like status');
        }
        const data = await response.json();
        if (data === null) {
            return null
        }
        return data
    } catch (error) {
        console.error('Error:', error);
        return null;
    }
}