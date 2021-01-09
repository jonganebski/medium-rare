import { usermenu, avatarFrame } from "./elements/header";

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

const init = () => {
  avatarFrame?.addEventListener("click", toggleUsermenu);
};

init();
