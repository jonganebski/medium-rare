import {
  followersPagefollowBtns,
  followersPagefollowingBtns,
} from "../elements/followers";
import {
  followersPageFollowBtnClick,
  followersPageFollowingBtnClick,
} from "../follow";

const init = () => {
  followersPagefollowBtns?.forEach((btn) => {
    btn.addEventListener("click", followersPageFollowBtnClick);
  });
  followersPagefollowingBtns?.forEach((btn) => {
    btn.addEventListener("click", followersPageFollowingBtnClick);
  });
};

init();
