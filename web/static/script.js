let form = document.querySelector("form");

form.addEventListener("submit", (e) => {
  handleSubmission(e)
}, false);

function handleSubmission(e) {
  e.preventDefault();

  let email = form.querySelector("#email").value;
  let password = form.querySelector("#password").value;

  fetch("/api/auth", {
      method: 'post',
      body: JSON.stringify({
        "email": email,
        "password": password
      })
    })
    .then(response => response.json())
    .then(data => authorizeWith(data));
}

function authorizeWith(data) {
  document.querySelector(".card").innerHTML = JSON.stringify(data);
}