export const mainHeaderEl = document.getElementById("main-header");
export const openSigninModalEl = mainHeaderEl?.querySelector(
  ".main-header__signin-anchor"
);
export const openSignupModalEl = mainHeaderEl?.querySelector(
  ".main-header__signup-btn"
);
export const publishBtn = mainHeaderEl?.querySelector(
  ".main-header__publish-btn"
);
const avatarFrame = mainHeaderEl?.querySelector(".main-header__avatar-frame");
const usermenu = mainHeaderEl?.querySelector(".header-usermenu");

const toggleUsermenu = () => {
  if (usermenu) {
    if (usermenu.className.includes("_none")) {
      usermenu.className = "header-usermenu _block";
      return;
    }
    if (usermenu.className.includes("_block")) {
      usermenu.className = "header-usermenu _none";
      return;
    }
  }
};

const useHeader = () => {
  avatarFrame?.addEventListener("click", toggleUsermenu);
};

useHeader();
