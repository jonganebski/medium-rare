import { followingsPagefollowBtns } from "./elements.followings";
import { followingsPageUnfollowBtnClick } from "./follow";

const initFollowings = () => {
  followingsPagefollowBtns?.forEach((btn) => {
    btn.addEventListener("click", followingsPageUnfollowBtnClick);
  });
};

initFollowings();
