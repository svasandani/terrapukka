var form = document.querySelector("form");
var card = document.querySelector(".sign-in-card");
let urlParams = new URLSearchParams(window.location.search);
let responseOK = false;

const API = "http://localhost:3000/api/";

const errorDict = {
  "password too short; minimum 8 alphanumeric characters": "Your password is too insecure. Please use a password longer than 8 characters.",
  "client could not be found": "Something went wrong. Try going back to the previous page and signing in again.",
  "user email or password is incorrect": "Email or password incorrect. Please try again.",
  "invalid field: email": "Your email appears to be invalid. Please try again.",
  "required field missing: response_type": "Something went wrong. Try going back to the previous page and signing in again.",
  "required field missing: client_id": "Something went wrong. Try going back to the previous page and signing in again.",
  "required field missing: redirect_uri": "Something went wrong. Try going back to the previous page and signing in again.",
  "required field missing: email": "You must enter an email. Please try again.",
  "required field missing: password": "You must enter a password. Please try again."
}

function doReady(method) {
  form = document.querySelector("form");
  card = document.querySelector(".sign-in-card");

  if (method === "sign-in") {
    urlParams.set("method", "register");
    document.querySelector(".register-link").href = "//" + location.host + location.pathname + "?" + urlParams.toString();

    urlParams.set("method","reset-token");
    document.querySelector(".forgot-password-link").href = "//" + location.host + location.pathname + "?" + urlParams.toString();
  } else if (method === "register" || method === "reset-token" || method === "reset") {
    urlParams.set("method", "sign-in");
    document.querySelector(".sign-in-link").href = "//" + location.host + location.pathname + "?" + urlParams.toString();
  }

  form.addEventListener("submit", (e) => {
    handleSubmission(e, method)
  }, false);

}


function handleSubmission(e, method) {
  e.preventDefault();

  let name = "";
  let user = {};
  
  if (method === "reset-token") {
    user.email = form.querySelector("#email").value;

    fetch(API + ("reset_token"), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        "redirect_uri": urlParams.get("redirect_uri"),
        "client_id": urlParams.get("client_id"),
        "user": user
      })
    })
    .then(response => {
      if (response.status === 200) {
        response = response.json();
        responseOK = true;
        return response;
      } else {
        response = response.json();
        return response;
      }
    })
    .then(data => console.log(data));

    return;
  }

  let password = form.querySelector("#password").value;
  let confirmPassword = "";

  if (method === "reset") {
    confirmPassword = form.querySelector("#confirm-password").value;

    if (confirmPassword !== password) {
      let el = appendError(createError("Your passwords need to match. Please try again."));
      setTimeout(() => {
        removeError(el);
      }, 4000);

      form.querySelector("#password").value = "";
      form.querySelector("#confirm-password").value = "";

      return;
    }

    user.password = password;

    fetch(API + ("reset"), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        "reset_token": urlParams.get("reset_token"),
        "user": user
      })
    })
    .then(response => {
      if (response.status === 200) {
        urlParams.set("method", "sign-in");
        urlParams.delete("reset_token");
        window.location = "//" + location.host + location.pathname + "?" + urlParams.toString();
      } else {
        let el = appendError(createError("Something went wrong. Please try again."));
        setTimeout(() => {
          removeError(el);
        }, 4000);
      }
    })

    return;
  }

  let email = form.querySelector("#email").value;

  if (method === "register") {
    name = form.querySelector("#name").value;
    confirmPassword = form.querySelector("#confirm-password").value;

    if (confirmPassword !== password) {
      let el = appendError(createError("Your passwords need to match. Please try again."));
      setTimeout(() => {
        removeError(el);
      }, 4000);

      form.querySelector("#password").value = "";
      form.querySelector("#confirm-password").value = "";

      return;
    }

    user.name = name;
  }

  user.email = email;
  user.password = password;


  fetch(API + (method === "register" ? method : "auth"), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        "response_type": "code",
        "redirect_uri": urlParams.get("redirect_uri"),
        "client_id": urlParams.get("client_id"),
        "state": urlParams.get("state"),
        "user": user
      })
    })
    .then(response => {
      if (response.status === 200) {
        response = response.json();
        responseOK = true;
        return response;
      } else {
        response = response.json();
        return response;
      }
    })
    .then(data => authorizeWith(data));
}

function authorizeWith(data) {
  if (responseOK) {
    let redirectParams = new URLSearchParams("");
    redirectParams.set("auth_code", data.auth_code);
    redirectParams.set("state", data.state);

    window.location = data.redirect_uri + "?" + redirectParams.toString();
  } else {
    if (errorDict[data.message] == undefined) {
      let el = appendError(createError("An error occurred. Please try signing in again."));
      setTimeout(() => {
        removeError(el);
      }, 4000);
    } else {
      let el = appendError(createError(errorDict[data.message]));
      setTimeout(() => {
        removeError(el);
      }, 4000);
    }
  }
}

function removeErrors() {
  let errors = document.querySelectorAll(".errors");
  errors.forEach(error => {
    card.removeChild(error);
  });
}

function removeError(el) {
  el.style.opacity = "50%";
  setTimeout(() => {
    el.style.opacity = "0%";
    setTimeout(() => {
      removeElement(el);
    }, 300);
  }, 3500);
}

function removeElement(el) {
  el.parentNode.removeChild(el);
}

function createError(text) {
  let el = document.createElement("div");
  el.classList.add("errors", "card");
  el.innerHTML = text;
  return el;
}

function appendError(el) {
  let errors = document.querySelectorAll(".errors");
  if (errors.length > 2) {
    removeElement(errors[errors.length - 1]);
  }
  let head = document.querySelector("h1");
  head.after(el);
  return el;
}

function loadMainElement(mode, service) {
  let url = "";

  if (mode == "register") {
    url = "register.html";
  } else if (mode == "reset-token") {
    url = "reset_token.html";
  } else if (mode == "reset") {
    url = "reset.html";
  } else if (mode == "sign-in") {
    url = "sign_in.html";
  } 

  fetch(url).then(response => response.text()).then(data => {
    let header = document.querySelector("header");
    let main = document.querySelector("main");

    data = data.replace("{{service}}", service);

    if (main == null) {
      header.insertAdjacentHTML("afterend", data);
    } else {
      document.body.removeChild(main);
      header.insertAdjacentHTML("afterend", data);
    }
  }).then(() => {
    doReady(urlParams.get("method"))
  });
}

window.addEventListener('DOMContentLoaded', () => {
  fetch(API + "client/identify", {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        "client": {

          "redirect_uri": urlParams.get("redirect_uri"),
          "id": urlParams.get("client_id")
        }
      })
    })
    .then(response => {
      if (response.status === 200) {
        response = response.json();
        return response;
      } else {
        let section = document.createElement("section");
        let span = document.createElement("span");
        span.setAttribute("aria-hidden", "true");
        section.appendChild(span);
        section.id = "form";
        let el = document.createElement("div");
        el.classList.add("card", "form-card");
        el.appendChild(createError("Something went wrong. Please try signing in again."));
        section.appendChild(el);
        document.querySelector("header").after(section);
      }
    }).then(data => {
      if (data)
        loadMainElement(urlParams.get("method"), data.client.name);
    });
});
