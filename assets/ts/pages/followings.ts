import { followingsPagefollowBtns } from "../elements/followings";
import { followingsPageUnfollowBtnClick } from "../follow";

const init = () => {
  followingsPagefollowBtns?.forEach((btn) => {
    btn.addEventListener("click", followingsPageUnfollowBtnClick);
  });
};

init();
