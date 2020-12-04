const filter = document.getElementById("filter");

const mainHeaderEl = document.getElementById("main-header");
const signinLink = mainHeaderEl?.querySelector(".main-header__signin-anchor");
const signupSubmit = mainHeaderEl?.querySelector(".main-header__signup-btn");

const modals = document.querySelectorAll(".modal");
const signinModal = document.getElementById("signin-modal");
const signupModal = document.getElementById("signup-modal");

const closeAuthModals = () => {
  filter && (filter.className = "_none");
  signinModal && (signinModal.className = "modal _none");
  signupModal && (signupModal.className = "modal _none");
};

const openAuthModal = (type: "signin" | "signup") => {
  filter && (filter.className = "_block");
  if (type === "signin") {
    signinModal && (signinModal.className = "modal _flex-c-cb");
    signupModal && (signupModal.className = "modal _none");
  } else if (type === "signup") {
    signinModal && (signinModal.className = "modal _none");
    signupModal && (signupModal.className = "modal _flex-c-cb");
  }
};

const auth = () => {
  filter?.addEventListener("click", closeAuthModals);
  signinLink?.addEventListener("click", () => openAuthModal("signin"));
  signupSubmit?.addEventListener("click", () => openAuthModal("signup"));
  modals.forEach((modal) => {
    const modalCloseIcon = modal.querySelector(".modal__close-icon");
    modalCloseIcon?.addEventListener("click", closeAuthModals);
  });
};

auth();
