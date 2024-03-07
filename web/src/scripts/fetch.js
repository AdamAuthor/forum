import { DisplayPosts } from '../scripts/displayPosts.js';
import { ShowAlert } from "./showAlert.js";


const categoriesContainer = document.getElementById('categories');
const errorMessage = document.getElementById('error-message');


export function FetchCategories() {
    categoriesContainer.style.display = 'flex';

    const options = {
        mode: 'cors',
        method: 'GET'
    };
    fetch('http://localhost:8080/categories', options)
        .then(response => response.json())
        .then(categories => {
            categoriesContainer.innerHTML = '';
            const allCategoryButton = document.createElement('div');
            allCategoryButton.classList.add('category', 'active');
            allCategoryButton.innerHTML = `<h1>All</h1>`;
            allCategoryButton.addEventListener('click', () => FetchPostsByCategory('All'));
            categoriesContainer.appendChild(allCategoryButton);
            categories.forEach(category => {
                const categoryElement = document.createElement('div');
                categoryElement.classList.add('category');
                categoryElement.innerHTML = `
                    <h1>${category.name}</h1>
                `;
                categoryElement.addEventListener('click', () => FetchPostsByCategory(category.name));
                categoriesContainer.appendChild(categoryElement);
            });
            FetchPostsByCategory('All');
        })
        .catch(error => {
            console.error('Error fetching categories:', error);
        });
}


export function FetchPostsByCategory(categoryName) {
    let url = 'http://localhost:8080/posts';
    if (categoryName !== 'All') {
        url = `http://localhost:8080/posts-by-category?category=${encodeURIComponent(categoryName)}`;
    }
    const options = {
        mode: 'cors',
        method: 'GET'
    };
    fetch(url, options)
        .then(response => response.json())
        .then(posts => DisplayPosts(posts, 'main'))
        .catch(error => console.error(`Error fetching posts for category ${categoryName}:`, error));
}

export function FetchLikedPosts() {
    let url = 'http://localhost:8080/posts-by-like';

    const options = {
        mode: 'cors',
        method: 'GET'
    };
    fetch(url, options)
        .then(response => response.json())
        .then(posts => {
            DisplayPosts(posts, 'like')
        })
        .catch(error => {
            if (error.message.includes('Unauthorized')) {
                ShowAlert('Unauthorized: Please login to create a post.');
            } else {
                errorMessage.textContent = 'Error: please, try again';
                errorMessage.style.display = 'block';
                setTimeout(function () {
                    errorMessage.style.display = 'none';
                }, 5000);
            }
        });
}


export function FetchMyPosts() {
    let url = 'http://localhost:8080/my-posts';

    const options = {
        mode: 'cors',
        method: 'GET'
    };
    fetch(url, options)
        .then(response => response.json())
        .then(posts => {
            DisplayPosts(posts, 'my');
        })
        .catch(error => {
            if (error.message.includes('Unauthorized')) {
                ShowAlert('Unauthorized: Please login to create a post.');
            } else {
                errorMessage.textContent = 'Error: please, try again';
                errorMessage.style.display = 'block';
                setTimeout(function () {
                    errorMessage.style.display = 'none';
                }, 5000);
            }
        });
}


export function DeletePost(postId) {
    const url = `http://localhost:8080/delete-post?post_id=${postId}`;
    const options = {
        method: 'DELETE',
    };

    fetch(url, options)
        .then(response => {
            if (!response.ok) {
                throw new Error('Delete failed');
            }
            FetchMyPosts();
        })
        .catch(error => {
            console.error('Error deleting post:', error);
            ShowAlert('Error deleting post. Please try again.');
        });
}
