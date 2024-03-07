const btnPopup = document.querySelector('.btnLogin-popup');
const btnOut = document.querySelector('.btnLogout-popup');

export function CheckAuthentication() {
    var xhr = new XMLHttpRequest();
    xhr.open('GET', 'http://localhost:8080/is-login', true);
    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
                btnPopup.style.display = 'none';
                btnOut.style.display = 'inline-block';
            } else if (xhr.status === 401) {
                logout();
            }
        }
    };
    xhr.send();
}

function logout() {
    var xhr = new XMLHttpRequest();
    xhr.open('POST', 'http://localhost:8080/logout', true);
    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
                document.cookie = 'session_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
                btnPopup.style.display = 'inline-block';
                btnOut.style.display = 'none';;
            } else {
                console.error('Error logging out');
            }
        }
    };
    xhr.send();
}
