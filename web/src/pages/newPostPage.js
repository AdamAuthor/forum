import { ShowAlert } from "../scripts/showAlert.js";
import { CheckAuthentication } from "../scripts/checkAuth.js";

const newPostButton = document.getElementById('new-post');
const postsContainer = document.getElementById('posts');
const categories = document.getElementById('categories');
const successMessage = document.getElementById('success-message');


document.addEventListener('DOMContentLoaded', function () {
    newPostButton.addEventListener('click', function () {
        CheckAuthentication();
        categories.innerHTML = ''
        postsContainer.innerHTML = `
            <div class="new-post-form">
                <h2>Create New Post</h2>
                <form id="new-post-form">
                    <div class="input-box">
                        <input type="text" id="post-title" name="post-title"  required>
                        <label>Title</label>
                    </div>
                    <div class="input-box" id="tt-area">
                        <textarea id="post-content" name="post-content" rows="5"required></textarea>
                        <label id="tt-label">Content</label>
                    </div>
                    <div id="categories-container"></div>
                    <button type="submit" class="btn">Create Post</button>
                </form>
            </div>
        `;

        const categoriesContainer = document.getElementById('categories-container');
        fetchCategories(categoriesContainer);

        const newPostForm = document.getElementById('new-post-form');
        newPostForm.addEventListener('submit', function (event) {
            event.preventDefault();

            const categoryCheckboxes = document.querySelectorAll('.category-new-post input[type="checkbox"]:checked');
            if (categoryCheckboxes.length === 0) {
                ShowAlert('Please select at least one category.');
                return;
            }


            const formData = new FormData(newPostForm);
            const post = {
                title: formData.get('post-title'),
                content: formData.get('post-content'),
                categories: getSelectedCategories()
            };

            fetch('http://localhost:8080/create-post', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Cookie': document.cookie
                },
                body: JSON.stringify(post)
            })
                .then(response => {
                    if (!response.ok) {
                        throw new Error(response.statusText);
                    }
                    return;
                })
                .then(data => {
                    successMessage.textContent = 'New post created successful!';
                    successMessage.style.display = 'block';
                    setTimeout(function () {
                        successMessage.style.display = 'none';
                    }, 5000);
                    newPostForm.reset();
                })
                .catch(error => {
                    if (error.message.includes('Unauthorized')) {
                        ShowAlert('Unauthorized: Please login to create a post.');
                    } else {
                        ShowAlert('Error creating post: ' + error.message);
                    }
                });
        });
    });

    function fetchCategories(categoriesContainer) {
        const options = {
            mode: 'cors',
            method: 'GET'
        };
        fetch('http://localhost:8080/categories', options)
            .then(response => response.json())
            .then(categories => {
                displayCategories(categories, categoriesContainer);
            })
            .catch(error => {
                console.error('Error fetching categories:', error);
            });
    }

    function displayCategories(categories, categoriesContainer) {
        categories.forEach(category => {
            const categoryId = `category-${category.id}`;
            const categoryCheckbox = document.createElement('input');
            categoryCheckbox.type = 'checkbox';
            categoryCheckbox.name = 'category-new-post';
            categoryCheckbox.value = category.name;
            categoryCheckbox.id = categoryId;

            const categoryLabel = document.createElement('label');
            categoryLabel.textContent = category.name;
            categoryLabel.htmlFor = categoryId;
            const categoryElement = document.createElement('div');
            categoryElement.classList.add('category-new-post');
            categoryElement.appendChild(categoryCheckbox);
            categoryElement.appendChild(categoryLabel);

            categoriesContainer.appendChild(categoryElement);
        });
    }


    function getSelectedCategories() {
        const selectedCategories = [];
        const categoryCheckboxes = document.querySelectorAll('.category-new-post input[type="checkbox"]:checked');
        categoryCheckboxes.forEach(checkbox => {
            selectedCategories.push(checkbox.value);
        });
        return selectedCategories;
    }
});


