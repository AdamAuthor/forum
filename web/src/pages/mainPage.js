import { FetchCategories, FetchPostsByCategory } from "../scripts/fetch.js";
import { CheckAuthentication } from "../scripts/checkAuth.js";

const categoriesContainer = document.getElementById('categories');

window.onload = function () {
    CheckAuthentication();
    FetchCategories();
};

document.addEventListener('DOMContentLoaded', function () {
    const homeBtn = document.getElementById('home');
    homeBtn.addEventListener('click', () => {
        CheckAuthentication();
        FetchCategories();
    });

    categoriesContainer.addEventListener('click', function (event) {
        const clickedCategory = event.target.closest('.category');
        if (!clickedCategory) return;

        const categoryName = clickedCategory.textContent.trim();
        const activeCategory = document.querySelector('.category.active');
        activeCategory.classList.remove('active');
        clickedCategory.classList.add('active');
        CheckAuthentication();
        FetchPostsByCategory(categoryName);
    });
});




