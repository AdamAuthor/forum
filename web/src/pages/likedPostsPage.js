import { FetchLikedPosts } from "../scripts/fetch.js";

document.addEventListener('DOMContentLoaded', function () {
    const likeMenu = document.getElementById('like-menu')
    likeMenu.addEventListener('click', () => {
        FetchLikedPosts();
    });
});