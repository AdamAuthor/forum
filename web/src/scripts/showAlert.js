const customAlert = document.getElementById('custom-alert');
const customAlertMessage = document.getElementById('custom-alert-message');
const customAlertButton = document.getElementById('custom-alert-button');

export function ShowAlert(message) {
    customAlertMessage.textContent = message;
    customAlert.classList.add('active');

    customAlertButton.addEventListener('click', function () {
        customAlert.classList.remove('active');
    });
}