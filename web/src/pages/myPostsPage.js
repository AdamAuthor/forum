import { FetchMyPosts } from "../scripts/fetch.js";
import { CheckAuthentication } from "../scripts/checkAuth.js";

document.addEventListener('DOMContentLoaded', function () {
    const myPosts = document.getElementById('my-posts');

    myPosts.addEventListener('click', () => {
        CheckAuthentication();
        FetchMyPosts();
    });
});
