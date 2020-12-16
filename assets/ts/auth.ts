import Axios from "axios";
import { openSigninModalEl, openSignupModalEl } from "./elements.header";

const filter = document.getElementById("filter");

const modals = document.querySelectorAll(".modal");
const signinModal = document.getElementById(
  "signin-modal"
) as HTMLFormElement | null;
const signupModal = document.getElementById(
  "signup-modal"
) as HTMLFormElement | null;
const signinError = signinModal?.querySelector(
  ".modal__error-message"
) as HTMLSpanElement | null;
const signupEmailError = signupModal?.querySelector(
  ".signup-emailError"
) as HTMLSpanElement | null;
const signupPasswordError = signupModal?.querySelector(
  ".signup-passwordError"
) as HTMLSpanElement | null;

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

const signinUser = async (e: Event) => {
  e.preventDefault();
  if (signinModal) {
    const data = new FormData(signinModal);
    const { pathname } = document.location;
    try {
      const { status } = await Axios.post("/signin", data);
      if (status === 200) {
        document.location.href = pathname;
      }
    } catch (err) {
      const { status } = err.response;
      if (status === 400 && signinError) {
        signinError.innerText = "Wrong email or password";
      } else if (status === 500 && signinError) {
        signinError.innerText = "Sorry.. server has a problem";
      }
    }
  }
};

const signupUser = async (e: Event) => {
  e.preventDefault();
  if (signupModal) {
    const data = new FormData(signupModal);
    const { pathname } = document.location;
    const passwordLen = data.get("password")?.toString().length;
    if (!passwordLen) {
      return;
    }
    if (passwordLen < 6) {
      signupPasswordError &&
        (signupPasswordError.innerText =
          "Password must be at least 6 characters");
      return;
    } else {
      signupPasswordError && (signupPasswordError.innerText = "");
    }
    try {
      const { status } = await Axios.post("/signup", data);
      if (status === 201) {
        document.location.href = pathname;
      }
    } catch (err) {
      const { status, data: errMessage } = err.response;
      if (status === 400 || status === 500) {
        signupEmailError && (signupEmailError.innerText = errMessage);
      }
    }
  }
};

const init = () => {
  filter?.addEventListener("click", closeAuthModals);
  openSigninModalEl?.addEventListener("click", () => openAuthModal("signin"));
  openSignupModalEl?.addEventListener("click", () => openAuthModal("signup"));
  modals.forEach((modal) => {
    const modalCloseIcon = modal.querySelector(".modal__close-icon");
    modalCloseIcon?.addEventListener("click", closeAuthModals);
  });
  signinModal?.addEventListener("submit", signinUser);
  signupModal?.addEventListener("submit", signupUser);
};

init();
