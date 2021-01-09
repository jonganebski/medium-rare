export const mainHeaderEl = document.getElementById("main-header");
export const openSigninModalEl = mainHeaderEl?.querySelector(
  ".main-header__signin-anchor"
);
export const openSignupModalEl = mainHeaderEl?.querySelector(
  ".main-header__signup-btn"
);
export const saveStatusEl = mainHeaderEl?.querySelector(
  ".main-header__save-status"
) as HTMLElement | null;
export const publishBtn = mainHeaderEl?.querySelector(
  ".main-header__publish-btn"
) as HTMLButtonElement | null;
export const unpublishBtn = mainHeaderEl?.querySelector(
  ".main-header__unpublish-btn"
) as HTMLButtonElement | null;
export const avatarFrame = mainHeaderEl?.querySelector(
  ".main-header__avatar-frame"
);
export const usermenu = mainHeaderEl?.querySelector(".header-usermenu");
