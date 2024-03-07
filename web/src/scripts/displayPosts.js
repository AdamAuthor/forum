import { PostInfo } from '../pages/postInfoPage.js';
import { DeletePost } from "../scripts/fetch.js";

const postsContainer = document.getElementById('posts');
const categoriesContainer = document.getElementById('categories');


export function DisplayPosts(posts, name) {
    postsContainer.innerHTML = '';
    if (name !== 'main') {
        categoriesContainer.style.display = 'none';
    }

    if (posts !== null) {
        posts.forEach(post => {
            const postElement = document.createElement('div');
            postElement.classList.add('post');

            const createdAt = new Date(post.post.created_at);
            const formattedDate = `${createdAt.toLocaleDateString()} ${createdAt.getHours()}:${createdAt.getMinutes()}`;
            const categories = post.post.categories.join(', ');

            if (name == 'my') {
                postElement.innerHTML = `
                    <h2>${post.post.title}</h2>
                    <p>Posted by: ${post.post.author_name}</p>
                    <p>Likes: ${post.likes.like}</p>
                    <p>Posted at: ${formattedDate}</p>
                    <p>Dislikes: ${post.likes.dislike}</p>
                    <p>Categories: ${categories}</p>
                    <button class="btnDelete" data-post-id="${post.post.id}">Delete</button>
                `;
                postElement.querySelector('.btnDelete').addEventListener('click', function (event) {
                    event.stopPropagation();
                    const postId = event.target.dataset.postId;
                    DeletePost(postId);
                });
            } else {
                postElement.innerHTML = `
                    <h2>${post.post.title}</h2>
                    <p>Posted by: ${post.post.author_name}</p>
                    <p>Likes: ${post.likes.like}</p>
                    <p>Posted at: ${formattedDate}</p>
                    <p>Dislikes: ${post.likes.dislike}</p>
                    <p>Categories: ${categories}</p>
                `;
            }
            postElement.addEventListener('click', () => {
                categoriesContainer.style.display = 'none';
                postsContainer.innerHTML = '';
                let id = post.post.id
                if (name === 'main') {
                    PostInfo(postsContainer, 'main', id)
                } else if (name === 'my') {
                    PostInfo(postsContainer, 'my', id)
                } else if (name === 'like') {
                    PostInfo(postsContainer, 'like', id)
                }
            });
            postsContainer.appendChild(postElement);
        });
    }
}